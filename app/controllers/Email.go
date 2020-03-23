package controllers

import (
	models "hanyny/app/models"
	v "hanyny/app/utils/view"
	"net/http"
)

// AdminEmailGet ...
func AdminEmailGet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")

	db := models.OpenDB()
	var emails []models.Email
	db.Offset(offset).Limit(limit).Order("updated_at desc").Find(&emails)

	data := v.Message(true, "Lấy danh email thành công.")
	data["emails"] = emails
	v.Respond(w, data)
}

// AdminEmailTotal ...
func AdminEmailTotal(w http.ResponseWriter, r *http.Request) {
	db := models.OpenDB()

	total := 0
	var emails []models.Email
	db.Find(&emails).Count(&total)

	data := v.Message(true, "")
	data["total"] = total
	v.Respond(w, data)
}
