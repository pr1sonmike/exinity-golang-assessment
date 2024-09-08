package models

import "time"

type Transaction struct {
	ID        uint `gorm:"primaryKey"`
	Gateway   string
	Amount    float64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
