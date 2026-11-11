// Larsens World Launcher: öffnet die aktuelle Version von GitHub in Opera GX
// (oder, falls Opera GX fehlt, im Standard-Browser). Updates kommen automatisch,
// weil immer die Seite von GitHub geladen wird.
package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const url = "https://larsen140.github.io/Larsens-World/"

func main() {
	la, pf, pf86 := os.Getenv("LOCALAPPDATA"), os.Getenv("ProgramFiles"), os.Getenv("ProgramFiles(x86)")
	cands := []string{
		filepath.Join(la, "Programs", "Opera GX", "opera.exe"),
		filepath.Join(la, "Programs", "Opera GX", "launcher.exe"),
		filepath.Join(pf, "Opera GX", "opera.exe"),
		filepath.Join(pf, "Opera GX", "launcher.exe"),
		filepath.Join(pf86, "Opera GX", "launcher.exe"),
	}
	for _, c := range cands {
		if c == "" {
			continue
		}
		if _, err := os.Stat(c); err == nil {
			cmd := exec.Command(c, "--new-window", url)
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			if cmd.Start() == nil {
				return
			}
		}
	}
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_ = cmd.Start()
}
