package config

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8081"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8081"`
	DatabaseAddress string `env:"DATABASE_DSN" envDefault:""`
}
