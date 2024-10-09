package global

import (
	"github.com/spf13/viper"
)

const (
	CONFIF_FILE = "config/config.yaml"
)

func initConfiguration() {
	configOnce.Do(func() {
		viper.AutomaticEnv()
		viper.SetDefault(KEY_LOG_LEVEL, -1) // -1 ~ 5

		viper.SetConfigFile(CONFIF_FILE)
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}

		err = viper.Unmarshal(&Config)
		if err != nil {
			panic(err)
		}
	})
}
