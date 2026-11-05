using System;
using System.IO;
using System.Linq;
using System.Net.Security;
using System.Net.Sockets;
using System.Security.Authentication;
using System.Text;
using System.Threading;
using Avalonia;
using Avalonia.Controls;
using Avalonia.Input;
using Avalonia.Layout;
using Avalonia.Media;
using Avalonia.Media.Imaging;
using Avalonia.Platform;
using Avalonia.Threading;
using DotNetVnc;
using XenAPI;

namespace XenAdmin.Desktop
{
    /// <summary>Console tab: RFB over an authenticated, certificate-checked HTTP CONNECT tunnel to the VM's host.</summary>
    public class ConsoleView : DockPanel, IVNCGraphicsClient
    {
        public static CertificateTrust Trust;
        /// <summary>Raised (on the UI thread) when a tunnel was refused because its certificate needs the user's decision.</summary>
        public static event Action<ConsoleView> CertificatePending;

        private readonly VM _vm;
        private readonly Screen _screen;
        private readonly TextBlock _status = new TextBlock { VerticalAlignment = VerticalAlignment.Center, Margin = new Thickness(8, 0) };
        private readonly CheckBox _scale = new CheckBox { Content = "Scale", IsChecked = false, VerticalAlignment = VerticalAlignment.Center };
        private readonly object _fbLock = new object();
        private byte[] _fb = Array.Empty<byte>();
        private int _width, _height;
        private VNCStream _vnc;
        private Session _session;
        private int _buttons;
        private volatile bool _closed;
        private int _redrawQueued;

        public static Control For(VM vm) =>
            vm == null ? new TextBlock() : new ConsoleView(vm);

        private ConsoleView(VM vm)
        {
            _vm = vm;
            Background = Brushes.White;
            _screen = new Screen(this) { Focusable = true };

            // top: teal band with the DVD link and "Launch SSH Console", as in VNCTabView
            var top = new DockPanel { Height = 36, Background = Views.Gradient("#396D8C", "#3F8B89", true) };
            // VNCTabView: the console switch (Remote Desktop) at the far right, Launch SSH Console to its left
            _rdp.Click += async (_, _) => await SwitchToRdp();
            DockPanel.SetDock(_rdp, Dock.Right);
            top.Children.Add(_rdp);
            _ssh.Click += (_, _) => LaunchSsh();
            DockPanel.SetDock(_ssh, Dock.Right);
            top.Children.Add(_ssh);
            var dvd = Views.LinkText("Click here to create a DVD drive", null);
            dvd.Foreground = Brushes.White;
            dvd.FontWeight = FontWeight.Bold;
            dvd.TextDecorations = TextDecorations.Underline;
            dvd.HorizontalAlignment = HorizontalAlignment.Center;
            dvd.VerticalAlignment = VerticalAlignment.Center;
            top.Children.Add(dvd);

            // bottom: Send Ctrl+Alt+Del | Scale, Undock, Fullscreen
            var bottom = new DockPanel { Height = 34, Background = Brush.Parse("#F0F0F0") };
            var cad = new Button { Content = "Send Ctrl+Alt+Del (Ctrl+Alt+Insert)", Margin = new Thickness(4) };
            cad.Click += (_, _) => SendCtrlAltDel();
            DockPanel.SetDock(cad, Dock.Left);
            bottom.Children.Add(cad);
            var right = new StackPanel { Orientation = Orientation.Horizontal, HorizontalAlignment = HorizontalAlignment.Right, Spacing = 6, Margin = new Thickness(4) };
            right.Children.Add(_scale);
            var undock = new StackPanel { Orientation = Orientation.Horizontal, Spacing = 4 };
            undock.Children.Add(Views.Icon("detach_24", 16));
            undock.Children.Add(new TextBlock { Text = "Undock (Alt+Shift+U)", VerticalAlignment = VerticalAlignment.Center });
            right.Children.Add(new Button { Content = undock, IsEnabled = false });
            right.Children.Add(new Button { Content = "Fullscreen (Ctrl+Enter)", IsEnabled = false });
            bottom.Children.Add(right);

            var status = new Border { Child = _status, Height = 22 };
            DockPanel.SetDock(top, Dock.Top);
            DockPanel.SetDock(bottom, Dock.Bottom);
            DockPanel.SetDock(status, Dock.Bottom);
            Children.Add(top);
            Children.Add(bottom);
            Children.Add(status);
            Children.Add(_screen);
            _scale.IsCheckedChanged += (_, _) => _screen.InvalidateVisual();
            AttachedToVisualTree += (_, _) => { Start(); UpdateButtons(); _buttonsTimer.Start(); };
            DetachedFromVisualTree += (_, _) => { Stop(); _buttonsTimer.Stop(); };
            // the guest agent reports addresses and the RDP state later: follow its metrics while the tab is open
            _buttonsTimer.Tick += (_, _) => UpdateButtons();
        }

