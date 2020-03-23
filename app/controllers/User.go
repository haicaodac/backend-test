package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	jobs "hanyny/app/jobs"
	library "hanyny/app/library"
	models "hanyny/app/models"
	auth "hanyny/app/utils/auth"
	system "hanyny/app/utils/system"
	v "hanyny/app/utils/view"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"

	fb "github.com/huandu/facebook"
	// fb "golang.org/x/oauth2/facebook"
)

// SignUp created user
func SignUp(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&account)
	if err != nil {
		v.Respond(w, v.Message(false, err.Error()))
		return
	}
	defer r.Body.Close()

	// Validator format
	if _, err := govalidator.ValidateStruct(account); err != nil {
		v.Respond(w, v.Message(false, err.Error()))
		return
	}

	// Validator username
	db := models.OpenDB()

	// Check email user
	var user models.User
	db.Select("id").Where("email = ?", account.Email).First(&user)
	if user.ID != 0 {
		v.Respond(w, v.Message(false, "Email đã tồn tại. Vui lòng truy cập phần lấy lại mật khẩu."))
		return
	}
	// Created User
	user.Email = account.Email
	user.Password = library.HashAndSalt(account.Password)
	user.FirstName = account.FirstName
	user.LastName = account.LastName
	user.Username = library.RandomString(10)
	db.Create(&user)

	tokenString, expiresAt := auth.CreateToken(user)

	data := v.Message(true, "Đăng ký tài khoản mới thành công.")
	data["token"] = tokenString
	data["expires"] = expiresAt
	v.Respond(w, data)
	return
}

//SignIn ...
func SignIn(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		v.Respond(w, v.Message(false, err.Error()))
		return
	}
	defer r.Body.Close()
	var user models.User
	db := models.OpenDB()

	// Kiểm tra xem ng dùng nhập username hay email
	if govalidator.IsEmail(account.Username) { // Email
		db.Where("email = ?", account.Username).First(&user)
		if user.ID == 0 {
			v.Respond(w, v.Message(false, "Địa chỉ Email không chính xác."))
			return
		}
	} else { // Username
		db.Where("username = ?", account.Username).First(&user)
		if user.ID == 0 {
			v.Respond(w, v.Message(false, "Tên đăng nhập không chính xác."))
			return
		}
	}

	if user.Status == system.GetStatus().Block {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Tài khoản của bạn đã bị khoá. Vui lòng liên hệ với quản trị viên để được khôi phục."))
		return
	}

	pwdMatch := library.ComparePasswords(user.Password, account.Password)
	if !pwdMatch {
		v.Respond(w, v.Message(false, "Mật khẩu không chính xác."))
		return
	}

	tokenString, expiresAt := auth.CreateToken(user)

	data := v.Message(true, "Đăng nhập thành công.")
	data["token"] = tokenString
	data["expires"] = expiresAt
	v.Respond(w, data)
	return
}

// SignInFacebook ...
func SignInFacebook(w http.ResponseWriter, r *http.Request) {
	var account models.AccountFacebook
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		v.Respond(w, v.Message(false, err.Error()))
		return
	}
	defer r.Body.Close()

	appFB := fb.New("985944054831668", "289ecf6d4c60e8c32ea6402f20bda6ec")

	session := appFB.Session(account.AccessToken)
	session.EnableAppsecretProof(true)

	// Use session.
	res, err := session.Get("/me", fb.Params{
		"fields": "first_name, last_name, email, birthday, picture",
	})
	if err != nil { // Email
		v.Respond(w, v.Message(false, err.Error()))
		return
	}

	var user models.User
	db := models.OpenDB()
	res.DecodeField("id", &user.FacebookID)

	if user.FacebookID == "" {
		v.Respond(w, v.Message(false, "Chúng tôi không thể đăng nhập facebook của bạn."))
		return
	}

	if user.Status == system.GetStatus().Block {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Tài khoản của bạn đã bị khoá. Vui lòng liên hệ với quản trị viên để được khôi phục."))
		return
	}

	db.Where("facebook_id = ?", user.FacebookID).Last(&user)
	res.DecodeField("first_name", &user.FirstName)
	res.DecodeField("last_name", &user.LastName)
	res.DecodeField("id", &user.FacebookID)
	res.DecodeField("picture.data.url", &user.Avatar)
	res.DecodeField("email", &user.Email)

	if !govalidator.IsEmail(user.Email) { // Email
		v.Respond(w, v.Message(false, "Chúng tôi không tìm thấy Email của bạn."))
		return
	}

	err = db.Save(&user).Error
	if err != nil { // Email
		v.Respond(w, v.Message(false, err.Error()))
		return
	}

	tokenString, expiresAt := auth.CreateToken(user)

	data := v.Message(true, "Đăng nhập thành công.")
	data["token"] = tokenString
	data["expires"] = expiresAt
	v.Respond(w, data)
	return
}

