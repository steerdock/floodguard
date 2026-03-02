# Changelog

All notable changes to FloodGuard will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [v1.0.3] - 2026-03-02

### Added

- **Blacklist enforcement** — IPs listed under `blacklist` in `config.yaml` are now
  blocked immediately on startup. The config field existed previously but was never
  read by any code. (`internal/monitor/monitor.go`)
  - Added `applyBlacklist()` — called at start of `Monitor.Start()` to block all
    configured entries via the active firewall backend.
  - Added `isBlacklisted()` — mirrors `isWhitelisted()`, supports both plain IPs and
    CIDR notation (e.g. `5.6.7.0/24`).

- **Webhook alert notifications** — when an IP is successfully blocked, a notification
  is sent to the configured webhook. (`internal/notifier/`)
  - New `internal/notifier` package with a `Notifier` interface and four
    implementations: `disabled`, `dingtalk`, `slack`, `wecom`.
  - Factory function `notifier.New(cfg, logger)` selects the backend based on
    `notification.type`; unknown values fall back to DingTalk.
  - HTTP timeout is 10 s; failures are logged but never crash the main loop.
  - No new third-party dependencies — uses only the Go standard library.

### Fixed

- Previously, the `blacklist` config field was parsed into `Config.Blacklist` but
  never acted upon, so configured entries had no effect.
- Previously, the `notification` config section was parsed but no HTTP request was
  ever sent, making the feature entirely non-functional.

---

## [v1.0.2] - Initial release

### Added

- Core connection monitoring via `/proc/net/tcp` and `/proc/net/tcp6`.
- Per-IP connection count threshold (`max_connections`) and QPS threshold (`max_qps`).
- Automatic firewall backend selection: nftables → iptables → dummy (fallback).
- IP whitelist with CIDR support; server public and local IPs auto-added on `init`.
- Ban duration policy (`ban.duration`); `0` = permanent ban.
- Structured logging (JSON / text) via `go.uber.org/zap`.
- `floodguard init` — generates a default `config.yaml` with auto-detected IPs.
- `floodguard start` — starts the protection loop.
- `floodguard status` — prints current status.
- systemd service support (see `scripts/install.sh`).
- Multi-language README: English, 简体中文, 日本語, 한국어, Deutsch, Français, Русский.