        private readonly Button _ssh = new Button { Content = "Launch SSH Console", Margin = new Thickness(6, 4), Background = Brushes.White };
        private readonly Button _rdp = new Button { Margin = new Thickness(6, 4, 0, 4), Background = Brushes.White };
        private volatile int _generation;
        private bool _rdpBusy; // the 3389 probe runs; WinForms keeps the button disabled meanwhile
        private readonly Avalonia.Threading.DispatcherTimer _buttonsTimer = new Avalonia.Threading.DispatcherTimer { Interval = TimeSpan.FromSeconds(2) };

        /// <summary>Set by the main window to run actions in its task list.</summary>
        public static Action<XenAdmin.Actions.AsyncAction> RunAction;

        private void UpdateButtons()
        {
            _ssh.IsVisible = ConsoleRules.SshVisible(_vm);
            _ssh.IsEnabled = ConsoleRules.SshEnabled(_vm);
            _rdp.IsVisible = ConsoleRules.RdpVisible(_vm);
            _rdp.IsEnabled = !_rdpBusy && ConsoleRules.RdpEnabled(_vm);
            _rdp.Content = Dialogs.M(ConsoleRules.RdpLabel(_vm));
        }

        /// <summary>VNCTabView.buttonSSH_Click: the desktop's terminal with OpenSSH instead of PuTTY.</summary>
        private void LaunchSsh()
        {
            if (!ConsoleRules.SshEnabled(_vm)) return;
            var target = ConsoleRules.SshTarget(_vm);
            if (target == null)
                _status.Text = "The address reported by the guest is not valid for SSH.";
            else if (!Platform.PlatformServices.Current.LaunchSsh(target))
                _status.Text = "No terminal with an ssh client was found (install OpenSSH and a terminal emulator, or set $TERMINAL).";
        }

        /// <summary>
        /// toggleConsoleButton (RDP): offer to turn RDP on in the guest, then open the desktop's Remote Desktop client on
        /// the first address answering on 3389 (XSVNCScreen.PollPort). There is no embedded RDP control here.
        /// </summary>
        private async System.Threading.Tasks.Task SwitchToRdp()
        {
            if (!ConsoleRules.RdpEnabled(_vm)) return;
            if (ConsoleRules.CanTurnOnRdp(_vm))
            {
                if (!await Dialogs.Confirm(TopLevel.GetTopLevel(this) as Window, "XCP-ng Center", Messages.FORCE_ENABLE_RDP, Messages.YES_BUTTON_CAPTION, Messages.NO_BUTTON_CAPTION))
                    return;
                RunAction?.Invoke(ConsoleRules.TurnOnRdp(_vm));
                _status.Text = "Asked the guest to turn on Remote Desktop; try again once it is enabled.";
                return;
            }
            _rdpBusy = true;
            UpdateButtons();
            var generation = _generation;
            var candidates = ConsoleRules.RdpAddresses(_vm);
            string open;
            try { open = await System.Threading.Tasks.Task.Run(() => candidates.FirstOrDefault(ip => PortOpen(ip.Trim('[', ']'), 3389))); }
            finally { _rdpBusy = false; UpdateButtons(); }
            if (_closed || generation != _generation)
                return; // the tab was left meanwhile
            if (open == null)
                _status.Text = "Remote Desktop is not answering on any address of this VM.";
            else if (!Platform.PlatformServices.Current.LaunchRdp(open))
                _status.Text = "No Remote Desktop client was found (install FreeRDP or Remmina).";
        }

