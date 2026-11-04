using System;
using System.Diagnostics;
using System.Runtime.InteropServices;

// shortcut: the WinForms VNC protocol sources (XenAdmin/VNC/VNCStream.cs etc.) are compiled here unchanged;
// these types live in their namespace so they win over the System.Drawing / XenCenterLib.Win32 imports,
// which are Windows-only at runtime. Replace with the shared XenAdmin.Vnc library (openspec task 5.1).
namespace DotNetVnc
{
    public enum PixelFormat { Format32bppRgb, Format32bppArgb, Format16bppRgb565, Format16bppRgb555, Format16bppArgb1555 }

    public readonly struct Color
    {
        public readonly byte A, R, G, B;
        private Color(byte a, byte r, byte g, byte b) { A = a; R = r; G = g; B = b; }
        public static Color FromArgb(int r, int g, int b) => new Color(255, (byte)r, (byte)g, (byte)b);
        public static Color White => new Color(255, 255, 255, 255);
        public static Color Black => new Color(255, 0, 0, 0);
    }

    /// <summary>Pixels converted to BGRA (premultiplied for opaque formats) at construction.</summary>
    public sealed class Bitmap : IDisposable
    {
        public int Width { get; }
        public int Height { get; }
        public byte[] Bgra { get; }

        public Bitmap(int width, int height, int stride, PixelFormat format, IntPtr scan0)
        {
            Width = width;
            Height = height;
            Bgra = new byte[width * height * 4];
            var row = new byte[Math.Abs(stride)];
            for (int y = 0; y < height; y++)
            {
                Marshal.Copy(scan0 + y * stride, row, 0, row.Length);
                int o = y * width * 4;
                for (int x = 0; x < width; x++, o += 4)
                {
                    switch (format)
                    {
                        case PixelFormat.Format32bppRgb:
                            Bgra[o] = row[x * 4]; Bgra[o + 1] = row[x * 4 + 1]; Bgra[o + 2] = row[x * 4 + 2]; Bgra[o + 3] = 255;
                            break;
                        case PixelFormat.Format32bppArgb:
                            Buffer.BlockCopy(row, x * 4, Bgra, o, 4);
                            break;
                        default:
                            int p = row[x * 2] | row[x * 2 + 1] << 8;
                            bool is565 = format == PixelFormat.Format16bppRgb565;
                            int r = is565 ? (p >> 11) & 31 : (p >> 10) & 31, g = is565 ? (p >> 5) & 63 : (p >> 5) & 31, b = p & 31;
                            Bgra[o] = (byte)(b * 255 / 31);
                            Bgra[o + 1] = (byte)(g * 255 / (is565 ? 63 : 31));
                            Bgra[o + 2] = (byte)(r * 255 / 31);
                            Bgra[o + 3] = format == PixelFormat.Format16bppArgb1555 && (p & 0x8000) == 0 ? (byte)0 : (byte)255;
                            break;
                    }
                }
            }
        }

        public void Dispose() { }
    }

    internal static class Win32
    {
        public static bool QueryPerformanceFrequency(out long x) { x = Stopwatch.Frequency; return true; }
        public static bool QueryPerformanceCounter(out long x) { x = Stopwatch.GetTimestamp(); return true; }
    }
}
