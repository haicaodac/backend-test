package controllers

import (
	"bytes"
	"encoding/json"
	library "hanyny/app/library"
	models "hanyny/app/models"
	system "hanyny/app/utils/system"
	v "hanyny/app/utils/view"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/fogleman/gg"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	"github.com/nfnt/resize"
)

// PostGet ...
func PostGet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")

	db := models.OpenDB()

	key := "PostGet" + offset + limit
	var posts []models.Post
	err := library.CacheGet(key, &posts)
	if err != nil {
		db.Offset(offset).Limit(limit).Select("id,title,slug,description,thumbnail,user_id,updated_at,created_at").Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, slug, name, level, category_posts.primary, category_posts.category_id, category_posts.post_id").Where("categories.level = ? AND category_posts.primary = ?", 1, true)
		}).Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, first_name, last_name")
		}).Where("status = ?", system.GetStatus().Active).Order("updated_at desc").Find(&posts)

		if len(posts) > 0 {
			err = library.CacheSet(key, posts)
			if err != nil {
				v.Respond(w, v.Message(false, "Quá trình tạo Cache thất bại."))
				return
			}
		}
	}

	data := v.Message(true, "Lấy bài viết thành công.")
	data["posts"] = posts
	v.Respond(w, data)
	return
}

// PostGetSitemap ...
func PostGetSitemap(w http.ResponseWriter, r *http.Request) {
	db := models.OpenDB()

	key := "PostGetSitemap"
	var posts []models.Post
	err := library.CacheGet(key, &posts)
	if err != nil {
		db.Select("id,slug,updated_at,created_at").Where("status = ?", system.GetStatus().Active).Order("created_at desc").Find(&posts)
		if len(posts) > 0 {
			err = library.CacheSet(key, posts)
			if err != nil {
				v.Respond(w, v.Message(false, "Quá trình tạo Cache thất bại."))
				return
			}
		}
	}
	data := v.Message(true, "Lấy bài viết thành công.")
	data["posts"] = posts
	v.Respond(w, data)
}

// PostGetRelated ...
func PostGetRelated(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")
	if limit == "" {
		limit = "10"
	}

	key := "PostGetRelated" + slug + offset + limit
	var posts []models.Post
	err := library.CacheGet(key, &posts)
	if err != nil {
		db := models.OpenDB()

		var post models.Post
		db.Where("slug = ? AND status = ?", slug, system.GetStatus().Active).First(&post)
		if post.ID == 0 {
			v.Respond(w, v.Message(false, "Đường dẫn bài viết không chính xác."))
			return
		}

		words := strings.Split(post.Title, " ")
		sqlLike := ""
		sqlOrderBy := ""
		for _, word := range words {
			sqlLike += "title LIKE '%" + word + "%' OR "
			sqlOrderBy += "(CASE WHEN title LIKE '%" + word + "%' THEN 1 ELSE 0 END) +"
		}
		sqlLike = strings.TrimRight(sqlLike, " OR ")
		sqlOrderBy = strings.TrimRight(sqlOrderBy, " +")
		db.Offset(offset).Limit(limit).Select("id,title,slug,thumbnail,user_id").Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, first_name, last_name")
		}).Where(
			"status = ? AND ("+sqlLike+")",
			system.GetStatus().Active,
		).Not("id", post.ID).Order(sqlOrderBy + " DESC").Find(&posts)

		if len(posts) > 0 {
			err = library.CacheSet(key, posts)
			if err != nil {
				v.Respond(w, v.Message(false, "Quá trình tạo Cache thất bại."))
				return
			}
		}
	}

	data := v.Message(true, "Lấy bài viết liên quan thành công.")
	data["posts"] = posts
	v.Respond(w, data)
}

// PostGetOne ...
func PostGetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	key := "PostGetOne" + slug
	var post models.Post
	err := library.CacheGet(key, &post)
	if err != nil {
		db := models.OpenDB()
		db.Where("slug = ? AND status = ?", slug, system.GetStatus().Active).Preload("Attachment", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, url")
		}).Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, slug, name, level, category_posts.primary, category_posts.category_id, category_posts.post_id").Where("category_posts.primary = ?", true).Order("level asc")
		}).Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, first_name, last_name")
		}).Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Where("parent_id = ?", "").Preload("User", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, username, first_name, last_name, avatar, level")
			}).Preload("Comments", func(db *gorm.DB) *gorm.DB {
				return db.Preload("User", func(db *gorm.DB) *gorm.DB {
					return db.Select("id, username, first_name, last_name, avatar, level")
				})
			})
		}).First(&post)
		if post.ID == 0 {
			v.Respond(w, v.Message(false, "Đường dẫn bài viết không chính xác."))
			return
		}

		if post.ID != 0 {
			err = library.CacheSet(key, post)
			if err != nil {
				v.Respond(w, v.Message(false, "Quá trình tạo Cache thất bại."))
				return
			}
		}
	}

	go countViewPost(slug)

	data := v.Message(true, "Lấy bài viết thành công.")
	data["post"] = post
	v.Respond(w, data)
	return
}

