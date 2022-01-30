package service

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type MysqlConfig struct {
	Dsn string `yaml:"dsn"`
}
type RedisConfig struct {
	Host     string `yaml: "host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
	Addr     string `yaml:"addr"`
}

type Configuration struct {
	Server ServerConfig `yaml:"server"`
	Mysql  MysqlConfig  `yaml:"mysql"`
	Redis  RedisConfig  `yaml:"redis"`
}

var (
	configuration *Configuration
)

func ParseConfiguration(configPath string) error {
	v := viper.New()
	v.SetConfigFile(configPath)
	// v.SetConfigType(configType)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	var config Configuration
	if err := v.Unmarshal(&config); err != nil {
		log.Println("解析v.unmarshal出错" + err.Error())
		return err
	}

	configuration = &config
	return nil
}

// GetConfiguration ...
func GetConfiguration() *Configuration {
	err := ParseConfiguration("config.yaml")
	if err != nil {
		log.Println("解析yaml出错" + err.Error())
		return nil
	}
	return configuration
}
