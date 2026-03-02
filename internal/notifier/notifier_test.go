package notifier

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yourusername/floodguard/internal/config"
	"go.uber.org/zap"
)

// captureServer returns a test HTTP server that records the last request body,
// along with a pointer to that body string.
func captureServer(t *testing.T) (*httptest.Server, *string) {
	t.Helper()
	var lastBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		lastBody = string(data)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	return srv, &lastBody
}

func TestDisabled_NeverSends(t *testing.T) {
	srv, body := captureServer(t)

	// enabled=false — should produce a no-op notifier regardless of URL
	cfg := config.NotificationConfig{
		Enabled:    false,
		WebhookURL: srv.URL,
		Type:       "dingtalk",
	}
	logger, _ := zap.NewDevelopment()
	n := New(cfg, logger)
	n.Send("1.2.3.4", 200)

	if *body != "" {
		t.Errorf("expected no request when disabled, got body: %s", *body)
	}
}

func TestEmptyURL_NeverSends(t *testing.T) {
	cfg := config.NotificationConfig{
		Enabled:    true,
		WebhookURL: "",
		Type:       "slack",
	}
	logger, _ := zap.NewDevelopment()
	n := New(cfg, logger)
	// Should not panic or send anything
	n.Send("5.6.7.8", 50)
}

func TestDingTalk_PayloadFormat(t *testing.T) {
	srv, body := captureServer(t)

	cfg := config.NotificationConfig{
		Enabled:    true,
		WebhookURL: srv.URL,
		Type:       "dingtalk",
	}
	logger, _ := zap.NewDevelopment()
	New(cfg, logger).Send("10.0.0.1", 300)

	var payload map[string]any
	if err := json.Unmarshal([]byte(*body), &payload); err != nil {
		t.Fatalf("invalid JSON: %v\nbody: %s", err, *body)
	}
	if payload["msgtype"] != "text" {
		t.Errorf("expected msgtype=text, got %v", payload["msgtype"])
	}
	text, ok := payload["text"].(map[string]any)
	if !ok {
		t.Fatalf("expected text object, got %T", payload["text"])
	}
	content, _ := text["content"].(string)
	if !strings.Contains(content, "10.0.0.1") {
		t.Errorf("expected IP in content, got: %s", content)
	}
}

func TestSlack_PayloadFormat(t *testing.T) {
	srv, body := captureServer(t)

	cfg := config.NotificationConfig{
		Enabled:    true,
		WebhookURL: srv.URL,
		Type:       "slack",
	}
	logger, _ := zap.NewDevelopment()
	New(cfg, logger).Send("10.0.0.2", 150)

	var payload map[string]any
	if err := json.Unmarshal([]byte(*body), &payload); err != nil {
		t.Fatalf("invalid JSON: %v\nbody: %s", err, *body)
	}
	text, ok := payload["text"].(string)
	if !ok {
		t.Errorf("expected text string, got %T", payload["text"])
	}
	if !strings.Contains(text, "10.0.0.2") {
		t.Errorf("expected IP in text, got: %s", text)
	}
}

func TestWeCom_PayloadFormat(t *testing.T) {
	srv, body := captureServer(t)

	cfg := config.NotificationConfig{
		Enabled:    true,
		WebhookURL: srv.URL,
		Type:       "wecom",
	}
	logger, _ := zap.NewDevelopment()
	New(cfg, logger).Send("10.0.0.3", 99)

	var payload map[string]any
	if err := json.Unmarshal([]byte(*body), &payload); err != nil {
		t.Fatalf("invalid JSON: %v\nbody: %s", err, *body)
	}
	if payload["msgtype"] != "text" {
		t.Errorf("expected msgtype=text, got %v", payload["msgtype"])
	}
	text, ok := payload["text"].(map[string]any)
	if !ok {
		t.Fatalf("expected text object, got %T", payload["text"])
	}
	content, _ := text["content"].(string)
	if !strings.Contains(content, "10.0.0.3") {
		t.Errorf("expected IP in content, got: %s", content)
	}
}

func TestUnknownType_FallsBackToDingTalk(t *testing.T) {
	srv, body := captureServer(t)

	cfg := config.NotificationConfig{
		Enabled:    true,
		WebhookURL: srv.URL,
		Type:       "unknown_type",
	}
	logger, _ := zap.NewDevelopment()
	New(cfg, logger).Send("1.1.1.1", 10)

	// Should still send with dingtalk format
	var payload map[string]any
	if err := json.Unmarshal([]byte(*body), &payload); err != nil {
		t.Fatalf("invalid JSON: %v\nbody: %s", err, *body)
	}
	if payload["msgtype"] != "text" {
		t.Errorf("expected dingtalk fallback with msgtype=text, got %v", payload["msgtype"])
	}
}

func TestServerError_DoesNotPanic(t *testing.T) {
	// Server that always returns 500
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	cfg := config.NotificationConfig{
		Enabled:    true,
		WebhookURL: srv.URL,
		Type:       "slack",
	}
	logger, _ := zap.NewDevelopment()
	// Should log error but not panic
	New(cfg, logger).Send("2.2.2.2", 500)
}
