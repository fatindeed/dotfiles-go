package app

import (
	"fmt"

	"github.com/spf13/viper"
)

func SetConfigFile(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		// viper.SetConfigType("yaml")
	}
}

type handler struct {
	config config
}

func (h *handler) init() error {
	// Use config file from the flag.
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&h.config); err != nil {
		return err
	}
	if h.config.Name != viper.GetString("name") {
		return fmt.Errorf("invalid app name (%s) in the config file", h.config.Name)
	}

	// init backup config
	return h.config.init()
}
