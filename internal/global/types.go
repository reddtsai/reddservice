package global

type Configuration struct {
	LogOpts  LogOptions                   `mapstructure:"log-options"`
	SQLOpts  map[string]SQLOptions        `mapstructure:"sql-options"`
	GrpcOpts map[string]GrpcClientOptions `mapstructure:"grpc-client-options"`
}

type LogOptions struct {
	Level int `mapstructure:"level"` // -1,0,1,2,3,4,5
}

type GrpcClientOptions struct {
	Addr string `mapstructure:"addr"`
}

type SQLOptions struct {
	Addr        string `mapstructure:"addr"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	DB          string `mapstructure:"db"`
	MaxLifetime int    `mapstructure:"max-lifetime"` // 分鐘
	MaxOpenConn int    `mapstructure:"max-open-conn"`
	MaxIdleConn int    `mapstructure:"max-idle-conn"`
}

type PostgresqlConnSetting struct {
	DSN         string
	MaxOpenConn int
	MaxIdleConn int
	MaxLifetime int
}
