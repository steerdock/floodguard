# FloodGuard

[English](README.md) | [简体中文](README_CN.md) | [日本語](README_JA.md) | [한국어](README_KO.md) | [Deutsch](README_DE.md) | [Français](README_FR.md) | [Русский](README_RU.md) | [Changelog](CHANGELOG.md)

Ein modernes, leichtgewichtiges Linux-Firewall-Tool zum Schutz vor CC-Angriffen und DDoS-Angriffen.

## Funktionen

- 🚀 **Leichtgewichtig & Schnell**: In Go geschrieben, einzelne Binärdatei, minimaler Ressourcenverbrauch
- 🛡️ **Intelligenter Schutz**: Mehrdimensionale Erkennung abnormaler Verbindungen, automatische IP-Sperrung
- 🔧 **Flexible Konfiguration**: YAML-basierte Konfiguration mit anpassbaren Schwellwerten und Richtlinien
- 📊 **Echtzeit-Überwachung**: Verbindungsstatistiken, Angriffsprotokolle und Sperr-Aufzeichnungen
- 🔔 **Alarmbenachrichtigungen**: Webhook-Unterstützung (DingTalk, WeCom, Slack)
- 🌐 **Multi-Backend-Unterstützung**: Automatische Erkennung von iptables, nftables, firewalld
- 📝 **Detailliertes Logging**: Strukturierte Protokollausgabe in mehreren Formaten

## Schnellstart

### Installation

```bash
# Installation über Go
go install github.com/steerdock/floodguard/cmd/floodguard@latest
```

> **Hinweis**: Bei der Initialisierung erkennt FloodGuard automatisch die öffentliche und lokale IP des Servers und fügt sie der Whitelist hinzu, um versehentliche Sperrungen zu verhindern.

### Verwendung

```bash
# Standard-Konfigurationsdatei erstellen
sudo floodguard init

# Schutz starten
sudo floodguard start

# Status anzeigen
sudo floodguard status

# Sperrliste anzeigen
sudo floodguard list

# IP-Sperre aufheben
sudo floodguard unban 1.2.3.4
```

## Konfiguration

Konfigurationsdatei: `/etc/floodguard/config.yaml`

```yaml
# Überwachungseinstellungen
monitor:
  interval: 10s              # Prüfintervall
  max_connections: 100       # Max. Verbindungen pro IP
  max_qps: 50                # Max. QPS pro IP

# Sperrrichtlinie
ban:
  duration: 3600            # Sperrdauer (Sekunden), 0 für permanent
  mode: "auto"              # auto/iptables/nftables/firewalld

# Whitelist
whitelist:
  - "127.0.0.1"
  - "192.168.0.0/16"

# Benachrichtigungen
notification:
  enabled: true
  webhook_url: "https://your-webhook-url"
```

## Systemanforderungen

- Linux (Kernel 3.10+)
- Root-Rechte
- iptables oder nftables

## Entwicklung

```bash
# Repository klonen
git clone https://github.com/steerdock/floodguard.git
cd floodguard

# Abhängigkeiten installieren
go mod download

# Bauen
make build

# Tests ausführen
make test
```

## Deployment (systemd)

```bash
# Binärdatei installieren
sudo cp build/floodguard /usr/local/bin/
sudo chmod +x /usr/local/bin/floodguard

# SELinux-Kontext reparieren (RHEL/CentOS/Fedora)
sudo restorecon -v /usr/local/bin/floodguard

# Konfiguration initialisieren (zuerst ausführen!)
sudo floodguard init

# systemd-Dienst erstellen
sudo tee /etc/systemd/system/floodguard.service > /dev/null <<EOF
[Unit]
Description=FloodGuard - DDoS Protection Service
After=network.target

[Service]
Type=exec
ExecStart=/usr/local/bin/floodguard start --config /etc/floodguard/config.yaml
Restart=on-failure
RestartSec=5s
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

# Aktivieren und starten
sudo systemctl daemon-reload
sudo systemctl enable floodguard
sudo systemctl start floodguard
sudo systemctl status floodguard
```

## Dienstverwaltung

```bash
sudo systemctl start floodguard
sudo systemctl stop floodguard
sudo systemctl restart floodguard
sudo systemctl status floodguard
sudo journalctl -u floodguard -f
```

## Changelog

Die vollständige Versionshistorie finden Sie in [CHANGELOG.md](CHANGELOG.md).

## Lizenz

MIT License
