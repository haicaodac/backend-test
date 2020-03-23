package models

// Post in database posts
type Post struct {
	Model
	Title        string      `json:"title,omitempty" gorm:"type:varchar(150);not null" valid:"required~Tiêu đề không thể trống.,runelength(30|150)~Tiêu đề phải có độ dài từ 30 đến 150 ký tự."`
	Slug         string      `json:"slug,omitempty" gorm:"type:varchar(150)" valid:"runelength(10|150)~Dường dẫn tĩnh phải có độ dài từ 10 đến 150 ký tự."`
	Description  string      `json:"description,omitempty" gorm:"type:text"`
	Content      string      `json:"content,omitempty" gorm:"type:longtext;not null" valid:"required~Nội dung không thể trống."`
	Thumbnail    string      `json:"thumbnail,omitempty" gorm:"type:varchar(1000);" valid:"runelength(1|1000)~Ảnh đại diện phải có độ dài từ 1 đến 1000 ký tự.,url~Ảnh đại diện phải là một đường dẫn(URL)."`
	View         int         `json:"view,omitempty" gorm:"not null;default:0"`
	AttachmentID uint        `json:"attachment_id,omitempty" valid:"int~Tài liệu đính kèm phải là ID của một tài liệu nào đó"`
	UserID       uint        `json:"user_id,omitempty" gorm:"not null"`
	Status       string      `json:"status,omitempty" gorm:"type:varchar(20);not null;default:'active'"`
	User         *User       `json:"user,omitempty"`
	Categories   []*Category `json:"categories,omitempty" gorm:"many2many:category_posts;"`
	Attachment   *Attachment `json:"attachment,omitempty"`
	Comments     []*Comment  `json:"comments,omitempty"`
}

// ValidatePost ...
type ValidatePost struct {
	Title        string      `json:"title" valid:"required~Tiêu đề không thể trống.,runelength(30|150)~Tiêu đề phải có độ dài từ 30 đến 150 ký tự."`
	Slug         string      `json:"slug" valid:"required~Đường dẫn tĩnh không thể trống.,runelength(10|150)~Dường dẫn tĩnh phải có độ dài từ 10 đến 150 ký tự."`
	Description  string      `json:"description" valid:"required~Mô tả không thể trống.,runelength(100|400)~Mô tả phải có độ dài từ 100 đến 400 ký tự."`
	Content      string      `json:"content" valid:"required~Nội dung không thể trống."`
	Thumbnail    string      `json:"thumbnail" valid:"required~Ảnh đại diện không thể trống.,runelength(1|1000)~Ảnh đại diện phải có độ dài từ 1 đến 1000 ký tự.,url~Ảnh đại diện phải là một đường dẫn(URL)."`
	AttachmentID uint        `json:"attachment_id" valid:"int~Tài liệu đính kèm phải là ID của một tài liệu nào đó"`
	Status       string      `json:"status" valid:"required~Trạng thái không thể trống."`
	Categories   []*Category `json:"categories" valid:"required~Danh sách chuyên mục không thể trống."`
}
