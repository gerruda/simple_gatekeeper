package welcome

import (
	"context"
	"fmt"
	"testing"

	"github.com/gerruda/simple_gatekeeper/internal/adapters/memory"
	"github.com/gerruda/simple_gatekeeper/internal/ports"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

var (
	supergroupID int64 = -100123
	privateID    int64 = 123
	groupID      int64 = -123
	userID       int64 = 456
)

func TestWelcomeFullPath(t *testing.T) {
	testCases := []struct {
		name                  string
		chatID                int64
		userID                int64
		acceptRulesID         int64
		admin                 bool
		newUser               bool
		wantErrorSendGeatings bool
		wantErrorAcceptRules  bool
	}{
		{
			name:                  "Игнорим админа",
			admin:                 true,
			newUser:               true,
			wantErrorSendGeatings: false,
			chatID:                groupID,
			userID:                userID,
			acceptRulesID:         userID,
		},
		{
			name:                  "Обычный новый пользователь",
			admin:                 false,
			newUser:               true,
			wantErrorSendGeatings: false,
			chatID:                groupID,
			userID:                userID,
			acceptRulesID:         userID,
		},
		{
			name:                  "Обычный новый пользователь в приватном чате с ботом",
			admin:                 false,
			newUser:               true,
			wantErrorSendGeatings: false,
			chatID:                privateID,
			userID:                userID,
			acceptRulesID:         userID,
		},
		{
			name:                 "Обычный новый пользователь, но принимает правила другой",
			admin:                false,
			newUser:              true,
			wantErrorAcceptRules: true,
			chatID:               groupID,
			userID:               userID,
			acceptRulesID:        123456,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(tt *testing.T) {
			tt.Parallel()

			pending := memory.NewMemAdapter()
			mock := NewMock(testCase.newUser, testCase.admin)

			w := NewWelcomeService(mock, mock, mock, pending, mock)

			err := w.FirstMeet(context.Background(), testCase.chatID, testCase.userID)
			if testCase.wantErrorSendGeatings {
				assert.Error(tt, err, "Должны получить ошибку")

				return
			}

			if testCase.admin {
				baned := false
				for _, ids := range mock.bunedUsers {
					if ids == [2]int64{testCase.chatID, testCase.userID} {
						baned = true
					}
				}
				assert.False(tt, baned, "Админа не должны банить")
				assert.Empty(tt, mock.msges, "Админам сообщение не отправляем")

				return
			}

			errAccept := w.OnAcceptRules(context.Background(), testCase.chatID, testCase.userID, testCase.acceptRulesID)
			if testCase.wantErrorAcceptRules {
				assert.Error(tt, errAccept, "Должны получить ошибку")
				assert.ErrorIs(tt, errAccept, ErrNotYou, "Ошибка должна совпадать")
				assert.Equal(tt, [][2]int64{{testCase.chatID, testCase.userID}}, mock.bunedUsers, "Должен быть разбан")

				return
			}

			baned := false
			for _, ids := range mock.bunedUsers {
				if ids == [2]int64{testCase.chatID, testCase.userID} {
					baned = true
				}
			}
			assert.True(tt, baned, "Подтверждение бана")

			unbaned := false
			for _, ids := range mock.unbunedUsers {
				if ids == [2]int64{testCase.chatID, testCase.userID} {
					unbaned = true
				}
			}
			assert.True(tt, unbaned, "Подтверждение разбана")

			assert.NotEmpty(tt, mock.msges, "Сообщения отправлены")
		})
	}
}

type MockWelcome struct {
	isNewUser    bool
	isAdmin      bool
	msges        []tgbotapi.MessageConfig
	bunedUsers   [][2]int64
	unbunedUsers [][2]int64
}

var (
	_ ports.Moderator    = (*MockWelcome)(nil)
	_ ports.Messenger    = (*MockWelcome)(nil)
	_ ports.AdminChecker = (*MockWelcome)(nil)
	_ ports.Greater      = (*MockWelcome)(nil)
)

func NewMock(isNewUser, isAdmin bool) *MockWelcome {
	return &MockWelcome{
		isNewUser: isNewUser,
		isAdmin:   isAdmin,
	}
}

func (m *MockWelcome) IsAdmin(ctx context.Context, chatID, userID int64) (bool, error) {
	return m.isAdmin, nil
}

func (m *MockWelcome) SendGreetingMsg(ctx context.Context, chatID, userID int64, msgText string) error {
	keyboardButtonUrl := fmt.Sprintf("Accept:%d", userID)

	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{
			tgbotapi.NewInlineKeyboardButtonData("Принимаю", keyboardButtonUrl),
		}},
	}

	m.msges = append(m.msges, msg)

	return nil
}

func (m *MockWelcome) BunUser(ctx context.Context, chatID, userID int64) error {
	m.bunedUsers = append(m.bunedUsers, [2]int64{chatID, userID})

	return nil
}

func (m *MockWelcome) UnbunUser(ctx context.Context, chatID, userID int64) error {
	m.unbunedUsers = append(m.unbunedUsers, [2]int64{chatID, userID})

	return nil
}

func (m *MockWelcome) SndMsg(ctx context.Context, chatID int64, msgText string) error {
	msg := tgbotapi.NewMessage(chatID, msgText)

	m.msges = append(m.msges, msg)

	return nil
}
