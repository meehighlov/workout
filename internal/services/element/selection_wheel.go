package element

import (
	"context"
	"strings"

	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/repositories/element"
)

func (s *Service) ElementsSelectionWheel(ctx context.Context, update *telegram.Update) error {
	msg := s.constants.ELEMENTS_SELECTION_WHEEL_MESSAGE

	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	keyboard, err := s.BuildElementsKeyboard(
		ctx,
		update,
		s.constants.COMMAND_ADD_ELEMENT_TO_WORKOUT,
		s.constants.COMMAND_ADD_ELEMENT_TO_WORKOUT_CONTROL,
	)
	if err != nil {
		return err
	}

	keyboard.
		AppendAsLine(
			keyboard.NewButton(s.constants.BUTTON_TEXT_SAVE, s.builders.CallbackDataBuilder.Build(update.GetChatIdStr(), s.constants.COMMAND_SAVE_WORKOUT, "0").String()),
		).
		PutFirstAsLine(
			keyboard.NewButton(s.constants.BUTTON_TEXT_ERASE_ELEMENT, s.builders.CallbackDataBuilder.Build(update.GetChatIdStr(), s.constants.COMMAND_ADD_ELEMENT_TO_WORKOUT_RM_EL, params.Offset).String()),
		)

	if params.Command == s.constants.COMMAND_ADD_ELEMENT_TO_WORKOUT_RM_EL {
		s.clients.Cache.PopWorkoutElement(update.GetChatIdStr())

		drills := s.clients.Cache.GetWorkoutElements(update.GetChatIdStr())
		msg += strings.Join(drills, "\n")

		_, err = s.clients.Telegram.Edit(ctx, msg, update, telegram.WithReplyMurkup(keyboard.Murkup()))

		return err
	}

	if params.Command == s.constants.COMMAND_ADD_ELEMENT_TO_WORKOUT_CONTROL {
		msg := update.GetMessage().Text
		_, err = s.clients.Telegram.Edit(ctx, msg, update, telegram.WithReplyMurkup(keyboard.Murkup()))
		return err
	}

	if params.Command == s.constants.COMMAND_ADD_ELEMENT_TO_WORKOUT {
		element, _ := s.repositories.Element.Get(ctx, &element.Filter{ID: params.ID}, nil)
		s.clients.Cache.AppendWorkoutElement(update.GetChatIdStr(), element.Name)

		drills := s.clients.Cache.GetWorkoutElements(update.GetChatIdStr())
		msg += strings.Join(drills, "\n")

		_, err = s.clients.Telegram.Edit(ctx, msg, update, telegram.WithReplyMurkup(keyboard.Murkup()))

		return nil
	}

	drills := s.clients.Cache.GetWorkoutElements(update.GetChatIdStr())
	msg += strings.Join(drills, "\n")

	_, err = s.clients.Telegram.Reply(ctx, msg, update, telegram.WithReplyMurkup(keyboard.Murkup()))
	return err
}
