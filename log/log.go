package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Format selects the output encoding of the logger.
type Format int

const (
	// FormatJSON produces JSON-encoded log records (one object per line).
	FormatJSON Format = iota
	// FormatText produces human-readable key=value log records.
	FormatText
)

// Options configures the logger produced by New.
type Options struct {
	// Level is the minimum level at which records are emitted.
	Level slog.Level
	// Format selects between JSON and text output. Defaults to FormatJSON.
	Format Format
	// Output is the destination for log records. Defaults to os.Stdout.
	Output io.Writer
	// AddSource appends the source file and line number to every record.
	AddSource bool
}

// New creates a *slog.Logger configured by opts. Timestamps are written in UTC.
// Attributes stored in the context via ContextWith are automatically included in
// every record produced by the returned logger.
func New(opts Options) *slog.Logger {
	out := opts.Output
	if out == nil {
		out = os.Stdout
	}

	ho := &slog.HandlerOptions{
		Level:       opts.Level,
		AddSource:   opts.AddSource,
		ReplaceAttr: utcTime,
	}

	var base slog.Handler
	switch opts.Format {
	case FormatText:
		base = slog.NewTextHandler(out, ho)
	default:
		base = slog.NewJSONHandler(out, ho)
	}

	return slog.New(&contextHandler{base})
}

// ParseLevel converts a level name ("debug", "info", "warn", "error") to the
// corresponding slog.Level. Comparison is case-insensitive. An unrecognised
// string returns an error.
func ParseLevel(s string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("log: unknown level %q", s)
	}
}

// utcTime is a slog ReplaceAttr function that converts the top-level timestamp
// attribute to UTC so all log records carry a canonical timezone.
func utcTime(groups []string, a slog.Attr) slog.Attr {
	if len(groups) == 0 && a.Key == slog.TimeKey {
		a.Value = slog.TimeValue(a.Value.Time().UTC())
	}
	return a
}