// UserForgotPassword ...
func UserForgotPassword(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		v.Respond(w, v.Message(false, err.Error()))
		return
	}
	defer r.Body.Close()

	// Cần chặn IP nếu truy cập quá nhiều

	var user models.User
	db := models.OpenDB()

	if account.Username == "" {
		v.Respond(w, v.Message(false, "Địa chỉ Email không tồn tại."))
		return
	}

	// Kiểm tra xem ng dùng nhập username hay email
	if govalidator.IsEmail(account.Username) { // Email
		db.Where("email = ?", account.Username).First(&user)
		if user.ID == 0 {
			v.Respond(w, v.Message(false, "Địa chỉ Email không tồn tại."))
			return
		}
	} else { // Username
		db.Where("username = ?", account.Username).First(&user)
		if user.ID == 0 {
			v.Respond(w, v.Message(false, "Tên đăng nhập không tồn tại."))
			return
		}
	}

	randomString := library.RandomString(32)
	var userSession models.UserSession
	userSession.Type = "forgot-password"
	userSession.UserID = user.ID
	userSession.Data = randomString
	db.Save(&userSession)

	if userSession.ID == 0 {
		v.Respond(w, v.Message(false, "Có lỗi xảy ra trong quá trình lưu trữ."))
		return
	}

	/* Gửi mail lấy lại mật khẩu */
	go jobs.ForgotPassword("forgot-password.html", user, randomString)

	data := v.Message(true, "Yêu cầu lấy lại mật khẩu thành công. Vui lòng truy cập Email để nhận hướng dẫn.")
	v.Respond(w, data)
	return
}

// UserCheckForgotPassword ...
func UserCheckForgotPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	token = strings.Trim(token, " ")
	if token == "" {
		v.Respond(w, v.Message(false, "Mã xác thực không chính xác."))
		return
	}

	var userSession models.UserSession
	db := models.OpenDB()

	db.Where("type = ? AND data = ?", "forgot-password", token).First(&userSession)
	if userSession.ID == 0 || userSession.UserID == 0 {
		v.Respond(w, v.Message(false, "Mã xác thực không chính xác."))
		return
	}

	var user models.User
	db.Where("id = ?", userSession.UserID).First(&user)
	if user.ID == 0 {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(user)
		logger.Close()
		v.Respond(w, v.Message(false, "Mã xác thực có vấn đề. Vui lòng thử lại sau."))
		return
	}

	tokenString, expiresAt := auth.CreateToken(user)

	data := v.Message(true, "Mã xác thực chính xác.")

	data["token"] = tokenString
	data["expires"] = expiresAt
	v.Respond(w, data)
	return
}

