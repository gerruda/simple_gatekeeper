package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gerruda/simple_gatekeeper/internal/adapters/memory"
	"github.com/gerruda/simple_gatekeeper/internal/adapters/telegram"
	"github.com/gerruda/simple_gatekeeper/internal/config"
	"github.com/gerruda/simple_gatekeeper/internal/logger"
	"github.com/gerruda/simple_gatekeeper/internal/services/welcome"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer func() {
		fmt.Println("[main] stopped gracfgully shuttdown")
		stop()
	}()

	// Infra
	log := logger.NewLogger(ctx)
	path, errGetPwd := os.Getwd()
	if errGetPwd != nil {
		path = "."
	}

	configs, errConfig := config.NewConfig(ctx, log, path)
	if errConfig != nil {
		log.Logger.Sugar().Errorf("[main] Ошибка получения конфигов: %v", errConfig)
	}

	// Client
	bot, err := telegram.NewBotClient(configs, log)
	if err != nil {
		panic(err)
	}
	memAdapter := memory.NewMemAdapter()

	// Services
	welcomeSvc := welcome.NewWelcomeService(bot, bot, bot, memAdapter, bot)

	// Adapter
	telegramAdapter := telegram.NewUpdatesRouter(bot, welcomeSvc)
	errUpdate := telegramAdapter.UpdatesRouter(ctx)
	if errUpdate != nil {
		log.Logger.Sugar().Debugf("[main] Ошибка получения обновлений: %v", errUpdate)
	}
}
