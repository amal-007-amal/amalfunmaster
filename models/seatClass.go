package models

import "database/sql"

type SeatClass struct {
	ID          uint `gorm:"primarykey"`
	Title       string
	MaxPrice    sql.NullFloat64
	NormalPrice sql.NullFloat64
	MinPrice    sql.NullFloat64
	SeatCount   int
	BookedCount int
	Seats       []Seat
}
