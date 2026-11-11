// Larsens World für Windows: installiert sich beim ersten Start selbst
// (%LOCALAPPDATA%\Programs\Larsens World, Verknüpfungen auf Desktop und im Startmenü,
// Eintrag unter "Apps & Features") und öffnet dann ein eigenes App-Fenster (WebView2),
// das immer die aktuelle Version von GitHub lädt. Daten: %APPDATA%\LarsensWorld.
package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	webview2 "github.com/jchv/go-webview2"
	"golang.org/x/sys/windows/registry"
)

const (
	appURL  = "https://larsen140.github.io/Larsens-World/?app=1"
	appName = "Larsens World"
	uninKey = `Software\Microsoft\Windows\CurrentVersion\Uninstall\LarsensWorld`
)

var user32 = syscall.NewLazyDLL("user32.dll")

func msg(text string, flags uintptr) int {
	t, _ := syscall.UTF16PtrFromString(appName)
	m, _ := syscall.UTF16PtrFromString(text)
	r, _, _ := user32.NewProc("MessageBoxW").Call(0, uintptr(unsafe.Pointer(m)), uintptr(unsafe.Pointer(t)), flags)
	return int(r)
}

// schwarze Titelleiste (Windows 10: Dark Mode, Windows 11: komplett schwarz mit weißer Schrift)
func darkTitle(hwnd uintptr) {
	dwm := syscall.NewLazyDLL("dwmapi.dll").NewProc("DwmSetWindowAttribute")
	if dwm.Find() != nil || hwnd == 0 {
		return
	}
	set := func(attr uintptr, v uint32) { dwm.Call(hwnd, attr, uintptr(unsafe.Pointer(&v)), 4) }
	set(20, 1)          // DWMWA_USE_IMMERSIVE_DARK_MODE
	set(19, 1)          // ältere Windows-10-Versionen
	set(35, 0x00000000) // DWMWA_CAPTION_COLOR: schwarz
	set(34, 0x00000000) // DWMWA_BORDER_COLOR: schwarz
	set(36, 0x00FFFFFF) // DWMWA_TEXT_COLOR: weiß
}

func hidden(name string, args ...string) *exec.Cmd {
	c := exec.Command(name, args...)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	return c
}

func shortcut(lnk, target string) {
	ps := "$s=(New-Object -ComObject WScript.Shell).CreateShortcut('" + strings.ReplaceAll(lnk, "'", "''") + "');" +
		"$s.TargetPath='" + strings.ReplaceAll(target, "'", "''") + "';$s.WorkingDirectory='" + strings.ReplaceAll(filepath.Dir(target), "'", "''") + "';" +
		"$s.IconLocation='" + strings.ReplaceAll(target, "'", "''") + ",0';$s.Description='Larsens World';$s.Save()"
	_ = hidden("powershell", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", ps).Run()
}

func paths() (dir, exe, desk, start string) {
	dir = filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", appName)
	exe = filepath.Join(dir, appName+".exe")
	home := os.Getenv("USERPROFILE")
	desk = filepath.Join(home, "Desktop", appName+".lnk")
	if k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\User Shell Folders`, registry.QUERY_VALUE); err == nil {
		if v, _, err := k.GetStringValue("Desktop"); err == nil && v != "" {
			if p, err := registry.ExpandString(v); err == nil {
				desk = filepath.Join(p, appName+".lnk")
			}
		}
		k.Close()
	}
	start = filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", appName+".lnk")
	return
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err = io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

func install(self string) bool {
	dir, exe, desk, start := paths()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false
	}
	if err := copyFile(self, exe); err != nil {
		return false
	}
	shortcut(desk, exe)
	shortcut(start, exe)
	if k, _, err := registry.CreateKey(registry.CURRENT_USER, uninKey, registry.SET_VALUE); err == nil {
		k.SetStringValue("DisplayName", appName)
		k.SetStringValue("DisplayIcon", exe+",0")
		k.SetStringValue("Publisher", "Larsen")
		k.SetStringValue("DisplayVersion", "1.0")
		k.SetStringValue("InstallLocation", dir)
		k.SetStringValue("UninstallString", `"`+exe+`" --uninstall`)
		k.SetDWordValue("NoModify", 1)
		k.SetDWordValue("NoRepair", 1)
		k.Close()
	}
	return true
}

func uninstall() {
	if msg("Larsens World wirklich deinstallieren?\nDeine Daten auf diesem PC werden dabei auch gelöscht.", 0x24) != 6 {
		return
	}
	dir, _, desk, start := paths()
	os.Remove(desk)
	os.Remove(start)
	registry.DeleteKey(registry.CURRENT_USER, uninKey)
	os.RemoveAll(filepath.Join(os.Getenv("APPDATA"), "LarsensWorld"))
	_ = hidden("cmd", "/c", "ping 127.0.0.1 -n 3 >nul & rmdir /s /q \""+dir+"\"").Start()
	msg("Larsens World wurde deinstalliert.", 0x40)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--uninstall" {
		uninstall()
		return
	}
	self, _ := os.Executable()
	_, exe, _, _ := paths()
	if !strings.EqualFold(filepath.Clean(self), filepath.Clean(exe)) {
		if install(self) {
			msg("Larsens World ist installiert!\nDu findest die App ab jetzt auf dem Desktop und im Startmenü.\nDie heruntergeladene Datei kannst du löschen.", 0x40)
			if hidden(exe).Start() == nil {
				return
			}
		}
	}
	data := filepath.Join(os.Getenv("APPDATA"), "LarsensWorld")
	_ = os.MkdirAll(data, 0o755)
	// gestochen scharf bei 125 %, 150 %, 4K usw.: Windows nicht mehr strecken lassen
	if p := user32.NewProc("SetProcessDpiAwarenessContext"); p.Find() == nil {
		p.Call(^uintptr(3)) // DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 (-4)
	} else {
		user32.NewProc("SetProcessDPIAware").Call()
	}
	dpi := uintptr(96)
	if p := user32.NewProc("GetDpiForSystem"); p.Find() == nil {
		if d, _, _ := p.Call(); d > 0 {
			dpi = d
		}
	}
	// dunkler Hintergrund schon während des Ladens (kein weißes Aufblitzen)
	os.Setenv("WEBVIEW2_DEFAULT_BACKGROUND_COLOR", "FF0B0D12")
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		AutoFocus: true,
		DataPath:  data,
		WindowOptions: webview2.WindowOptions{Title: appName, Width: uint(1280 * dpi / 96), Height: uint(860 * dpi / 96), IconId: 1, Center: true},
	})
	if w == nil {
		msg("Für Larsens World fehlt die Microsoft WebView2-Komponente.\nSie öffnet sich jetzt zum Herunterladen. Danach Larsens World einfach nochmal starten.", 0x40)
		_ = hidden("rundll32", "url.dll,FileProtocolHandler", "https://go.microsoft.com/fwlink/p/?LinkId=2124703").Start()
		return
	}
	defer w.Destroy()
	darkTitle(uintptr(w.Window()))
	user32.NewProc("ShowWindow").Call(uintptr(w.Window()), 3) // maximiert starten
	w.Navigate(appURL)
	w.Run()
}
