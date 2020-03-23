package models

// Botchat in database botchats
type Botchat struct {
	Model
	Message  string `json:"message,omitempty" gorm:"type:text"`
	Tag      string `json:"tag,omitempty" gorm:"type:text"`
	Response string `json:"response,omitempty" gorm:"type:text"`
}
