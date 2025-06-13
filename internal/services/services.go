package services

import (
	"log/slog"

	"github.com/meehighlov/workout/internal/builders"
	"github.com/meehighlov/workout/internal/clients"
	"github.com/meehighlov/workout/internal/config"
	"github.com/meehighlov/workout/internal/constants"
	"github.com/meehighlov/workout/internal/pagination"
	"github.com/meehighlov/workout/internal/parsers"
	"github.com/meehighlov/workout/internal/repositories"
	"github.com/meehighlov/workout/internal/services/element"
	"github.com/meehighlov/workout/internal/services/user"
	"github.com/meehighlov/workout/internal/services/workout"
	"github.com/meehighlov/workout/internal/validators"
)

type Services struct {
	User    *user.Service
	Element *element.Service
	Workout *workout.Service
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	repositories *repositories.Repositories,
	clients *clients.Clients,
	builders *builders.Builders,
	validators *validators.Validators,
	constants *constants.Constants,
	parsers *parsers.Parsers,
	pagination *pagination.Pagination,
) *Services {
	return &Services{
		User:    user.New(cfg, logger, repositories, clients, builders, validators, constants, parsers),
		Element: element.New(cfg, logger, repositories, clients, builders, validators, constants, parsers, pagination),
		Workout: workout.New(cfg, logger, repositories, clients, builders, validators, constants, parsers, pagination),
	}
}
