// Package config содержит функции и структуры для работы с конфигурацией приложения.
package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

// Config хранит все конфигурации приложения.
type Config struct {
	AppConfig AppConfig `mapstructure:"app"`
	DBConfig  DBConfig  `mapstructure:"db"`
}

// AppConfig хранит конфигурацию приложения.
type AppConfig struct {
	Env         string        `mapstructure:"env"`
	Port        string        `mapstructure:"port"`
	TimeOut     time.Duration `mapstructure:"timeout"`
	Idletimeout time.Duration `mapstructure:"idle_timeout"`
}

// DBConfig хранит конфигурацию подключения к базе данных.
type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

// InitConfig загружает конфигурацию из файла и окружения и возвращает структуру Config.
func InitConfig() *Config {

	configPath := getConfigPath()

	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading configfile: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	return &cfg

}

func getConfigPath() string {

	err := gotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file: ", err)
	}

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH IS NOT SET")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG FILE DOES NOT EXIST %s", configPath)
	}
	return configPath
}
