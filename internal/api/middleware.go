package api

import (
	"context"
	"fajntvajb/internal/repository"
	"net/http"
	"strings"
)

func (handlers *handlers) authenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static") {
			next.ServeHTTP(w, r)
			return
		}

		session, _ := handlers.session.Get(r, "session")
		userId, ok := session.Values["userId"].(int)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		user, err := handlers.db.Users.GetUserByID(userId)
		if err != nil {
			handleWebError(w, err)
			return
		}
		if user == nil {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", user)))
	})
}

func (handlers *handlers) requireAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value("user").(*repository.User)
		if !ok {
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (handlers *handlers) requireNoAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value("user").(*repository.User)
		if ok {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
