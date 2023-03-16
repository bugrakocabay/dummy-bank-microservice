package main

import (
	"github.com/spf13/viper"
)

type EnvConfig struct {
	MongoURI string `mapstructure:"MONGO_URI"`
	Username string `mapstructure:"MONGO_USERNAME"`
	Password string `mapstructure:"MONGO_PASSWORD"`
}

func LoadConfig() (config EnvConfig, err error) {
	viper.SetConfigFile("./app/config.json")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
