package config

type Config struct {
	AppPort     string
	DatabaseURL string
}

func Load() Config {
	port := "8081"
	databaseURL := "postgres://postgres:postgres@localhost:5433/mini_wms?sslmode=disable"

	return Config{AppPort: port, DatabaseURL: databaseURL}
}
