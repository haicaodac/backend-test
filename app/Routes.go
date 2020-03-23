package routes

import (
	"encoding/json"
	controllers "hanyny/app/controllers"
	middleware "hanyny/app/middleware"
	auth "hanyny/app/utils/auth"
	system "hanyny/app/utils/system"
	"net/http"

	botchat "hanyny/app/botchat"

	"github.com/gorilla/mux"
)

// NewRouter function configures a new router to the API
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Use(middleware.TimingRouter)

	/*==================================== Client =============================================================*/
	// User
	router.HandleFunc("/api/user/signin", controllers.SignIn).Methods("POST")
	router.HandleFunc("/api/user/signin/facebook", controllers.SignInFacebook).Methods("POST")
	router.HandleFunc("/api/user/signup", controllers.SignUp).Methods("POST")
	router.HandleFunc("/api/user/forgot-password", controllers.UserForgotPassword).Methods("POST")
	router.HandleFunc("/api/user/check-token-forgot-password/{token}", controllers.UserCheckForgotPassword).Methods("GET")
	router.HandleFunc("/api/user/new-password", auth.Validate(controllers.UserNewPassword, system.GetLevel().User)).Methods("POST")
	router.HandleFunc("/api/user/subscribe/email", controllers.UserSubscribeEmail).Methods("GET")

	// Media
	// router.HandleFunc("/api/media", auth.Validate(controllers.MediaGet, system.GetLevel().User)).Methods("GET")

	// Media
	router.HandleFunc("/api/media", auth.Validate(controllers.MediaAdd, system.GetLevel().Editor)).Methods("POST")
	router.HandleFunc("/api/media", auth.Validate(controllers.MediaGet, system.GetLevel().Editor)).Methods("GET")
	router.HandleFunc("/api/media/thumbnail", auth.Validate(controllers.AdminMediaCreateThumbnail, system.GetLevel().Editor)).Methods("POST")

	// Post
	router.HandleFunc("/api/post", controllers.PostGet).Methods("GET")
	router.HandleFunc("/api/post/sitemap", controllers.PostGetSitemap).Methods("GET")
	router.HandleFunc("/api/post/related/{slug}", controllers.PostGetRelated).Methods("GET")
	router.HandleFunc("/api/post/search", controllers.PostSearch).Methods("GET")
	router.HandleFunc("/api/post/search/speed", controllers.PostSearchSpeed).Methods("GET")
	router.HandleFunc("/api/post/{slug}", controllers.PostGetOne).Methods("GET")
	router.HandleFunc("/api/post", auth.Validate(controllers.PostAdd, system.GetLevel().User)).Methods("POST")
	// router.HandleFunc("/api/post/{slug}", controllers.PostGetOne).Methods("GET")

	// Comment
	router.HandleFunc("/api/comment", auth.Validate(controllers.CommentAdd, system.GetLevel().User)).Methods("POST")

	// Attachment
	router.HandleFunc("/api/attachment", auth.Validate(controllers.AttachmentGet, system.GetLevel().User)).Methods("GET")
	router.HandleFunc("/api/attachment/{id}", controllers.AttachmentGetOne).Methods("GET")
	router.HandleFunc("/api/attachment", auth.Validate(controllers.AttachmentAdd, system.GetLevel().User)).Methods("POST")

	// Category
	router.HandleFunc("/api/category", controllers.CategoryGet).Methods("GET")
	router.HandleFunc("/api/category/sitemap", controllers.CategoryGetSitemap).Methods("GET")
	router.HandleFunc("/api/category/{slug}", controllers.CategoryGetOne).Methods("GET")

	/*==================================== ADMIN =============================================================*/
	// User
	router.HandleFunc("/api/admin/user", auth.Validate(controllers.AdminUserGet, system.GetLevel().Editor)).Methods("GET")
	router.HandleFunc("/api/admin/user/total", auth.Validate(controllers.AdminUserTotal, system.GetLevel().Editor)).Methods("GET")
	router.HandleFunc("/api/admin/user/search", auth.Validate(controllers.AdminUserSearch, system.GetLevel().Editor)).Methods("GET")
	router.HandleFunc("/api/admin/user/{id}", auth.Validate(controllers.AdminUserGetOne, system.GetLevel().Editor)).Methods("GET")
	router.HandleFunc("/api/admin/user/{id}", auth.Validate(controllers.AdminUserUpdate, system.GetLevel().Admin)).Methods("PUT")

	// Post
	router.HandleFunc("/api/admin/post", auth.Validate(controllers.AdminPostGet, system.GetLevel().Editor)).Methods("GET")
	router.HandleFunc("/api/admin/post/private", auth.Validate(controllers.AdminPostGetPrivate, system.GetLevel().Editor)).Methods("GET")
	router.HandleFunc("/api/admin/post/search", auth.Validate(controllers.AdminPostSearch, system.GetLevel().Editor)).Methods("GET")
	router.HandleFunc("/api/admin/post/total", auth.Validate(controllers.AdminPostTotal, system.GetLevel().Editor)).Methods("GET")
	router.HandleFunc("/api/admin/post/{id}", auth.Validate(controllers.AdminPostGetOne, system.GetLevel().Editor)).Methods("GET")
	router.HandleFunc("/api/admin/post/{id}", auth.Validate(controllers.AdminPostUpdate, system.GetLevel().Editor)).Methods("PUT")
	router.HandleFunc("/api/admin/post", auth.Validate(controllers.AdminPostAdd, system.GetLevel().Editor)).Methods("POST")
	router.HandleFunc("/api/admin/post/{id}", auth.Validate(controllers.AdminPostDelete, system.GetLevel().Admin)).Methods("DELETE")

	// Category
	router.HandleFunc("/api/admin/category", auth.Validate(controllers.AdminCategoryGet, system.GetLevel().Editor)).Methods("GET")
	router.HandleFunc("/api/admin/category/{id}", auth.Validate(controllers.AdminCategoryGetOne, system.GetLevel().Editor)).Methods("GET")
	router.HandleFunc("/api/admin/category", auth.Validate(controllers.AdminCategoryAdd, system.GetLevel().Admin)).Methods("POST")
	router.HandleFunc("/api/admin/category/{id}", auth.Validate(controllers.AdminCategoryUpdate, system.GetLevel().Admin)).Methods("PUT")
	router.HandleFunc("/api/admin/category/{slug}", auth.Validate(controllers.AdminCategoryDelete, system.GetLevel().Admin)).Methods("DELETE")
	router.HandleFunc("/api/admin/category/total", auth.Validate(controllers.AdminCategoryTotal, system.GetLevel().Editor)).Methods("GET")

	// Attachment
	router.HandleFunc("/api/admin/attachment", auth.Validate(controllers.AdminAttachmentGet, system.GetLevel().Editor)).Methods("GET")

	// Media
	router.HandleFunc("/api/admin/media", auth.Validate(controllers.MediaAdd, system.GetLevel().Editor)).Methods("POST")
	router.HandleFunc("/api/admin/media", auth.Validate(controllers.MediaGet, system.GetLevel().Editor)).Methods("GET")
	router.HandleFunc("/api/media/thumbnail", auth.Validate(controllers.AdminMediaCreateThumbnail, system.GetLevel().Editor)).Methods("POST")

	// Email
	router.HandleFunc("/api/admin/email", auth.Validate(controllers.AdminEmailGet, system.GetLevel().Editor)).Methods("GET")
	router.HandleFunc("/api/admin/email/total", auth.Validate(controllers.AdminEmailTotal, system.GetLevel().Editor)).Methods("GET")

	// API thumbnail
	router.HandleFunc("/uploads/thumbnail/{link}", controllers.PostThumbnail).Methods("GET")
	router.HandleFunc("/uploads/resize/thumbnail/{link}", controllers.PostThumbnail).Methods("GET")

	// Nyny
	var nyny botchat.Nyny
	nyny.New()
	router.HandleFunc("/api/nyny", nyny.QuestionAnswer).Methods("POST")

	/*==================================== Static =============================================================*/
	// Router uploads
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("public/uploads/"))))
	router.PathPrefix("/public/thumbnail/store/").Handler(http.StripPrefix("/public/thumbnail/store/", http.FileServer(http.Dir("public/thumbnail/store/"))))

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "404 not found",
			"status":  false,
		})
	})

	//// CHECK LIST
	// err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	// 	pathTemplate, err := route.GetPathTemplate()
	// 	if err == nil {
	// 		fmt.Println("ROUTE:", pathTemplate)
	// 	}
	// 	pathRegexp, err := route.GetPathRegexp()
	// 	if err == nil {
	// 		fmt.Println("Path regexp:", pathRegexp)
	// 	}
	// 	queriesTemplates, err := route.GetQueriesTemplates()
	// 	if err == nil {
	// 		fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
	// 	}
	// 	queriesRegexps, err := route.GetQueriesRegexp()
	// 	if err == nil {
	// 		fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
	// 	}
	// 	methods, err := route.GetMethods()
	// 	if err == nil {
	// 		fmt.Println("Methods:", strings.Join(methods, ","))
	// 	}
	// 	fmt.Println()
	// 	return nil
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// }

	return router
}
