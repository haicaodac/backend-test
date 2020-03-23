package controllers

import (
	library "hanyny/app/library"
	models "hanyny/app/models"
	v "hanyny/app/utils/view"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// AttachmentGet ...
func AttachmentGet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")

	db := models.OpenDB()
	user := r.Context().Value("user").(models.User)

	var attachments []models.Attachment
	db.Offset(offset).Limit(limit).Where("user_id = ?", user.ID).Find(&attachments)

	data := v.Message(true, "Lấy tệp tin thành công.")
	data["attachments"] = attachments
	v.Respond(w, data)
}

// AttachmentGetOne ...
func AttachmentGetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var attachment models.Attachment
	db := models.OpenDB()
	db.Where("id = ?", id).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, first_name, last_name")
	}).Preload("Post", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, slug")
	}).First(&attachment)
	if attachment.ID == 0 {
		v.Respond(w, v.Message(false, "ID của tệp tin không chính xác."))
		return
	}

	data := v.Message(true, "Lấy tệp tin thành công.")
	data["attachment"] = attachment
	v.Respond(w, data)
	return
}

// AttachmentAdd ...
func AttachmentAdd(w http.ResponseWriter, r *http.Request) {
	file, handle, err := r.FormFile("file")
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
	url := ""

	arrayType := []string{
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/pdf",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"application/vnd.ms-powerpoint",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		// "application/zip",
		// "application/x-rar-compressed",
	}
	if !library.ArrayStringContains(arrayType, mimeType) {
		v.Respond(w, v.Message(false, "Tài liệu đính kèm phải là dạng .doc/.docx/.pdf/.pptx/.xls/.xlsx"))
		// v.Respond(w, v.Message(false, "Tài liệu đính kèm phải là dạng .doc/.docx/.pdf/.pptx/.xls/.xlsx/.zip/.rar"))
		return
	}

	db := models.OpenDB()
	user := r.Context().Value("user").(models.User)

	// Kiểm tra tệp đã tồn tại chưa
	var checkAttachment models.Attachment
	db.Where("name = ? AND type = ? AND size = ? AND user_id = ?", handle.Filename, mimeType, handle.Size, user.ID).First(&checkAttachment)
	if checkAttachment.ID != 0 {
		data := v.Message(true, "Thêm tài liệu mới thành công.")
		data["attachment"] = checkAttachment
		v.Respond(w, data)
		return
	}

	url, err = library.UploadFileToServer(file, handle)
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Có lỗi xảy ra trong quá trình lưu file, vui lòng thử lại sau."))
		return
	}

	var attachment models.Attachment
	attachment.Name = handle.Filename
	attachment.Type = mimeType
	attachment.Size = handle.Size
	attachment.URL = url
	attachment.UserID = user.ID

	// Lưu vào cơ sở dữ liệu
	err = db.Create(&attachment).Error
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Có lỗi xảy ra khi lưu dữ liệu, vui lòng thử lại sau."))
		return
	}

	data := v.Message(true, "Thêm tài liệu mới thành công.")
	data["attachment"] = attachment
	v.Respond(w, data)
}

// AttachmentUpdate ...
func AttachmentUpdate(w http.ResponseWriter, r *http.Request) {

}

// AttachmentDelete ...
func AttachmentDelete(w http.ResponseWriter, r *http.Request) {

}

// AdminAttachmentGet ...
func AdminAttachmentGet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")

	db := models.OpenDB()

	var attachments []models.Attachment
	db.Offset(offset).Limit(limit).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, first_name, last_name")
	}).Order("created_at DESC").Find(&attachments)

	data := v.Message(true, "Lấy tệp tin thành công.")
	data["attachments"] = attachments
	v.Respond(w, data)
}
