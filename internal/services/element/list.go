package element

import (
	"context"
	"strconv"

	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/repositories/element"
	"github.com/meehighlov/workout/internal/repositories/user"
)

func (s *Service) List(ctx context.Context, update *telegram.Update) error {
	user, err := s.repositories.User.Get(ctx, &user.Filter{TgChatID: update.GetChatIdStr()}, nil)
	if err != nil {
		return err
	}

	filter := &element.Filter{
		UserID: user.ID.String(),
	}

	count, err := s.repositories.Element.Count(ctx, &element.Filter{UserID: user.ID.String()})
	if err != nil {
		return err
	}

	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	offset, err := strconv.Atoi(params.Offset)
	if err != nil {
		return err
	}

	filter.Limit = s.pagination.Limit
	filter.Offset = offset

	elements, err := s.repositories.Element.List(ctx, filter)
	if err != nil {
		return err
	}

	keyboard := s.builders.KeyboardBuilder.Keyboard()

	for _, element := range elements {
		keyboard.AppendAsLine(keyboard.NewButton(element.Name, s.builders.CallbackDataBuilder.Build(element.ID.String(), s.constants.COMMAND_ELEMENT_INFO, params.Offset).String()))
	}

	keyboard.
	OptimizeView().
	Append(
		s.pagination.BuildControls(
			int(count),
			s.constants.COMMAND_LIST_ELEMENT,
			update.GetChatIdStr(),
			offset,
		),
	).
	PutFirstAsLine(
		keyboard.NewButton(s.constants.BUTTON_TEXT_ADD_ELEMENT, s.builders.CallbackDataBuilder.Build(user.ID.String(), s.constants.COMMAND_ADD_ELEMENT, params.Offset).String()),
	)

	if update.IsCallback() {
		_, err = s.clients.Telegram.Edit(ctx, "Мои элементы", update, telegram.WithReplyMurkup(keyboard.Murkup()))
	} else {
		_, err = s.clients.Telegram.Reply(ctx, "Мои элементы", update, telegram.WithReplyMurkup(keyboard.Murkup()))
	}

	if err != nil {
		return err
	}

	return nil
}
