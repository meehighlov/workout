package workout

import (
	"context"
	"errors"

	"github.com/meehighlov/workout/internal/repositories/models"
	"gorm.io/gorm"
)

func (r *Repository) Get(ctx context.Context, filter *Filter, tx *gorm.DB) (*models.Workout, error) {
	db := r.db
	if tx != nil {
		db = tx
	}

	workout := &models.Workout{}

	if filter == nil {
		return nil, errors.New("filter is nil")
	}

	q := db.WithContext(ctx).Order("created_at DESC")

	if filter.ID != "" {
		q = q.Where("id = ?", filter.ID)
	}

	if filter.UserID != "" {
		q = q.Where("user_id = ?", filter.UserID)
	}

	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}

	err := q.First(workout).Error
	if err != nil {
		return nil, err
	}

	return workout, nil
}