        private static bool PortOpen(string host, int port)
        {
            try
            {
                using var client = new TcpClient();
                return client.ConnectAsync(host, port).Wait(TimeSpan.FromSeconds(2)) && client.Connected;
            }
            catch (Exception)
            {
                return false;
            }
        }

        public string VmName => _vm.Name();
        public string UUID => _vm.uuid;

        public void Retry()
        {
            Stop();
            Start();
        }

        private void Start()
        {
            // the view is reused when the tab comes back: each attempt has its generation, and only the current one may connect
            _closed = false;
            var generation = ++_generation;
            if (_vm.power_state != vm_power_state.Running)
            {
                _status.Text = "This VM is not running.";
                return;
            }
            var console = _vm.Connection.ResolveAll(_vm.consoles).FirstOrDefault(c => c.protocol == console_protocol.rfb);
            if (console == null)
            {
                _status.Text = "This VM has no graphical console.";
                return;
            }
            _status.Text = "Connecting...";
            var location = console.location;
            ThreadPool.QueueUserWorkItem(_ =>
            {
                try
                {
                    var session = _vm.Connection.DuplicateSession();
                    var stream = OpenTunnel(new Uri(location), session.opaque_ref);
                    if (_closed || generation != _generation) { stream.Dispose(); return; }
                    _session = session;
                    _vnc = new VNCStream(this, stream, false);
                    _vnc.ErrorOccurred += (_, e) => { if (KeyLog) System.Console.Error.WriteLine("console vnc error " + e); };
                    _vnc.ErrorOccurred += (_, e) => Ui(() => _status.Text = _closed ? "" : "Disconnected: " + e.Message);
                    _vnc.ConnectionSuccess += (_, _) => Ui(() => _status.Text = $"Connected to {VmName}");
                    _vnc.Connect(null);
                }
                catch (AuthenticationException)
                {
                    Ui(() => { _status.Text = "Certificate not trusted."; CertificatePending?.Invoke(this); });
                }
                catch (Exception e)
                {
                    Ui(() => _status.Text = "Unable to connect to the console: " + e.Message);
                }
            });
        }

        private void Stop()
        {
            _closed = true;
            try { _vnc?.Close(); } catch { }
            _vnc = null;
            // never log out: DuplicateSession shares the connection's session ref, so logout would drop the whole
            // connection (XSVNCScreen only forgets it too)
            _session = null;
        }

        /// <summary>CONNECT tunnel as in XenAPI.HTTP.HttpConnectStream, but with the app's certificate policy.</summary>
        private static Stream OpenTunnel(Uri uri, string sessionRef)
        {
            var port = uri.IsDefaultPort ? 443 : uri.Port;
            var tcp = new TcpClient(uri.Host, port) { NoDelay = true };
            var ssl = new SslStream(tcp.GetStream(), false, (_, cert, _, errors) => Trust.ValidateHost($"{uri.Host}:{port}", cert, errors));
            ssl.AuthenticateAsClient(uri.Host);
            var request = $"CONNECT {uri.PathAndQuery} HTTP/1.0\r\nHost: {uri.Host}\r\nCookie: session_id={sessionRef}\r\n\r\n";
            ssl.Write(Encoding.ASCII.GetBytes(request));
            ssl.Flush();
            var status = ReadLine(ssl);
            while (ReadLine(ssl).Length > 0) { }
            if (!status.Contains(" 200"))
                throw new IOException("Console tunnel refused: " + status);
            return ssl;
        }

