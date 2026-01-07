package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Chains   ChainsConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port        string
	Host        string
	Environment string
}

type DatabaseConfig struct {
	Host           string
	Port           string
	Name           string
	User           string
	Password       string
	SSLMode        string
	MaxConnections int
	MaxIdleConns   int
}

type LoggingConfig struct {
	Level  string
	Format string
	Output string
}

type ChainsConfig struct {
	DefaultChainID int64
	ConfigPath     string
}

func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("DB_MAX_CONNECTIONS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNECTIONS", 5)
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_FORMAT", "json")
	viper.SetDefault("LOG_OUTPUT", "stdout")
	viper.SetDefault("CHAINS_CONFIG_PATH", "internal/config/chains.json")
	viper.SetDefault("DEFAULT_CHAIND_ID", 1337)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found. using defaults and environment variables")
	}

	config := &Config{
		Server: ServerConfig{
			Port:        viper.GetString("SERVER_PORT"),
			Host:        viper.GetString("SERVER_HOST"),
			Environment: viper.GetString("ENVIRONMENT"),
		},
		Database: DatabaseConfig{
			Host:           viper.GetString("DB_HOST"),
			Port:           viper.GetString("DB_PORT"),
			Name:           viper.GetString("DB_NAME"),
			User:           viper.GetString("DB_USER"),
			Password:       viper.GetString("DB_PASSWORD"),
			SSLMode:        viper.GetString("DB_SSL_MODE"),
			MaxConnections: viper.GetInt("DB_MAX_CONNECTIONS"),
			MaxIdleConns:   viper.GetInt("DB_MAX_IDLE_CONNECTIONS"),
		},
		Logging: LoggingConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
			Output: viper.GetString("LOG_OUTPUT"),
		},
		Chains: ChainsConfig{
			DefaultChainID: viper.GetInt64("DEFAULT_CHAIN_ID"),
			ConfigPath:     viper.GetString("CHAINS_CONFIG_PATH"),
		},
	}
	return config, nil
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}
