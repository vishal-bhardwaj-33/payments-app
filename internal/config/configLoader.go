package config

import (
	"fmt"
	"log"
	"net/http"

	er "payments-app/internal/utils"

	"github.com/spf13/viper"
)

const DefaultConfigFile = "./config/default.toml"

// LoadConfig loads the configuration from the default TOML file into the given struct.
func LoadConfig(config interface{}) error {
	v := viper.New()
	v.SetConfigFile(DefaultConfigFile)
	v.SetConfigType("toml")
	v.AutomaticEnv() // Allow environment variables to override config values

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		return er.NewError(http.StatusInternalServerError, fmt.Sprintf("error loading config file: %v", err))
	}

	// Unmarshal into the provided struct
	if err := v.Unmarshal(config); err != nil {
		return er.NewError(http.StatusInternalServerError, fmt.Sprintf("error parsing config file: %v", err))
	}

	log.Printf("Config loaded from %s", DefaultConfigFile)
	return nil
}
