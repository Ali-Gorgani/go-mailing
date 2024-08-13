package config

import (
	"github.com/spf13/viper"
)

// Config holds the application wide configurations.
// The values are read by viper from the config file or environment variables.
type Config struct {
	Server    ServerConfig
	Kavenegar KavenegarConfig
}

type ServerConfig struct {
	Address             string   `mapstructure:"address"`
	LoadbalancerAddresses []string `mapstructure:"loadbalancer_address"`
}

type KavenegarConfig struct {
	Sender string `mapstructure:"sender"`
	URL    string `mapstructure:"url"`
	APIKey string `mapstructure:"api_key"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AutomaticEnv()

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
