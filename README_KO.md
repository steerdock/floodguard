# FloodGuard

[English](README.md) | [简体中文](README_CN.md) | [日本語](README_JA.md) | [한국어](README_KO.md) | [Deutsch](README_DE.md) | [Français](README_FR.md) | [Русский](README_RU.md) | [변경 로그](CHANGELOG.md)

CC 공격 및 DDoS 공격을 방어하기 위한 현대적이고 가벼운 Linux 방화벽 도구입니다.

## 특징

- 🚀 **경량 고성능**: Go 언어로 작성된 단일 바이너리, 적은 리소스 사용
- 🛡️ **스마트 보호**: 다차원 비정상 연결 감지, 공격 IP 자동 차단
- 🔧 **유연한 설정**: YAML 기반 설정, 임계값 및 정책 커스터마이징 가능
- 📊 **실시간 모니터링**: 연결 수 통계, 공격 로그, 차단 기록
- 🔔 **알림 경보**: Webhook 알림 지원 (DingTalk, WeCom, Slack)
- 🌐 **멀티 백엔드 지원**: iptables, nftables, firewalld 자동 감지
- 📝 **상세 로그**: 다양한 형식의 구조화된 로그 출력

## 빠른 시작

### 설치

```bash
# Go를 사용하여 설치
go install github.com/steerdock/floodguard/cmd/floodguard@latest
```

> **주의**: `init` 실행 시 서버의 공인 IP와 로컬 IP를 자동으로 감지하여 화이트리스트에 추가합니다. 실수로 차단되는 것을 방지합니다.

### 사용법

```bash
# 기본 설정 파일 생성
sudo floodguard init

# 보호 시작
sudo floodguard start

# 상태 확인
sudo floodguard status

# 차단 목록 확인
sudo floodguard list

# IP 차단 해제
sudo floodguard unban 1.2.3.4
```

## 설정

설정 파일 위치: `/etc/floodguard/config.yaml`

```yaml
# 모니터링 설정
monitor:
  interval: 10s              # 체크 간격
  max_connections: 100       # IP당 최대 연결 수
  max_qps: 50                # IP당 최대 QPS

# 차단 정책
ban:
  duration: 3600            # 차단 기간 (초), 0은 영구 차단
  mode: "auto"              # auto/iptables/nftables/firewalld

# 화이트리스트
whitelist:
  - "127.0.0.1"
  - "192.168.0.0/16"

# 알림
notification:
  enabled: true
  webhook_url: "https://your-webhook-url"
```

## 시스템 요구 사항

- Linux (커널 3.10+)
- root 권한
- iptables 또는 nftables

## 개발

```bash
# 저장소 복제
git clone https://github.com/steerdock/floodguard.git
cd floodguard

# 의존성 설치
go mod download

# 빌드
make build

# 테스트 실행
make test
```

## 배포 (systemd)

```bash
# 바이너리 설치
sudo cp build/floodguard /usr/local/bin/
sudo chmod +x /usr/local/bin/floodguard

# SELinux 컨텍스트 수정 (RHEL/CentOS/Fedora)
sudo restorecon -v /usr/local/bin/floodguard

# 설정 초기화 (먼저 실행하세요!)
sudo floodguard init

# systemd 서비스 생성
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

# 활성화 및 시작
sudo systemctl daemon-reload
sudo systemctl enable floodguard
sudo systemctl start floodguard
sudo systemctl status floodguard
```

## 서비스 관리

```bash
sudo systemctl start floodguard
sudo systemctl stop floodguard
sudo systemctl restart floodguard
sudo systemctl status floodguard
sudo journalctl -u floodguard -f
```

## 변경 로그

전체 릴리스 기록은 [CHANGELOG.md](CHANGELOG.md)를 참조하세요.

## 라이선스

MIT License
