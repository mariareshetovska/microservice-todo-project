package config

import "github.com/spf13/viper"

type Config struct {
	Db     DbConfig     `mapstructure:"db"`
	Server ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DbConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBname   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

var vp *viper.Viper

func LoadConfig() (Config, error) {
	vp = viper.New()
	var config Config

	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath("./config")
	vp.AddConfigPath(".")
	err := vp.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	err = vp.Unmarshal((&config))
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