        private static string ReadLine(Stream s)
        {
            var sb = new StringBuilder();
            int c;
            while ((c = s.ReadByte()) >= 0 && c != '\n')
                if (c != '\r') sb.Append((char)c);
            if (c < 0) throw new EndOfStreamException("Console tunnel closed");
            return sb.ToString();
        }

        private static void Ui(Action a) => Dispatcher.UIThread.Post(a);

        // ---------- IVNCGraphicsClient (called on the VNC thread) ----------

        public void ClientBell() { }
        public void ClientCutText(string text) { }
        public void ClientSetCursor(DotNetVnc.Bitmap image, int x, int y, int width, int height) { }

        public void ClientDesktopSize(int width, int height)
        {
            if (KeyLog) System.Console.Error.WriteLine($"console vnc desktop size {width}x{height}");
            lock (_fbLock)
            {
                _width = width;
                _height = height;
                _fb = new byte[width * height * 4];
            }
            Ui(() => _screen.Resize(width, height));
        }

        public void ClientDrawImage(DotNetVnc.Bitmap image, int x, int y, int width, int height)
        {
            lock (_fbLock)
            {
                for (int row = 0; row < height && y + row < _height; row++)
                {
                    var n = Math.Min(width, _width - x) * 4;
                    if (n > 0)
                        Buffer.BlockCopy(image.Bgra, row * width * 4, _fb, ((y + row) * _width + x) * 4, n);
                }
            }
        }

        public void ClientFillRectangle(int x, int y, int width, int height, DotNetVnc.Color color)
        {
            lock (_fbLock)
            {
                for (int row = y; row < Math.Min(y + height, _height); row++)
                    for (int col = x; col < Math.Min(x + width, _width); col++)
                    {
                        var o = (row * _width + col) * 4;
                        _fb[o] = color.B; _fb[o + 1] = color.G; _fb[o + 2] = color.R; _fb[o + 3] = 255;
                    }
            }
        }

        public void ClientCopyRectangle(int x, int y, int width, int height, int dx, int dy)
        {
            lock (_fbLock)
            {
                var tmp = new byte[width * height * 4];
                for (int row = 0; row < height; row++)
                    Buffer.BlockCopy(_fb, ((y + row) * _width + x) * 4, tmp, row * width * 4, width * 4);
                for (int row = 0; row < height; row++)
                    Buffer.BlockCopy(tmp, row * width * 4, _fb, ((dy + row) * _width + dx) * 4, width * 4);
            }
        }

        public void ClientFrameBufferUpdate()
        {
            // one repaint in flight at a time; the next update catches up
            if (Interlocked.Exchange(ref _redrawQueued, 1) == 0)
                Dispatcher.UIThread.Post(() => { _redrawQueued = 0; _screen.Upload(); }, DispatcherPriority.Render);
        }

        // ---------- input ----------

        private void SendCtrlAltDel()
        {
            var vnc = _vnc;
            if (vnc == null) return;
            if (_vm.IsHVM())
            {
                foreach (var (k, sym) in new[] { (0x1d, 0xffe3), (0x38, 0xffe9), (0xd3, 0xffff) }) vnc.keyScanEvent(true, k, sym, true);
                foreach (var (k, sym) in new[] { (0xd3, 0xffff), (0x38, 0xffe9), (0x1d, 0xffe3) }) vnc.keyScanEvent(false, k, sym, true);
            }
            else
            {
                foreach (var k in new[] { 0xffe3, 0xffe9, 0xffff }) vnc.keyCodeEvent(true, k);
                foreach (var k in new[] { 0xffff, 0xffe9, 0xffe3 }) vnc.keyCodeEvent(false, k);
            }
            _screen.Focus();
        }

