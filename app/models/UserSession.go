package models

// UserSession in database user_sessions
type UserSession struct {
	Model
	Type   string `json:"type" gorm:"type:varchar(100);not null"`
	UserID uint   `json:"user_id" gorm:"not null"`
	Data   string `json:"data" gorm:"type:varchar(1000);not null"`
}
