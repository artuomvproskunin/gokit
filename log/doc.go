// Package log provides a structured [slog.Logger] with JSON and text output,
// UTC timestamps, and request-scoped field propagation via context.
//
// # Creating a logger
//
// Call [New] with an [Options] value. The result is a standard *[slog.Logger]
// and works with all slog APIs:
//
//	logger := log.New(log.Options{
//		Level:  slog.LevelInfo,
//		Format: log.FormatJSON,  // or log.FormatText for local development
//	})
//
//	logger.Info("server started", slog.Int("port", 8080))
//	// {"time":"2024-01-15T12:00:00Z","level":"INFO","msg":"server started","port":8080}
//
// [Options.Output] defaults to os.Stdout. Timestamps are always in UTC.
//
// # Reading the log level from the environment
//
// [ParseLevel] converts "debug", "info", "warn", or "error" (case-insensitive)
// to a [slog.Level]:
//
//	lvl, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
//	if err != nil {
//		lvl = slog.LevelInfo // sensible default for unknown / missing value
//	}
//
//	logger := log.New(log.Options{Level: lvl, Format: log.FormatJSON})
//
// # Request-scoped fields
//
// [ContextWith] attaches [slog.Attr] values to a context. Any logger created
// with [New] automatically includes those attributes in every record that is
// logged with that context:
//
//	// In HTTP middleware — stamp a request-id on the context once.
//	ctx = log.ContextWith(ctx,
//		slog.String("request_id", requestID),
//		slog.String("method", r.Method),
//	)
//
//	// Deep in a handler or service — no logger argument required.
//	logger.InfoContext(ctx, "order created", slog.Int("order_id", id))
//	// {"time":"...Z","level":"INFO","msg":"order created",
//	//  "request_id":"abc-123","method":"POST","order_id":42}
//
// Multiple calls to [ContextWith] accumulate: new attributes are appended and
// earlier ones are never overwritten.
//
// # Using logger.With
//
// The standard logger.With / logger.WithGroup methods are fully supported —
// context injection is preserved across the derived loggers:
//
//	svcLog := logger.With(slog.String("service", "orders"))
//	svcLog.InfoContext(ctx, "started") // includes both "service" and context attrs
package log