        private void Pointer(PointerEventArgs e)
        {
            var p = _screen.ToFramebuffer(e.GetPosition(_screen));
            var props = e.GetCurrentPoint(_screen).Properties;
            _buttons = (props.IsLeftButtonPressed ? 1 : 0) | (props.IsMiddleButtonPressed ? 2 : 0) | (props.IsRightButtonPressed ? 4 : 0);
            _vnc?.PointerEvent(_buttons, p.X, p.Y);
        }

        private void Key(KeyEventArgs e, bool down)
        {
            if (down && e.Key == Avalonia.Input.Key.Insert && e.KeyModifiers.HasFlag(KeyModifiers.Control) && e.KeyModifiers.HasFlag(KeyModifiers.Alt))
            {
                SendCtrlAltDel();
                e.Handled = true;
                return;
            }
            var sym = KeySyms.Track(_pressed, e, down);
            // HVM (QEMU): scancodes like XSVNCScreen.SendScanCodes, so the guest's layout picks the character;
            // the PV text console (vncterm) only understands keysyms
            var scan = _vm.IsHVM() && !ForceKeysyms ? KeySyms.Scancode(e.PhysicalKey) : 0;
            if (KeyLog)
                System.Console.Error.WriteLine($"console key {(down ? "down" : "up  ")} key={e.Key} physical={e.PhysicalKey} symbol={(e.KeySymbol == null ? "null" : "'" + e.KeySymbol + "'")} modifiers={e.KeyModifiers} -> scancode 0x{scan:x} keysym 0x{sym:x}");
            if (scan != 0)
            {
                _vnc?.keyScanEvent(down, scan, sym, true);
                e.Handled = true;
            }
            else if (sym != 0)
            {
                _vnc?.keyCodeEvent(down, sym);
                e.Handled = true;
            }
        }

        private readonly System.Collections.Generic.Dictionary<(Avalonia.Input.PhysicalKey, Avalonia.Input.Key), int> _pressed = new();

        // XCP_CONSOLE_KEYLOG=1 prints every key event and the keysym sent, to diagnose keyboard layouts.
        // Debug builds only: it would print passwords typed in the guests.
#if DEBUG
        private static readonly bool KeyLog = Environment.GetEnvironmentVariable("XCP_CONSOLE_KEYLOG") == "1";
#else
        private const bool KeyLog = false;
#endif
        // XCP_CONSOLE_KEYSYMS=1 sends keysyms to HVM guests too (diagnostics; the guest then depends on QEMU's keymap)
        private static readonly bool ForceKeysyms = Environment.GetEnvironmentVariable("XCP_CONSOLE_KEYSYMS") == "1";

        /// <summary>Draws the framebuffer, scaled to fit (keeping the aspect ratio) or at real size.</summary>
        private class Screen : Control
        {
            private readonly ConsoleView _owner;
            private WriteableBitmap _bitmap;

            public Screen(ConsoleView owner)
            {
                _owner = owner;
                ClipToBounds = true;
                PointerMoved += (_, e) => _owner.Pointer(e);
                PointerPressed += (_, e) => { Focus(); _owner.Pointer(e); };
                PointerReleased += (_, e) => _owner.Pointer(e);
                PointerWheelChanged += (_, e) =>
                {
                    var p = ToFramebuffer(e.GetPosition(this));
                    _owner._vnc?.PointerWheelEvent(p.X, p.Y, e.Delta.Y > 0 ? -1 : 1);
                };
                KeyDown += (_, e) => _owner.Key(e, true);
                KeyUp += (_, e) => _owner.Key(e, false);
            }

            public void Resize(int w, int h)
            {
                _bitmap = w > 0 && h > 0 ? new WriteableBitmap(new PixelSize(w, h), new Vector(96, 96), Avalonia.Platform.PixelFormat.Bgra8888, AlphaFormat.Opaque) : null;
                InvalidateVisual();
            }

