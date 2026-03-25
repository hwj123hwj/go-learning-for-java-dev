package config

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:     "localhost",
		DBPort:     "5433",
		DBUser:     "root",
		DBPassword: "15671040800q",
		DBName:     "go_user_api",
		ServerPort: "8080",
	}
}
