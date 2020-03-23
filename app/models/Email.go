package models

// Email in database emails
type Email struct {
	Model
	Email          string `json:"email,omitempty" gorm:"type:varchar(200);not null"`
	IP             string `json:"-" gorm:"type:varchar(200);"`
	CountSendError int    `json:"count_send_error"`
}
