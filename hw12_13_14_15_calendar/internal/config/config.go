package config

import (
	"github.com/mhersimonyan2003/otus_go_home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.

type LoggerConf struct {
	Level logger.Level
	// TODO
}

type Config struct {
	Logger     LoggerConf
	HTTPServer struct {
		Host string
		Port int
	}
	GRPCServer struct {
		Host string
		Port int
	}
	Storage struct {
		Type string
	}
	API struct {
		Mode string
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
}

func NewConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
