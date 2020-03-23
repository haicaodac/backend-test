package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	models "hanyny/app/models"
	system "hanyny/app/utils/system"
	v "hanyny/app/utils/view"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

// TokenClaim This is the cliam object which gets parsed from the authorization header
type TokenClaim struct {
	*jwt.StandardClaims
	models.BaseUser
}

// CreateToken ...
func CreateToken(user models.User) (string, int64) {
	expiresAt := time.Now().Add(time.Hour * 24 * 30).Unix()

	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = &TokenClaim{
		&jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		models.BaseUser{
			user.ID,
			user.Username,
			user.FirstName,
			user.LastName,
			user.Level,
			user.Avatar,
		},
	}

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	return tokenString, expiresAt
}

// Validate validator auth req
func Validate(next http.HandlerFunc, level string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		// Kiểm tra xem có tồn tại token không
		if authorizationHeader == "" {
			v.Respond(w, v.Message(false, "Bạn chưa đăng nhập."))
			return

		}
		// Kiểm tra xem token có đúng định dạng Bearer token không
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) != 2 {
			v.Respond(w, v.Message(false, "Có lỗi xảy ra, vui lòng đăng xuất và đăng nhập trở lại."))
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Phiên đăng nhập đã hết hạn.")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			v.Respond(w, v.Message(false, err.Error()))
			return
		}
		if !token.Valid {
			v.Respond(w, v.Message(false, "Phiên đăng nhập đã hết hạn."))
			return
		}

		var user models.User
		db := models.OpenDB()
		var baseUser models.BaseUser
		mapstructure.Decode(token.Claims, &baseUser)
		db.Where("id = ? AND username = ?", baseUser.ID, baseUser.Username).First(&user)

		if user.ID == 0 {
			v.Respond(w, v.Message(false, "Tài khoản đăng nhập dương như có vấn đề gì đó."))
			return
		}
		if user.Status == system.GetStatus().Block {
			v.Respond(w, v.Message(false, "Tài khoản của bạn đã bị khoá do vi phạm điều khoản sử dụng. Vui lòng liên hệ với quản trị viên để được khôi phục tài khoản."))
			return
		}

		if !system.ValidLevel(user.Level, level) {
			v.Respond(w, v.Message(false, "Bạn không có quyền truy cập."))
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)
		next(w, r)

		// v.Vars["status"] = true
		// v.Vars["user"] = user
		// v.JSON(w)
		// return
	})
}
