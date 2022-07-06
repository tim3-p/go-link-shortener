package configs

type Config struct {
	ServerAddress   string `json:"server_address" env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `json:"base_url" env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `json:"file_storage_path" env:"FILE_STORAGE_PATH" envDefault:""`
	DatabaseDSN     string `json:"database_dsn" env:"DATABASE_DSN" envDefault:""`
	EnableHTTPS     bool   `json:"enable_https" env:"ENABLE_HTTPS" envDefault:"false"`
	ConfigJson      string `json:"" env:"CONFIG"`
}

var EnvConfig Config