            public void Upload()
            {
                if (_bitmap == null) return;
                using (var fb = _bitmap.Lock())
                lock (_owner._fbLock)
                {
                    if (_owner._fb.Length == fb.RowBytes * fb.Size.Height)
                        System.Runtime.InteropServices.Marshal.Copy(_owner._fb, 0, fb.Address, _owner._fb.Length);
                }
                InvalidateVisual();
            }

            private Rect Target()
            {
                if (_bitmap == null) return default;
                var w = _bitmap.PixelSize.Width;
                var h = _bitmap.PixelSize.Height;
                if (_owner._scale.IsChecked != true)
                    return new Rect(Math.Max(0, (Bounds.Width - w) / 2), Math.Max(0, (Bounds.Height - h) / 2), w, h);
                var k = Math.Min(Bounds.Width / w, Bounds.Height / h);
                return new Rect((Bounds.Width - w * k) / 2, (Bounds.Height - h * k) / 2, w * k, h * k);
            }

            public PixelPoint ToFramebuffer(Point p)
            {
                var t = Target();
                if (_bitmap == null || t.Width <= 0) return default;
                var x = (int)((p.X - t.X) * _bitmap.PixelSize.Width / t.Width);
                var y = (int)((p.Y - t.Y) * _bitmap.PixelSize.Height / t.Height);
                return new PixelPoint(Math.Clamp(x, 0, _bitmap.PixelSize.Width - 1), Math.Clamp(y, 0, _bitmap.PixelSize.Height - 1));
            }

            public override void Render(DrawingContext context)
            {
                context.FillRectangle(Brushes.White, new Rect(Bounds.Size));
                if (_bitmap == null) return;
                var t = Target();
                context.DrawImage(_bitmap, new Rect(_bitmap.Size), t);
                // focus frame around the screen, as the WinForms console draws it
                context.DrawRectangle(new Pen(Brush.Parse(IsFocused ? "#3399FF" : "#7A7A7A"), 1), t.Inflate(1));
            }
        }
    }

    /// <summary>Avalonia key → X11 keysym (what RFB key events carry).</summary>
    public static class KeySyms
    {
        /// <summary>
        /// The keysym to send for a press or release. A release repeats its press's keysym: the symbol can differ on key up
        /// (Shift released first turns '&lt;' into ',', or X11 reports none), which left keys stuck in the guest.
        /// </summary>
        public static int Track(System.Collections.Generic.Dictionary<(Avalonia.Input.PhysicalKey, Avalonia.Input.Key), int> pressed, KeyEventArgs e, bool down)
        {
            var id = (e.PhysicalKey, e.Key);
            if (down)
            {
                var sym = From(e);
                if (sym != 0) pressed[id] = sym;
                return sym;
            }
            return pressed.Remove(id, out var held) ? held : From(e);
        }

