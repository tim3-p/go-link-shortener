package configs

/*
const (
	DefaultAddress string = "http://localhost:8080/"
)
*/
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080/"`
}

var EnvConfig Config
