# FloodGuard

[English](README.md) | [简体中文](README_CN.md) | [日本語](README_JA.md) | [한국어](README_KO.md) | [Deutsch](README_DE.md) | [Français](README_FR.md) | [Русский](README_RU.md)

Современный легковесный инструмент межсетевого экрана для Linux, защищающий от CC-атак и DDoS-атак.

## Возможности

- 🚀 **Лёгкий и быстрый**: Написан на Go, единый бинарный файл, минимальное потребление ресурсов
- 🛡️ **Умная защита**: Многомерное обнаружение аномальных подключений, автоматическая блокировка IP
- 🔧 **Гибкая настройка**: Конфигурация на основе YAML с настраиваемыми порогами и политиками
- 📊 **Мониторинг в реальном времени**: Статистика подключений, журналы атак и записи блокировок
- 🔔 **Уведомления**: Поддержка Webhook (DingTalk, WeCom, Slack)
- 🌐 **Мультибэкенд**: Автоматическое определение iptables, nftables, firewalld
- 📝 **Детальное логирование**: Структурированный вывод журнала в различных форматах

## Быстрый старт

### Установка

```bash
# Установка через Go
go install github.com/steerdock/floodguard/cmd/floodguard@latest
```

> **Примечание**: При инициализации FloodGuard автоматически определяет публичный и локальный IP сервера и добавляет их в белый список, чтобы избежать случайной блокировки.

### Использование

```bash
# Создать файл конфигурации по умолчанию
sudo floodguard init

# Запустить защиту
sudo floodguard start

# Проверить статус
sudo floodguard status

# Показать список заблокированных IP
sudo floodguard list

# Разблокировать IP
sudo floodguard unban 1.2.3.4
```

## Конфигурация

Файл конфигурации: `/etc/floodguard/config.yaml`

```yaml
# Настройки мониторинга
monitor:
  interval: 10s              # Интервал проверки
  max_connections: 100       # Макс. подключений на IP
  max_qps: 50                # Макс. QPS на IP

# Политика блокировки
ban:
  duration: 3600            # Длительность блокировки (сек), 0 — навсегда
  mode: "auto"              # auto/iptables/nftables/firewalld

# Белый список
whitelist:
  - "127.0.0.1"
  - "192.168.0.0/16"

# Уведомления
notification:
  enabled: true
  webhook_url: "https://your-webhook-url"
```

## Системные требования

- Linux (ядро 3.10+)
- Права root
- iptables или nftables

## Разработка

```bash
# Клонировать репозиторий
git clone https://github.com/steerdock/floodguard.git
cd floodguard

# Установить зависимости
go mod download

# Собрать
make build

# Запустить тесты
make test
```

## Развёртывание (systemd)

```bash
# Установить бинарный файл
sudo cp build/floodguard /usr/local/bin/
sudo chmod +x /usr/local/bin/floodguard

# Исправить контекст SELinux (RHEL/CentOS/Fedora)
sudo restorecon -v /usr/local/bin/floodguard

# Инициализировать конфигурацию (выполните в первую очередь!)
sudo floodguard init

# Создать службу systemd
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

# Включить и запустить
sudo systemctl daemon-reload
sudo systemctl enable floodguard
sudo systemctl start floodguard
sudo systemctl status floodguard
```

## Управление службой

```bash
sudo systemctl start floodguard
sudo systemctl stop floodguard
sudo systemctl restart floodguard
sudo systemctl status floodguard
sudo journalctl -u floodguard -f
```

## Лицензия

MIT License