        public static int From(KeyEventArgs e)
        {
            switch (e.Key)
            {
                case Avalonia.Input.Key.Back: return 0xff08;
                case Avalonia.Input.Key.Tab: return 0xff09;
                case Avalonia.Input.Key.Enter: return 0xff0d;
                case Avalonia.Input.Key.Escape: return 0xff1b;
                case Avalonia.Input.Key.Delete: return 0xffff;
                case Avalonia.Input.Key.Home: return 0xff50;
                case Avalonia.Input.Key.Left: return 0xff51;
                case Avalonia.Input.Key.Up: return 0xff52;
                case Avalonia.Input.Key.Right: return 0xff53;
                case Avalonia.Input.Key.Down: return 0xff54;
                case Avalonia.Input.Key.PageUp: return 0xff55;
                case Avalonia.Input.Key.PageDown: return 0xff56;
                case Avalonia.Input.Key.End: return 0xff57;
                case Avalonia.Input.Key.Insert: return 0xff63;
                case Avalonia.Input.Key.LeftShift: return 0xffe1;
                case Avalonia.Input.Key.RightShift: return 0xffe2;
                case Avalonia.Input.Key.LeftCtrl: return 0xffe3;
                case Avalonia.Input.Key.RightCtrl: return 0xffe4;
                case Avalonia.Input.Key.LeftAlt: return 0xffe9;
                case Avalonia.Input.Key.RightAlt: return 0xfe03; // AltGr
                case Avalonia.Input.Key.LWin: return 0xffeb;
                case Avalonia.Input.Key.CapsLock: return 0xffe5;
                case Avalonia.Input.Key.Space: return 0x20;
            }
            if (e.Key >= Avalonia.Input.Key.F1 && e.Key <= Avalonia.Input.Key.F12)
                return 0xffbe + (e.Key - Avalonia.Input.Key.F1);
            // printable keys: the symbol the user's layout produced (ABNT2 accents included)
            var s = e.KeySymbol;
            if (!string.IsNullOrEmpty(s) && s.Length == 1 && !char.IsControl(s[0]))
                return s[0] < 0x100 ? s[0] : 0x01000000 | s[0];
            // X11 reports no symbol for some shifted keys (Shift+5): use the US character of the physical key
            var us = UsCharacter(e.PhysicalKey, e.KeyModifiers.HasFlag(KeyModifiers.Shift));
            return us == '\0' ? 0 : us;
        }

        private const string UsRow = "`1234567890-=[]\\;',./", UsShiftRow = "~!@#$%^&*()_+{}|:\"<>?";

        private static char UsCharacter(PhysicalKey k, bool shift)
        {
            if (k >= PhysicalKey.A && k <= PhysicalKey.Z)
                return (char)((shift ? 'A' : 'a') + (k - PhysicalKey.A));
            var i = k switch
            {
                PhysicalKey.Backquote => 0, PhysicalKey.Digit1 => 1, PhysicalKey.Digit2 => 2, PhysicalKey.Digit3 => 3, PhysicalKey.Digit4 => 4,
                PhysicalKey.Digit5 => 5, PhysicalKey.Digit6 => 6, PhysicalKey.Digit7 => 7, PhysicalKey.Digit8 => 8, PhysicalKey.Digit9 => 9,
                PhysicalKey.Digit0 => 10, PhysicalKey.Minus => 11, PhysicalKey.Equal => 12, PhysicalKey.BracketLeft => 13, PhysicalKey.BracketRight => 14,
                PhysicalKey.Backslash => 15, PhysicalKey.Semicolon => 16, PhysicalKey.Quote => 17, PhysicalKey.Comma => 18, PhysicalKey.Period => 19,
                PhysicalKey.Slash => 20, _ => -1
            };
            return i < 0 ? '\0' : (shift ? UsShiftRow : UsRow)[i];
        }

