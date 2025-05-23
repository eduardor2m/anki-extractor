package config

type Config struct {
	ServerPort string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort: "8080",
	}
}
