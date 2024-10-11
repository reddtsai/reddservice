package global

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	CONFIG_FILE = "config.yaml"
)

func initConfiguration() {
	configOnce.Do(func() {
		viper.AutomaticEnv()
		viper.AddConfigPath(CONFIG_PATH)
		viper.SetConfigType("yaml")

		viper.SetConfigName(CONFIG_FILE)
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

func GetPostgresqlConnSetting(name string) PostgresqlConnSetting {
	opts := Config.SQLOpts[name]
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", opts.Username, opts.Password, opts.Addr, opts.DB)
	return PostgresqlConnSetting{
		DSN:         dsn,
		MaxOpenConn: opts.MaxOpenConn,
		MaxIdleConn: opts.MaxIdleConn,
		MaxLifetime: opts.MaxLifetime,
	}
}