func countViewPost(slug string) {
	db := models.OpenDB()
	var post models.Post
	err := db.Where("slug = ? AND status = ?", slug, system.GetStatus().Active).Last(&post).Error
	if err == nil {
		post.View = post.View + 1
		db.Save(&post)
	}
}

// PostAdd ...
func PostAdd(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
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
	post.UserID = user.ID
	post.Status = system.GetStatus().Private

	if _, err := govalidator.ValidateStruct(post); err != nil {
		v.Respond(w, v.Message(false, err.Error()))
		return
	}

	db := models.OpenDB()

	// Kiểm tra title đã đăng chưa?
	var checkTitle models.Post
	db.Where("title = ?", post.Title).First(&checkTitle)
	if checkTitle.ID != 0 {
		v.Respond(w, v.Message(false, "Bài viết của bạn đã được chia sẻ rồi."))
		return
	}

	// Kiểm tra lịch sử
	var history models.History
	db.Where("type = ? AND user_id = ?", system.GetTypeHistory().AddPost, user.ID).First(&history)

	if history.ID == 0 { // Chưa có history
		history.Count = 1
		history.Type = system.GetTypeHistory().AddPost
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
			v.Respond(w, v.Message(false, "Tài khoản của bạn đã bị khoá do đăng bài quá nhiều. Vui lòng liên hệ với admin để khôi phục."))
			return
		}
	}

	// Tạo bài viết mới
	err = db.Create(&post).Error
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Có lỗi xảy ra khi lưu dữ liệu, vui lòng thử lại sau."))
		return
	}

	data := v.Message(true, "Tạo bài viết mới thành công.")
	data["post"] = post
	v.Respond(w, data)
}

// PostUpdate ...
func PostUpdate(w http.ResponseWriter, r *http.Request) {

}

// PostDelete ...
func PostDelete(w http.ResponseWriter, r *http.Request) {

}

// AdminPostGet ...
func AdminPostGet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")

	db := models.OpenDB()
	var posts []models.Post
	db.Where("status = ?", system.GetStatus().Active).Offset(offset).Limit(limit).Select("id, title, slug, created_at, updated_at, view, status, user_id").Preload("Categories", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, slug, name, category_posts.category_id, category_posts.post_id").Where("categories.level = ? AND category_posts.primary = ?", 1, true)
	}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, first_name, last_name")
	}).Order("id desc").Find(&posts)

	data := v.Message(true, "Lấy danh sách bài viết thành công.")
	data["posts"] = posts
	v.Respond(w, data)
}

// AdminPostGetPrivate ...
func AdminPostGetPrivate(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")

	db := models.OpenDB()
	var posts []models.Post
	db.Offset(offset).Limit(limit).Select("id, title, slug, created_at, updated_at, view, status, user_id").Where("status = ?", system.GetStatus().Private).Preload("Categories", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, slug, name, category_posts.category_id, category_posts.post_id").Where("categories.level = ? AND category_posts.primary = ?", 1, true)
	}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, first_name, last_name")
	}).Order("id desc").Find(&posts)

	data := v.Message(true, "Lấy danh sách bài viết thành công.")
	data["posts"] = posts
	v.Respond(w, data)
}

// AdminPostGetOne ...
func AdminPostGetOne(w http.ResponseWriter, r *http.Request) {
	db := models.OpenDB()

	vars := mux.Vars(r)
	id := vars["id"]

	var post models.Post
	db.Where("id = ?", id).Preload("Attachment", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, url")
	}).Preload("Categories", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, slug, name, level, category_posts.primary, category_posts.category_id, category_posts.post_id")
	}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, first_name, last_name")
	}).Order("id desc").Find(&post)
	if post.ID == 0 {
		v.Respond(w, v.Message(false, "ID bài viết không chính xác."))
		return
	}
	data := v.Message(true, "Lấy bài viết thành công.")
	data["post"] = post
	v.Respond(w, data)
}

