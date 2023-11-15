package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port        int      `mapstructure:"port"`
	Environment string   `mapstructure:"environment"`
	ScraperURL  string   `mapstructure:"scraper_url"`
	DB          DBConfig `mapstructure:"mongodb"`
}

type DBConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

func NewConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file. err = %v", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config file. err = %v", err)
	}

	return &config, nil
}
