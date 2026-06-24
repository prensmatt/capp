package config

import(
	"fmt"
	"os"
)

type Config struct{
	Port string
	DBURL string
	JWTSecret string
}

func Load() (*Config, error){
	cfg := &Config{
		Port: os.Getenv("PORT")
		DBURL: os.Getenv("DB_URL")
		JWTSecret: os.Getenv("JWT_SECRET")
	}
	if cfg.Port == ""{
		cfg.Port = "8080"
	}
	if cfg.DBURL == ""{
		return nil, fmt.Errorf("DB_URL is required")
	}
	if cfg.JWTSecret == ""{
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}