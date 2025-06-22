package workout

import (
	"context"

	"github.com/google/uuid"
	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/repositories/models"
	"github.com/meehighlov/workout/internal/repositories/user"
	"github.com/meehighlov/workout/internal/repositories/workout"
)

func (s *Service) SaveWorkout(ctx context.Context, update *telegram.Update) error {
	user, err := s.repositories.User.Get(ctx, &user.Filter{TgChatID: update.GetChatIdStr()}, nil)
	if err != nil {
		return err
	}

	elements := s.clients.Cache.GetWorkoutElements(ctx, update.GetChatIdStr())
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
		Status: models.WORKOUT_STATUS_ACTIVE,
	}

	workoutId := s.clients.Cache.GetWorkoutID(ctx, update.GetChatIdStr())
	if workoutId != "" {
		workout, err := s.repositories.Workout.Get(ctx, &workout.Filter{ID: workoutId}, nil)
		if err != nil {
			return err
		}
		workoutUpdate = workout
		workoutUpdate.Drills = s.mergeDrills(workout.Drills, drills)
	}

	err = s.repositories.Workout.Save(ctx, workoutUpdate, nil)
	if err != nil {
		return err
	}

	s.clients.Cache.Reset(ctx, update.GetChatIdStr())

	keyboard := s.builders.KeyboardBuilder.Keyboard()
	keyboard.AppendAsLine(keyboard.NewButton(s.constants.BUTTON_TEXT_OPEN, s.builders.CallbackDataBuilder.Build(workoutUpdate.ID.String(), s.constants.COMMAND_INFO_WORKOUT, "0").String()))
	keyboard.AppendAsLine(keyboard.NewButton(s.constants.BUTTON_TEXT_ADD, s.builders.CallbackDataBuilder.Build(user.ID.String(), s.constants.COMMAND_NEW_WORKOUT, "0").String()))

	_, err = s.clients.Telegram.Edit(ctx, s.constants.WORKOUT_SAVED_MESSAGE, update, telegram.WithReplyMurkup(keyboard.Murkup()))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) mergeDrills(workoutDrills []models.Drill, newDrills []models.Drill) models.Drills {
	var mergedDrills models.Drills

	seen := make(map[string]*models.Drill)
	for _, drill := range workoutDrills {
		seen[drill.ElementName] = &drill
	}

	for _, drill := range newDrills {
		if existingDrill, ok := seen[drill.ElementName]; ok {
			mergedDrills = append(mergedDrills, *existingDrill)
		} else {
			mergedDrills = append(mergedDrills, drill)
		}
	}

	return mergedDrills
}
