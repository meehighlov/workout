package element

import (
	"log/slog"

	"github.com/meehighlov/workout/internal/builders"
	"github.com/meehighlov/workout/internal/clients"
	"github.com/meehighlov/workout/internal/config"
	"github.com/meehighlov/workout/internal/constants"
	"github.com/meehighlov/workout/internal/pagination"
	"github.com/meehighlov/workout/internal/parsers"
	"github.com/meehighlov/workout/internal/repositories"
	"github.com/meehighlov/workout/internal/validators"
)

type Service struct {
	logger       *slog.Logger
	repositories *repositories.Repositories
	clients      *clients.Clients
	builders     *builders.Builders
	validators   *validators.Validators
	constants    *constants.Constants
	parsers      *parsers.Parsers
	pagination   *pagination.Pagination
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
) *Service {
	return &Service{
		logger:       logger,
		repositories: repositories,
		clients:      clients,
		builders:     builders,
		validators:   validators,
		constants:    constants,
		parsers:      parsers,
		pagination:   pagination,
	}
}
