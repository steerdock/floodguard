package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/floodguard/internal/config"
	"go.uber.org/zap"
)

// Notifier sends alert messages when an IP is blocked.
type Notifier interface {
	Send(ip string, connections int)
}

// New returns the appropriate Notifier based on configuration.
// If notifications are disabled or webhook_url is empty, a no-op is returned.
func New(cfg config.NotificationConfig, logger *zap.Logger) Notifier {
	if !cfg.Enabled || cfg.WebhookURL == "" {
		return &disabled{}
	}

	client := &http.Client{Timeout: 10 * time.Second}

	switch cfg.Type {
	case "slack":
		return &slackNotifier{url: cfg.WebhookURL, client: client, logger: logger}
	case "wecom":
		return &wecomNotifier{url: cfg.WebhookURL, client: client, logger: logger}
	default: // "dingtalk" or any unrecognized value defaults to dingtalk
		return &dingtalkNotifier{url: cfg.WebhookURL, client: client, logger: logger}
	}
}

// postJSON sends a JSON payload to the given URL via HTTP POST.
func postJSON(client *http.Client, url string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

// ---------------------------------------------------------------------------
// disabled — no-op implementation
// ---------------------------------------------------------------------------

type disabled struct{}

func (d *disabled) Send(_ string, _ int) {}

// ---------------------------------------------------------------------------
// dingtalk
// ---------------------------------------------------------------------------

type dingtalkNotifier struct {
	url    string
	client *http.Client
	logger *zap.Logger
}

func (n *dingtalkNotifier) Send(ip string, connections int) {
	msg := fmt.Sprintf(
		"[FloodGuard Alert] Malicious IP blocked\nIP: %s\nConnections: %d\nTime: %s",
		ip, connections, time.Now().Format(time.RFC3339),
	)

	// DingTalk text message format
	payload := map[string]any{
		"msgtype": "text",
		"text":    map[string]string{"content": msg},
	}

	if err := postJSON(n.client, n.url, payload); err != nil {
		n.logger.Error("DingTalk notification failed",
			zap.String("ip", ip),
			zap.Error(err),
		)
		return
	}
	n.logger.Info("DingTalk notification sent", zap.String("ip", ip))
}

// ---------------------------------------------------------------------------
// slack
// ---------------------------------------------------------------------------

type slackNotifier struct {
	url    string
	client *http.Client
	logger *zap.Logger
}

func (n *slackNotifier) Send(ip string, connections int) {
	msg := fmt.Sprintf(
		":rotating_light: *FloodGuard Alert* — Malicious IP blocked\n*IP:* `%s`\n*Connections:* %d\n*Time:* %s",
		ip, connections, time.Now().Format(time.RFC3339),
	)

	// Slack incoming webhook format
	payload := map[string]string{"text": msg}

	if err := postJSON(n.client, n.url, payload); err != nil {
		n.logger.Error("Slack notification failed",
			zap.String("ip", ip),
			zap.Error(err),
		)
		return
	}
	n.logger.Info("Slack notification sent", zap.String("ip", ip))
}

// ---------------------------------------------------------------------------
// wecom (WeCom / 企业微信)
// ---------------------------------------------------------------------------

type wecomNotifier struct {
	url    string
	client *http.Client
	logger *zap.Logger
}

func (n *wecomNotifier) Send(ip string, connections int) {
	msg := fmt.Sprintf(
		"[FloodGuard Alert] Malicious IP blocked\nIP: %s\nConnections: %d\nTime: %s",
		ip, connections, time.Now().Format(time.RFC3339),
	)

	// WeCom robot webhook format
	payload := map[string]any{
		"msgtype": "text",
		"text":    map[string]string{"content": msg},
	}

	if err := postJSON(n.client, n.url, payload); err != nil {
		n.logger.Error("WeCom notification failed",
			zap.String("ip", ip),
			zap.Error(err),
		)
		return
	}
	n.logger.Info("WeCom notification sent", zap.String("ip", ip))
}
