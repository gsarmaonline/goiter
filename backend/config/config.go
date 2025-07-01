package config

const (
	ModeTest    ModeT = "test"
	ModeDev     ModeT = "dev"
	ModeStaging ModeT = "staging"
	ModeProd    ModeT = "prod"
)

type (
	ModeT  string
	Config struct {
		Mode ModeT

		Port string
	}
)

func DefaultConfig() *Config {
	return &Config{
		Mode: ModeProd,
		Port: "8080",
	}
}
