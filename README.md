# PHOENIX Feuer OS

PHOENIX Feuer OS ist eine schlanke, rein clientseitige Kommandozentrale für Solo-Unternehmer*innen, die ihr Projekt "PHOENIX" strukturiert dokumentieren möchten. Die Anwendung läuft vollständig im Browser, speichert Daten lokal und unterstützt dich dabei, Fortschritt, Compliance-Leitplanken und Tag-X-Vorbereitungen im Blick zu behalten.

## Feature-Überblick

- **Compliance-Monitor & 80-Tage-Kompass** – behalte die Meilensteine zwischen Projektstart und Tag X im Blick, inklusive Schnell-Checks für P-Konto, kaltes Geschäftskonto und vorbereitete Insolvenzanträge.
- **Tägliches Schutz-Protokoll** – protokolliere Arbeitsminuten, Modus (Struktur vs. Übungs-Audits) und halte Gesundheits-, Executive- sowie Fakten-Logs getrennt und konsistent.【F:index.html†L692-L772】【F:index.html†L820-L870】
- **Primär-Direktiven (MIT)** – setze bis zu drei Tagesziele, hake sie ab und halte Fokus ohne Overhead.【F:index.html†L772-L812】
- **Mini-Betrieb · Finanz-Snapshot** – erfasse Einnahmen/Ausgaben, erhalte Jahres-Overviews und dokumentiere Mini-EÜR-Daten für Nachweise.【F:index.html†L883-L960】
- **90-Minuten-Fokus-Sprint** – plane dokumentierte Fokusblöcke inklusive Timer, Restlaufzeit und Leitplankenwarnung für die 180-Minuten-Grenze.【F:index.html†L1047-L1126】
- **Prototypen-/Evidence-Verwaltung** – erstelle Übungs-Audits, sammle Evidence und baue Tag-X-Pakete anhand rechtlicher Mindestanforderungen strukturiert auf.【F:index.html†L1132-L1608】
- **Export/Import & Backups** – sichere deinen Zustand als JSON, erstelle Tag-X-Pakete und importiere Backups bei Bedarf.【F:index.html†L995-L1046】【F:index.html†L1602-L1706】

## Lokale Einrichtung

> **Hinweis:** Die Anwendung benötigt keinen Build-Prozess. Alles ist statisch und kann mit jedem Webserver ausgeliefert werden.

1. **Repository klonen**
   ```bash
   git clone https://github.com/<dein-account>/phoenix-feuer-os.git
   cd phoenix-feuer-os
   ```
2. **Lokalen Server starten** – wähle eine Option:
   - Python 3: `python3 -m http.server 4173`
   - Node.js (serve): `npx serve . -l 4173`
   - PHP: `php -S localhost:4173`
3. **Browser öffnen** – rufe `http://localhost:4173` auf und erlaube dem Browser, `localStorage` zu nutzen.
4. **PWA installieren (optional)** – klicke im Browser auf „Installieren“, um PHOENIX Feuer OS als Progressive Web App offline verfügbar zu machen.

## Deployment-Optionen

Da es sich um eine statische Web-App handelt, sind Deployments sehr flexibel:

- **GitHub Pages** – lege das Repo als Pages-Projekt (Branch `main` → `/`-Verzeichnis) an; GitHub liefert HTML/CSS/JS direkt aus.
- **Netlify / Vercel** – importiere das Repo und setze das Build-Kommando auf „(leer)“; Output-Verzeichnis `.`.
- **Cloud Storage (S3, R2, Firebase Hosting)** – lade die Dateien `index.html`, `404.html`, `sw.js`, `manifest.webmanifest` und den `assets/`-Ordner hoch und aktiviere HTTPS.
- **Self-Hosted / On-Premise** – kopiere den Projektordner auf deinen Server und konfiguriere nginx/Apache mit einem simplen „try_files“-Fallback auf `index.html`, um die PWA-Assets auszuliefern.

### Deployment-Checks

- Stelle sicher, dass der MIME-Type für `manifest.webmanifest` und `sw.js` korrekt gesetzt ist (`application/manifest+json` bzw. `application/javascript`).
- Aktiviere HTTPS (z. B. über Let’s Encrypt), damit Service Worker, PWA-Installation und `localStorage` ohne Einschränkungen funktionieren.
- Teste das Offline-Verhalten, indem du den Service Worker in den DevTools aktivierst und anschließend in den Offline-Modus wechselst.

## Datenschutz & Datensouveränität

- **Speicherort** – alle Eingaben werden ausschließlich in `localStorage` des Browsers hinterlegt; es gibt keinen externen Sync.【F:index.html†L354-L406】
- **Exports** – JSON-Exports werden nur auf deinem Gerät erstellt. Verteile sie ausschließlich an Personen, die gemäß Datenschutzrecht berechtigt sind (z. B. Anwält*innen, Insolvenzverwalter*innen).
- **Sensitive Bereiche** – Buddha-Notizen und andere private Felder sind für externe Reports vorgesehen, aber werden nicht automatisch exportiert; prüfe vor Weitergabe den Exportinhalt.【F:index.html†L844-L870】【F:index.html†L995-L1046】
- **Rechtsgrundlagen** – beachte bei Verarbeitung personenbezogener Daten insbesondere:
  - [Art. 6, 17 DSGVO – Rechtmäßigkeit & Löschung](https://eur-lex.europa.eu/eli/reg/2016/679/oj)
  - [§ 26 BDSG – Beschäftigtendaten](https://www.gesetze-im-internet.de/bdsg_2018/__26.html)
  - [§ 35 Abs. 2 InsO – Freigabe von Vermögen](https://www.gesetze-im-internet.de/inso/__35.html)
  - [SGB XII / SGB II – Einkommensanrechnung](https://www.gesetze-im-internet.de/sgb_12/__82.html)

> **Disclaimer:** Diese App ersetzt keine Rechtsberatung. Hole dir bei Unsicherheiten Unterstützung durch Fachanwält*innen oder Schuldnerberatungen.

## FAQ

**Brauche ich einen Server oder ein Backend?**  
Nein. Alles läuft im Browser. Nutze bei Bedarf nur einen simplen Webserver, um die Dateien auszuliefern.

**Kann ich mehrere Profile oder Haushalte anlegen?**  
Die App verwaltet aktuell genau einen Zustand pro Browserprofil. Für getrennte Setups verwende unterschiedliche Browser-Profile oder exportiere/importiere JSON-Dateien.

**Wie sichere ich meine Daten?**  
Nutze den Button „Export JSON“ im Bereich „Backup“, speichere die Datei an einem sicheren Ort (verschlüsselt) und importiere sie bei Bedarf wieder über „Import“.

**Sind Tag-X-Pakete rechtssicher?**  
Die App liefert Vorlagen und Checklisten, die an rechtliche Leitplanken angelehnt sind, ersetzt aber keine Prüfung durch Rechtsprofis. Lass finale Dokumente verifizieren, bevor du sie einreichst.【F:index.html†L1602-L1706】

**Was passiert bei einem Browser-Update oder Gerätewechsel?**  
`localStorage` kann durch Browser-Resets gelöscht werden. Regelmäßige Exports sind Pflicht, wenn du Geräte wechselst oder den Browser neu installierst.

## Weiterführende Dokumentation

- [Nutzer-Workflows](docs/WORKFLOWS.md)
- [Changelog & Release Notes](CHANGELOG.md)

---

Für Fragen oder Beiträge erstelle gerne Issues bzw. Pull Requests im Repository. 
