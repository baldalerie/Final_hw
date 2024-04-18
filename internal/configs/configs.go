package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	DataBase Database `yaml:"database"`
	Restapi  Api      `yaml:"api"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	TimeZone string `yaml:"timezone"`
}

type Api struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port" default:"8080"`
	Sslmode string `yaml:"sslmode"`
}

func LoadConfig(configFilePath string) (*Config, error) {
	viper.SetConfigFile(configFilePath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
