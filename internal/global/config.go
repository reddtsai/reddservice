package global

import (
	"fmt"

	"github.com/spf13/viper"
)

// this value is set by ldflags
// build -ldflags "-X github.com/reddtsai/reddservice/internal/global.CONFIG_PATH=conf.d"
// var CONFIG_PATH = "conf.d"

const (
	CONFIG_FILE = "config.yaml"
)

func loadConfiguration(configPath string) {
	configOnce.Do(func() {
		viper.AutomaticEnv()
		viper.AddConfigPath(configPath)
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

func GetGrpcClientOptions(name string) GrpcClientOptions {
	return Config.GrpcOpts[name]
}
