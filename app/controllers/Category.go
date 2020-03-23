package controllers

import (
	"encoding/json"
	library "hanyny/app/library"
	models "hanyny/app/models"
	system "hanyny/app/utils/system"
	v "hanyny/app/utils/view"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// CategoryGet ...
func CategoryGet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")

	key := "CategoryGet" + offset + limit
	var categories []models.Category
	err := library.CacheGet(key, &categories)
	if err != nil {
		db := models.OpenDB()
		db.Offset(offset).Limit(limit).Where("status = ?", system.GetStatus().Active).Find(&categories)

		err = library.CacheSet(key, categories)
		if err != nil {
			v.Respond(w, v.Message(false, "Quá trình tạo Cache thất bại."))
			return
		}
	}
	data := v.Message(true, "Lấy chuyên mục thành công.")
	data["categories"] = categories
	v.Respond(w, data)
}

// CategoryGetSitemap ...
func CategoryGetSitemap(w http.ResponseWriter, r *http.Request) {

	key := "CategoryGetSitemap"
	var categories []models.Category
	err := library.CacheGet(key, &categories)
	if err != nil {
		db := models.OpenDB()
		db.Select("id,slug,created_at,updated_at").Where("status = ?", system.GetStatus().Active).Find(&categories)
		err = library.CacheSet(key, categories)
		if err != nil {
			v.Respond(w, v.Message(false, "Quá trình tạo Cache thất bại."))
			return
		}
	}
	data := v.Message(true, "Lấy chuyên mục thành công.")
	data["categories"] = categories
	v.Respond(w, data)
}

// CategoryGetOne ...
func CategoryGetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")

	key := "CategoryGetOne" + slug + offset + limit
	var category models.Category
	err := library.CacheGet(key, &category)
	if err != nil {
		db := models.OpenDB()
		db.Where("slug = ?", slug).Last(&category)
		if category.ID == 0 {
			v.Respond(w, v.Message(false, "Không tìm thấy chuyên mục."))
			return
		}

		var categoryPosts []models.CategoryPost
		db.Where("category_id = ?", category.ID).Find(&categoryPosts)

		listID := []uint{0}
		for _, categoryPost := range categoryPosts {
			listID = append(listID, categoryPost.PostID)
		}

		var posts []*models.Post
		db.Offset(offset).Limit(limit).Select("id,title,slug,description,thumbnail,user_id,updated_at,created_at").Preload("Categories", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, slug, name, level, category_posts.primary, category_posts.category_id, category_posts.post_id").Where("categories.level = ? AND category_posts.primary = ?", 1, true)
		}).Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, first_name, last_name")
		}).Where("status = ? AND id in (?)", system.GetStatus().Active, listID).Order("updated_at desc").Find(&posts)

		category.Posts = posts

		if category.ID == 0 {
			v.Respond(w, v.Message(false, "Đường dẫn chuyên mục không chính xác."))
			return
		}

		err = library.CacheSet(key, category)
		if err != nil {
			v.Respond(w, v.Message(false, "Quá trình tạo Cache thất bại."))
			return
		}
	}

	data := v.Message(true, "success")
	data["category"] = category
	v.Respond(w, data)
	return
}

