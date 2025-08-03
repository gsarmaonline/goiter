package config

import "os"

const (
	// Server modes
	ModeDev  ModeT = "dev"
	ModeProd ModeT = "prod"

	// DBType
	PostgresDbType DbTypeT = iota + 1
	SqliteDbType
)

type (
	ModeT   string
	DbTypeT uint8
	Config  struct {
		Mode    ModeT
		GinMode string

		DBHost     string
		DBPort     string
		DBUser     string
		DBPassword string
		DBName     string
		DBType     DbTypeT

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

func (cfg *Config) GetKey(cfgKey string) (cfgVal string) {
	cfgVal = os.Getenv(cfgKey)
	return
}
