package crawl

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	library "hanyny/app/library"
	"hanyny/app/models"

	"github.com/gosimple/slug"
)

// Page ...
type Page struct {
	Total int
	Posts []*Post
}

// Post ...
type Post struct {
	Title string
	Link  string
	Desc  string
	File  string
}

// Push ...
func Push() {

	pushOne("lop-lon", 8)
}

func pushOne(keyword string, number int) {
	for i := 1; i <= number; i++ {
		fileJSON := fmt.Sprintf("crawl/"+keyword+"/page-%d.json", i)
		pageJSON, err := ioutil.ReadFile(fileJSON)
		if err != nil {
			panic(err.Error())
		}
		var page Page
		err = json.Unmarshal(pageJSON, &page)
		if err != nil {
			panic(err.Error())
		}

		db := models.OpenDB()

		for _, postJSON := range page.Posts {

			tx := db.Begin()
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()

			fmt.Println("New: ", postJSON.Title)
			var post models.Post
			post.Title = postJSON.Title
			post.Slug = slug.Make(post.Title)
			post.Content = postJSON.Desc
			post.UserID = 1
			post.Status = "private"

			var checkPost models.Post
			db.Where("title = ? OR slug = ?", post.Title, post.Slug).Last(&checkPost)

			if checkPost.ID == 0 {
				fi, err := os.Stat("crawl/" + postJSON.File)
				if err != nil {
					log.Fatal(err)
				}
				// get the size
				size := fi.Size()

				oldLocation := "crawl/" + postJSON.File
				fileName := strings.Replace(postJSON.File, keyword+"/posts/", "", -1)
				fileName = library.RandomString(10) + "-" + fileName
				newLocation := "public/uploads/files/" + fileName

				var attachment models.Attachment
				attachment.Name = post.Title
				attachment.Type = "application/pdf"
				attachment.Size = size
				attachment.URL = fileName
				attachment.UserID = 1
				err = tx.Create(&attachment).Error
				if err == nil {
					err = copyFile(oldLocation, newLocation)
					if err != nil {
						tx.Rollback()
						log.Fatal(err)
					}

					post.AttachmentID = attachment.ID
					err = tx.Create(&post).Error
					if err != nil {
						tx.Rollback()
						log.Fatal(err)
					}
				} else {
					tx.Rollback()
					log.Fatal(err)
				}

				tx.Commit()
			}
		}
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// func migrateAttachment() {
// 	user, err := ioutil.ReadFile("app/library/migrateOldVersion/media.json")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	var arrayData []map[string]interface{}
// 	err = json.Unmarshal(user, &arrayData)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	db := models.OpenDB()

// 	for _, data := range arrayData {
// 		var attachment models.Attachment
// 		// Name          string `json:"name,omitempty" gorm:"type:varchar(4000);not null" valid:"required~Tên tài liệu không thể trống.,runelength(1|4000)~Tên tài liệu phải có độ dài từ 1 đến 4000 ký tự."`
// 		// Type          string `json:"type,omitempty" gorm:"type:varchar(200);not null"`
// 		// Size          int64  `json:"size,omitempty" gorm:"not null"`
// 		// URL           string `json:"url,omitempty" gorm:"type:varchar(4000);not null" valid:"required~Đường dẫn tài liệu không thể trống.,runelength(1|4000)~Đường dẫn tài liệu phải có độ dài từ 1 đến 4000 ký tự."`
// 		// View          uint   `json:"view,omitempty" gorm:"not null;default:0"`
// 		// CountDownload uint   `json:"count_download,omitempty" gorm:"not null;default:0"`
// 		// UserID        uint   `json:"user_id,omitempty" gorm:"not null"`
// 		// Status        string `json:"status,omitempty" gorm:"type:varchar(20);not null;default:'active'"`

// 		attachment.OldID = data["_id"].(map[string]interface{})["$oid"].(string)

// 		if data["fieldname"] != "" && data["fieldname"] != nil {
// 			fileN := data["fieldname"].(string)
// 			if fileN == "doc" {
// 				if data["originalname"] != "" && data["originalname"] != nil {
// 					attachment.Name = data["originalname"].(string)
// 				}
// 				if data["type"] != "" && data["type"] != nil {
// 					attachment.Type = data["type"].(string)
// 				}
// 				if data["size"] != "" && data["size"] != nil {
// 					i, _ := strconv.Atoi(data["size"].(string))
// 					attachment.Size = int64(i)
// 				}
// 				if data["link"] != "" && data["link"] != nil {
// 					attachment.URL = data["link"].(string)
// 				}
// 				if data["download"] != "" && data["download"] != nil {
// 					attachment.CountDownload = int(data["download"].(float64))
// 				}

// 				db.Where("old_id = ?", attachment.OldID).Find(&attachment)
// 				if err := db.Save(&attachment).Error; err != nil {
// 					fmt.Println(attachment.OldID)
// 					panic(err.Error())
// 				}

// 				if data["userId"] != "" && data["userId"] != nil {
// 					userID := data["userId"].(map[string]interface{})["$oid"].(string)

// 					var checkUser models.User
// 					db.Where("old_id = ?", userID).Last(&checkUser)
// 					fmt.Println(checkUser.ID)
// 					attachment.UserID = checkUser.ID
// 					db.Save(&attachment)
// 				}
// 			}
// 		}
// 	}
// }
