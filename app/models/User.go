package models

// User in database users
type User struct {
	Model
	Password   string `json:"-" gorm:"type:varchar(200);not null"`
	Username   string `json:"username,omitempty" gorm:"type:varchar(50);not null"`
	FirstName  string `json:"first_name,omitempty" gorm:"type:varchar(15);not null"`
	LastName   string `json:"last_name,omitempty" gorm:"type:varchar(45);not null"`
	Email      string `json:"email,omitempty" gorm:"type:varchar(200);not null"`
	Avatar     string `json:"avatar,omitempty" gorm:"type:text;" valid:"url"`
	Level      string `json:"level,omitempty" gorm:"type:varchar(20);not null;default:'user'"`
	IP         string `json:"-" gorm:"type:varchar(200);"`
	FacebookID string `json:"-" gorm:"type:varchar(100);"`
	GoogleID   string `json:"-" gorm:"type:varchar(100);"`
	Status     string `json:"status,omitempty" gorm:"type:varchar(20);not null;default:'active'"`
}

// BaseUser ...
type BaseUser struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Level     string `json:"level"`
	Avatar    string `json:"avatar"`
}

// Account ...
type Account struct {
	Password  string `json:"password" valid:"required~Mật khẩu không thể trống.,runelength(6|50)~Mật khẩu phải từ 6 đến 50 ký tự."`
	FirstName string `json:"first_name" valid:"required~Tên không thể trống.,runelength(1|15)~Tên phải có độ dài từ 1 đến 15 ký tự."`
	LastName  string `json:"last_name" valid:"required~Họ không thể trống,runelength(1|45)~Họ phải có độ dài từ 1 đến 45 ký tự."`
	Email     string `json:"email" valid:"required~Địa chỉ Email không thể trống,runelength(1|200)~Địa chỉ Email từ 1 đến 200 ký tự.,email~Không đúng định dạng Email."`
	Username  string `json:"username"`
}

// AccountFacebook ...
type AccountFacebook struct {
	FacebookID  string `json:"facebook_id"`
	AccessToken string `json:"access_token"`
}

// AccountPassword update password
type AccountPassword struct {
	OldPassword string `json:"old_password" valid:"required~Mật khẩu cũ không thể để trống.,runelength(6|50)~Mật khẩu cũ không chính xác."`
	Password    string `json:"password" valid:"required~Mật khẩu mới không thể để trống.,runelength(6|50)~Mật khẩu mới phải có độ dài từ 6 đến 50 ký tự."`
}
