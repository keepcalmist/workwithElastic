package config

import (
	"github.com/spf13/viper"
	"log"
)

type ConfI interface {
	GetPort() string
}

type ConfElastic struct {
	port string
}

func (c *ConfElastic) GetPort() string {
	return c.port
}

func GetPort() string {
	return viper.GetString("PORT")
}

func New() error {
	openEnv()
	return nil
}

func openEnv(){
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Config file doesn't exists")
	}
}



