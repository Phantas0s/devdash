package main

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
)

type config struct {
	GoogleAnalytics GoogleAnalytics `mapstructure:"google_analytics"`
}

type GoogleAnalytics struct {
	Keyfile string `mapstructure:"keyfile"`
	ViewID  string `mapstructure:"view_id"`
}

func mapConfig(data []byte) config {
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		panic(fmt.Errorf("could not read config data %s: %s", string(data), err))
	}

	var cfg config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
