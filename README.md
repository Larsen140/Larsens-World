# Larsens World

Notizen, Trading, Casino (Blackjack, Roulette, Slots, Okey), Flying Bird und Glücksrad.

**Live:** https://larsen140.github.io/Larsens-World/

## PC-App
[`pc/Larsens-World.exe`](pc/Larsens-World.exe) herunterladen und starten. Sie öffnet Larsens World in einem eigenen App-Fenster (WebView2, in Windows eingebaut) und lädt immer die aktuelle Version von GitHub, Updates kommen also automatisch. Ohne Internet startet die zuletzt geladene Version. Daten liegen in `%APPDATA%\LarsensWorld`.

Windows zeigt beim ersten Start evtl. „Der Computer wurde durch Windows geschützt“, weil die Datei nicht signiert ist: **Weitere Informationen → Trotzdem ausführen**.

## Hinweise
- Diese Version läuft ohne Claude-Konto: alles wird nur im jeweiligen Browser gespeichert.
- KI-Chat und Angebote-Suche gibt es nur in der Claude-Version.
- Quelltext der PC-App: `pc/launcher/main.go` (Go + github.com/jchv/go-webview2, gebaut mit `GOOS=windows go build -ldflags "-H windowsgui -s -w"`).
