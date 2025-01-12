package config

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Hostname string
	Port     string
}

type DatabaseConfig struct {
	Hostname string
	Port     string
	Username string
	Password string
	Database string
}

func NewConfig() Config {
	return Config{
		App: AppConfig{
			Hostname: "localhost",
			Port:     "9000",
		},
		Database: DatabaseConfig{
			Hostname: "db.local",
			Port:     "5432",
			Username: "user",
			Password: "password",
			Database: "shortUrl",
		},
	}
}