// AdminPostTotal ...
func AdminPostTotal(w http.ResponseWriter, r *http.Request) {
	db := models.OpenDB()

	total := 0
	var posts []models.Post
	db.Where("status = ?", system.GetStatus().Active).Find(&posts).Count(&total)

	data := v.Message(true, "")
	data["total"] = total
	v.Respond(w, data)
}

// AdminPostUpdate ...
func AdminPostUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var post models.Post
	db := models.OpenDB()
	db.Where("id = ?", id).First(&post)
	if post.ID == 0 {
		v.Respond(w, v.Message(false, "ID bài viết không chính xác."))
		return
	}

	var adminPost models.ValidatePost
	err := json.NewDecoder(r.Body).Decode(&adminPost)
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, err.Error()))
		return
	}
	defer r.Body.Close()

	if _, err := govalidator.ValidateStruct(adminPost); err != nil {
		v.Respond(w, v.Message(false, err.Error()))
		return
	}

	// Kiểm tra slug đã tồn tại chưa
	var checkSlug models.Post
	adminPost.Slug = slug.Make(adminPost.Slug)
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Không thể chuyển đường dẫn tĩnh sang dạng đơn giản."))
		return
	}
	db.Where("slug = ?", adminPost.Slug).First(&checkSlug)
	if checkSlug.ID != 0 && checkSlug.ID != post.ID {
		v.Respond(w, v.Message(false, "Đường dẫn tĩnh đã tồn tại."))
		return
	}

	if adminPost.Status == system.GetStatus().Active && strings.Contains(adminPost.Thumbnail, "store/") {
		// Cập nhật đường dẫn hoàn thiện của thumbnail
		url := adminPost.Thumbnail
		url = strings.Split(url, "?")[0]
		arrayURL := strings.Split(url, "/")
		thumbnail := arrayURL[len(arrayURL)-1]
		adminPost.Thumbnail = thumbnail

		fileInput := "public/thumbnail/store/" + adminPost.Thumbnail
		fileOutput := "public/uploads/thumbnail/" + adminPost.Thumbnail
		err = library.MoveFile(fileInput, fileOutput)
		if err != nil {
			logger := library.Logger{Type: "ERROR"}
			log := logger.Open()
			log.Println(err.Error())
			logger.Close()
			v.Respond(w, v.Message(false, "Có lỗi xảy ra khi tạo ảnh đại diện, vui lòng thử lại sau."))
			return
		}
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Cập nhật bài viết
	if err := tx.Model(&post).Select("title", "slug", "description", "content", "thumbnail", "attachment_id", "status").Updates(adminPost).Error; err != nil {
		tx.Rollback()

		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Có lỗi xảy ra khi cập nhật dữ liệu, vui lòng thử lại sau."))
		return
	}

	// Cập nhật danh sách chuyên mục
	// Array Remove Duplicates
	var newCategories []models.Category
	for _, category := range adminPost.Categories {
		flag := true
		for _, newCategory := range newCategories {
			if category.ID == newCategory.ID {
				flag = true
				break
			}
		}
		if flag {
			newCategories = append(newCategories, *category)
		}
	}
	var categoryPosts []models.CategoryPost
	db.Where("post_id = ?", post.ID).Find(&categoryPosts)
	for _, categoryPost := range categoryPosts { // Danh sách cũ
		flag := false
		for i, newCategory := range adminPost.Categories { // Danh sachs mới
			if newCategory.ID == categoryPost.CategoryID {
				flag = true                                                                            // Nếu đã tồn tại trong cơ sở dữ liệu
				adminPost.Categories = append(adminPost.Categories[:i], adminPost.Categories[i+1:]...) //Nếu đã tồn tại thì xoá khỏi list mới
				if categoryPost.Primary != newCategory.Primary {                                       // Nếu khác trạng thái thì cập nhật
					if err := tx.Model(&categoryPost).Update("primary", newCategory.Primary).Error; err != nil {
						tx.Rollback()
						logger := library.Logger{Type: "ERROR"}
						log := logger.Open()
						log.Println(err.Error())
						logger.Close()
						v.Respond(w, v.Message(false, "Có lỗi xảy ra khi cập nhật dữ liệu, vui lòng thử lại sau."))
						return
					}
				}
			}
		}
		if !flag { // Không nằm trong list mới thì xoá
			if err := db.Where("post_id = ? AND category_id =?", categoryPost.PostID, categoryPost.CategoryID).Delete(&categoryPost).Error; err != nil {
				logger := library.Logger{Type: "ERROR"}
				log := logger.Open()
				log.Println(err.Error())
				logger.Close()
				v.Respond(w, v.Message(false, "Có lỗi xảy ra khi cập nhật dữ liệu, vui lòng thử lại sau."))
				return
			}
		}
	}
	if len(adminPost.Categories) > 0 { // Danh sách còn lại thì cập nhật thêm vào DB
		for _, category := range adminPost.Categories {
			var categoryPost models.CategoryPost
			categoryPost.Primary = category.Primary
			categoryPost.CategoryID = category.ID
			categoryPost.PostID = post.ID
			if err := tx.Save(&categoryPost).Error; err != nil {
				tx.Rollback()
				logger := library.Logger{Type: "ERROR"}
				log := logger.Open()
				log.Println(err.Error())
				logger.Close()
				v.Respond(w, v.Message(false, "Có lỗi xảy ra khi cập nhật dữ liệu, vui lòng thử lại sau."))
				return
			}
		}
	}
	tx.Commit()
	library.CacheCleanAll()

	data := v.Message(true, "Cập nhật bài viết thành công.")
	data["post"] = post
	v.Respond(w, data)
}

