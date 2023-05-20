package models

import (
	"database/sql"
	"time"
)

type Seat struct {
	ID          string `gorm:"primaryKey"`
	SeatClassID int
	SeatClass   SeatClass
	BookingID   sql.NullInt64
	UpdatedAt   time.Time
}
