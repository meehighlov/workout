package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	WORKOUT_STATUS_ACTIVE = "active"
	WORKOUT_STATUS_COMPLETED = "completed"
)

type Workout struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	Name      string    `gorm:"column:name;type:varchar(255);not null"`
	Drills    Drills    `gorm:"column:drills;type:json"`
	Status    string    `gorm:"column:status;type:varchar(255);not null;default:active"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp with time zone;not null;default:now()"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp with time zone;not null;default:now()"`
}

func (*Workout) TableName() string {
	return "workout"
}

func (w *Workout) GetID() string {
	return w.ID.String()
}

func (w *Workout) GetDrills() []string {
	drills := []string{}
	for _, drill := range w.Drills {
		drills = append(drills, drill.ElementName)
	}
	return drills
}

type DrillSet struct {
	RepetitionCount int    `json:"repetition_count"`
	Weight          string `json:"weight"`
}

type Drill struct {
	ElementName            string     `json:"element_name"`
	CurrentlyObesrvableSet int        `json:"currently_obesrvable_set"`
	Sets                   []DrillSet `json:"sets"`
}

type Drills []Drill

func (d Drills) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	return json.Marshal(d)
}

func (d *Drills) Scan(value interface{}) error {
	if value == nil {
		*d = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot scan non-bytes into Drills")
	}

	return json.Unmarshal(bytes, d)
}