// AdminCategoryAdd ...
func AdminCategoryAdd(w http.ResponseWriter, r *http.Request) {
	var validateCategory models.ValidateCategory
	err := json.NewDecoder(r.Body).Decode(&validateCategory)
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Dữ liệu không chính xác."))
		return
	}
	defer r.Body.Close()

	if _, err := govalidator.ValidateStruct(validateCategory); err != nil {
		v.Respond(w, v.Message(false, err.Error()))
		return
	}

	db := models.OpenDB()

	// Kiểm tra slug đã tồn tại chưa
	var checkSlug models.Category
	db.Where("slug = ?", validateCategory.Slug).First(&checkSlug)
	if checkSlug.ID != 0 {
		v.Respond(w, v.Message(false, "Đường dẫn tĩnh đã tồn tại."))
		return
	}

	var category models.Category
	// Cập nhật thêm người tạo
	user := r.Context().Value("user").(models.User)
	category.UserID = user.ID

	category.Name = validateCategory.Name
	category.Title = validateCategory.Title
	category.Slug = validateCategory.Slug
	category.Description = validateCategory.Description
	category.Thumbnail = validateCategory.Thumbnail
	category.Icon = validateCategory.Icon
	category.Position = validateCategory.Position
	category.ParentID = validateCategory.ParentID
	category.Level = validateCategory.Level
	category.Status = validateCategory.Status

	// Tạo chuyên mục mới
	err = db.Create(&category).Error
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Có lỗi xảy ra khi lưu dữ liệu, vui lòng thử lại sau."))
		return
	}

	library.CacheCleanAll()

	data := v.Message(true, "Tạo chuyên mục mới thành công.")
	data["category"] = category
	v.Respond(w, data)
}

// AdminCategoryUpdate ...
func AdminCategoryUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var category models.Category
	db := models.OpenDB()
	db.Where("id = ?", id).First(&category)
	if category.ID == 0 {
		v.Respond(w, v.Message(false, "ID của chuyên mục không chính xác."))
		return
	}

	var validateCategory models.ValidateCategory
	err := json.NewDecoder(r.Body).Decode(&validateCategory)
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Dữ liệu không chính xác."))
		return
	}
	defer r.Body.Close()

	if _, err := govalidator.ValidateStruct(validateCategory); err != nil {
		v.Respond(w, v.Message(false, err.Error()))
		return
	}

	// Kiểm tra slug đã tồn tại chưa
	var checkSlug models.Category
	db.Where("slug = ?", validateCategory.Slug).First(&checkSlug)
	if checkSlug.ID != 0 && checkSlug.ID != category.ID {
		v.Respond(w, v.Message(false, "Đường dẫn tĩnh đã tồn tại."))
		return
	}

	// Cập nhật lại chuyên mục
	err = db.Model(&category).Select("name", "title", "slug", "description", "thumbnail", "icon", "position", "parent_id", "level", "status").Updates(validateCategory).Error
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Có lỗi xảy ra khi cập nhật dữ liệu, vui lòng thử lại sau."))
		return
	}

	library.CacheCleanAll()

	data := v.Message(true, "Cập nhật chuyên mục thành công.")
	data["category"] = category
	v.Respond(w, data)
}

// AdminCategoryDelete ...
func AdminCategoryDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	var category models.Category
	db := models.OpenDB()
	db.Where("slug = ?", slug).First(&category)
	if category.ID == 0 {
		v.Respond(w, v.Message(false, "Đường dẫn của chuyên mục không chính xác."))
		return
	}

	err := db.Delete(&category).Error
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Có lỗi xảy ra khi cập nhật dữ liệu, vui lòng thử lại sau."))
		return
	}

	library.CacheCleanAll()

	data := v.Message(true, "Xoá chuyên mục thành công.")
	v.Respond(w, data)
	return
}

// AdminCategoryTotal ...
func AdminCategoryTotal(w http.ResponseWriter, r *http.Request) {
	db := models.OpenDB()

	total := 0
	var category models.Category
	db.Find(&category).Count(&total)

	data := v.Message(true, "")
	data["total"] = total
	v.Respond(w, data)
}

// AdminCategoryGetOne ...
func AdminCategoryGetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var category models.Category
	db := models.OpenDB()
	db.Where("id = ?", id).First(&category)
	if category.ID == 0 {
		v.Respond(w, v.Message(false, "ID chuyên mục không chính xác."))
		return
	}

	data := v.Message(true, "")
	data["category"] = category
	v.Respond(w, data)
	return
}

// AdminCategoryGet ...
func AdminCategoryGet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")

	db := models.OpenDB()

	var categories []models.Category
	db.Offset(offset).Limit(limit).Find(&categories)

	data := v.Message(true, "Lấy chuyên mục thành công.")
	data["categories"] = categories
	v.Respond(w, data)
}
