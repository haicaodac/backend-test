package models

// Attachment in database attachments
type Attachment struct {
	Model
	Name          string `json:"name,omitempty" gorm:"type:varchar(4000);not null" valid:"required~Tên tài liệu không thể trống.,runelength(1|4000)~Tên tài liệu phải có độ dài từ 1 đến 4000 ký tự."`
	Type          string `json:"type,omitempty" gorm:"type:varchar(200);not null"`
	Size          int64  `json:"size,omitempty" gorm:"not null"`
	URL           string `json:"url,omitempty" gorm:"type:varchar(4000);not null" valid:"required~Đường dẫn tài liệu không thể trống.,runelength(1|4000)~Đường dẫn tài liệu phải có độ dài từ 1 đến 4000 ký tự."`
	View          int    `json:"view,omitempty" gorm:"not null;default:0"`
	CountDownload int    `json:"count_download,omitempty" gorm:"not null;default:0"`
	UserID        uint   `json:"user_id,omitempty" gorm:"not null"`
	Status        string `json:"status,omitempty" gorm:"type:varchar(20);not null;default:'active'"`
	User          *User  `json:"user,omitempty"`
	Post          *Post  `json:"post,omitempty"`
}
