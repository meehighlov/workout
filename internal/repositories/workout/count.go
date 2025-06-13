package workout

import (
	"context"

	"github.com/meehighlov/workout/internal/repositories/models"
)

func (r *Repository) Count(ctx context.Context, filter *Filter) (int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Workout{})

	if filter != nil {
		if filter.UserID != "" {
			query = query.Where("user_id = ?", filter.UserID)
		}
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
