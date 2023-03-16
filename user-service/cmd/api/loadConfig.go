package main

import (
	"github.com/spf13/viper"
)

type EnvConfig struct {
	UserDbConnString string `mapstructure:"USER_DB_CONN_STRING"`
	SymmetricKey     string `mapstructure:"SYMMETRIC_KEY"`
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
