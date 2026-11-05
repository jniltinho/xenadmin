//go:build linux && !nogui && cgo

package gui

/*
#cgo pkg-config: gtk4 gtk4-x11 x11
#include <gtk/gtk.h>
#include <gdk/x11/gdkx.h>
#include <X11/Xatom.h>
#include <stdlib.h>

// The window icon. GTK 4 has no per-window icon and, on X11, does not write _NET_WM_ICON from the icon theme, so the
// window manager would show a generic icon: the icon is written to the X11 window when it is mapped (as Avalonia
// does for the C# app). On Wayland the compositor takes it from the .desktop file matching the app id.
static unsigned long *xcp_icon;
static int xcp_icon_len;

static void xcp_apply(GtkWidget *widget, gpointer unused) {
	GdkSurface *s = gtk_native_get_surface(GTK_NATIVE(widget));
	if (s == NULL || !GDK_IS_X11_SURFACE(s) || xcp_icon == NULL)
		return;
	Display *dpy = GDK_DISPLAY_XDISPLAY(gdk_surface_get_display(s));
	XChangeProperty(dpy, gdk_x11_surface_get_xid(s), XInternAtom(dpy, "_NET_WM_ICON", False), XA_CARDINAL, 32,
		PropModeReplace, (unsigned char *)xcp_icon, xcp_icon_len);
}

static char *xcp_icon_name;

static void xcp_watch(GtkWidget *w) {
	gtk_window_set_icon_name(GTK_WINDOW(w), xcp_icon_name);
	g_signal_connect(w, "map", G_CALLBACK(xcp_apply), NULL);
	if (gtk_widget_get_mapped(w))
		xcp_apply(w, NULL);
}

// windows Wails creates after the application started get the icon too
static void xcp_toplevels_changed(GListModel *list, guint position, guint removed, guint added, gpointer unused) {
	for (guint i = position; i < position + added; i++) {
		GtkWidget *w = g_list_model_get_item(list, i);
		xcp_watch(w);
		g_object_unref(w);
	}
}

static void xcp_set_icon(unsigned long *data, int len, const char *dir, const char *name) {
	xcp_icon_name = g_strdup(name);
	xcp_icon = data;
	xcp_icon_len = len;
	GdkDisplay *d = gdk_display_get_default();
	if (d != NULL)
		gtk_icon_theme_add_search_path(gtk_icon_theme_get_for_display(d), dir);
	gtk_window_set_default_icon_name(name);
	GListModel *wins = gtk_window_get_toplevels();
	for (guint i = 0; i < g_list_model_get_n_items(wins); i++) {
		GtkWidget *w = g_list_model_get_item(wins, i);
		xcp_watch(w);
		g_object_unref(w);
	}
	g_signal_connect(wins, "items-changed", G_CALLBACK(xcp_toplevels_changed), NULL);
}
*/
import "C"

import (
	"bytes"
	"image"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

const iconName = "xcp-ng-center-go"

// setWindowIcon gives the windows the XCP-ng Center icon (icons/app.png of the embedded assets): as _NET_WM_ICON on
// X11, and by name from a copy in the user's cache directory (nothing is installed in the system).
func setWindowIcon(app *application.App, assets fs.FS) {
	data, err := fs.ReadFile(assets, "icons/app.png")
	if err != nil {
		return
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return
	}
	cardinals := netWMIcon(img, 32, 48, 64)
	dir := ""
	if base, err := os.UserCacheDir(); err == nil {
		dir = filepath.Join(base, "xcp-ng-center-go", "icons")
		if os.MkdirAll(dir, 0o700) == nil {
			_ = os.WriteFile(filepath.Join(dir, iconName+".png"), data, 0o600)
		}
	}
	// C memory: X keeps no pointer, but GTK calls xcp_apply later, so the array must outlive this function
	buf := (*C.ulong)(C.malloc(C.size_t(len(cardinals)) * C.size_t(unsafe.Sizeof(C.ulong(0)))))
	copy(unsafe.Slice(buf, len(cardinals)), cardinals)
	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
		application.InvokeSync(func() {
			cdir, cname := C.CString(dir), C.CString(iconName)
			defer C.free(unsafe.Pointer(cdir))
			defer C.free(unsafe.Pointer(cname))
			C.xcp_set_icon(buf, C.int(len(cardinals)), cdir, cname)
		})
	})
}

// netWMIcon is the _NET_WM_ICON value: for each size, width, height, then ARGB pixels (an X "long" each), the image
// scaled by sampling the area under each pixel.
func netWMIcon(img image.Image, sizes ...int) []C.ulong {
	var out []C.ulong
	b := img.Bounds()
	for _, n := range sizes {
		out = append(out, C.ulong(n), C.ulong(n))
		for y := 0; y < n; y++ {
			for x := 0; x < n; x++ {
				var r, g, bl, a, count uint32
				x0, x1 := b.Min.X+x*b.Dx()/n, b.Min.X+(x+1)*b.Dx()/n
				y0, y1 := b.Min.Y+y*b.Dy()/n, b.Min.Y+(y+1)*b.Dy()/n
				for sy := y0; sy < max(y1, y0+1); sy++ {
					for sx := x0; sx < max(x1, x0+1); sx++ {
						cr, cg, cb, ca := img.At(sx, sy).RGBA() // premultiplied 16-bit
						r, g, bl, a, count = r+cr, g+cg, bl+cb, a+ca, count+1
					}
				}
				r, g, bl, a = r/count>>8, g/count>>8, bl/count>>8, a/count>>8
				if a > 0 { // un-premultiply: X wants straight alpha
					r, g, bl = min(255, r*255/a), min(255, g*255/a), min(255, bl*255/a)
				}
				out = append(out, C.ulong(a<<24|r<<16|g<<8|bl))
			}
		}
	}
	return out
}