        /// <summary>
        /// The XT (set 1) scancode of a physical key, as QEMU's extended key event wants it: 0xE0-prefixed keys
        /// are sent with the high bit set (Right Ctrl 0xE0 0x1D = 0x9D). 0 when the key has none.
        /// </summary>
        public static int Scancode(PhysicalKey k)
        {
            if (k >= PhysicalKey.F1 && k <= PhysicalKey.F10) return 0x3b + (k - PhysicalKey.F1);
            return k switch
            {
                PhysicalKey.Escape => 0x01, PhysicalKey.Digit1 => 0x02, PhysicalKey.Digit2 => 0x03, PhysicalKey.Digit3 => 0x04, PhysicalKey.Digit4 => 0x05,
                PhysicalKey.Digit5 => 0x06, PhysicalKey.Digit6 => 0x07, PhysicalKey.Digit7 => 0x08, PhysicalKey.Digit8 => 0x09, PhysicalKey.Digit9 => 0x0a,
                PhysicalKey.Digit0 => 0x0b, PhysicalKey.Minus => 0x0c, PhysicalKey.Equal => 0x0d, PhysicalKey.Backspace => 0x0e, PhysicalKey.Tab => 0x0f,
                PhysicalKey.Q => 0x10, PhysicalKey.W => 0x11, PhysicalKey.E => 0x12, PhysicalKey.R => 0x13, PhysicalKey.T => 0x14, PhysicalKey.Y => 0x15,
                PhysicalKey.U => 0x16, PhysicalKey.I => 0x17, PhysicalKey.O => 0x18, PhysicalKey.P => 0x19, PhysicalKey.BracketLeft => 0x1a,
                PhysicalKey.BracketRight => 0x1b, PhysicalKey.Enter => 0x1c, PhysicalKey.ControlLeft => 0x1d, PhysicalKey.A => 0x1e, PhysicalKey.S => 0x1f,
                PhysicalKey.D => 0x20, PhysicalKey.F => 0x21, PhysicalKey.G => 0x22, PhysicalKey.H => 0x23, PhysicalKey.J => 0x24, PhysicalKey.K => 0x25,
                PhysicalKey.L => 0x26, PhysicalKey.Semicolon => 0x27, PhysicalKey.Quote => 0x28, PhysicalKey.Backquote => 0x29, PhysicalKey.ShiftLeft => 0x2a,
                PhysicalKey.Backslash => 0x2b, PhysicalKey.Z => 0x2c, PhysicalKey.X => 0x2d, PhysicalKey.C => 0x2e, PhysicalKey.V => 0x2f, PhysicalKey.B => 0x30,
                PhysicalKey.N => 0x31, PhysicalKey.M => 0x32, PhysicalKey.Comma => 0x33, PhysicalKey.Period => 0x34, PhysicalKey.Slash => 0x35,
                PhysicalKey.ShiftRight => 0x36, PhysicalKey.NumPadMultiply => 0x37, PhysicalKey.AltLeft => 0x38, PhysicalKey.Space => 0x39,
                PhysicalKey.CapsLock => 0x3a, PhysicalKey.NumLock => 0x45, PhysicalKey.ScrollLock => 0x46, PhysicalKey.NumPad7 => 0x47,
                PhysicalKey.NumPad8 => 0x48, PhysicalKey.NumPad9 => 0x49, PhysicalKey.NumPadSubtract => 0x4a, PhysicalKey.NumPad4 => 0x4b,
                PhysicalKey.NumPad5 => 0x4c, PhysicalKey.NumPad6 => 0x4d, PhysicalKey.NumPadAdd => 0x4e, PhysicalKey.NumPad1 => 0x4f,
                PhysicalKey.NumPad2 => 0x50, PhysicalKey.NumPad3 => 0x51, PhysicalKey.NumPad0 => 0x52, PhysicalKey.NumPadDecimal => 0x53,
                PhysicalKey.IntlBackslash => 0x56, PhysicalKey.F11 => 0x57, PhysicalKey.F12 => 0x58, PhysicalKey.IntlRo => 0x73, PhysicalKey.IntlYen => 0x7d,
                PhysicalKey.NumPadComma => 0x7e,
                // 0xE0-prefixed
                PhysicalKey.NumPadEnter => 0x9c, PhysicalKey.ControlRight => 0x9d, PhysicalKey.NumPadDivide => 0xb5, PhysicalKey.PrintScreen => 0xb7,
                PhysicalKey.AltRight => 0xb8, PhysicalKey.Home => 0xc7, PhysicalKey.ArrowUp => 0xc8, PhysicalKey.PageUp => 0xc9,
                PhysicalKey.ArrowLeft => 0xcb, PhysicalKey.ArrowRight => 0xcd, PhysicalKey.End => 0xcf, PhysicalKey.ArrowDown => 0xd0,
                PhysicalKey.PageDown => 0xd1, PhysicalKey.Insert => 0xd2, PhysicalKey.Delete => 0xd3, PhysicalKey.MetaLeft => 0xdb,
                PhysicalKey.MetaRight => 0xdc, PhysicalKey.ContextMenu => 0xdd,
                _ => 0
            };
        }
    }
}
