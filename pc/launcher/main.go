// Larsens World für Windows: installiert sich beim ersten Start selbst
// (%LOCALAPPDATA%\Programs\Larsens World, Verknüpfungen auf Desktop und im Startmenü,
// Eintrag unter "Apps & Features") und öffnet dann ein eigenes App-Fenster (WebView2),
// das immer die aktuelle Version von GitHub lädt. Daten: %APPDATA%\LarsensWorld.
package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	webview2 "github.com/jchv/go-webview2"
	"golang.org/x/sys/windows/registry"
)

const (
	appVersion = "1.3"
	verURL     = "https://raw.githubusercontent.com/Larsen140/Larsens-World/main/pc/version.txt"
	exeURL     = "https://github.com/Larsen140/Larsens-World/raw/main/pc/Larsens-World.exe"
	appURL     = "https://larsen140.github.io/Larsens-World/?app=1"
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
	if _, err := os.Stat(exe); err == nil {
		old := exe + ".old"
		os.Remove(old)
		if err := os.Rename(exe, old); err != nil {
			return false
		}
	}
	if err := copyFile(self, exe); err != nil {
		return false
	}
	shortcut(desk, exe)
	shortcut(start, exe)
	registerProtocol(exe)
	if k, _, err := registry.CreateKey(registry.CURRENT_USER, uninKey, registry.SET_VALUE); err == nil {
		k.SetStringValue("DisplayName", appName)
		k.SetStringValue("DisplayIcon", exe+",0")
		k.SetStringValue("Publisher", "Larsen")
		k.SetStringValue("DisplayVersion", appVersion)
		k.SetStringValue("InstallLocation", dir)
		k.SetStringValue("UninstallString", `"`+exe+`" --uninstall`)
		k.SetDWordValue("NoModify", 1)
		k.SetDWordValue("NoRepair", 1)
		k.Close()
	}
	return true
}

// larsensworld://import/<code> öffnet die App und übernimmt die Daten aus dem Link
func registerProtocol(exe string) {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Classes\larsensworld`, registry.SET_VALUE)
	if err != nil {
		return
	}
	k.SetStringValue("", "URL:Larsens World")
	k.SetStringValue("URL Protocol", "")
	k.Close()
	if k, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Classes\larsensworld\DefaultIcon`, registry.SET_VALUE); err == nil {
		k.SetStringValue("", exe+",0")
		k.Close()
	}
	if k, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Classes\larsensworld\shell\open\command`, registry.SET_VALUE); err == nil {
		k.SetStringValue("", `"`+exe+`" "%1"`)
		k.Close()
	}
}

func newer(remote, local string) bool {
	r, err1 := strconv.ParseFloat(strings.TrimSpace(remote), 64)
	l, err2 := strconv.ParseFloat(local, 64)
	return err1 == nil && err2 == nil && r > l
}

// Selbst-Update: neue Programmversion von GitHub holen, austauschen und neu starten
func selfUpdate(exe string) bool {
	cl := &http.Client{Timeout: 3 * time.Second}
	res, err := cl.Get(verURL + "?t=" + strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		return false
	}
	b, _ := io.ReadAll(io.LimitReader(res.Body, 64))
	res.Body.Close()
	if res.StatusCode != 200 || !newer(string(b), appVersion) {
		return false
	}
	cl.Timeout = 90 * time.Second
	res, err = cl.Get(exeURL + "?t=" + strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		return false
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return false
	}
	tmp := exe + ".new"
	f, err := os.Create(tmp)
	if err != nil {
		return false
	}
	n, err := io.Copy(f, res.Body)
	f.Close()
	head := make([]byte, 2)
	if hf, e := os.Open(tmp); e == nil {
		hf.Read(head)
		hf.Close()
	}
	if err != nil || n < 1_000_000 || string(head) != "MZ" {
		os.Remove(tmp)
		return false
	}
	old := exe + ".old"
	os.Remove(old)
	if os.Rename(exe, old) != nil {
		os.Remove(tmp)
		return false
	}
	if os.Rename(tmp, exe) != nil {
		os.Rename(old, exe)
		return false
	}
	return hidden(exe, os.Args[1:]...).Start() == nil
}

func uninstall() {
	if msg("Larsens World wirklich deinstallieren?\nDeine Daten auf diesem PC werden dabei auch gelöscht.", 0x24) != 6 {
		return
	}
	dir, _, desk, start := paths()
	os.Remove(desk)
	os.Remove(start)
	registry.DeleteKey(registry.CURRENT_USER, uninKey)
	for _, k := range []string{`larsensworld\shell\open\command`, `larsensworld\shell\open`, `larsensworld\shell`, `larsensworld\DefaultIcon`, `larsensworld`} {
		registry.DeleteKey(registry.CURRENT_USER, `Software\Classes\`+k)
	}
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
	os.Remove(exe + ".old")
	if !strings.EqualFold(filepath.Clean(self), filepath.Clean(exe)) {
		if install(self) {
			msg("Larsens World ist installiert!\nDu findest die App ab jetzt auf dem Desktop und im Startmenü.\nUpdates kommen ab jetzt automatisch.\nDie heruntergeladene Datei kannst du löschen.", 0x40)
			if hidden(exe, os.Args[1:]...).Start() == nil {
				return
			}
		}
	} else {
		registerProtocol(exe)
		if selfUpdate(exe) {
			return
		}
	}
	nav := appURL
	for _, a := range os.Args[1:] {
		if i := strings.Index(a, "import/"); strings.HasPrefix(strings.ToLower(a), "larsensworld:") && i >= 0 {
			code := strings.Trim(a[i+len("import/"):], "/\" ")
			ok := code != ""
			for _, c := range code {
				if !(c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-' || c == '_') {
					ok = false
					break
				}
			}
			if ok {
				nav = appURL + "#lwimport=" + code
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
		WindowOptions: webview2.WindowOptions{Title: appName, Width: uint(1280 * dpi / 96), Height: uint(860 * dpi / 96), IconId: 2, Center: true},
	})
	if w == nil {
		msg("Für Larsens World fehlt die Microsoft WebView2-Komponente.\nSie öffnet sich jetzt zum Herunterladen. Danach Larsens World einfach nochmal starten.", 0x40)
		_ = hidden("rundll32", "url.dll,FileProtocolHandler", "https://go.microsoft.com/fwlink/p/?LinkId=2124703").Start()
		return
	}
	defer w.Destroy()
	darkTitle(uintptr(w.Window()))
	user32.NewProc("ShowWindow").Call(uintptr(w.Window()), 3) // maximiert starten
	w.Navigate(nav)
	w.Run()
}
