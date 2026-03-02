# FloodGuard

[English](README.md) | [简体中文](README_CN.md) | [日本語](README_JA.md) | [한국어](README_KO.md) | [Deutsch](README_DE.md) | [Français](README_FR.md) | [Русский](README_RU.md) | [更新履歴](CHANGELOG.md)

CC 攻撃および DDoS 攻撃から保護するための、モダンで軽量な Linux ファイアウォールツールです。

## 特徴

- 🚀 **軽量・高速**：Go 言語製、単一バイナリ、リソース消費が少ない
- 🛡️ **スマート防護**：異常接続を多次元で検出し、攻撃 IP を自動ブロック
- 🔧 **柔軟な設定**：YAML 設定対応、各種閾値やポリシーをカスタマイズ可能
- 📊 **リアルタイム監視**：接続数統計、攻撃ログ、ブロック記録
- 🔔 **アラート通知**：Webhook 通知対応（DingTalk、WeCom、Slack）
- 🌐 **マルチバックエンド**：iptables、nftables、firewalld を自動検出
- 📝 **詳細ログ**：構造化ログ出力、複数フォーマット対応

## クイックスタート

### インストール

```bash
# Go を使ってインストール
go install github.com/steerdock/floodguard/cmd/floodguard@latest
```

> **注意**：`init` 実行時に、サーバーのパブリック IP およびローカル IP を自動検出してホワイトリストに追加します。誤ったブロックを防ぐために必ず実行してください。

### 使い方

```bash
# デフォルト設定ファイルを生成
sudo floodguard init

# 保護を開始
sudo floodguard start

# ステータス確認
sudo floodguard status

# ブロックリストを表示
sudo floodguard list

# IP のブロックを解除
sudo floodguard unban 1.2.3.4
```

## 設定

設定ファイルのパス：`/etc/floodguard/config.yaml`

```yaml
# 監視設定
monitor:
  interval: 10s              # チェック間隔
  max_connections: 100       # IP あたりの最大接続数
  max_qps: 50                # IP あたりの最大 QPS

# ブロックポリシー
ban:
  duration: 3600            # ブロック期間（秒）、0 で永久
  mode: "auto"              # auto/iptables/nftables/firewalld

# ホワイトリスト
whitelist:
  - "127.0.0.1"
  - "192.168.0.0/16"

# 通知
notification:
  enabled: true
  webhook_url: "https://your-webhook-url"
```

## システム要件

- Linux（カーネル 3.10+）
- root 権限
- iptables または nftables

## 開発

```bash
# リポジトリをクローン
git clone https://github.com/steerdock/floodguard.git
cd floodguard

# 依存関係をインストール
go mod download

# ビルド
make build

# テスト実行
make test
```

## デプロイ（systemd）

```bash
# バイナリをインストール
sudo cp build/floodguard /usr/local/bin/
sudo chmod +x /usr/local/bin/floodguard

# SELinux コンテキストを修正（RHEL/CentOS/Fedora）
sudo restorecon -v /usr/local/bin/floodguard

# 設定を初期化（最初に必ず実行！）
sudo floodguard init

# systemd サービスを作成
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

# 有効化して起動
sudo systemctl daemon-reload
sudo systemctl enable floodguard
sudo systemctl start floodguard
sudo systemctl status floodguard
```

## サービス管理

```bash
sudo systemctl start floodguard
sudo systemctl stop floodguard
sudo systemctl restart floodguard
sudo systemctl status floodguard
sudo journalctl -u floodguard -f
```

## 更新履歴

完全なリリース履歴は [CHANGELOG.md](CHANGELOG.md) をご覧ください。

## ライセンス

MIT License
