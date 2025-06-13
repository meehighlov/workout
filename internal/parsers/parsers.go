package parsers

import (
	"log/slog"

	"github.com/meehighlov/workout/internal/config"
)

type Parsers struct {
}

func New(cfg *config.Config, logger *slog.Logger) *Parsers {
	return &Parsers{}
}
