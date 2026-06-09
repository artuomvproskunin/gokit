package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Loader reads configuration values from environment variables, accumulating
// parse errors rather than panicking. Callers assemble their own config struct
// from typed getter calls and check Err() once at the end.
type Loader struct {
	lookup func(string) (string, bool)
	errs   []error
}

// Option is a functional option for Loader.
type Option func(*Loader)

// WithLookup replaces the default os.LookupEnv with fn.
// Use this in tests to provide a hermetic in-memory environment.
func WithLookup(fn func(string) (string, bool)) Option {
	return func(l *Loader) {
		l.lookup = fn
	}
}

// New creates a Loader that reads from os.LookupEnv by default.
// Pass Option values to override behaviour.
func New(opts ...Option) *Loader {
	l := &Loader{lookup: os.LookupEnv}
	for _, o := range opts {
		o(l)
	}
	return l
}

func (l *Loader) get(key string) (string, bool) {
	return l.lookup(key)
}

func (l *Loader) addErr(key string, err error) {
	l.errs = append(l.errs, fmt.Errorf("config: key %q: %w", key, err))
}

// Err returns all accumulated errors joined into a single error, or nil.
func (l *Loader) Err() error {
	return errors.Join(l.errs...)
}

// String returns the string value of key, or def when the key is absent or empty.
func (l *Loader) String(key, def string) string {
	v, ok := l.get(key)
	if !ok || v == "" {
		return def
	}
	return v
}

// Int returns the value of key parsed as an int, or def when the key is absent,
// empty, or cannot be parsed. A parse error is accumulated in Err.
func (l *Loader) Int(key string, def int) int {
	v, ok := l.get(key)
	if !ok || v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		l.addErr(key, err)
		return def
	}
	return n
}

// Int64 returns the value of key parsed as an int64, or def when the key is
// absent, empty, or cannot be parsed. A parse error is accumulated in Err.
func (l *Loader) Int64(key string, def int64) int64 {
	v, ok := l.get(key)
	if !ok || v == "" {
		return def
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		l.addErr(key, err)
		return def
	}
	return n
}

// Bool returns the value of key parsed as a bool, or def when the key is absent,
// empty, or cannot be parsed. A parse error is accumulated in Err.
func (l *Loader) Bool(key string, def bool) bool {
	v, ok := l.get(key)
	if !ok || v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		l.addErr(key, err)
		return def
	}
	return b
}

// Float64 returns the value of key parsed as a float64, or def when the key is
// absent, empty, or cannot be parsed. A parse error is accumulated in Err.
func (l *Loader) Float64(key string, def float64) float64 {
	v, ok := l.get(key)
	if !ok || v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		l.addErr(key, err)
		return def
	}
	return f
}

// Duration returns the value of key parsed as a time.Duration, or def when the
// key is absent, empty, or cannot be parsed. A parse error is accumulated in Err.
func (l *Loader) Duration(key string, def time.Duration) time.Duration {
	v, ok := l.get(key)
	if !ok || v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		l.addErr(key, err)
		return def
	}
	return d
}

// Strings returns the value of key split by comma as a []string, or def when the
// key is absent or empty. Each element is trimmed; blank elements are dropped.
// If all elements are blank after trimming, def is returned.
func (l *Loader) Strings(key string, def []string) []string {
	v, ok := l.get(key)
	if !ok || v == "" {
		return def
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return def
	}
	return out
}

// RequiredString returns the value of key as a string.
// If the key is absent or empty, an error is accumulated in Err and "" is returned.
func (l *Loader) RequiredString(key string) string {
	v, ok := l.get(key)
	if !ok || v == "" {
		l.addErr(key, errors.New("required but not set"))
		return ""
	}
	return v
}
