package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	AppPort     string `env:"APP_PORT" env-default:"8080"`
	LogLevel    string `env:"LOG_LEVEL" env-default:"info"`
	DatabaseDir string `env:"DB_DIR" env-default:"/app/data"` // For internal SQLite
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(".env", cfg)
	if err != nil {
		// If .env missing, try reading purely from environment vars
		err = cleanenv.ReadEnv(cfg)
	}
	return cfg, err
}