// AdminPostAdd ...
func AdminPostAdd(w http.ResponseWriter, r *http.Request) {
	db := models.OpenDB()

	var adminPost models.ValidatePost
	err := json.NewDecoder(r.Body).Decode(&adminPost)
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, err.Error()))
		return
	}
	defer r.Body.Close()

	if _, err := govalidator.ValidateStruct(adminPost); err != nil {
		v.Respond(w, v.Message(false, err.Error()))
		return
	}

	// Kiểm tra slug đã tồn tại chưa
	var checkSlug models.Post
	adminPost.Slug = slug.Make(adminPost.Slug)
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Không thể chuyển đường dẫn tĩnh sang dạng đơn giản."))
		return
	}
	db.Where("slug = ?", adminPost.Slug).First(&checkSlug)
	if checkSlug.ID != 0 {
		v.Respond(w, v.Message(false, "Đường dẫn tĩnh đã tồn tại."))
		return
	}

	if adminPost.Status == system.GetStatus().Active && !strings.Contains(adminPost.Thumbnail, "uploads/") {
		// Cập nhật đường dẫn hoàn thiện của thumbnail
		url := adminPost.Thumbnail
		url = strings.Split(url, "?")[0]
		arrayURL := strings.Split(url, "/")
		thumbnail := arrayURL[len(arrayURL)-1]
		adminPost.Thumbnail = thumbnail

		fileInput := "public/thumbnail/store/" + adminPost.Thumbnail
		fileOutput := "public/uploads/thumbnail/" + adminPost.Thumbnail
		err = library.MoveFile(fileInput, fileOutput)
		if err != nil {
			logger := library.Logger{Type: "ERROR"}
			log := logger.Open()
			log.Println(err.Error())
			logger.Close()
			v.Respond(w, v.Message(false, "Có lỗi xảy ra khi tạo ảnh đại diện, vui lòng thử lại sau."))
			return
		}
	}

	var post models.Post
	post.Title = adminPost.Title
	post.Slug = adminPost.Slug
	post.Description = adminPost.Description
	post.Content = adminPost.Content
	post.Thumbnail = adminPost.Thumbnail
	post.AttachmentID = adminPost.AttachmentID
	post.Status = adminPost.Status
	user := r.Context().Value("user").(models.User)
	post.UserID = user.ID

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Cập nhật bài viết
	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Có lỗi xảy ra khi cập nhật dữ liệu, vui lòng thử lại sau."))
		return
	}

	// Cập nhật danh sách chuyên mục
	// Array Remove Duplicates
	var newCategories []models.Category
	for _, category := range adminPost.Categories {
		flag := true
		for _, newCategory := range newCategories {
			if category.ID == newCategory.ID {
				flag = true
				break
			}
		}
		if flag {
			newCategories = append(newCategories, *category)
		}
	}
	var categoryPosts []models.CategoryPost
	db.Where("post_id = ?", post.ID).Find(&categoryPosts)
	for _, categoryPost := range categoryPosts { // Danh sách cũ
		flag := false
		for i, newCategory := range adminPost.Categories { // Danh sachs mới
			if newCategory.ID == categoryPost.CategoryID {
				flag = true                                                                            // Nếu đã tồn tại trong cơ sở dữ liệu
				adminPost.Categories = append(adminPost.Categories[:i], adminPost.Categories[i+1:]...) //Nếu đã tồn tại thì xoá khỏi list mới
				if categoryPost.Primary != newCategory.Primary {                                       // Nếu khác trạng thái thì cập nhật
					if err := tx.Model(&categoryPost).Update("primary", newCategory.Primary).Error; err != nil {
						tx.Rollback()
						logger := library.Logger{Type: "ERROR"}
						log := logger.Open()
						log.Println(err.Error())
						logger.Close()
						v.Respond(w, v.Message(false, "Có lỗi xảy ra khi cập nhật dữ liệu, vui lòng thử lại sau."))
						return
					}
				}
			}
		}
		if !flag { // Không nằm trong list mới thì xoá
			if err := tx.Delete(&categoryPost).Error; err != nil {
				tx.Rollback()
				logger := library.Logger{Type: "ERROR"}
				log := logger.Open()
				log.Println(err.Error())
				logger.Close()
				v.Respond(w, v.Message(false, "Có lỗi xảy ra khi cập nhật dữ liệu, vui lòng thử lại sau."))
				return
			}
		}
	}
	if len(adminPost.Categories) > 0 { // Danh sách còn lại thì cập nhật thêm vào DB
		for _, category := range adminPost.Categories {
			var categoryPost models.CategoryPost
			categoryPost.Primary = category.Primary
			categoryPost.CategoryID = category.ID
			categoryPost.PostID = post.ID
			if err := tx.Save(&categoryPost).Error; err != nil {
				tx.Rollback()
				logger := library.Logger{Type: "ERROR"}
				log := logger.Open()
				log.Println(err.Error())
				logger.Close()
				v.Respond(w, v.Message(false, "Có lỗi xảy ra khi cập nhật dữ liệu, vui lòng thử lại sau."))
				return
			}
		}
	}
	tx.Commit()
	library.CacheCleanAll()

	data := v.Message(true, "Thêm bài viết thành công.")
	data["post"] = post
	v.Respond(w, data)
}

