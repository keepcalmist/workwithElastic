package config

import (
	"github.com/spf13/viper"
)

func New() *viper.Viper {
	cfg := viper.New()
	cfg.SetDefault("hui","pizda")
}
