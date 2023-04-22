package entity

import "time"

type Tag struct {
	ID        int64     `gorm:"primarykey"`
	UserID    int64     `gorm:"index;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Title     string    `gorm:"type:text;not null"`
}