// UserNewPassword ...
func UserNewPassword(w http.ResponseWriter, r *http.Request) {
	type structBodyReq struct {
		Token       string `json:"token" valid:"required~Mã xác thực là trường bắt buộc."`
		NewPassword string `json:"new_password" valid:"required~Mật khẩu mới không thể để trống.,runelength(6|50)~Mật khẩu mới phải có độ dài từ 6 đến 50 ký tự."`
	}
	bodyReq := structBodyReq{}
	err := json.NewDecoder(r.Body).Decode(&bodyReq)
	if err != nil {
		v.Respond(w, v.Message(false, err.Error()))
		return
	}
	defer r.Body.Close()

	// Validator format
	if _, err := govalidator.ValidateStruct(bodyReq); err != nil {
		v.Respond(w, v.Message(false, err.Error()))
		return
	}

	user := r.Context().Value("user").(models.User)

	db := models.OpenDB()

	var userSession models.UserSession
	db.Where("data = ?", bodyReq.Token).First(&userSession)
	if userSession.ID == 0 || userSession.UserID == 0 {
		v.Respond(w, v.Message(false, "Mã xác thực không chính xác."))
		return
	}
	if userSession.UserID != user.ID {
		v.Respond(w, v.Message(false, "Mã xác thực không chính xác."))
		return
	}

	user.Password = library.HashAndSalt(bodyReq.NewPassword)
	db.Save(&user)

	db.Delete(&userSession)

	data := v.Message(true, "Cập nhật mật khẩu mới thành công.")
	v.Respond(w, data)
	return
}

// AdminUserTotal ...
func AdminUserTotal(w http.ResponseWriter, r *http.Request) {
	db := models.OpenDB()

	total := 0
	var users []models.User
	db.Find(&users).Count(&total)

	data := v.Message(true, "Lấy tổng số thành viên thành công.")
	data["total"] = total
	v.Respond(w, data)
}

// AdminUserGet ...
func AdminUserGet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")

	db := models.OpenDB()
	var users []models.User
	db.Offset(offset).Limit(limit).Order("id desc").Find(&users)

	data := v.Message(true, "Lấy danh sách thành viên thành công.")
	data["users"] = users
	v.Respond(w, data)
}

// AdminUserGetOne ...
func AdminUserGetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user models.User
	db := models.OpenDB()
	db.Where("id = ?", id).First(&user)

	data := v.Message(true, "Lấy thành viên thành công.")
	data["user"] = user
	v.Respond(w, data)
}

// AdminUserUpdate ...
func AdminUserUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	db := models.OpenDB()
	var user models.User
	db.Where("id = ?", id).First(&user)
	if user.ID == 0 {
		v.Respond(w, v.Message(false, "ID thành viên không chính xác."))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		v.Respond(w, v.Message(false, err.Error()))
		return
	}
	defer r.Body.Close()

	// Cập nhật thành viên
	if err := db.Model(&user).Select("first_name", "last_name", "username", "avatar", "level", "status").Updates(user).Error; err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Có lỗi xảy ra khi cập nhật dữ liệu, vui lòng thử lại sau."))
		return
	}

	data := v.Message(true, "Cập nhật thành viên thành công.")
	data["user"] = user
	v.Respond(w, data)
}

// AdminUserSearch ...
func AdminUserSearch(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	offset := params.Get("offset")
	limit := params.Get("limit")
	str := params.Get("s")

	db := models.OpenDB()
	var users []models.User
	db.Offset(offset).Limit(limit).Where("email LIKE ?", "%"+str+"%").Find(&users)

	data := v.Message(true, "Tìm thành viên thành công.")
	data["users"] = users
	v.Respond(w, data)
}

// UserSubscribeEmail ...
func UserSubscribeEmail(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	email := params.Get("email")

	ip, err := library.GetIP()
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Có lỗi xảy ra khi kiểm tra email của bạn. Vui lòng thử lại sau"))
		return
	}

	email = strings.TrimSpace(email)
	if !govalidator.IsEmail(email) { // !Email
		v.Respond(w, v.Message(false, "Email không đúng."))
		return
	}

	var emailModel models.Email
	emailModel.Email = email
	emailModel.IP = ip.String()

	var checkEmail models.Email
	db := models.OpenDB()
	db.Where("email = ?", emailModel.Email).Last(&checkEmail)
	if checkEmail.ID != 0 {
		v.Respond(w, v.Message(false, "Email này đã được đăng ký."))
		return
	}

	err = db.Create(&emailModel).Error
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
		v.Respond(w, v.Message(false, "Có lỗi xảy ra khi kiểm tra email của bạn. Vui lòng thử lại sau"))
		return
	}

	data := v.Message(true, "Đăng ký thành công.")
	v.Respond(w, data)
}
