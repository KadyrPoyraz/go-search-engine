package utils

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
    Port int `mapstructure:"port"`
}

type Config struct {
    Server ServerConfig `mapstructure:"server"`
}

var vp *viper.Viper

func LoadConfig() (Config, error) {
   vp = viper.New()
   var config Config
   
   vp.SetConfigName("config")
   vp.SetConfigType("json")
   vp.AddConfigPath(".")
   vp.AddConfigPath("./utils")

    err := vp.ReadInConfig()
    if err != nil {
        return Config{}, err
    }

    err = vp.Unmarshal(&config)
    if err != nil {
        return Config{}, err
    }

   return config, nil
}
