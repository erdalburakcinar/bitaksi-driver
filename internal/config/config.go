package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Ports struct {
		HTTP int `mapstructure:"http"`
	} `mapstructure:"ports"`
	Server struct {
		APIKey string `mapstructure:"api_key"`
	} `mapstructure:"server"`
	MongoDB struct {
		Username   string `mapstructure:"username"`
		Password   string `mapstructure:"password"`
		Host       string `mapstructure:"host"`
		Port       int    `mapstructure:"port"`
		Database   string `mapstructure:"database"`
		Collection string `mapstructure:"collection"`
	} `mapstructure:"mongodb"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("Unable to decode into struct: %v", err)
		return nil, err
	}

	return &config, nil
}
