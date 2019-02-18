package middleware

import (
	"log"
	"net/http"
	"time"
)

type Session struct {
	tokenUsers map[string]string
}

func NewSession() *Session {
	return &Session{
		tokenUsers: make(map[string]string),
	}
}

// Initialize it somewhere
func (s *Session) Populate() {
	s.tokenUsers["00000000"] = "user0"
	s.tokenUsers["aaaaaaaa"] = "userA"
	s.tokenUsers["05f717e5"] = "randomUser"
	s.tokenUsers["deadbeef"] = "user0"
}

func (s *Session) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")

		if user, found := s.tokenUsers[token]; found {
			log.Printf("Authenticated user %s\n", user)
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func (s *Session) SetCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "hello",
		Value:    "hello",
		HttpOnly: true,
		Expires:  time.Now().Add(1 * time.Hour),
	})
}
