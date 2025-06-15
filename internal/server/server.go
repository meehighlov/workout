package server

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/meehighlov/workout/internal/builders"
	"github.com/meehighlov/workout/internal/clients"
	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/config"
	"github.com/meehighlov/workout/internal/constants"
	"github.com/meehighlov/workout/internal/services"
)

type Server struct {
	logger        *slog.Logger
	handleTimeout time.Duration
	constants     *constants.Constants
	services      *services.Services
	clients       *clients.Clients
	builders      *builders.Builders
	allowedUsers  []string
	webServer     *http.Server
	wgWebServer   sync.WaitGroup
	shutdownChan  chan struct{}
	cfg           *config.Config

	wgWorkerPool  sync.WaitGroup
	workerCount  int
	updatesChan  *telegram.UpdatesChannel
	workerCtx    context.Context
	workerCancel context.CancelFunc
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	services *services.Services,
	clients *clients.Clients,
	constants *constants.Constants,
	builders *builders.Builders,
) *Server {
	return &Server{
		logger:        logger,
		services:      services,
		clients:       clients,
		constants:     constants,
		builders:      builders,
		handleTimeout: time.Duration(cfg.TelegramHandlerTimeoutSec) * time.Second,
		allowedUsers:  cfg.AllowedUsers(),
		cfg:           cfg,
		shutdownChan:  make(chan struct{}),
		wgWebServer:   sync.WaitGroup{},
		wgWorkerPool:  sync.WaitGroup{},
		workerCount:   cfg.WorkerCount,
		workerCtx:     context.Background(),
		workerCancel:  func() {},
	}
}
