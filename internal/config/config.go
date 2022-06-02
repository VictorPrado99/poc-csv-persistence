package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config file
type Config struct {
	Port             string `mapstructure:"port"`
	ConnectionString string `mapstructure:"connection_string"`
}

// Instance of configuration
var AppConfig *Config

func LoadAppConfig() {
	PERSISTENCE_HOME := os.Getenv("PERSISTENCE_HOME") // Get envrioment variable with the path

	if PERSISTENCE_HOME == "" { // If do not exist
		PERSISTENCE_HOME = os.Getenv("HOME") // Consider self
	}

	log.Println("Loading Server Configurations...")
	viper.AddConfigPath(PERSISTENCE_HOME)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatal(err)
	}
}
