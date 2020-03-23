package models

// Media in database medias
type Media struct {
	Model
	Name   string `json:"name,omitempty" gorm:"type:varchar(4000);not null" valid:"required~Tên tài liệu không thể trống.,runelength(1|4000)~Tên tài liệu phải có độ dài từ 1 đến 4000 ký tự."`
	Type   string `json:"type,omitempty" gorm:"type:varchar(100);not null"`
	Width  uint   `json:"width,omitempty"`
	Height uint   `json:"height,omitempty"`
	URL    string `json:"url,omitempty" gorm:"type:varchar(4000);not null" valid:"required~Đường dẫn tài liệu không thể trống.,runelength(1|4000)~Đường dẫn tài liệu phải có độ dài từ 1 đến 4000 ký tự."`
	UserID uint   `json:"user_id,omitempty" gorm:"not null"`
}
