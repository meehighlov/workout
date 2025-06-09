package clients

import (
	"log/slog"

	"github.com/meehighlov/workout/internal/clients/cache"
	"github.com/meehighlov/workout/internal/clients/telegram"
	"github.com/meehighlov/workout/internal/config"
)

type Clients struct {
	Telegram *telegram.Client
	Cache    *cache.Client
}

func New(cfg *config.Config, logger *slog.Logger) *Clients {
	return &Clients{
		Telegram: telegram.New(cfg, logger),
		Cache:    cache.New(cfg, logger),
	}
}

func (c *Clients) Close() error {
	c.Cache.Close()
	return nil
}
