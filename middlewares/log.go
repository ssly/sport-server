package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

// LogMiddleware 日志中间件
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		next.ServeHTTP(w, r)

		timeElapsed := time.Since(timeStart)
		url := r.RequestURI
		time := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println(time, url, timeElapsed)
	})
}
