package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	AccountDbConnString string `mapstructure:"ACCOUNT_DB_CONN_STRING"`
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigFile("./app/config.json")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
