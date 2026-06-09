package config_test

import (
	"fmt"
	"time"

	"github.com/artuomvproskunin/gokit/config"
)

// env is a helper that returns a lookup function backed by a static map.
func env(m map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := m[key]
		return v, ok
	}
}

func Example() {
	l := config.New(config.WithLookup(env(map[string]string{
		"APP_HOST":    "db.example.com",
		"APP_PORT":    "5432",
		"APP_TIMEOUT": "10s",
		"APP_TAGS":    "prod, eu-west",
	})))

	host := l.String("APP_HOST", "localhost")
	port := l.Int("APP_PORT", 5432)
	timeout := l.Duration("APP_TIMEOUT", 30*time.Second)
	tags := l.Strings("APP_TAGS", nil)

	if err := l.Err(); err != nil {
		fmt.Println("config error:", err)
		return
	}

	fmt.Printf("host=%s port=%d timeout=%s tags=%v\n", host, port, timeout, tags)
	// Output:
	// host=db.example.com port=5432 timeout=10s tags=[prod eu-west]
}

func ExampleLoader_RequiredString_missing() {
	l := config.New(config.WithLookup(env(nil))) // empty environment

	_ = l.RequiredString("DATABASE_URL")

	fmt.Println(l.Err())
	// Output:
	// config: key "DATABASE_URL": required but not set
}

func ExampleLoader_Err_multipleErrors() {
	l := config.New(config.WithLookup(env(map[string]string{
		"PORT":    "not-a-number",
		"TIMEOUT": "not-a-duration",
	})))

	l.Int("PORT", 0)
	l.Duration("TIMEOUT", 0)

	// Err returns all problems joined; each line names the offending key.
	fmt.Println(l.Err())
	// Output:
	// config: key "PORT": strconv.Atoi: parsing "not-a-number": invalid syntax
	// config: key "TIMEOUT": time: invalid duration "not-a-duration"
}
