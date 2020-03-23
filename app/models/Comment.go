package models

// Comment in database comments
type Comment struct {
	Model
	ParentID  uint       `json:"parent_id,omitempty" valid:"int~ParentID phải là ID của một bình luận nào đó."`
	Content   string     `json:"content,omitempty" gorm:"not null" valid:"required~Nội dung bình luận không thể trống."`
	CommentID uint       `json:"comment_id,omitempty" valid:"int~CommentID phải là ID của một comment nào đó."`
	PostID    uint       `json:"post_id,omitempty" valid:"int~PostID phải là ID của một bài viết nào đó."`
	UserID    uint       `json:"user_id,omitempty" gorm:"not null"`
	User      *User      `json:"user,omitempty"`
	Comments  []*Comment `json:"comments,omitempty" gorm:"foreignkey:ParentID"`
}
