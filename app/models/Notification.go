package models

// Notification in database notifications
type Notification struct {
	Model
	Type      string `json:"type,omitempty" gorm:"type:varchar(100);not null" valid:"required~Kiểu của thông báo không thể trống.,runelength(1|100)~Kiểu của thôbg báo phải có độ dài từ 1 đến 100 ký tự."`
	PostID    uint   `json:"post_id,omitempty" valid:"int~PostID phải là ID của một bài viết nào đó."`
	ReviewID  uint   `json:"review_id,omitempty" valid:"int~ReviewID phải là ID của một đánh giá nào đó."`
	CommentID uint   `json:"comment_id,omitempty" valid:"int~CommentID phải là ID của một comment nào đó."`
	UserID    uint   `json:"user_id,omitempty" gorm:"not null"`
}
