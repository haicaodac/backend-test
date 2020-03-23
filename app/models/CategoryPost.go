package models

// CategoryPost in database category_posts
type CategoryPost struct {
	CategoryID uint `json:"category_id" gorm:"primary_key"`
	PostID     uint `json:"post_id" gorm:"primary_key"`
	Primary    bool `json:"primary,omitempty"`
}
