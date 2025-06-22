package element

import (
	"context"
	"strings"

	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/repositories/element"
)

func (s *Service) ElementsSelectionWheelOnWrokoutEdit(ctx context.Context, update *telegram.Update) error {
	msg := s.constants.ELEMENTS_SELECTION_WHEEL_MESSAGE

	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	keyboard, err := s.BuildElementsKeyboard(
		ctx,
		update,
		s.constants.COMMAND_EDIT_WORKOUT_DRILLS_ADD_EL,
		s.constants.COMMAND_ADD_ELEMENT_TO_EDIT_WORKOUT_CONTROL,
	)
	if err != nil {
		return err
	}

	keyboard.
		AppendAsLine(
			keyboard.NewButton(s.constants.BUTTON_TEXT_CANCEL, s.builders.CallbackDataBuilder.Build(update.GetChatIdStr(), s.constants.COMMAND_CANCEL, "0").String()),
			keyboard.NewButton(s.constants.BUTTON_TEXT_SAVE, s.builders.CallbackDataBuilder.Build(update.GetChatIdStr(), s.constants.COMMAND_SAVE_WORKOUT, "0").String()),
		).
		PutFirstAsLine(
			keyboard.NewButton(s.constants.BUTTON_TEXT_ERASE_ELEMENT, s.builders.CallbackDataBuilder.Build(update.GetChatIdStr(), s.constants.COMMAND_EDIT_WORKOUT_DRILLS_RM_EL, params.Offset).String()),
		)

	if params.Command == s.constants.COMMAND_EDIT_WORKOUT_DRILLS_RM_EL {
		s.clients.Cache.PopWorkoutElement(ctx, update.GetChatIdStr())

		drills := s.clients.Cache.GetWorkoutElements(ctx, update.GetChatIdStr())
		msg += strings.Join(drills, "\n")

		_, err = s.clients.Telegram.Edit(ctx, msg, update, telegram.WithReplyMurkup(keyboard.Murkup()))

		return err
	}

	if params.Command == s.constants.COMMAND_ADD_ELEMENT_TO_EDIT_WORKOUT_CONTROL {
		msg := update.GetMessage().Text
		_, err = s.clients.Telegram.Edit(ctx, msg, update, telegram.WithReplyMurkup(keyboard.Murkup()))
		return err
	}

	if params.Command == s.constants.COMMAND_ADD_ELEMENT_TO_EDIT_WORKOUT {
		element, err := s.repositories.Element.Get(ctx, &element.Filter{ID: params.ID}, nil)
		if err != nil {
			return err
		}

		pickedDrill := element.Name

		drills := s.clients.Cache.GetWorkoutElements(ctx, update.GetChatIdStr())

		drills = append(drills, pickedDrill)

		seen := make(map[string]bool)
		repeated := false
		for _, drill := range drills {
			if !seen[drill] {
				seen[drill] = true
			} else {
				repeated = true
			}
		}

		if repeated {
			msg += strings.Join(drills[:len(drills)-1], "\n")
		} else {
			s.clients.Cache.AppendWorkoutElement(ctx, update.GetChatIdStr(), pickedDrill)
			msg += strings.Join(drills, "\n")
		}

		_, err = s.clients.Telegram.Edit(ctx, msg, update, telegram.WithReplyMurkup(keyboard.Murkup()))

		return nil
	}

	drills := s.clients.Cache.GetWorkoutElements(ctx, update.GetChatIdStr())
	msg += strings.Join(drills, "\n")

	_, err = s.clients.Telegram.Reply(ctx, msg, update, telegram.WithReplyMurkup(keyboard.Murkup()))
	return err
}
