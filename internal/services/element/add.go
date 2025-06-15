package element

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/repositories/models"
	"github.com/meehighlov/workout/internal/repositories/user"
)

func (s *Service) Add(ctx context.Context, update *telegram.Update) error {
	s.clients.Telegram.Reply(
		ctx,
		"Введите название элемента",
		update,
		telegram.WithReplyMurkup(s.builders.KeyboardBuilder.Keyboard().Murkup()),
	)

	s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), s.constants.COMMAND_ADD_ELEMENT_SAVE)

	return nil
}

func (s *Service) AddSave(ctx context.Context, update *telegram.Update) error {
	user, err := s.repositories.User.Get(ctx, &user.Filter{TgChatID: update.GetChatIdStr()}, nil)
	if err != nil {
		return err
	}

	name := strings.ToLower(strings.TrimSpace(update.Message.Text))

	element := &models.Element{
		ID:     uuid.New(),
		UserID: user.ID,
		Name:   name,
	}

	err = s.repositories.Element.Save(ctx, element, nil)
	if err != nil {
		return err
	}

	s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), "")

	keyboard := s.builders.KeyboardBuilder.Keyboard()
	keyboard.AppendAsLine(keyboard.NewButton(s.constants.BUTTON_TEXT_OPEN, s.builders.CallbackDataBuilder.Build(element.ID.String(), s.constants.COMMAND_INFO_ELEMENT, "0").String()))
	keyboard.AppendAsLine(keyboard.NewButton(s.constants.BUTTON_TEXT_ADD, s.builders.CallbackDataBuilder.Build(user.ID.String(), s.constants.COMMAND_ADD_ELEMENT, "0").String()))

	s.clients.Telegram.Reply(ctx, "Элемент добавлен", update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}
