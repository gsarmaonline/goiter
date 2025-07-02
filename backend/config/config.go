package config

import "os"

const (
	ModeTest    ModeT = "test"
	ModeDev     ModeT = "dev"
	ModeStaging ModeT = "staging"
	ModeProd    ModeT = "prod"
)

type (
	ModeT  string
	Config struct {
		Mode    ModeT
		GinMode string

		DBHost     string
		DBPort     string
		DBUser     string
		DBPassword string
		DBName     string

		Port string
	}
)

func DefaultConfig() *Config {
	return &Config{
		Mode:    ModeT(os.Getenv("MODE")),
		GinMode: os.Getenv("GIN_MODE"),

		Port: os.Getenv("PORT"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}
}
