package workout

import (
	"context"

	"github.com/meehighlov/workout/internal/repositories/models"
)

type Filter struct {
	ID     string
	UserID string
	Name   string
	Status string
	Limit  int
	Offset int
}

func (r *Repository) List(ctx context.Context, filter *Filter) ([]*models.Workout, error) {
	var workouts []*models.Workout

	query := r.DB().WithContext(ctx)

	if filter != nil {
		if filter.UserID != "" {
			query = query.Where("user_id = ?", filter.UserID)
		}
		if filter.Name != "" {
			query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
		}
		if filter.Limit > 0 {
			query = query.Limit(filter.Limit)
		}
		if filter.Offset > 0 {
			query = query.Offset(filter.Offset)
		}
		if filter.Status != "" {
			query = query.Where("status = ?", filter.Status)
		}
	}

	query = query.Order("created_at DESC")

	err := query.Find(&workouts).Error
	if err != nil {
		return nil, err
	}
	return workouts, nil
}
