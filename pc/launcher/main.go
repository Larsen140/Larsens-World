// Larsens World für Windows: eigenes App-Fenster (WebView2, in Windows eingebaut),
// lädt immer die aktuelle Version von GitHub. Daten liegen in %APPDATA%\LarsensWorld.
package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"

	webview2 "github.com/jchv/go-webview2"
)

const appURL = "https://larsen140.github.io/Larsens-World/"

func msg(title, text string) {
	u := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW")
	t, _ := syscall.UTF16PtrFromString(title)
	m, _ := syscall.UTF16PtrFromString(text)
	u.Call(0, uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(t)), 0x40)
}

func main() {
	data := filepath.Join(os.Getenv("APPDATA"), "LarsensWorld")
	_ = os.MkdirAll(data, 0o755)
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     false,
		AutoFocus: true,
		DataPath:  data,
		WindowOptions: webview2.WindowOptions{
			Title:  "Larsens World",
			Width:  1280,
			Height: 860,
			IconId: 1,
			Center: true,
		},
	})
	if w == nil {
		msg("Larsens World", "Für Larsens World fehlt die Microsoft WebView2-Komponente.\nSie öffnet sich jetzt zum Herunterladen. Danach Larsens World einfach nochmal starten.")
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", "https://go.microsoft.com/fwlink/p/?LinkId=2124703").Start()
		return
	}
	defer w.Destroy()
	w.Navigate(appURL)
	w.Run()
}
