package workout

import (
	"context"

	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/repositories/workout"
)

func (s *Service) Edit(ctx context.Context, update *telegram.Update) error {
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	workout, err := s.repositories.Workout.Get(ctx, &workout.Filter{ID: params.ID}, nil)
	if err != nil {
		s.clients.Telegram.Reply(ctx, "Возникла непредвиденная ошибка", update)
		return err
	}

	s.clients.Cache.SetWorkout(ctx, update.GetChatIdStr(), workout)

	keyboard := s.builders.KeyboardBuilder.Keyboard()
	keyboard.AppendAsLine(
		keyboard.NewButton(s.constants.BUTTON_TEXT_NAME, s.builders.CallbackDataBuilder.Build("name", s.constants.COMMAND_EDIT_WORKOUT_REQUEST, "0").String()),
		keyboard.NewButton(s.constants.BUTTON_TEXT_ELEMENTS_IN_WORKOUT, s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_EDIT_WORKOUT_DRILLS, "0").String()),
	)

	s.clients.Telegram.Reply(ctx, "Что будем редактировать?", update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}

func (s *Service) EditRequest(ctx context.Context, update *telegram.Update) error {
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	nextHandler := ""
	switch params.ID {
	case "name":
		s.clients.Telegram.Reply(ctx, "Введите новое название", update)
		nextHandler = s.constants.COMMAND_EDIT_WORKOUT_NAME_SAVE
	}

	s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), nextHandler)

	return nil
}

func (s *Service) EditNameSave(ctx context.Context, update *telegram.Update) error {
	workoutId := s.clients.Cache.GetWorkoutID(ctx, update.GetChatIdStr())

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

	s.clients.Cache.SetNextHandler(ctx, update.GetChatIdStr(), "")

	keyboard := s.builders.KeyboardBuilder.Keyboard()
	backButton := keyboard.NewButton(s.constants.BUTTON_TEXT_BACK, s.builders.CallbackDataBuilder.Build(workoutToEdit.ID.String(), s.constants.COMMAND_INFO_WORKOUT, "0").String())
	keyboard.AppendAsLine(backButton)

	s.clients.Telegram.Reply(ctx, "Тренировка обновлена", update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}
