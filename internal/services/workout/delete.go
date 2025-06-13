package workout

import (
	"context"
	"fmt"

	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/repositories/workout"
)

func (s *Service) Delete(ctx context.Context, update *telegram.Update) error {
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	workoutToDelete, err := s.repositories.Workout.Get(ctx, &workout.Filter{ID: params.ID}, nil)
	if err != nil {
		return err
	}

	keyboard := s.builders.KeyboardBuilder.Keyboard()
	keyboard.AppendAsLine(
		keyboard.NewButton(s.constants.BUTTON_TEXT_BACK, s.builders.CallbackDataBuilder.Build(workoutToDelete.ID.String(), s.constants.COMMAND_INFO_WORKOUT, params.Offset).String()),
		keyboard.NewButton(s.constants.BUTTON_TEXT_DELETE, s.builders.CallbackDataBuilder.Build(workoutToDelete.ID.String(), s.constants.COMMAND_DELETE_WORKOUT_CONFIRM, params.Offset).String()),
	)

	header := fmt.Sprintf("Вы уверены, что хотите удалить тренировку %s?", workoutToDelete.Name)
	s.clients.Telegram.Edit(ctx, header, update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}

func (s *Service) DeleteConfirm(ctx context.Context, update *telegram.Update) error {
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	workoutToDelete, err := s.repositories.Workout.Get(ctx, &workout.Filter{ID: params.ID}, nil)
	if err != nil {
		return err
	}

	err = s.repositories.Workout.Delete(ctx, workoutToDelete.ID, nil)
	if err != nil {
		return err
	}

	keyboard := s.builders.KeyboardBuilder.Keyboard()
	keyboard.AppendAsLine(keyboard.NewButton(s.constants.BUTTON_TEXT_BACK, s.builders.CallbackDataBuilder.Build(workoutToDelete.ID.String(), s.constants.COMMAND_LIST_WORKOUT, params.Offset).String()))

	header := fmt.Sprintf("Тренировка %s удалена", workoutToDelete.Name)
	s.clients.Telegram.Edit(ctx, header, update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}