// AdminPostDelete ...
func AdminPostDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var post models.Post
	db := models.OpenDB()
	db.Where("id = ?", id).First(&post)
	if post.ID == 0 {
		v.Respond(w, v.Message(false, "ID bài viết không chính xác."))
		return
	}

	err := db.Delete(&post).Error
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Có lỗi xảy ra khi cập nhật dữ liệu, vui lòng thử lại sau."))
		return
	}

	library.CacheCleanAll()

	data := v.Message(true, "Xoá bài viết thành công.")
	v.Respond(w, data)
	return
}

// AdminPostSearch ...
func AdminPostSearch(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")
	str := params.Get("s")

	db := models.OpenDB()
	var posts []models.Post

	words := strings.Split(str, " ")
	sqlLike := ""
	sqlOrderBy := ""
	for _, word := range words {
		sqlLike += "title LIKE '%" + word + "%' OR "
		sqlOrderBy += "(CASE WHEN title LIKE '%" + word + "%' THEN 1 ELSE 0 END) +"
	}
	sqlLike = strings.TrimRight(sqlLike, " OR ")
	sqlOrderBy = strings.TrimRight(sqlOrderBy, " +")

	db.Offset(offset).Limit(limit).Select("id, title, slug, created_at, updated_at, status, user_id").Preload("Categories", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, slug, name, category_posts.category_id, category_posts.post_id").Where("categories.level = ? AND category_posts.primary = ?", 1, true)
	}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, first_name, last_name")
	}).Where(
		"(" + sqlLike + ")",
	).Order(sqlOrderBy + " DESC").Find(&posts)

	data := v.Message(true, "Tìm bài viết thành công.")
	data["posts"] = posts
	v.Respond(w, data)
}

