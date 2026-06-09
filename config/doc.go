// Package config provides type-safe loading of application configuration from
// environment variables.
//
// # Basic usage
//
// Create a [Loader] with [New], call typed getters to read each variable, then
// check [Loader.Err] once at the end to surface all parse failures in one place:
//
//	l := config.New()
//
//	type AppConfig struct {
//		Host    string
//		Port    int
//		Debug   bool
//		Timeout time.Duration
//		Tags    []string
//		DSN     string
//	}
//
//	cfg := AppConfig{
//		Host:    l.String("APP_HOST", "localhost"),
//		Port:    l.Int("APP_PORT", 8080),
//		Debug:   l.Bool("APP_DEBUG", false),
//		Timeout: l.Duration("APP_TIMEOUT", 30*time.Second),
//		Tags:    l.Strings("APP_TAGS", nil),          // comma-separated: "a, b, c"
//		DSN:     l.RequiredString("DATABASE_URL"),    // error if absent or empty
//	}
//
//	if err := l.Err(); err != nil {
//		log.Fatal(err) // err names every offending key
//	}
//
// # Available getters
//
// All optional getters accept a default value that is returned when the key is
// absent or empty. Parse failures are accumulated — the getter returns the
// default and the error is surfaced later through [Loader.Err].
//
//	l.String(key, def)    // raw string
//	l.Int(key, def)       // strconv.Atoi
//	l.Int64(key, def)     // strconv.ParseInt base 10
//	l.Bool(key, def)      // strconv.ParseBool ("true"/"1"/"false"/"0" …)
//	l.Float64(key, def)   // strconv.ParseFloat
//	l.Duration(key, def)  // time.ParseDuration ("30s", "5m", "1h30m" …)
//	l.Strings(key, def)   // comma-separated CSV, elements trimmed
//
// [Loader.RequiredString] accumulates an error when the key is absent or empty
// and returns "".
//
// # Error handling
//
// [Loader.Err] returns all saved errors joined via [errors.Join]. Each
// sub-error mentions the offending key, for example:
//
//	config: key "DATABASE_URL": required but not set
//	config: key "APP_PORT": strconv.Atoi: parsing "abc": invalid syntax
//
// # Testing
//
// Use [WithLookup] to replace [os.LookupEnv] with an in-memory map so tests
// are hermetic and never call os.Setenv:
//
//	l := config.New(config.WithLookup(func(key string) (string, bool) {
//		m := map[string]string{"APP_PORT": "9090"}
//		v, ok := m[key]
//		return v, ok
//	}))
package config
