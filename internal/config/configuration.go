package config

import (
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
	v := viper.New()

	v.SetDefault("AppName", AppName)
	v.SetDefault("ConfigPath", "config.json")
	v.SetDefault("Collection", "")
	v.SetDefault("Database", "")
	v.SetDefault("DBURL", "")
	v.SetDefault("Port", "8080")
	v.SetDefault("Version", Version)

	v.AutomaticEnv()
	v.SetConfigType("json")
	v.SetConfigFile(v.GetString("ConfigPath"))

	err := v.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	var c Config
	err = v.Unmarshal(&c)

	return c, err
}
