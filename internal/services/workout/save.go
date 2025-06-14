package workout

import (
	"context"

	"github.com/google/uuid"
	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/repositories/models"
	"github.com/meehighlov/workout/internal/repositories/user"
)

func (s *Service) SaveWorkout(ctx context.Context, update *telegram.Update) error {
	user, err := s.repositories.User.Get(ctx, &user.Filter{TgChatID: update.GetChatIdStr()}, nil)
	if err != nil {
		return err
	}

	elements := s.clients.Cache.GetWorkoutElements(update.GetChatIdStr())
	drills := models.Drills{}
	for _, element := range elements {
		drills = append(drills, models.Drill{
			ElementName: element,
			Sets:        []models.DrillSet{},
		})
	}

	if len(elements) == 0 {
		_, err = s.clients.Telegram.Edit(ctx, s.constants.WORKOUT_NOT_SAVED_MESSAGE, update)
		if err != nil {
			return err
		}
		return nil
	}

	workoutUpdate := &models.Workout{
		ID:     uuid.New(),
		Name:   s.builders.ShortIdBuilder.Build(),
		UserID: user.ID,
		Drills: drills,
	}

	err = s.repositories.Workout.Save(ctx, workoutUpdate, nil)
	if err != nil {
		return err
	}

	s.clients.Cache.Reset(update.GetChatIdStr())

	keyboard := s.builders.KeyboardBuilder.Keyboard()
	keyboard.AppendAsLine(keyboard.NewButton(s.constants.BUTTON_TEXT_OPEN, s.builders.CallbackDataBuilder.Build(workoutUpdate.ID.String(), s.constants.COMMAND_INFO_WORKOUT, "0").String()))
	keyboard.AppendAsLine(keyboard.NewButton(s.constants.BUTTON_TEXT_ADD, s.builders.CallbackDataBuilder.Build(user.ID.String(), s.constants.COMMAND_NEW_WORKOUT, "0").String()))

	_, err = s.clients.Telegram.Edit(ctx, s.constants.WORKOUT_SAVED_MESSAGE, update, telegram.WithReplyMurkup(keyboard.Murkup()))
	if err != nil {
		return err
	}

	return nil
}
