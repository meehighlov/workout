package element

import (
	"context"

	"github.com/meehighlov/workout/internal/repositories/models"
)

type Filter struct {
	ID     string
	UserID string
	Status string
	Limit  int
	Offset int
}

func (r *Repository) List(ctx context.Context, filter *Filter) ([]*models.Element, error) {
	var elements []*models.Element

	query := r.DB().WithContext(ctx)

	if filter != nil {
		if filter.UserID != "" {
			query = query.Where("user_id = ?", filter.UserID)
		}
		if filter.Status != "" {
			query = query.Where("status = ?", filter.Status)
		}
		if filter.Limit > 0 {
			query = query.Limit(filter.Limit)
		}
		if filter.Offset > 0 {
			query = query.Offset(filter.Offset)
		}
	}

	query = query.Order("created_at DESC")

	err := query.Find(&elements).Error
	if err != nil {
		return nil, err
	}
	return elements, nil
}
