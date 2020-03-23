package models

// Review in database reviews
type Review struct {
	Model
	Content string `json:"content,omitempty" gorm:"type:varchar(4000);"`
	Point   int    `json:"point,omitempty" gorm:"not null" valid:"required~Điểm của đánh giá không thể trống.,int~Điểm phải là một giá trị số."`
	PostID  uint   `json:"post_id,omitempty" valid:"int~PostID phải là ID của một bài viết nào đó."`
	UserID  uint   `json:"user_id,omitempty" gorm:"not null"`
}
