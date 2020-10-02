package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// AddHSTS sets HSTS header
func AddHSTS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(strictTransportSecurity, "max-age=31536000")
		next.ServeHTTP(w, r)
	})
}

// AddCSP sets HSTS header
func AddCSP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		next.ServeHTTP(w, r)
	})
}

// JSONify sets content type to application json
func JSONify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentType, "application/json")
		next.ServeHTTP(w, r)
	})
}

// Logging middleware hanlde all incoming request
// and logs it to log/request.log file
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// log incoming request details
		now := time.Now().Format(formatDateYMD)

		if _, err := os.Stat("log"); os.IsNotExist(err) {
			if err = os.Mkdir("log", os.ModeDir); err != nil {
				log.Fatal(err)
			}
		}

		filepath := fmt.Sprintf("log/REQUEST_%s.log", now)
		file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer file.Close()
		log.SetOutput(file)
		log.Println(r)

		next.ServeHTTP(w, r)
	})
}
