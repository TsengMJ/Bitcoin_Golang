package config

import (
	"errors"

	"github.com/spf13/viper"
)

var config *viper.Viper

func Init(env string) (error) {
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("./src/config/")
	err = config.ReadInConfig()
	if err != nil {
		return errors.New("failed to read config file: " + err.Error())
	}

	return nil
}

func GetConfig() *viper.Viper {
	return config
}
