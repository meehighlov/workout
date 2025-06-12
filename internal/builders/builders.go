package builders

import (
	"log/slog"

	"github.com/meehighlov/workout/internal/builders/callback_data"
	"github.com/meehighlov/workout/internal/builders/inline_keyboard"
	"github.com/meehighlov/workout/internal/builders/short_id"
	"github.com/meehighlov/workout/internal/config"
)

type Builders struct {
	ShortIdBuilder      *short_id.Builder
	CallbackDataBuilder *callbackdata.Builder
	KeyboardBuilder     *inlinekeyboard.Builder
}

func New(cfg *config.Config, logger *slog.Logger) *Builders {
	return &Builders{
		ShortIdBuilder:      short_id.New(6),
		CallbackDataBuilder: callbackdata.New(),
		KeyboardBuilder:     inlinekeyboard.New(),
	}
}
