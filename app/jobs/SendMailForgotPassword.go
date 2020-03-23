package jobs

import (
	"bytes"
	"hanyny/app/library"
	"hanyny/app/models"
	"html/template"
	"os"
)

func parseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ForgotPassword ...
func ForgotPassword(pathTemp string, user models.User, token string) {
	url := os.Getenv("DOMAIN") + "tao-mat-khau-moi?token=" + token
	templateData := struct {
		URL string
	}{
		URL: url,
	}

	strBody, err := parseTemplate("app/template-email/"+pathTemp, templateData)
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
	}
	err = library.SendMailGun("[ Hanyny ] Lấy lại mật khẩu", strBody, user.Email)
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
	}
}

// NotifiNewComment ...
func NotifiNewComment(post models.Post, comment models.Comment, user models.User) {
	templateData := struct {
		Title    string
		Content  string
		URL      string
		Level    string
		FullName string
		Time     string
	}{}

	templateData.Title = post.Title
	templateData.FullName = user.LastName + " " + user.FirstName
	templateData.Content = comment.Content
	templateData.URL = post.Slug
	templateData.Level = comment.User.Level
	templateData.Time = comment.CreatedAt.Format("15:04:05 02-01-2006")

	strBody, err := parseTemplate("app/template-email/comment.html", templateData)
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
	}
	err = library.SendMailGun("[ Hanyny ] Bình luận mới", strBody, user.Email)
	if err != nil {
		logger := library.Logger{Type: "ERROR"}
		log := logger.Open()
		log.Println(err.Error())
		logger.Close()
	}
}
