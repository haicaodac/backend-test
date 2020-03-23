/*
 * Created by Dac Hai on 23/10/2018
 */

package models

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // Library connect mysql gorm
)

// Model ...
type Model struct {
	ID        uint       `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

var db *gorm.DB

// Init ...
func Init() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	fmt.Println(dbURI)
	conn, err := gorm.Open("mysql", dbURI)
	if err != nil {
		fmt.Println("-------------------Database ERROR--------------------------")
		fmt.Println(err)
	}
	db = conn

	// db.DropTable(
	// 	&User{},
	// 	&Post{},
	// 	&CategoryPost{},
	// 	&Category{},
	// 	&Comment{},
	// 	&Attachment{},
	// 	&Media{},
	// 	&Review{},
	// 	&Notification{},
	// 	&History{},
	// 	&UserSession{},
	// )

	db.AutoMigrate(
		&User{},
		&Post{},
		&CategoryPost{},
		&Category{},
		&Comment{},
		&Attachment{},
		&Media{},
		&Review{},
		&Notification{},
		&History{},
		&UserSession{},
		&Email{},
		&Botchat{},
	)

	// ALTER TABLE posts ADD FULLTEXT (`title`, `description`, `content`)
	// SELECT * FROM posts WHERE MATCH (`title`,`description`,`content`) AGAINST ('vẽ quả chuối' IN NATURAL LANGUAGE MODE)
}

// OpenDB created connect database
func OpenDB() *gorm.DB {
	return db
}
