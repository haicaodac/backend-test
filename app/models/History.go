package models

// History in database histories
type History struct {
	Model
	Type   string `json:"type"`
	Data   string `json:"data"`
	Count  uint   `json:"count" gorm:"not null;default:1"`
	IP     uint   `json:"ip"`
	UserID uint   `json:"user_id"`
}
