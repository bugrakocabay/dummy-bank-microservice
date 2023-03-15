package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	AccountDbConnString string `mapstructure:"ACCOUNT_DB_CONN_STRING"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
