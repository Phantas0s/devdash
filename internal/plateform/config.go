package plateform

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
)

type config struct {
	Test string `mapstructure:"test"`
}

func MapConfig(data []byte) config {
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		panic(fmt.Errorf("could not read config data %s: %s \n", string(data), err))
	}

	var cfg config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
