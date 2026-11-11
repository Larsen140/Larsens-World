# Larsens World

Notizen, Trading, Casino (Blackjack, Roulette, Slots, Okey), Flying Bird und Glücksrad.

**Live:** https://larsen140.github.io/Larsens-World/

## PC-App für Windows
**[⬇ Larsens World herunterladen](https://github.com/Larsen140/Larsens-World/raw/main/pc/Larsens-World.exe)**

Datei starten, fertig: Sie installiert sich selbst (Verknüpfung auf dem Desktop und im Startmenü) und öffnet Larsens World in einem eigenen App-Fenster. Updates kommen automatisch von GitHub, ohne Internet startet die zuletzt geladene Version. Deinstallieren: Windows-Einstellungen → Apps → Larsens World.

Windows zeigt beim ersten Start evtl. „Der Computer wurde durch Windows geschützt“, weil die Datei nicht signiert ist: **Weitere Informationen → Trotzdem ausführen**.

## Hinweise
- Diese Version läuft ohne Claude-Konto: alles wird nur im jeweiligen Browser gespeichert.
- KI-Chat und Angebote-Suche gibt es nur in der Claude-Version.
- Quelltext der PC-App: `pc/launcher/main.go` (Go + github.com/jchv/go-webview2, gebaut mit `GOOS=windows go build -ldflags "-H windowsgui -s -w"`).
