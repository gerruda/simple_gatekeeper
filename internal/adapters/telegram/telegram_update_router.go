package telegram

import (
	"context"
	"fmt"

	"github.com/gerruda/simple_gatekeeper/internal/services/welcome"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramRouter struct {
	BotClient      *BotClient
	welcomeService *welcome.WelcomeService
}

func NewUpdatesRouter(botApi *BotClient, welcomeservice *welcome.WelcomeService) *TelegramRouter {
	return &TelegramRouter{
		BotClient:      botApi,
		welcomeService: welcomeservice,
	}
}

func (a *TelegramRouter) UpdatesRouter(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 8

	updates := a.BotClient.api.GetUpdatesChan(u)
	defer a.BotClient.api.StopReceivingUpdates()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("[adapters:telegram:StartUpdates] ошибка контекста: %w", ctx.Err())
		case update, ok := <-updates:
			if !ok {
				return nil
			}

			switch {
			case update.Message != nil && len(update.Message.NewChatMembers) > 0:
				a.onNewUser(ctx, update)

			case update.CallbackQuery != nil:
				a.onCallback(ctx, update)

			case update.Message != nil && update.Message.IsCommand():
				a.onCommand(ctx, update)
			}
		}
	}
}
