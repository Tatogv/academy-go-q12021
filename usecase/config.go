package usecase

import (
	"api/entities"

	"github.com/spf13/viper"
)

func LoadConfig() (config entities.Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("api")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
