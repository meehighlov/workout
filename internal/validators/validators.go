package validators

import (
	"log/slog"

	"github.com/meehighlov/workout/internal/config"
)

type Validators struct {
}

func New(cfg *config.Config, logger *slog.Logger) *Validators {
	return &Validators{}
}
