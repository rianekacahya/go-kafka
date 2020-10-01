package goconf

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

var (
	configuration *viper.Viper
	mutex         sync.Once
)

func Config() *viper.Viper {
	mutex.Do(func() {
		configuration = new()
	})

	return configuration
}

func new() *viper.Viper {
	configuration := viper.New()
	configuration.SetConfigType("yaml")
	configuration.SetConfigName("config")
	configuration.AddConfigPath("files/config")
	if err := configuration.ReadInConfig(); err != nil {
		log.Fatalf("got an error reading file config, error: %s", err)
	}
	return configuration
}
