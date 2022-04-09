package configs

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:"host=localhost user=postgres password=postgres dbname=postgres sslmode=disable"`
}

var EnvConfig Config
