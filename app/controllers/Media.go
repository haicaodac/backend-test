package controllers

import (
	"encoding/json"
	library "hanyny/app/library"
	libthumb "hanyny/app/library/libthumb"
	models "hanyny/app/models"
	v "hanyny/app/utils/view"
	"net/http"
)

// MediaGet ...
func MediaGet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")

	db := models.OpenDB()

	var media []models.Media
	db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&media)

	data := v.Message(true, "Get media successully")
	data["media"] = media
	v.Respond(w, data)
}

// MediaGetOne ...
func MediaGetOne(w http.ResponseWriter, r *http.Request) {

}

// MediaAdd ...
func MediaAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	fhs := r.MultipartForm.File["files"]

	var listMedia []models.Media
	for _, handle := range fhs {
		file, err := handle.Open()
		if err != nil {
			logger := library.Logger{Type: "ERROR"}
			log := logger.Open()
			log.Println(err.Error())
			logger.Close()
			v.Respond(w, v.Message(false, "Có lỗi xảy ra trong quá trình mở tài liệu, vui lòng thử lại sau."))
			return
		}
		defer file.Close()
		mimeType := handle.Header.Get("Content-Type")
		arrayType := []string{
			"image/jpeg",
			"image/jpg",
			"image/png",
			"image/gif",
			"image/webp",
		}
		if !library.ArrayStringContains(arrayType, mimeType) {
			v.Respond(w, v.Message(false, "Hình ảnh phải có dạng .jpeg/.png/.jpg"))
			return
		}

		db := models.OpenDB()
		user := r.Context().Value("user").(models.User)

		// Kiểm tra tệp đã tồn tại chưa
		var checkMedia models.Media
		db.Where("name = ? AND type = ? AND user_id = ?", handle.Filename, mimeType, user.ID).First(&checkMedia)
		if checkMedia.ID != 0 {
			listMedia = append(listMedia, checkMedia)
			continue
		}

		url, err := library.UploadImageToServer(file, handle)
		if err != nil {
			logger := library.Logger{Type: "ERROR"}
			log := logger.Open()
			log.Println(err.Error())
			logger.Close()
			v.Respond(w, v.Message(false, "Có lỗi xảy ra trong quá trình lưu file, vui lòng thử lại sau."))
			return
		}

		var media models.Media
		media.Name = handle.Filename
		media.Type = mimeType
		media.URL = url
		media.UserID = user.ID

		// Lưu vào cơ sở dữ liệu
		err = db.Create(&media).Error
		if err != nil {
			logger := library.Logger{Type: "ERROR"}
			log := logger.Open()
			log.Println(err.Error())
			logger.Close()
			v.Respond(w, v.Message(false, "Có lỗi xảy ra khi lưu dữ liệu, vui lòng thử lại sau."))
			return
		}
		listMedia = append(listMedia, media)
	}

	data := v.Message(true, "Upload images success.")
	data["media"] = listMedia
	v.Respond(w, data)
}

// MediaUpdate ...
func MediaUpdate(w http.ResponseWriter, r *http.Request) {

}

// MediaDelete ...
func MediaDelete(w http.ResponseWriter, r *http.Request) {

}

// AdminMediaCreateThumbnail ...
func AdminMediaCreateThumbnail(w http.ResponseWriter, r *http.Request) {
	var media models.Media
	err := json.NewDecoder(r.Body).Decode(&media)
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Dữ liệu không chính xác."))
		return
	}
	defer r.Body.Close()
	if media.Name == "" {
		v.Respond(w, v.Message(false, "Bắt buộc phải nhập tên."))
		return
	}
	path, err := libthumb.New(media.Name, 650, 365)
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Quá trình tạo ảnh gặp sự cố."))
		return
	}
	url := path
	data := v.Message(true, "Tạo ảnh đại diện thành công.")
	data["url"] = url
	v.Respond(w, data)
}
