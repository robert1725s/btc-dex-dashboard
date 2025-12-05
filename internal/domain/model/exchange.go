package model

import "time"

// Exchange は取引所マスタ
type Exchange struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Key         string    `gorm:"uniqueIndex;size:50;not null" json:"key"`
	DisplayName string    `gorm:"size:100;not null" json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
