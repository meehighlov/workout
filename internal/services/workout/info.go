package workout

import (
	"context"
	"fmt"
	"strconv"

	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/repositories/models"
	"github.com/meehighlov/workout/internal/repositories/workout"
)

func (s *Service) InfoWorkout(ctx context.Context, update *telegram.Update) error {
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	workout, err := s.repositories.Workout.Get(ctx, &workout.Filter{ID: params.ID}, nil)
	if err != nil {
		return err
	}

	keyboard := s.builders.KeyboardBuilder.Keyboard()

	offset, _ := strconv.Atoi(params.Offset)
	nextOffset := offset + 1
	prevOffset := offset - 1
	if nextOffset >= len(workout.Drills) {
		nextOffset = 0
	}
	if prevOffset < 0 {
		prevOffset = len(workout.Drills) - 1
	}

	buttonNextElement := keyboard.NewButton(s.constants.BUTTON_TEXT_NEXT_ELEMENT, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_INFO_WORKOUT, strconv.Itoa(nextOffset)).String())
	buttonPrevElement := keyboard.NewButton(s.constants.BUTTON_TEXT_PREV_ELEMENT, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_INFO_WORKOUT, strconv.Itoa(prevOffset)).String())
	keyboard.
		AppendAsLine(buttonPrevElement, buttonNextElement)

	header := fmt.Sprintf("🏃 %s, всего упражнений %d\n\n", workout.Name, len(workout.Drills))
	drill := workout.Drills[offset]

	drillText := fmt.Sprintf("💪 %s\n\n", drill.ElementName)
	header += drillText

	if params.Command == s.constants.COMMAND_WORKOUT_PLUS_SET {
		newSet := models.DrillSet{
			RepetitionCount: 0,
		}
		drill.Sets = append(drill.Sets, newSet)
		drill.CurrentlyObesrvableSet = len(drill.Sets) - 1
	}

	if params.Command == s.constants.COMMAND_WORKOUT_MINUS_SET {
		indexToRemove := drill.CurrentlyObesrvableSet
		drill.Sets = append(drill.Sets[:indexToRemove], drill.Sets[indexToRemove+1:]...)
		drill.CurrentlyObesrvableSet--
		if drill.CurrentlyObesrvableSet < 0 {
			drill.CurrentlyObesrvableSet = 0
		}
	}

	if params.Command == s.constants.COMMAND_WORKOUT_NEXT_SET {
		drill.CurrentlyObesrvableSet++
		if drill.CurrentlyObesrvableSet >= len(drill.Sets) {
			drill.CurrentlyObesrvableSet = 0
		}
	}

	if params.Command == s.constants.COMMAND_WORKOUT_PREV_SET {
		drill.CurrentlyObesrvableSet--
		if drill.CurrentlyObesrvableSet < 0 {
			drill.CurrentlyObesrvableSet = len(drill.Sets) - 1
		}
	}

	if params.Command == s.constants.COMMAND_WORKOUT_PLUS_REPS {
		drill.Sets[drill.CurrentlyObesrvableSet].RepetitionCount++
	}

	if params.Command == s.constants.COMMAND_WORKOUT_MINUS_REPS {
		drill.Sets[drill.CurrentlyObesrvableSet].RepetitionCount--
		if drill.Sets[drill.CurrentlyObesrvableSet].RepetitionCount < 0 {
			drill.Sets[drill.CurrentlyObesrvableSet].RepetitionCount = 0
		}
	}

	if len(drill.Sets) > 0 {
		currentSet := drill.Sets[drill.CurrentlyObesrvableSet]
		drillSet := fmt.Sprintf("Текущий подход %d\n", drill.CurrentlyObesrvableSet+1)

		reps := fmt.Sprintf("Повторения в подходе %d\n", currentSet.RepetitionCount)
		header += drillSet
		header += reps

		newSetButton := keyboard.NewButton("+ подход", s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_WORKOUT_PLUS_SET, strconv.Itoa(offset)).String())
		removeSetButton := keyboard.NewButton("- подход", s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_WORKOUT_MINUS_SET, strconv.Itoa(offset)).String())
		nextSetButton := keyboard.NewButton("подход >", s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_WORKOUT_NEXT_SET, strconv.Itoa(offset)).String())
		prevSetButton := keyboard.NewButton("< подход", s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_WORKOUT_PREV_SET, strconv.Itoa(offset)).String())
		plusRepsButton := keyboard.NewButton("+ повторения", s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_WORKOUT_PLUS_REPS, strconv.Itoa(offset)).String())
		minusRepsButton := keyboard.NewButton("- повторения", s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_WORKOUT_MINUS_REPS, strconv.Itoa(offset)).String())
		keyboard.
			AppendAsLine(minusRepsButton, plusRepsButton).
			AppendAsLine(prevSetButton, removeSetButton, newSetButton, nextSetButton)
	} else {
		keyboard.AppendAsLine(keyboard.NewButton("+ подход", s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_WORKOUT_PLUS_SET, strconv.Itoa(offset)).String()))
	}

	workout.Drills[offset] = drill
	err = s.repositories.Workout.Save(ctx, workout, nil)
	if err != nil {
		return err
	}

	header += fmt.Sprintf("Подходов в упражнении %d", len(drill.Sets))

	buttonBack := keyboard.NewButton(s.constants.BUTTON_TEXT_BACK, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_LIST_WORKOUT, "0").String())
	buttonEdit := keyboard.NewButton(s.constants.BUTTON_TEXT_EDIT, s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_EDIT_WORKOUT, strconv.Itoa(offset)).String())
	buttonDelete := keyboard.NewButton(s.constants.BUTTON_TEXT_DELETE, s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_DELETE_WORKOUT, strconv.Itoa(offset)).String())

	keyboard.
		AppendAsLine(buttonBack, buttonEdit, buttonDelete)

	_, err = s.clients.Telegram.Edit(ctx, header, update, telegram.WithReplyMurkup(keyboard.Murkup()))
	return err
}
