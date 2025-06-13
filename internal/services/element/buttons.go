package element

import (
	"context"
	"strconv"

	"github.com/meehighlov/workout/internal/builders/inline_keyboard"
	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/repositories/element"
	"github.com/meehighlov/workout/internal/repositories/user"
)

func (s *Service) BuildElementsKeyboard(
	ctx context.Context,
	update *telegram.Update,
	elementCommand string,
	controlCommand string,
) (*inlinekeyboard.Builder, error) {
	user, err := s.repositories.User.Get(ctx, &user.Filter{TgChatID: update.GetChatIdStr()}, nil)
	if err != nil {
		return nil, err
	}

	filter := &element.Filter{
		UserID: user.ID.String(),
	}

	count, err := s.repositories.Element.Count(ctx, &element.Filter{UserID: user.ID.String()})
	if err != nil {
		return nil, err
	}

	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	offset, err := strconv.Atoi(params.Offset)
	if err != nil {
		return nil, err
	}

	filter.Limit = s.pagination.Limit
	filter.Offset = offset

	elements, err := s.repositories.Element.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	keyboard := s.builders.KeyboardBuilder.Keyboard()

	for _, element := range elements {
		keyboard.AppendAsLine(keyboard.NewButton(element.Name, s.builders.CallbackDataBuilder.Build(element.ID.String(), elementCommand, params.Offset).String()))
	}

	keyboard.
		OptimizeView().
		Append(
			s.pagination.BuildControls(
				int(count),
				controlCommand,
				update.GetChatIdStr(),
				offset,
			),
		)

	return keyboard, nil
}
