# Larsens World

Notizen, Trading, Casino (Blackjack, Roulette, Slots, Okey), Flying Bird und Glücksrad.

**Live:** https://larsen140.github.io/Larsens-World/

## PC-App
[`pc/Larsens-World.exe`](pc/Larsens-World.exe) herunterladen und starten. Sie öffnet immer die aktuelle Version von GitHub in Opera GX (ohne Opera GX im Standard-Browser), Updates kommen also automatisch.

Windows zeigt beim ersten Start evtl. „Der Computer wurde durch Windows geschützt“, weil die Datei nicht signiert ist: **Weitere Informationen → Trotzdem ausführen**.

## Hinweise
- Diese Version läuft ohne Claude-Konto: alles wird nur im jeweiligen Browser gespeichert.
- KI-Chat und Angebote-Suche gibt es nur in der Claude-Version.
- Daten aus der Claude-Version holen: dort **Optionen → Speicher → „Mit allen Daten auf GitHub öffnen“**.
- Quelltext der PC-App: `pc/launcher/main.go` (Go, gebaut mit `GOOS=windows go build -ldflags "-H windowsgui"`).
