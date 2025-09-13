package platform

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type DbConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	Name string `mapstructure:"name"`
}

type Config struct {
	Port    string   `mapstructure:"port"`
	BaseUrl string   `mapstructure:"baseUrl"`
	DB      DbConfig `mapstructure:"db"`
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("Port", "8080")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No config file found, using env/defaults")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to decode config: %v", err)
	}

	return &cfg
}
