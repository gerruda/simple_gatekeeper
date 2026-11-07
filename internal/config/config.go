package config

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gerruda/simple_gatekeeper/internal/logger"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	logger *logger.Logger

	Token string
	Proxy string
}

func NewConfig(ctx context.Context, logger *logger.Logger, path string) (*Config, error) {
	c := &Config{}
	c.logger = logger

	err := c.LoadConfig(ctx, path)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) LoadConfig(ctx context.Context, path string) error {
	if err := godotenv.Load(filepath.Join(path, ".env")); err != nil {
		c.logger.Logger.Sugar().Errorf("[configs:LoadConfig] не удалось считать файл .env file, %s \n", err)
	}

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(path)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		c.logger.Logger.Sugar().Errorf("[configs:LoadConfig] ене удалось считать файл config.yaml, %s \n", err)

		return fmt.Errorf("[configs:LoadConfig] ошибка получения конфигов из файла")
	}

	c.Token = v.GetString("bot.token")
	c.Proxy = v.GetString("bot.proxy")

	return nil
}
