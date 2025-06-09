package app

import (
	"github.com/meehighlov/workout/internal/builders"
	"github.com/meehighlov/workout/internal/clients"
	"github.com/meehighlov/workout/internal/config"
	"github.com/meehighlov/workout/internal/constants"
	"github.com/meehighlov/workout/internal/parsers"
	"github.com/meehighlov/workout/internal/repositories"
	"github.com/meehighlov/workout/internal/server"
	"github.com/meehighlov/workout/internal/services"
	"github.com/meehighlov/workout/internal/validators"
)

func Run() {
	cfg := config.MustLoad()
	logger := MustSetupLogging(cfg)

	repositories := repositories.New(cfg, logger)
	clients := clients.New(cfg, logger)
	builders := builders.New(cfg, logger)
	validators := validators.New(cfg, logger)
	constants := constants.New(cfg)
	parsers := parsers.New(cfg, logger)
	services := services.New(cfg, logger, repositories, clients, builders, validators, constants, parsers)

	server := server.New(cfg, logger, services, clients, constants, builders)
	server.Serve()
}
