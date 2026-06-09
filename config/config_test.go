package config_test

import (
	"strings"
	"testing"
	"time"

	"github.com/artuomvproskunin/gokit/config"
)

// mapLookup returns a lookup function backed by the given map.
func mapLookup(m map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := m[key]
		return v, ok
	}
}

func TestString_Default(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(nil)))
	if got := l.String("FOO", "bar"); got != "bar" {
		t.Fatalf("got %q, want %q", got, "bar")
	}
	if err := l.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestString_Value(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(map[string]string{"FOO": "baz"})))
	if got := l.String("FOO", "bar"); got != "baz" {
		t.Fatalf("got %q, want %q", got, "baz")
	}
}

func TestInt_Default(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(nil)))
	if got := l.Int("N", 42); got != 42 {
		t.Fatalf("got %d, want 42", got)
	}
	if err := l.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInt_Value(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(map[string]string{"N": "99"})))
	if got := l.Int("N", 0); got != 99 {
		t.Fatalf("got %d, want 99", got)
	}
}

func TestInt_ParseError(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(map[string]string{"N": "abc"})))
	if got := l.Int("N", 7); got != 7 {
		t.Fatalf("got %d, want default 7 on parse error", got)
	}
	err := l.Err()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "N") {
		t.Fatalf("error %q does not mention key N", err.Error())
	}
}

func TestInt64_Default(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(nil)))
	if got := l.Int64("X", 100); got != 100 {
		t.Fatalf("got %d, want 100", got)
	}
}

func TestInt64_Value(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(map[string]string{"X": "9000000000"})))
	if got := l.Int64("X", 0); got != 9_000_000_000 {
		t.Fatalf("got %d, want 9000000000", got)
	}
}

func TestBool_Default(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(nil)))
	if got := l.Bool("B", true); !got {
		t.Fatalf("got %v, want true", got)
	}
}

func TestBool_Value(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(map[string]string{"B": "false"})))
	if got := l.Bool("B", true); got {
		t.Fatalf("got %v, want false", got)
	}
}

func TestFloat64_Default(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(nil)))
	if got := l.Float64("F", 1.5); got != 1.5 {
		t.Fatalf("got %v, want 1.5", got)
	}
}

func TestFloat64_Value(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(map[string]string{"F": "3.14"})))
	if got := l.Float64("F", 0); got != 3.14 {
		t.Fatalf("got %v, want 3.14", got)
	}
}

func TestDuration_Default(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(nil)))
	if got := l.Duration("D", 5*time.Second); got != 5*time.Second {
		t.Fatalf("got %v, want 5s", got)
	}
}

func TestDuration_Value(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(map[string]string{"D": "10m"})))
	if got := l.Duration("D", time.Second); got != 10*time.Minute {
		t.Fatalf("got %v, want 10m", got)
	}
}

func TestStrings_Default(t *testing.T) {
	def := []string{"a", "b"}
	l := config.New(config.WithLookup(mapLookup(nil)))
	got := l.Strings("S", def)
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("got %v, want %v", got, def)
	}
}

func TestStrings_Value(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(map[string]string{"S": "x, y , z"})))
	got := l.Strings("S", nil)
	want := []string{"x", "y", "z"}
	if len(got) != len(want) {
		t.Fatalf("len %d, want %d; got %v", len(got), len(want), got)
	}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("got[%d] = %q, want %q", i, got[i], v)
		}
	}
}

func TestRequiredString_Missing(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(nil)))
	if got := l.RequiredString("MUST"); got != "" {
		t.Fatalf("got %q, want empty string", got)
	}
	err := l.Err()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "MUST") {
		t.Fatalf("error %q does not mention key MUST", err.Error())
	}
}

func TestRequiredString_Value(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(map[string]string{"MUST": "present"})))
	if got := l.RequiredString("MUST"); got != "present" {
		t.Fatalf("got %q, want present", got)
	}
	if err := l.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestErr_AccumulatesMultipleKeys(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(map[string]string{
		"ALPHA": "bad",
		"BETA":  "also-bad",
	})))
	l.Int("ALPHA", 0)
	l.Int("BETA", 0)
	err := l.Err()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "ALPHA") {
		t.Errorf("error %q does not mention key ALPHA", msg)
	}
	if !strings.Contains(msg, "BETA") {
		t.Errorf("error %q does not mention key BETA", msg)
	}
}

func TestErr_NilOnClean(t *testing.T) {
	l := config.New(config.WithLookup(mapLookup(nil)))
	if err := l.Err(); err != nil {
		t.Fatalf("expected nil error on clean loader, got: %v", err)
	}
}
