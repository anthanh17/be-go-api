package configs

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
type Config struct {
	Database DatabaseConfig
	Cache    CacheConfig
	HTTP     HTTPConfig
}

// DatabaseConfig struct for database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Source   string
}

// CacheConfig struct for cache configuration
type CacheConfig struct {
	Type     string
	Address  string
	Username string
	Password string
}

// HTTPConfig struct for HTTP server configuration
type HTTPConfig struct {
	Address string
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
