# Larsens World

Notizen, Trading, Casino (Blackjack, Roulette, Slots, Okey), Flying Bird und Glücksrad.

**Live:** https://larsen140.github.io/Larsens-World/

## PC-App für Windows
**[⬇ Larsens World herunterladen](https://github.com/Larsen140/Larsens-World/raw/main/pc/Larsens-World.exe)**

Die exe ist nur ein kleiner Starter (125 KB): Sie öffnet Larsens World in Microsoft Edge (sonst Chrome) im App-Modus, also als eigenes Fenster ohne Adressleiste. Weil immer die Version von GitHub geladen wird, kommen Updates von selbst – nie neu installieren. Tipp: Rechtsklick auf die exe → „An Taskleiste anheften“ bzw. Verknüpfung auf den Desktop.

Windows zeigt beim ersten Start evtl. „Der Computer wurde durch Windows geschützt“, weil die Datei nicht signiert ist: **Weitere Informationen → Trotzdem ausführen**.

## Hinweise
- Diese Version läuft ohne Claude-Konto: alles wird nur im jeweiligen Browser gespeichert.
- KI-Chat und Angebote-Suche gibt es nur in der Claude-Version.
- Updates: `version.txt` muss zu `APP_VERSION` in `index.html` passen; weicht sie ab, lädt die App einmal pro Version neu.
- Quelltext des Starters: `pc/starter.c` + `pc/starter.rc`, gebaut mit `python3 -m ziglang cc -target x86_64-windows-gnu -municode -Wl,--subsystem,windows -Os -s starter.c starter.rc -o Larsens-World.exe`.
