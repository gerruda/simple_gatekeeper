package welcome

import (
	"context"
	"errors"
	"fmt"
	"github.com/gerruda/simple_gatekeeper/internal/ports"
)

var (
	msgText = "Для продолжения общения примите правила чата, нажав кнопку ниже"

	ErrNotYou = errors.New("Нажать принять правила могут только новые пользователи")
)

type WelcomeService struct {
	greater    ports.Greater
	moderator  ports.Moderator
	adminCheck ports.AdminChecker
	pending    ports.PendingRepository
	messanger  ports.Messenger
}

func NewWelcomeService(greater ports.Greater, moderator ports.Moderator, adminCheck ports.AdminChecker, pending ports.PendingRepository, messanger ports.Messenger) *WelcomeService {
	return &WelcomeService{
		greater:    greater,
		moderator:  moderator,
		adminCheck: adminCheck,
		pending:    pending,
		messanger:  messanger,
	}
}

func (w *WelcomeService) FirstMeet(ctx context.Context, chatID, userID int64) error {
	admin, errIsAdmin := w.adminCheck.IsAdmin(ctx, chatID, userID)
	if errIsAdmin != nil {
		return fmt.Errorf("[service:welcome:FirstMeet] Ошибка получения статуса админа у пользователя: %w", errIsAdmin)
	}
	if admin {
		return nil
	}

	waiting, errPending := w.pending.Exists(ctx, chatID, userID)
	if errPending != nil {
		return fmt.Errorf("[service:welcome:FirstMeet] Ошибка получения статуса ожидания принятия правил у пользователя: %w", errPending)
	}
	if waiting {
		return nil
	}

	errAdd := w.pending.Add(ctx, chatID, userID)
	if errAdd != nil {
		return fmt.Errorf("[service:welcome:FirstMeet] Ошибка добавления в очередь на ожидания принятия правил: %w", errAdd)
	}

	errBun := w.moderator.BunUser(ctx, chatID, userID)
	if errBun != nil {
		_ = w.pending.Remove(ctx, chatID, userID)

		return fmt.Errorf("[service:welcome:FirstMeet] Ошибка попытки ограничений парв пользователя: %w", errBun)
	}

	msgErr := w.greater.SendGreetingMsg(ctx, chatID, userID, msgText)
	if msgErr != nil {
		_ = w.moderator.UnbunUser(ctx, chatID, userID)
		_ = w.pending.Remove(ctx, chatID, userID)

		return fmt.Errorf("[service:welcome:FirstMeet] Ошибка отправки приветственного сообщения: %w", msgErr)
	}

	return nil
}

func (w *WelcomeService) OnAcceptRules(ctx context.Context, chatID, actorID, targetID int64) error {
	if actorID != targetID {
		return ErrNotYou
	}

	errUnbun := w.moderator.UnbunUser(ctx, chatID, actorID)
	if errUnbun != nil {
		return fmt.Errorf("[service:welcome:OnAcceptRules] Ошибка попытки снятия ограничений парв пользователя: %w", errUnbun)
	}

	errRemove := w.pending.Remove(ctx, chatID, actorID)
	if errRemove != nil {
		return errRemove
	}

	return nil
}
