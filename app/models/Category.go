package models

// Category in database categorys
type Category struct {
	Model
	Name        string  `json:"name,omitempty" gorm:"type:varchar(200);not null"`
	Title       string  `json:"title,omitempty" gorm:"type:varchar(200);not null"`
	Slug        string  `json:"slug,omitempty" gorm:"type:varchar(200);not null"`
	Description string  `json:"description,omitempty" gorm:"type:varchar(400);not null"`
	Thumbnail   string  `json:"thumbnail,omitempty" gorm:"type:varchar(1000);"`
	Icon        string  `json:"icon,omitempty" gorm:"type:varchar(1000)"`
	Position    int     `json:"position,omitempty"`
	ParentID    uint    `json:"parent_id,omitempty"`
	Level       int     `json:"level,omitempty" gorm:"not null;default:1"`
	Primary     bool    `json:"primary,omitempty"`
	UserID      uint    `json:"user_id,omitempty" gorm:"not null"`
	Status      string  `json:"status,omitempty" gorm:"type:varchar(20);not null;default:'active'"`
	Posts       []*Post `json:"posts,omitempty" gorm:"many2many:category_posts;"`
}

// ValidateCategory ...
type ValidateCategory struct {
	Name        string `json:"name" valid:"required~Tên chuyên mục không thể trống.,runelength(1|200)~Tên chuyển mục phải có độ dài từ 1 đến 200 ký tự."`
	Title       string `json:"title" valid:"required~Tiêu đề không thể trống.,runelength(10|200)~Tiêu đề phải có độ dài từ 10 đến 200 ký tự."`
	Slug        string `json:"slug" valid:"required~Đường dẫn tĩnh không thể trống.,runelength(5|200)~Dường dẫn tĩnh phải có độ dài từ 5 đến 200 ký tự."`
	Description string `json:"description" valid:"required~Mô tả không thể trống.,runelength(100|400)~Mô tả phải có độ dài từ 100 đến 400 ký tự."`
	Thumbnail   string `json:"thumbnail" valid:"required~Ảnh đại diện không thể trống.,runelength(1|1000)~Ảnh đại diện phải có độ dài từ 1 đến 1000 ký tự.,url~Ảnh đại diện phải là một đường dẫn(URL)."`
	Icon        string `json:"icon" valid:"runelength(0|1000)~Icon phải có độ dài từ 0 đến 1000 ký tự."`
	Position    int    `json:"position" valid:"int~Vị trí phải là số."`
	ParentID    uint   `json:"parent_id" valid:"int~ParentID phải là id của một chuyên mục nào đó."`
	Level       int    `json:"level" valid:"int~Cấp độ phải là số."`
	Status      string `json:"status" valid:"runelength(1|20)~Trạng thái phải có độ dài từ 1 đến 20 ký tự."`
}
