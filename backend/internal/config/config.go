package config

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port         int `mapstructure:"port"`
	IdleTimeout  int `mapstructure:"idle_timeout"`
	ReadTimeout  int `mapstructure:"read_timeout"`
	WriteTimeout int `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Password           string `mapstructure:"password"`
	PingTimeout        int    `mapstructure:"ping_timeout"`
	QueryTimeout       int    `mapstructure:"query_timeout"`
	TransactionTimeout int    `mapstructure:"transaction_timeout"`
}

type InitialAdminConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	FullName string `mapstructure:"full_name"`
	Email    string `mapstructure:"email"`
}

type Config struct {
	Server       ServerConfig       `mapstructure:"server"`
	Database     DatabaseConfig     `mapstructure:"database"`
	InitialAdmin InitialAdminConfig `mapstructure:"initial_admin"`
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
	if err := viper.BindEnv("database.password", "POSTGRES_PASSWORD"); err != nil {
		return nil, err
	}
	if err := viper.BindEnv("initial_admin.password", "INITIAL_ADMIN_PASSWORD"); err != nil {
		return nil, err
	}

	// unmarshal config into struct
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// return
	return cfg, nil
}
