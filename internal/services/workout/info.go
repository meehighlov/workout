package workout

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	inlinekeyboard "github.com/meehighlov/workout/internal/builders/inline_keyboard"
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

	offset, _ := strconv.Atoi(params.Offset)
	nextOffset := offset + 1
	prevOffset := offset - 1
	if nextOffset >= len(workout.Drills) {
		nextOffset = 0
	}
	if prevOffset < 0 {
		prevOffset = len(workout.Drills) - 1
	}

	header := fmt.Sprintf("🏃 %s\n\n", workout.Name)
	drill := workout.Drills[s.floorIndex(offset, len(workout.Drills))]

	drillText := fmt.Sprintf("💪 %d/%d %s\n\n", offset+1, len(workout.Drills), drill.ElementName)
	header += drillText

	if params.Command == s.constants.COMMAND_WORKOUT_PLUS_SET {
		newSet := models.DrillSet{
			Weight: "0.0",
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

	if strings.HasPrefix(params.Command, "tw_") {
		weight := s.ParseWeight(params.Command)
		prevWeight, _ := strconv.ParseFloat(
			drill.Sets[drill.CurrentlyObesrvableSet].Weight, 64)

		newWeight := prevWeight + weight
		if newWeight < 0 {
			newWeight = 0
		}

		drill.Sets[drill.CurrentlyObesrvableSet].Weight = strconv.FormatFloat(newWeight, 'f', -1, 64)
	}

	workout.Drills[s.floorIndex(offset, len(workout.Drills))] = drill
	err = s.repositories.Workout.Save(ctx, workout, nil)
	if err != nil {
		return err
	}

	header += fmt.Sprintf("Подходов в упражнении %d\n", len(drill.Sets))

	keyboard := s.builders.KeyboardBuilder.Keyboard()

	if params.Command != s.constants.COMMAND_INFO_WORKOUT {
		msg := fmt.Sprintf("💪 %s\n\n", drill.ElementName)

		newSetButton := keyboard.NewButton(s.constants.BUTTON_TEXT_WORKOUT_DRILL_SETS_INCREASE, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_WORKOUT_PLUS_SET, strconv.Itoa(offset)).String())
		removeSetButton := keyboard.NewButton(s.constants.BUTTON_TEXT_WORKOUT_DRILL_SETS_DECREASE, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_WORKOUT_MINUS_SET, strconv.Itoa(offset)).String())
		nextSetButton := keyboard.NewButton(s.constants.BUTTON_TEXT_NEXT, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_WORKOUT_NEXT_SET, strconv.Itoa(offset)).String())
		prevSetButton := keyboard.NewButton(s.constants.BUTTON_TEXT_PREV, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_WORKOUT_PREV_SET, strconv.Itoa(offset)).String())
		plusRepsButton := keyboard.NewButton(s.constants.BUTTON_TEXT_WORKOUT_DRILL_REPS_INCREASE, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_WORKOUT_PLUS_REPS, strconv.Itoa(offset)).String())
		minusRepsButton := keyboard.NewButton(s.constants.BUTTON_TEXT_WORKOUT_DRILL_REPS_DECREASE, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_WORKOUT_MINUS_REPS, strconv.Itoa(offset)).String())

		if len(drill.Sets) > 0 {
			currentSet := drill.Sets[drill.CurrentlyObesrvableSet]
			drillSet := fmt.Sprintf("Подход %d/%d\n", drill.CurrentlyObesrvableSet+1, len(drill.Sets))
			reps := fmt.Sprintf("Повторения %d\n", currentSet.RepetitionCount)
			weight := fmt.Sprintf("Доп вес %s(кг)\n", currentSet.Weight)
			msg += drillSet
			msg += reps
			msg += weight

			keyboard.
				AppendAsLine(prevSetButton, nextSetButton).
				AppendAsLine(removeSetButton, newSetButton).
				AppendAsLine(minusRepsButton, plusRepsButton).
				Append(s.WeightButtons(workout, offset))
		} else {
			keyboard.AppendAsLine(newSetButton)
		}

		keyboard.AppendAsLine(keyboard.NewButton(s.constants.BUTTON_TEXT_BACK, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_INFO_WORKOUT, strconv.Itoa(offset)).String()))

		_, err = s.clients.Telegram.Edit(ctx, msg, update, telegram.WithReplyMurkup(keyboard.Murkup()))
		return err
	} else {
		keyboard.AppendAsLine(keyboard.NewButton(s.constants.BUTTON_TEXT_EXEC, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_DRILL_EXEC, strconv.Itoa(offset)).String()))
	}

	buttonNextElement := keyboard.NewButton(s.constants.BUTTON_TEXT_NEXT, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_INFO_WORKOUT, strconv.Itoa(nextOffset)).String())
	buttonPrevElement := keyboard.NewButton(s.constants.BUTTON_TEXT_PREV, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_INFO_WORKOUT, strconv.Itoa(prevOffset)).String())
	keyboard.
		PutFirstAsLine(buttonPrevElement, buttonNextElement)

	buttonBack := keyboard.NewButton(s.constants.BUTTON_TEXT_BACK, s.builders.CallbackDataBuilder.Build(params.ID, s.constants.COMMAND_LIST_WORKOUT, "0").String())
	buttonEdit := keyboard.NewButton(s.constants.BUTTON_TEXT_EDIT, s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_EDIT_WORKOUT, strconv.Itoa(offset)).String())
	buttonDelete := keyboard.NewButton(s.constants.BUTTON_TEXT_DELETE, s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_DELETE_WORKOUT, strconv.Itoa(offset)).String())

	keyboard.
		AppendAsLine(buttonBack, buttonEdit, buttonDelete)

	_, err = s.clients.Telegram.Edit(ctx, header, update, telegram.WithReplyMurkup(keyboard.Murkup()))
	return err
}

