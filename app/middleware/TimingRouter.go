package middleware

import (
	"fmt"
	library "hanyny/app/library"
	"net/http"
	"os"
	"strings"
	"time"
)

// TimingRouter ...
func TimingRouter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/") {
			defer elapsed(r)()
		}

		next.ServeHTTP(w, r)
		return
	})
}

func elapsed(r *http.Request) func() {
	start := time.Now()
	return func() {
		elapsed := float64(time.Now().Sub(start)) / float64(time.Millisecond)

		msg := fmt.Sprintf("ROUTE: %s - %fms", r.URL.Path, elapsed)
		if time.Duration(elapsed)*time.Millisecond > time.Second*1 {
			logger := library.Logger{Type: "TIMMER"}
			log := logger.Open()
			log.Println(msg)
			logger.Close()
		}
		if os.Getenv("MODE") == "development" {
			fmt.Printf("%s\n", msg)
		}
	}
}
