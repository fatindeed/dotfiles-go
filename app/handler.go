package app

import (
	"fmt"

	"github.com/spf13/viper"
)

var cfgViper *viper.Viper

func init() {
	cfgViper = viper.New()
}

func SetConfigFile(cfgFile string) {
	if cfgFile != "" {
		cfgViper.SetConfigFile(cfgFile)
	} else {
		cfgViper.AddConfigPath(".")
		cfgViper.SetConfigName("config")
		// cfgViper.SetConfigType("yaml")
	}
}

type handler struct {
	config config
	// vault  *vault
}

func (h *handler) init() error {
	// Use config file from the flag.
	if err := cfgViper.ReadInConfig(); err != nil {
		return err
	}
	if err := cfgViper.Unmarshal(&h.config); err != nil {
		return err
	}
	if h.config.Name != viper.GetString("name") {
		return fmt.Errorf("invalid app name (%s) in the config file", h.config.Name)
	}

	// init backup config
	return h.config.init()
}