func (s *Service) floorIndex(index int, length int) int {
	if index >= length {
		return length - 1
	}
	if index < 0 {
		return 0
	}
	if index >= 0 && index < length {
		return index
	}
	return 0
}

func (s *Service) WeightButtons(workout *models.Workout, offset int) *inlinekeyboard.Builder {
	keyboard := s.builders.KeyboardBuilder.Keyboard()

	kg025plus := keyboard.NewButton("+0.25", s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_WORKOUT_TUNE_WEIGHT_0_25_PLUS, strconv.Itoa(offset)).String())
	kg025minus := keyboard.NewButton("-0.25", s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_WORKOUT_TUNE_WEIGHT_0_25_MINUS, strconv.Itoa(offset)).String())

	kg05plus := keyboard.NewButton("+0.5", s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_WORKOUT_TUNE_WEIGHT_0_5_PLUS, strconv.Itoa(offset)).String())
	kg05minus := keyboard.NewButton("-0.5", s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_WORKOUT_TUNE_WEIGHT_0_5_MINUS, strconv.Itoa(offset)).String())

	kg1plus := keyboard.NewButton("+1", s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_WORKOUT_TUNE_WEIGHT_1_PLUS, strconv.Itoa(offset)).String())
	kg1minus := keyboard.NewButton("-1", s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_WORKOUT_TUNE_WEIGHT_1_MINUS, strconv.Itoa(offset)).String())

	kg5plus := keyboard.NewButton("+5", s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_WORKOUT_TUNE_WEIGHT_5_PLUS, strconv.Itoa(offset)).String())
	kg5minus := keyboard.NewButton("-5", s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_WORKOUT_TUNE_WEIGHT_5_MINUS, strconv.Itoa(offset)).String())

	kg10plus := keyboard.NewButton("+10", s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_WORKOUT_TUNE_WEIGHT_10_PLUS, strconv.Itoa(offset)).String())
	kg10minus := keyboard.NewButton("-10", s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_WORKOUT_TUNE_WEIGHT_10_MINUS, strconv.Itoa(offset)).String())

	kg20plus := keyboard.NewButton("+20", s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_WORKOUT_TUNE_WEIGHT_20_PLUS, strconv.Itoa(offset)).String())
	kg20minus := keyboard.NewButton("-20", s.builders.CallbackDataBuilder.Build(workout.ID.String(), s.constants.COMMAND_WORKOUT_TUNE_WEIGHT_20_MINUS, strconv.Itoa(offset)).String())

	keyboard.
		AppendAsLine(kg025minus, kg05minus, kg025plus, kg05plus).
		AppendAsLine(kg1minus, kg5minus, kg1plus, kg5plus).
		AppendAsLine(kg10minus, kg20minus, kg10plus, kg20plus)

	return keyboard
}

func (s *Service) ParseWeight(rawWeightData string) float64 {
	weightData := strings.Split(rawWeightData, "_")
	weight := weightData[1]
	action := string(weight[len(weight)-1])
	weightWithoutAction := strings.Replace(weight, action, "", 1)

	kg, _ := strconv.ParseFloat(weightWithoutAction, 64)

	if action == "p" {
		return kg
	}

	if action == "m" {
		return (-1) * kg
	}

	return 0
}
