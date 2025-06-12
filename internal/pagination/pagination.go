package pagination

import (
	"github.com/meehighlov/workout/internal/builders"
)

type Pagination struct {
	builders *builders.Builders
	BaseOffset int
	Limit int
}

func New(builders *builders.Builders) *Pagination {
	return &Pagination{
		builders: builders,
		BaseOffset: 6,
		Limit: 6,
	}
}
