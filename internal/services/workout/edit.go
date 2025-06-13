package workout

import (
	"context"

	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/repositories/workout"
)

func (s *Service) Edit(ctx context.Context, update *telegram.Update) error {
	s.clients.Telegram.Reply(ctx, "Введите новое название тренировки", update)

	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)
	s.clients.Cache.AppendText(update.GetChatIdStr(), params.ID)
	s.clients.Cache.SetNextHandler(update.GetChatIdStr(), s.constants.COMMAND_EDIT_WORKOUT_NAME_SAVE)

	return nil
}

func (s *Service) EditNameSave(ctx context.Context, update *telegram.Update) error {
	workoutId := s.clients.Cache.GetTexts(update.GetChatIdStr())[0]

	workoutToEdit, err := s.repositories.Workout.Get(ctx, &workout.Filter{ID: workoutId}, nil)
	if err != nil {
		s.clients.Telegram.Reply(ctx, "Возникла непредвиденная ошибка", update)
		return err
	}

	workoutToEdit.Name = update.Message.Text

	err = s.repositories.Workout.Save(ctx, workoutToEdit, nil)
	if err != nil {
		s.clients.Telegram.Reply(ctx, "Возникла непредвиденная ошибка", update)
		return err
	}

	s.clients.Cache.SetNextHandler(update.GetChatIdStr(), "")

	keyboard := s.builders.KeyboardBuilder.Keyboard()
	backButton := keyboard.NewButton(s.constants.BUTTON_TEXT_BACK, s.builders.CallbackDataBuilder.Build(workoutToEdit.ID.String(), s.constants.COMMAND_INFO_WORKOUT, "0").String())
	keyboard.AppendAsLine(backButton)

	s.clients.Telegram.Reply(ctx, "Тренировка обновлена", update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}
