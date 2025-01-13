package config

import "github.com/spf13/viper"

type ServerConfig struct {
	Port         int `mapstructure:"port"`
	IdleTimeout  int `mapstructure:"idle_timeout"`
	ReadTimeout  int `mapstructure:"read_timeout"`
	WriteTimeout int `mapstructure:"write_timeout"`
}

type Config struct {
	Server ServerConfig `mapstructure:"server"`
}

func ReadConfig() (*Config, error) {
	// read config from config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// read config from environment variables
	viper.AutomaticEnv()

	// unmarshal config into struct
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// return
	return cfg, nil
}
