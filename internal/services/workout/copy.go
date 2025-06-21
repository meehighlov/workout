package workout

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/repositories/models"
	"github.com/meehighlov/workout/internal/repositories/workout"
)

func (s *Service) Copy(ctx context.Context, update *telegram.Update) error {
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	workoutToCopy, err := s.repositories.Workout.Get(ctx, &workout.Filter{ID: params.ID}, nil)
	if err != nil {
		return err
	}

	keyboard := s.builders.KeyboardBuilder.Keyboard()
	keyboard.AppendAsLine(
		keyboard.NewButton(s.constants.BUTTON_TEXT_BACK, s.builders.CallbackDataBuilder.Build(workoutToCopy.ID.String(), s.constants.COMMAND_INFO_WORKOUT, params.Offset).String()),
		keyboard.NewButton(s.constants.BUTTON_TEXT_COPY, s.builders.CallbackDataBuilder.Build(workoutToCopy.ID.String(), s.constants.COMMAND_COPY_WORKOUT_CONFIRM, params.Offset).String()),
	)

	header := fmt.Sprintf("Вы уверены, что хотите скопировать тренировку %s?", workoutToCopy.Name)
	s.clients.Telegram.Edit(ctx, header, update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}

func (s *Service) CopyConfirm(ctx context.Context, update *telegram.Update) error {
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	workoutToCopy, err := s.repositories.Workout.Get(ctx, &workout.Filter{ID: params.ID}, nil)
	if err != nil {
		return err
	}

	drills := models.Drills{}
	for _, drill := range workoutToCopy.Drills {
		drills = append(drills, models.Drill{
			ElementName:            drill.ElementName,
			CurrentlyObesrvableSet: 0,
			Sets:                   []models.DrillSet{},
		})
	}

	copiedWorkout := &models.Workout{
		ID:     uuid.New(),
		UserID: workoutToCopy.UserID,
		Name:   workoutToCopy.Name,
		Drills: drills,
		Status: workoutToCopy.Status,
	}

	err = s.repositories.Workout.Save(ctx, copiedWorkout, nil)
	if err != nil {
		return err
	}

	keyboard := s.builders.KeyboardBuilder.Keyboard()
	keyboard.AppendAsLine(keyboard.NewButton(s.constants.BUTTON_TEXT_OPEN, s.builders.CallbackDataBuilder.Build(copiedWorkout.ID.String(), s.constants.COMMAND_INFO_WORKOUT, "0").String()))
	keyboard.AppendAsLine(keyboard.NewButton(s.constants.BUTTON_TEXT_BACK, s.builders.CallbackDataBuilder.Build(workoutToCopy.ID.String(), s.constants.COMMAND_LIST_WORKOUT, params.Offset).String()))

	header := fmt.Sprintf("%s: %s", s.constants.WORKOUT_COPIED_MESSAGE, copiedWorkout.Name)
	s.clients.Telegram.Edit(ctx, header, update, telegram.WithReplyMurkup(keyboard.Murkup()))

	return nil
}
