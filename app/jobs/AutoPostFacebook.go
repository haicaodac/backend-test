package jobs

// import (
// 	"fmt"
// 	"hanyny/app/library"
// 	"net/http"
// 	"os"

// 	fb "github.com/huandu/facebook"
// 	"golang.org/x/oauth2"
// 	oauth2fb "golang.org/x/oauth2/facebook"
// )

// // AutoPostFacebook ...
// func AutoPostFacebook(w http.ResponseWriter, r *http.Request) {
// 	// Get Facebook access token.
// 	conf := &oauth2.Config{
// 		ClientID:     "985944054831668",
// 		ClientSecret: "289ecf6d4c60e8c32ea6402f20bda6ec",
// 		RedirectURL:  os.Getenv("DOMAIN") + "/facebook/callback",
// 		Scopes:       []string{"manage_pages"},
// 		Endpoint:     oauth2fb.Endpoint,
// 	}
// 	http.Redirect(w, r, conf.AuthCodeURL("state"), http.StatusMovedPermanently)
// }

// // CallbackAutoPostFacebook ...
// func CallbackAutoPostFacebook(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("VAOOOO")
// 	conf := &oauth2.Config{
// 		ClientID:     "985944054831668",
// 		ClientSecret: "289ecf6d4c60e8c32ea6402f20bda6ec",
// 		RedirectURL:  os.Getenv("DOMAIN") + "/facebook/callback",
// 		Scopes:       []string{"manage_pages"},
// 		Endpoint:     oauth2fb.Endpoint,
// 	}

// 	token, err := conf.Exchange(oauth2.NoContext, "code")
// 	fmt.Println(err)
// 	if err != nil {
// 		logger := library.Logger{Type: "ERROR"}
// 		log := logger.Open()
// 		log.Println(err.Error())
// 		logger.Close()
// 		return
// 	}

// 	// Create a client to manage access token life cycle.
// 	client := conf.Client(oauth2.NoContext, token)

// 	// Use OAuth2 client with session.
// 	session := &fb.Session{
// 		Version:    "v2.4",
// 		HttpClient: client,
// 	}
// 	session.EnableAppsecretProof(true)

// 	// Use session.
// 	idPage := "114648839196301"
// 	res, err := session.Post("/"+idPage+"/feed", fb.Params{
// 		"message": "Test post to page",
// 	})
// 	fmt.Println(err)
// 	fmt.Println(res)
// }
