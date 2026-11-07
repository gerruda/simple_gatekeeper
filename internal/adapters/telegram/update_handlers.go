package telegram

import (
	"context"
	"errors"
	"github.com/gerruda/simple_gatekeeper/internal/services/welcome"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *TelegramRouter) onNewUser(ctx context.Context, update tgbotapi.Update) {
	for _, member := range update.Message.NewChatMembers {
		if member.IsBot {
			continue
		}

		a.BotClient.logger.Logger.Sugar().Infof("пользователь %v %v", update.Message.Chat.ID, member.ID)

		errFirsMeet := a.welcomeService.FirstMeet(ctx, update.Message.Chat.ID, member.ID)
		if errFirsMeet != nil {
			a.BotClient.logger.Logger.Sugar().Errorf("[bot:NewBot] ошибка приветствия пользователя: %v", errFirsMeet.Error())
		}
	}
}

func (a *TelegramRouter) onCommand(ctx context.Context, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch update.Message.Command() {
	case "help":
		msg.Text = "Бот должен быть админом в вашей группе. Тогда он сможет приветсвовать новичков."
	case "start":
		msg.Text = "Привет, я бот привраптник. Я буду ограничивать работу новичков в группе, пока они не согласятся с правилами. Подробнее в команде /help"
	default:
		msg.Text = "Нет такой команды"
	}

	errS := a.BotClient.SndMsg(ctx, msg.ChatID, msg.Text)
	if errS != nil {
		a.BotClient.logger.Logger.Sugar().Errorf("[bot:NewBot] ошибка отправки сообщения: %v", errS.Error())
	}
}

func (a *TelegramRouter) onCallback(ctx context.Context, update tgbotapi.Update) {
	cb := update.CallbackQuery
	if cb.Message == nil {
		return
	}

	parsedDataId := parseDataFromCallback(cb.Data)
	errorCallback := a.welcomeService.OnAcceptRules(ctx, cb.Message.Chat.ID, cb.From.ID, parsedDataId)

	switch {
	case errors.Is(errorCallback, welcome.ErrNotYou):
		newCallback := tgbotapi.NewCallbackWithAlert(cb.ID, "Это не ваша кнопка")
		a.BotClient.api.Request(newCallback)

		return
	case errorCallback != nil:
		newCallback := tgbotapi.NewCallbackWithAlert(cb.ID, "")
		a.BotClient.api.Request(newCallback)

		a.BotClient.logger.Logger.Sugar().Errorf("[bot:NewBot] ошибка принятия правил: %v", errorCallback.Error())

		return
	}

	deleteConfig := tgbotapi.NewDeleteMessage(cb.Message.Chat.ID, cb.Message.MessageID)

	_, err := a.BotClient.api.Request(deleteConfig)
	if err != nil {
		a.BotClient.logger.Logger.Sugar().Errorf("[bot:NewBot] ошибка удаление сообщения: %v", err.Error())
	}
}

func parseDataFromCallback(data string) int64 {
	userIdString, ok := strings.CutPrefix(data, "Accept:")
	if !ok {
		return int64(0)
	}

	userIdInt, err := strconv.Atoi(userIdString)
	if err != nil {
		return int64(0)
	}

	return int64(userIdInt)
}