// PostSearch ...
func PostSearch(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")
	str := params.Get("s")

	if limit == "" {
		limit = "5"
	}

	key := "PostSearch" + str + offset + limit
	var posts []models.Post
	err := library.CacheGet(key, &posts)
	if err != nil {
		db := models.OpenDB()

		words := strings.Split(str, " ")
		sqlLike := ""
		sqlOrderBy := ""
		for _, word := range words {
			sqlLike += "title LIKE '%" + word + "%' OR "
			sqlOrderBy += "(CASE WHEN title LIKE '%" + word + "%' THEN 1 ELSE 0 END) +"
		}
		sqlLike = strings.TrimRight(sqlLike, " OR ")
		sqlOrderBy = strings.TrimRight(sqlOrderBy, " +")
		db.Offset(offset).Limit(limit).Select("id,title,slug,description,thumbnail,user_id,updated_at,created_at").Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, slug, name, level, category_posts.primary, category_posts.category_id, category_posts.post_id").Where("categories.level = ? AND category_posts.primary = ?", 1, true)
		}).Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, first_name, last_name")
		}).Where(
			"status = ? AND ("+sqlLike+")",
			system.GetStatus().Active,
		).Order(sqlOrderBy + " DESC").Find(&posts)

		if len(posts) > 0 {
			err = library.CacheSet(key, posts)
			if err != nil {
				v.Respond(w, v.Message(false, "Quá trình tạo Cache thất bại."))
				return
			}
		}
	}

	data := v.Message(true, "Tìm bài viết thành công.")
	data["posts"] = posts
	v.Respond(w, data)
}

// PostSearchSpeed ...
func PostSearchSpeed(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")
	str := params.Get("s")

	key := "PostSearchSpeed" + str + offset + limit
	var posts []models.Post
	err := library.CacheGet(key, &posts)
	if err != nil {
		db := models.OpenDB()

		words := strings.Split(str, " ")
		sqlLike := ""
		sqlOrderBy := ""
		for _, word := range words {
			sqlLike += "title LIKE '%" + word + "%' OR "
			sqlOrderBy += "(CASE WHEN title LIKE '%" + word + "%' THEN 1 ELSE 0 END) +"
		}
		sqlLike = strings.TrimRight(sqlLike, " OR ")
		sqlOrderBy = strings.TrimRight(sqlOrderBy, " +")
		db.Offset(offset).Limit(limit).Select("id, title, slug").Where(
			"status = ? AND ("+sqlLike+")",
			system.GetStatus().Active,
		).Order(sqlOrderBy + " DESC").Find(&posts)

		if len(posts) > 0 {
			err = library.CacheSet(key, posts)
			if err != nil {
				v.Respond(w, v.Message(false, "Quá trình tạo Cache thất bại."))
				return
			}
		}
	}

	data := v.Message(true, "Tìm bài viết thành công.")
	data["posts"] = posts
	v.Respond(w, data)
}

// PostThumbnail ...
func PostThumbnail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := vars["link"]

	var re = regexp.MustCompile(`(?m)\d*x\d*_`)
	var reOld = regexp.MustCompile(`(?m)\d*x\d*-`)
	match := re.FindAllString(link, -1)
	size := ""
	if len(match) > 0 {
		size = match[0]
		size = strings.Replace(size, "_", "", -1)
	} else {
		match := reOld.FindAllString(link, -1)
		size = ""
		if len(match) > 0 {
			size = match[0]
			size = strings.Replace(size, "-", "", -1)
		}
	}
	img, err := gg.LoadImage("public/uploads/thumbnail/" + link)
	buffer := new(bytes.Buffer)
	if err != nil {
		// Tạo mới
		sizeArray := strings.Split(size, "x")
		width, err := strconv.Atoi(sizeArray[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		height, err := strconv.Atoi(sizeArray[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newLink := re.ReplaceAllString(link, "")
		newLink = reOld.ReplaceAllString(newLink, "")
		im, err := gg.LoadImage("public/uploads/thumbnail/" + newLink)
		if err != nil {
			im, err = gg.LoadImage("public/uploads/images/" + newLink)
			if err != nil {
				img, _ = gg.LoadImage("public/banner.jpg")
				jpeg.Encode(buffer, img, nil)
				w.Header().Set("Content-Type", "image/jpeg")
				w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
				if _, err := w.Write(buffer.Bytes()); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}
		}
		img = resize.Resize(uint(width), uint(height), im, resize.Lanczos3)
		err = jpeg.Encode(buffer, img, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ioutil.WriteFile("public/uploads/thumbnail/"+link, buffer.Bytes(), 0644)
	} else {
		err = jpeg.Encode(buffer, img, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
