package config

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	DatabaseAddress string `env:"DATABASE_DSN" envDefault:"host=localhost user=postgres password=mysecretpassword dbname=yandex port=5432 sslmode=disable TimeZone=UTC"`
	SecretKey       string `env:"SECRET_KEY" envDefault:"secret key"`
}
