package models

import (
	"time"
)

type Log struct {
	ID          int        `gorm:"primaryKey" json:"id"`
	UserID      int        `json:"user_id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Url         string     `json:"url"`
	ViewedAt    *time.Time `json:"viewed_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

func NewLog(log *Log) error {
	return GetDB().Create(log).Error
}

func GetAllLogFromUserID(userId int) []Log {
	var logs []Log
	GetDB().Model(&Log{}).Where("user_id = ?", userId).Find(&logs)
	return logs
}
