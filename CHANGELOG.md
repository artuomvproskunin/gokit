# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

#### `config` package

- `Loader` — accumulates configuration values from environment variables without panicking on errors.
- `New(opts ...Option) *Loader` — constructor; defaults to `os.LookupEnv`.
- `WithLookup(fn func(string)(string,bool)) Option` — replaces the env source (useful for hermetic tests).
- Typed getters with defaults (absent/empty key returns the default; parse errors are accumulated):
  `String`, `Int`, `Int64`, `Bool`, `Float64`, `Duration`, `Strings` (comma-separated CSV).
- `RequiredString(key string) string` — accumulates an error when the key is absent or empty.
- `Err() error` — returns all accumulated errors joined with `errors.Join`.

#### `log` package

- `Format` (`FormatJSON` / `FormatText`) — output encoding selector.
- `Options{Level, Format, Output, AddSource}` — logger configuration.
- `New(opts Options) *slog.Logger` — builds a `*slog.Logger` with UTC timestamps and context-field injection.
- `ParseLevel(s string) (slog.Level, error)` — case-insensitive mapping of `debug`/`info`/`warn`/`error`.
- `ContextWith(ctx context.Context, attrs ...slog.Attr) context.Context` — stores request-scoped attributes
  in the context; they are automatically included in every log record produced by a logger from `New`.
