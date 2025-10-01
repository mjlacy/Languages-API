package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	AppName    string
	ConfigPath string
	Collection string
	Database   string
	DBURL      string
	Port       string
	Version    string
}

func New() (Config, error) {
	viper.SetDefault("AppName", AppName)
	viper.SetDefault("ConfigPath", "config.json")
	viper.SetDefault("Collection", "")
	viper.SetDefault("Database", "")
	viper.SetDefault("DBURL", "")
	viper.SetDefault("Port", "8080")
	viper.SetDefault("Version", Version)

	viper.SetConfigType("json")
	viper.SetConfigFile(viper.GetString("ConfigPath"))

	err := viper.ReadInConfig()
	if err != nil {
		log.Error().Err(err).Msg("Error reading in config file")
		return Config{}, err
	}

	var c Config
	err = viper.Unmarshal(&c)
	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling config file")
		return Config{}, err
	}

	return c, err
}
