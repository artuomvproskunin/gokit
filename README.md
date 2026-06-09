# gokit

Reusable Go libraries. Shared infrastructure primitives —
no domain logic, no external coupling.

## Packages

| Package    | Purpose |
|------------|---------|
| `config`   | Load and validate configuration from environment variables |
| `log`      | Structured slog-based logger (JSON/text, request-scoped fields) |
| `postgres` | pgxpool setup, health check, TxManager, goose migration runner |
| `httpx`    | HTTP mux, middleware (request-id/recovery/log/CORS/timeout), unified error format, decode/encode + validation |
| `auth`     | Session store, OTP/magic-link, session-middleware, API-key middleware |

---

## config

Type-safe configuration loading from environment variables. Errors are accumulated,
not panicked — a single `Err()` check at the end surfaces all problems at once.

```go
import (
    "log"
    "time"
    "github.com/artuomvproskunin/gokit/config"
)

l := config.New()

type AppConfig struct {
    Host    string
    Port    int
    Debug   bool
    Timeout time.Duration
    Tags    []string
    DSN     string
}

cfg := AppConfig{
    Host:    l.String("APP_HOST", "localhost"),
    Port:    l.Int("APP_PORT", 8080),
    Debug:   l.Bool("APP_DEBUG", false),
    Timeout: l.Duration("APP_TIMEOUT", 30*time.Second),
    Tags:    l.Strings("APP_TAGS", nil),         // comma-separated: "a, b, c"
    DSN:     l.RequiredString("DATABASE_URL"),   // error if absent or empty
}

if err := l.Err(); err != nil {
    log.Fatal(err) // names every offending key
}
```

**Available getters** — all optional ones return the default when the key is absent,
empty, or cannot be parsed:

| Getter | Type | Notes |
|--------|------|-------|
| `String(key, def)` | `string` | raw value |
| `Int(key, def)` | `int` | base-10 integer |
| `Int64(key, def)` | `int64` | base-10 integer |
| `Bool(key, def)` | `bool` | `"true"`, `"1"`, `"false"`, `"0"`, … |
| `Float64(key, def)` | `float64` | |
| `Duration(key, def)` | `time.Duration` | `"30s"`, `"5m"`, `"1h30m"` … |
| `Strings(key, def)` | `[]string` | comma-separated CSV, elements trimmed |
| `RequiredString(key)` | `string` | error if absent or empty |

### Testing

Use `WithLookup` to swap the environment source so tests are hermetic — no
`os.Setenv` required:

```go
l := config.New(config.WithLookup(func(key string) (string, bool) {
    m := map[string]string{"APP_PORT": "9090"}
    v, ok := m[key]
    return v, ok
}))
```

---

## log

A `*slog.Logger` (standard library) with JSON/text output, UTC timestamps, and
request-scoped field injection via `context.Context`.

### Creating a logger

```go
import (
    "log/slog"
    "os"
    "github.com/artuomvproskunin/gokit/log"
)

// Parse the log level from the environment.
lvl, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
if err != nil {
    lvl = slog.LevelInfo // default when the variable is absent or unknown
}

logger := log.New(log.Options{
    Level:  lvl,
    Format: log.FormatJSON,   // log.FormatText for local development
})

logger.Info("server started", slog.Int("port", 8080))
// {"time":"2024-01-15T12:00:00Z","level":"INFO","msg":"server started","port":8080}
```

`Options.Output` defaults to `os.Stdout`. Timestamps are always UTC.

### Request-scoped fields

`ContextWith` stamps attributes onto a `context.Context`. Any logger built with
`New` picks them up automatically — no need to thread the logger through every
function call:

```go
// In HTTP middleware — add request context once.
ctx = log.ContextWith(ctx,
    slog.String("request_id", requestID),
    slog.String("method", r.Method),
)

// Deep in a handler or service layer.
logger.InfoContext(ctx, "order created", slog.Int("order_id", id))
// {"time":"...Z","level":"INFO","msg":"order created",
//  "request_id":"abc-123","method":"POST","order_id":42}
```

Multiple calls to `ContextWith` accumulate: new attributes are appended and
earlier ones are never overwritten.

### Using `logger.With` / `logger.WithGroup`

The standard `slog` derivation methods work as expected — context injection is
preserved in any logger derived from the original:

```go
svcLog := logger.With(slog.String("service", "orders"))
svcLog.InfoContext(ctx, "started") // includes both "service" and all context attrs
```

---

## Installation

```bash
go get github.com/artuomvproskunin/gokit@vX.Y.Z
```

## Versioning

Releases follow [Semantic Versioning](https://semver.org/). Tags are pushed in the form
`vX.Y.Z`. Breaking changes bump the minor or major version and are documented in
[CHANGELOG.md](./CHANGELOG.md).

## Development

```bash
go test ./...
golangci-lint run
```
