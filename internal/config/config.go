package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type StorageConfig struct {
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

var (
	instance *StorageConfig
	once     sync.Once
)

func GetConfig() *StorageConfig {
	once.Do(func() {
		log.Print("read application configuration")
		instance = &StorageConfig{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
