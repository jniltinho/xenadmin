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
            var ssh = new Button { Content = "Launch SSH Console", IsEnabled = false, Margin = new Thickness(6, 4), Background = Brushes.White };
            DockPanel.SetDock(ssh, Dock.Right);
            top.Children.Add(ssh);
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
            AttachedToVisualTree += (_, _) => Start();
            DetachedFromVisualTree += (_, _) => Stop();
        }

        public string VmName => _vm.Name();
        public string UUID => _vm.uuid;

        public void Retry()
        {
            Stop();
            _closed = false;
            Start();
        }

        private void Start()
        {
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
                    _session = _vm.Connection.DuplicateSession();
                    var stream = OpenTunnel(new Uri(location), _session.opaque_ref);
                    if (_closed) { stream.Dispose(); return; }
                    _vnc = new VNCStream(this, stream, false);
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
            foreach (var k in new[] { 0xffe3, 0xffe9, 0xffff }) vnc.keyCodeEvent(true, k);
            foreach (var k in new[] { 0xffff, 0xffe9, 0xffe3 }) vnc.keyCodeEvent(false, k);
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
            if (sym != 0)
            {
                _vnc?.keyCodeEvent(down, sym);
                e.Handled = true;
            }
        }

        private readonly System.Collections.Generic.Dictionary<(Avalonia.Input.PhysicalKey, Avalonia.Input.Key), int> _pressed = new();

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
            return 0;
        }
    }
}
