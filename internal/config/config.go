package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string      `yaml:"env" env-default:"prod"`
	HTTP_Server HTTP_Server `yaml:"http_server"`
	DB          DB          `yaml:"db"`
}

type HTTP_Server struct {
	Address     string        `yaml"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DB struct {
	Host     string `yaml:"host" env-default:"board_db"`
	Port     string `yaml:"port" env-default:"5432"`
	Name     string `yaml:"name" env-default:"board_db"`
	User     string `yaml:"user" env-default:"board_user"`
	Password string `yaml:"password" env-default:"pwd123"`
	Sslmode  string `yaml:"sslmode" env-default:"disable"`
}

func (db DB) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		db.Host, db.Port, db.Name, db.User, db.Password, db.Sslmode)
}

func MustLoad() *Config {

	if err := godotenv.Load("local.env"); err != nil && !os.IsNotExist(err) {
		log.Fatalf("error loading .env file: %v", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(os.Getenv("CONFIG_PATH"), &cfg); err != nil {
		log.Fatalf("Cannot read config: %v:%v", err, os.Getenv("CONFIG_PATH"))
	}

	return &cfg
}
