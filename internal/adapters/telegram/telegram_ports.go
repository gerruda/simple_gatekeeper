package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *BotClient) SendGreetingMsg(ctx context.Context, chatID, userID int64, msgText string) error {
	buttonCallbackData := fmt.Sprintf("Accept:%d", userID)

	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{
			tgbotapi.NewInlineKeyboardButtonData("Принимаю", buttonCallbackData),
		}},
	}

	_, err := a.api.Send(msg)
	if err != nil {
		a.logger.Logger.Sugar().Errorf("[adapter:telegram:SendGreatingMsg] Ошибка при сообщения: %v", err)
		return err
	}

	return nil
}
