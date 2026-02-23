# FloodGuard

[English](README.md) | [简体中文](README_CN.md) | [日本語](README_JA.md) | [한국어](README_KO.md) | [Deutsch](README_DE.md) | [Français](README_FR.md) | [Русский](README_RU.md)

Un outil de pare-feu Linux moderne et léger pour se défendre contre les attaques CC et DDoS.

## Fonctionnalités

- 🚀 **Léger & Rapide** : Écrit en Go, binaire unique, faible consommation de ressources
- 🛡️ **Protection Intelligente** : Détection multi-dimensionnelle des connexions anormales, blocage automatique des IP
- 🔧 **Configuration Flexible** : Configuration YAML avec seuils et politiques personnalisables
- 📊 **Surveillance en Temps Réel** : Statistiques de connexions, journaux d'attaques et enregistrements de blocage
- 🔔 **Alertes et Notifications** : Support Webhook (DingTalk, WeCom, Slack)
- 🌐 **Multi-backend** : Détection automatique d'iptables, nftables, firewalld
- 📝 **Journalisation Détaillée** : Sortie de log structurée en plusieurs formats

## Démarrage Rapide

### Installation

```bash
# Installation via Go
go install github.com/steerdock/floodguard/cmd/floodguard@latest
```

> **Remarque** : Lors de l'initialisation, FloodGuard détecte automatiquement les IP publiques et locales du serveur et les ajoute à la liste blanche pour éviter les blocages accidentels.

### Utilisation

```bash
# Générer le fichier de configuration par défaut
sudo floodguard init

# Démarrer la protection
sudo floodguard start

# Vérifier le statut
sudo floodguard status

# Afficher la liste des IP bloquées
sudo floodguard list

# Débloquer une IP
sudo floodguard unban 1.2.3.4
```

## Configuration

Fichier de configuration : `/etc/floodguard/config.yaml`

```yaml
# Paramètres de surveillance
monitor:
  interval: 10s              # Intervalle de vérification
  max_connections: 100       # Connexions max par IP
  max_qps: 50                # QPS max par IP

# Politique de blocage
ban:
  duration: 3600            # Durée de blocage (secondes), 0 pour permanent
  mode: "auto"              # auto/iptables/nftables/firewalld

# Liste blanche
whitelist:
  - "127.0.0.1"
  - "192.168.0.0/16"

# Notifications
notification:
  enabled: true
  webhook_url: "https://your-webhook-url"
```

## Configuration Requise

- Linux (noyau 3.10+)
- Privilèges root
- iptables ou nftables

## Développement

```bash
# Cloner le dépôt
git clone https://github.com/steerdock/floodguard.git
cd floodguard

# Installer les dépendances
go mod download

# Compiler
make build

# Exécuter les tests
make test
```

## Déploiement (systemd)

```bash
# Installer le binaire
sudo cp build/floodguard /usr/local/bin/
sudo chmod +x /usr/local/bin/floodguard

# Corriger le contexte SELinux (RHEL/CentOS/Fedora)
sudo restorecon -v /usr/local/bin/floodguard

# Initialiser la configuration (à exécuter en premier !)
sudo floodguard init

# Créer le service systemd
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

# Activer et démarrer
sudo systemctl daemon-reload
sudo systemctl enable floodguard
sudo systemctl start floodguard
sudo systemctl status floodguard
```

## Gestion du Service

```bash
sudo systemctl start floodguard
sudo systemctl stop floodguard
sudo systemctl restart floodguard
sudo systemctl status floodguard
sudo journalctl -u floodguard -f
```

## Licence

MIT License
