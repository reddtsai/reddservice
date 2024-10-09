package global

type Configuration struct {
	LogOpts LogOptions `mapstructure:"log-options"`
}

type LogOptions struct {
	Level int `mapstructure:"level"` // -1,0,1,2,3,4,5
}
