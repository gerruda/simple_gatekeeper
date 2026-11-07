package telegram

import (
	"context"
	"fmt"
	"github.com/gerruda/simple_gatekeeper/internal/config"
	"github.com/gerruda/simple_gatekeeper/internal/logger"
	"github.com/gerruda/simple_gatekeeper/internal/ports"
	"net/http"
	"net/url"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap/zapcore"
)

var (
	_ ports.Moderator    = (*BotClient)(nil)
	_ ports.Messenger    = (*BotClient)(nil)
	_ ports.AdminChecker = (*BotClient)(nil)
	_ ports.Greater      = (*BotClient)(nil)
)

type BotClient struct {
	api    *tgbotapi.BotAPI
	logger *logger.Logger
}

func NewHTTPClient(proxyURL string) (*http.Client, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	if proxyURL != "" {
		parsed, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("parse telegram proxy URL: %v", err)
		}
		if parsed.Host == "" {
			return nil, fmt.Errorf("parse telegram proxy URL %q: host is required", proxyURL)
		}
		if parsed.Scheme != "http" && parsed.Scheme != "https" {
			return nil, fmt.Errorf("unsupported telegram proxy scheme %q: use http://host:port", parsed.Scheme)
		}

		transport.Proxy = http.ProxyURL(parsed)
	}

	return &http.Client{
		Transport: transport,
		Timeout:   8 * time.Second,
	}, nil
}

func NewBotClient(configs *config.Config, log *logger.Logger) (*BotClient, error) {
	client, err := NewHTTPClient(configs.Proxy)
	if err != nil {
		return nil, fmt.Errorf("[bot:NewBot] ошибка инициализации клиента для работы бота: %v", err.Error())
	}

	bot, err := tgbotapi.NewBotAPIWithClient(configs.Token, "https://api.telegram.org/bot%s/%s", client)
	if err != nil {
		return nil, fmt.Errorf("[bot:NewBot] ошибка инициализации бота: %v", err.Error())
	}

	if bot != nil {
		bot.Debug = true
	}

	log.Logger.Log(zapcore.DebugLevel, fmt.Sprintf("[bot:NewBot] аторизовались и запустили бота на аккуаунте %s", bot.Self.UserName))

	return &BotClient{
		api:    bot,
		logger: log,
	}, nil
}

func (a *BotClient) SndMsg(ctx context.Context, chatID int64, msgText string) error {
	msg := tgbotapi.NewMessage(chatID, msgText)

	_, err := a.api.Send(msg)
	if err != nil {
		a.logger.Logger.Sugar().Errorf("[bot:NewBot] ошибка отправки сообщения через сервис: %v", err.Error())

		return err
	}

	return nil
}

func (a *BotClient) BunUser(ctx context.Context, chatID, userID int64) error {
	config := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		UntilDate: 0,

		Permissions: &tgbotapi.ChatPermissions{
			CanSendMessages:       false,
			CanSendMediaMessages:  false,
			CanSendPolls:          false,
			CanAddWebPagePreviews: false,
		},
	}

	_, err := a.api.Request(config)
	if err != nil {
		a.logger.Logger.Sugar().Errorf("[adapter:telegram:BunUser] Ошибка при ограничении пользователя: %v", err)

		return err
	}

	return nil
}

func (a *BotClient) UnbunUser(ctx context.Context, chatID, userID int64) error {
	config := tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
		UntilDate: 0,

		Permissions: &tgbotapi.ChatPermissions{
			CanSendMessages:       true,
			CanSendMediaMessages:  true,
			CanSendPolls:          true,
			CanSendOtherMessages:  true,
			CanAddWebPagePreviews: true,
		},
	}

	_, err := a.api.Request(config)
	if err != nil {
		a.logger.Logger.Sugar().Errorf("[adapter:telegram:UnbunUser] Ошибка при снятии ограничении пользователя: %v", err)

		return err
	}

	return nil
}

func (a *BotClient) IsAdmin(ctx context.Context, chatID, userID int64) (bool, error) {
	reqConfig := tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	}

	member, err := a.api.GetChatMember(reqConfig)
	if err != nil {
		a.logger.Logger.Sugar().Errorf("[adapter:telegram:IsAdmin] Ошибка проверки пользователя на статус админа: %v", err)
		return false, err
	}

	return member.Status == "creator" || member.Status == "administrator", err
}
