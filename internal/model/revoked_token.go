package model

import "time"

type RevokedToken struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `gorm:"size:512;not null;uniqueIndex"`
	ExpiresAt time.Time
	CreatedAt time.Time
}
