package config



type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}