package workout

import (
	"context"
	"strconv"

	"github.com/meehighlov/workout/internal/builders/inline_keyboard"
	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/repositories/user"
	"github.com/meehighlov/workout/internal/repositories/workout"
)

func (s *Service) BuildWorkoutsKeyboard(
	ctx context.Context,
	update *telegram.Update,
	elementCommand string,
	controlCommand string,
) (*inlinekeyboard.Builder, error) {
	user, err := s.repositories.User.Get(ctx, &user.Filter{TgChatID: update.GetChatIdStr()}, nil)
	if err != nil {
		return nil, err
	}

	filter := &workout.Filter{
		UserID: user.ID.String(),
	}

	count, err := s.repositories.Workout.Count(ctx, &workout.Filter{UserID: user.ID.String()})
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

	workouts, err := s.repositories.Workout.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	keyboard := s.builders.KeyboardBuilder.Keyboard()

	for _, workout := range workouts {
		keyboard.AppendAsLine(keyboard.NewButton(workout.Name, s.builders.CallbackDataBuilder.Build(workout.ID.String(), elementCommand, params.Offset).String()))
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
