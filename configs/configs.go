package configs

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application wide configurations.
// The values are read by viper from the config file or environment variables.
type Config struct {
	Server      serverConfig
	Kavenegar   kavenegarConfig
	DB          dbConfig
	AccessToken accessTokenConfig
}

type serverConfig struct {
	Environment string `mapstructure:"environment"`
	Address     string `mapstructure:"address"`
}

type kavenegarConfig struct {
	Sender string `mapstructure:"sender"`
	URL    string `mapstructure:"url"`
	APIKey string `mapstructure:"api_key"`
}

type dbConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslmode"`
}

type accessTokenConfig struct {
	SecretKey string        `mapstructure:"secret_key"`
	Duration  time.Duration `mapstructure:"duration"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AutomaticEnv()
	viper.AddConfigPath(path)
	viper.SetConfigName("configs")
	viper.SetConfigType("json")

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("could not read config file: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("could not unmarshal config: %w", err)
	}
	config.AccessToken.SecretKey = viper.GetString("access_token.secret_key")
	config.AccessToken.Duration = viper.GetDuration("access_token.duration")
	// fmt.Println("Viper All Keys:", viper.AllKeys())

	return config, nil
}
