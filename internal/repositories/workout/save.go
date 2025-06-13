package workout

import (
	"context"

	"github.com/meehighlov/workout/internal/repositories/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) Save(ctx context.Context, workout *models.Workout, tx *gorm.DB) error {
	var db *gorm.DB
	if tx != nil {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	db = db.Session(&gorm.Session{
		SkipHooks: true,
	})

	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "drills", "status", "updated_at"}),
	}).Create(workout).Error
	if err != nil {
		r.logger.Error("SaveWorkout error", "reason", err)
		return err
	}

	r.logger.Debug("SaveWorkout done", "id", workout.ID)
	return nil
}
