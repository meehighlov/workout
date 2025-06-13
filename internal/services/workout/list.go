package workout

import (
	"context"

	"github.com/meehighlov/workout/internal/clients/telegram"
)

func (s *Service) ListWorkouts(ctx context.Context, update *telegram.Update) error {
	keyboard, err := s.BuildWorkoutsKeyboard(
		ctx,
		update,
		s.constants.COMMAND_INFO_WORKOUT,
		s.constants.COMMAND_LIST_WORKOUT,
	)
	if err != nil {
		return err
	}

	if update.IsCallback() {
		_, err = s.clients.Telegram.Edit(ctx, s.constants.WORKOUT_LIST_MESSAGE, update, telegram.WithReplyMurkup(keyboard.Murkup()))
	} else {
		_, err = s.clients.Telegram.Reply(ctx, s.constants.WORKOUT_LIST_MESSAGE, update, telegram.WithReplyMurkup(keyboard.Murkup()))
	}

	if err != nil {
		return err
	}

	return nil
}
