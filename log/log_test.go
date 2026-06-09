package log_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/artuomvproskunin/gokit/log"
)

func TestNew_JSONOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(log.Options{
		Level:  slog.LevelInfo,
		Format: log.FormatJSON,
		Output: &buf,
	})
	logger.Info("hello", slog.String("key", "val"))

	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v\nraw: %s", err, buf.String())
	}
	if m["level"] == nil {
		t.Error("JSON record missing level field")
	}
	if m["msg"] == nil {
		t.Error("JSON record missing msg field")
	}
	if got, _ := m["key"].(string); got != "val" {
		t.Errorf("key = %q, want %q", got, "val")
	}
}

func TestParseLevel(t *testing.T) {
	cases := []struct {
		input string
		want  slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"Info", slog.LevelInfo},
		{" info ", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"WARN", slog.LevelWarn},
		{"error", slog.LevelError},
		{"ERROR", slog.LevelError},
	}
	for _, c := range cases {
		got, err := log.ParseLevel(c.input)
		if err != nil {
			t.Errorf("ParseLevel(%q): unexpected error: %v", c.input, err)
			continue
		}
		if got != c.want {
			t.Errorf("ParseLevel(%q) = %v, want %v", c.input, got, c.want)
		}
	}

	if _, err := log.ParseLevel("verbose"); err == nil {
		t.Error(`ParseLevel("verbose"): expected error, got nil`)
	}
}

func TestContextWith_AddsFields(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(log.Options{
		Level:  slog.LevelInfo,
		Format: log.FormatJSON,
		Output: &buf,
	})

	ctx := log.ContextWith(context.Background(), slog.String("request_id", "abc123"))
	logger.InfoContext(ctx, "request handled")

	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v\nraw: %s", err, buf.String())
	}
	if got, _ := m["request_id"].(string); got != "abc123" {
		t.Errorf("request_id = %q, want %q", got, "abc123")
	}
}

func TestNew_LevelFilter(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(log.Options{
		Level:  slog.LevelInfo,
		Format: log.FormatJSON,
		Output: &buf,
	})
	logger.DebugContext(context.Background(), "should not appear")
	if buf.Len() != 0 {
		t.Errorf("expected empty output for debug at info level, got: %s", buf.String())
	}
}

func TestNew_UTCTimestamp(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(log.Options{
		Level:  slog.LevelInfo,
		Format: log.FormatJSON,
		Output: &buf,
	})
	logger.Info("ts test")

	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	ts, ok := m["time"].(string)
	if !ok {
		t.Fatal("time field missing or not a string")
	}
	if !strings.HasSuffix(ts, "Z") {
		t.Errorf("timestamp %q is not UTC (no Z suffix)", ts)
	}
}
