package platform

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type DbConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	Name string `mapstructure:"name"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Pass string `mapstructure:"pass"`
	DB   int    `mapstructure:"db"`
}

type CacheConfig struct {
	Enabled bool `mapstructure:"enabled"`

	HotURLTTL  time.Duration `mapstructure:"hotUrlTTL"`
	WarmURLTTL time.Duration `mapstructure:"warmUrlTTL"`
	ColdURLTTL time.Duration `mapstructure:"coldUrlTTL"`

	HotThreshold  int `mapstructure:"hotThreshold"`
	WarmThreshold int `mapstructure:"warmThreshold"`

	MaxTTL time.Duration `mapstructure:"maxTTL"`
	MinTTL time.Duration `mapstructure:"minTTL"`
}

type Config struct {
	Port    string      `mapstructure:"port"`
	BaseUrl string      `mapstructure:"baseUrl"`
	DB      DbConfig    `mapstructure:"db"`
	Redis   RedisConfig `mapstructure:"redis"`
	Cache   CacheConfig `mapstructure:"cache"`
}

func LoadConfig() *Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// APP
	viper.BindEnv("port", "PORT")
	viper.BindEnv("baseUrl", "BASEURL")

	// DB
	viper.BindEnv("db.host", "DB_HOST")
	viper.BindEnv("db.port", "DB_PORT")
	viper.BindEnv("db.user", "DB_USER")
	viper.BindEnv("db.pass", "DB_PASS")
	viper.BindEnv("db.name", "DB_NAME")

	// Redis
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.pass", "REDIS_PASS")
	viper.BindEnv("redis.db", "REDIS_DB")

	// Cache
	viper.BindEnv("cache.enabled", "CACHE_ENABLED")
	viper.BindEnv("cache.hotUrlTTL", "CACHE_HOT_TTL")
	viper.BindEnv("cache.warmUrlTTL", "CACHE_WARM_TTL")
	viper.BindEnv("cache.coldUrlTTL", "CACHE_COLD_TTL")
	viper.BindEnv("cache.hotThreshold", "CACHE_HOT_THRESHOLD")
	viper.BindEnv("cache.warmThreshold", "CACHE_WARM_THRESHOLD")
	viper.BindEnv("cache.maxTTL", "MAX_TTL")
	viper.BindEnv("cache.minTTL", "MIN_TTL")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No config file found, using env/defaults")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to decode config: %v", err)
	}

	fmt.Println("redis pass ________>>>", cfg.Redis.Host)

	return &cfg
}
