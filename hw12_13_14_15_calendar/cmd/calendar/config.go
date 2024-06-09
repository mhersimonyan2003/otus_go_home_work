package main

import "github.com/spf13/viper"

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.

type LoggerConf struct {
	Level string
	// TODO
}

type Config struct {
	Logger     LoggerConf
	HTTPServer struct {
		Host string
		Port int
	}
	Storage struct {
		Type string
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

func NewConfig() *Config {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic("failed to read config file: " + err.Error())
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic("failed to unmarshal config: " + err.Error())
	}

	return &config
}
