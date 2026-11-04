using System;
using System.Collections.Generic;
using System.Linq;
using Avalonia;
using Avalonia.Media.Imaging;
using Avalonia.Platform;

namespace XenAdmin.Desktop
{
    /// <summary>The WinForms app's own PNGs (XenAdmin/Images), embedded as Avalonia resources.</summary>
    public static class Images
    {
        private static readonly Dictionary<string, Bitmap> Cache = new Dictionary<string, Bitmap>();

        public static Bitmap Get(string file)
        {
            lock (Cache)
            {
                if (!Cache.TryGetValue(file, out var bmp))
                {
                    var uri = file == "Logo" ? "avares://xcp-ng-center/Assets/xcp-ng-center.png" : $"avares://xcp-ng-center/Images/{file}.png";
                    try { bmp = new Bitmap(AssetLoader.Open(new Uri(uri))); }
                    catch { bmp = null; }
                    Cache[file] = bmp;
                }
                return bmp;
            }
        }

        private static readonly Dictionary<string, Bitmap> GrayCache = new Dictionary<string, Bitmap>();

        /// <summary>Greyscale copy, as ToolStrip renders disabled buttons.</summary>
        public static Bitmap Gray(string file)
        {
            lock (GrayCache)
            {
                if (GrayCache.TryGetValue(file, out var gray))
                    return gray;
                var src = Get(file);
                if (src == null)
                    return null;
                var size = src.PixelSize;
                var wb = new WriteableBitmap(size, src.Dpi, Avalonia.Platform.PixelFormat.Bgra8888, AlphaFormat.Unpremul);
                using (var fb = wb.Lock())
                {
                    src.CopyPixels(new PixelRect(size), fb.Address, fb.RowBytes * size.Height, fb.RowBytes);
                    var px = new byte[fb.RowBytes * size.Height];
                    System.Runtime.InteropServices.Marshal.Copy(fb.Address, px, 0, px.Length);
                    for (int i = 0; i < px.Length; i += 4)
                    {
                        // lighter grey, like ControlPaint.DrawImageDisabled
                        var l = (byte)Math.Min(255, 70 + (px[i] * 11 + px[i + 1] * 59 + px[i + 2] * 30) / 100 * 2 / 3);
                        px[i] = px[i + 1] = px[i + 2] = l;
                    }
                    System.Runtime.InteropServices.Marshal.Copy(px, 0, fb.Address, px.Length);
                }
                return GrayCache[file] = wb;
            }
        }

    }
}
