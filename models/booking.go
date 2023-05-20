package models

import "gorm.io/gorm"

type Booking struct {
	gorm.Model
	Name  string
	Email string
	Phone string
	Seats []Seat
}
