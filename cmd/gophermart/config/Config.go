package config

type Config struct {
	ServerAddress   string `env:"RUN_ADDRESS" envDefault:":8080"`
	DatabaseAddress string `env:"DATABASE_URI" envDefault:"host=localhost user=postgres password=mysecretpassword dbname=yandex port=5432 sslmode=disable TimeZone=UTC"`
	SecretKey       string `env:"SECRET_KEY" envDefault:"secret key"`
	AccrualAddress  string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"8082"`
}
