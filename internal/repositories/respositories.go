package repositories

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/meehighlov/workout/internal/config"
	"github.com/meehighlov/workout/internal/repositories/element"
	"github.com/meehighlov/workout/internal/repositories/user"
	"github.com/meehighlov/workout/internal/repositories/workout"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Repositories struct {
	User    *user.Repository
	Element *element.Repository
	Workout *workout.Repository
}

func New(cfg *config.Config, logger *slog.Logger) *Repositories {
	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{
		Logger: WrapAppLogger(logger),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		NowFunc: func() time.Time {
			return time.Now()
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := RunMigrations(context.Background(), cfg, logger, db); err != nil {
		log.Fatal("Migration error:", err)
	}

	return &Repositories{
		User:    user.New(cfg, db, logger),
		Element: element.New(cfg, db, logger),
		Workout: workout.New(cfg, db, logger),
	}
}
