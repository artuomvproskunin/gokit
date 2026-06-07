# gokit

Reusable Go libraries for the car wash platform. Shared infrastructure primitives —
no domain logic, no external coupling.

## Packages

| Package    | Purpose |
|------------|---------|
| `config`   | Load and validate configuration from environment variables |
| `log`      | Structured slog-based logger (JSON/text, request-scoped fields) |
| `postgres` | pgxpool setup, health check, TxManager, goose migration runner |
| `httpx`    | HTTP mux, middleware (request-id/recovery/log/CORS/timeout), unified error format, decode/encode + validation |
| `auth`     | Session store, OTP/magic-link, session-middleware, API-key middleware |

## Usage

```bash
go get github.com/artuomvproskunin/gokit@vX.Y.Z
```

Import individual packages:

```go
import "github.com/artuomvproskunin/gokit/config"
import "github.com/artuomvproskunin/gokit/log"
import "github.com/artuomvproskunin/gokit/postgres"
import "github.com/artuomvproskunin/gokit/httpx"
import "github.com/artuomvproskunin/gokit/auth"
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
