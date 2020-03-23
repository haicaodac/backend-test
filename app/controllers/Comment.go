package controllers

import (
	"encoding/json"
	jobs "hanyny/app/jobs"
	library "hanyny/app/library"
	models "hanyny/app/models"
	system "hanyny/app/utils/system"
	v "hanyny/app/utils/view"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

// CommentGet ...
func CommentGet(w http.ResponseWriter, r *http.Request) {

}

// CommentGetOne ...
func CommentGetOne(w http.ResponseWriter, r *http.Request) {

}

// CommentAdd ...
func CommentAdd(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Dữ liệu không chính xác."))
		return
	}
	defer r.Body.Close()
	// db := models.OpenDB()
	user := r.Context().Value("user").(models.User)
	comment.UserID = user.ID

	if _, err := govalidator.ValidateStruct(comment); err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(comment)
		logger.Close()
		v.Respond(w, v.Message(false, err.Error()))
		return
	}

	db := models.OpenDB()

	// Kiểm tra comment đã đăng chưa?
	var checkContent models.Comment
	db.Where("content = ? AND post_id = ?", comment.Content, comment.PostID).First(&checkContent)
	if checkContent.ID != 0 {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(comment)
		logger.Close()
		v.Respond(w, v.Message(false, "Bình luận của bạn đã được chia sẻ rồi."))
		return
	}

	// Kiểm tra post có tồn tại không?
	var post models.Post
	db.Where("id = ?", comment.PostID).Last(&post)
	if post.ID == 0 {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(comment)
		logger.Close()
		v.Respond(w, v.Message(false, "Không tìm thấy bài viết bạn đang bình luận."))
		return
	}

	// Kiểm tra lịch sử
	var history models.History
	db.Where("type = ? AND user_id = ?", system.GetTypeHistory().AddComment, user.ID).First(&history)
	if history.ID == 0 { // Chưa có history
		history.Count = 1
		history.Type = system.GetTypeHistory().AddComment
		history.UserID = user.ID
		db.Create(&history)
	} else {
		remaining := time.Now().Sub(*history.CreatedAt)
		if remaining > (1 * time.Minute) { // Lớn hơn 1 phút rồi thì reset
			history.Count = 1
			*history.CreatedAt = time.Now()
			db.Save(&history)
		} else if history.Count < 10 { // Nhỏ hơn 10 lần câp nhật số lượng
			history.Count = history.Count + 1
			db.Save(&history)
		} else { // Lớn hơn 10 lần trong 1 phút tài khoản bị khoá
			user.Status = system.GetStatus().Block
			db.Save(&user)
			v.Respond(w, v.Message(false, "Tài khoản của bạn đã bị khoá do thêm bình luận quá nhiều. Vui lòng liên hệ với admin để khôi phục."))
			return
		}
	}

	db.Create(&comment)

	library.CacheCleanAll()

	db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, first_name, last_name, avatar, level")
	}).First(&comment)

	// Thông báo bằng email đến những nguời liên quan
	// NotifiNewComment
	var comments []models.Comment
	db.Where("post_id = ?", post.ID).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, first_name, last_name, avatar, level")
	}).Find(&comments)
	// Lấy user liên quan
	userIDs := []uint{}
	var users []models.User
	for _, comment := range comments {
		flag := true
		for _, user := range users {
			if comment.User != nil && user.ID == comment.User.ID {
				flag = false
				break
			}
		}
		if flag {
			users = append(users, *comment.User)
		}
	}
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
		go jobs.NotifiNewComment(post, comment, user)
	}
	// Thông báo cho admin
	var adminUsers []models.User
	db.Where("level in (?)", []string{system.GetLevel().Admin, system.GetLevel().Editor}).Not("id in (?)", userIDs).Find(&adminUsers)
	for _, user := range adminUsers {
		go jobs.NotifiNewComment(post, comment, user)
	}

	data := v.Message(true, "Tạo bình luận thành công.")
	data["comment"] = comment
	v.Respond(w, data)
}

// CommentUpdate ...
func CommentUpdate(w http.ResponseWriter, r *http.Request) {

}

// CommentDelete ...
func CommentDelete(w http.ResponseWriter, r *http.Request) {

}
