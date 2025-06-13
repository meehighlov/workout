package element

import (
	"context"

	"github.com/meehighlov/workout/internal/clients/telegram"
)

func (s *Service) List(ctx context.Context, update *telegram.Update) error {
	keyboard, err := s.BuildElementsKeyboard(
		ctx,
		update,
		s.constants.COMMAND_ELEMENT_INFO,
		s.constants.COMMAND_LIST_ELEMENT,
	)
	if err != nil {
		return err
	}

	if update.IsCallback() {
		_, err = s.clients.Telegram.Edit(ctx, s.constants.ELEMENTS_LIST_MESSAGE, update, telegram.WithReplyMurkup(keyboard.Murkup()))
	} else {
		_, err = s.clients.Telegram.Reply(ctx, s.constants.ELEMENTS_LIST_MESSAGE, update, telegram.WithReplyMurkup(keyboard.Murkup()))
	}

	if err != nil {
		return err
	}

	return nil
}
