package config

import (
	"os"
	"telegram-bot/internal/logger"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const configFile = "data/config.yaml"

type Config struct {
	// TODO лучше сделать всё в переменных среды?
	// Добавить переменные для cron и джобы (время выполнения, путь и т.д.)
	Token              string `yaml:"token"`              // Токен бота в телеграме.
	ConnectionStringDB string `yaml:"ConnectionStringDB"` // Строка подключения в базе данных.
	MaxAttempts        int    `yaml:"MaxAttempts"`
}

type Service struct {
	config Config
}

func New() (*Service, error) {
	s := &Service{}

	rawYAML, err := os.ReadFile(configFile)
	if err != nil {
		logger.Error("Ошибка reading config file", "ERROR", err)
		return nil, errors.Wrap(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &s.config)
	if err != nil {
		logger.Error("Ошибка parsing yaml", "ERROR", err)
		return nil, errors.Wrap(err, "parsing yaml")
	}

	return s, nil
}

func (s *Service) Token() string {
	return s.config.Token
}

func (s *Service) GetConfig() Config {
	return s.config
}
