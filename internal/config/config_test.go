package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_PORT", "")
	t.Setenv("DATABASE_URL", "")

	config := Load()
	if config.AppPort != "8081" {
		t.Fatalf("unexpected default port: %s", config.AppPort)
	}
	if config.DatabaseURL != "postgres://postgres:postgres@localhost:5433/mini_wms?sslmode=disable" {
		t.Fatalf("unexpected default database URL: %s", config.DatabaseURL)
	}
}

func TestLoadEnvironment(t *testing.T) {
	t.Setenv("APP_PORT", "9000")
	t.Setenv("DATABASE_URL", "postgres://example")

	config := Load()
	if config.AppPort != "9000" || config.DatabaseURL != "postgres://example" {
		t.Fatalf("unexpected config: %#v", config)
	}
}
