package api

import (
	"context"
	"net/http"
)

func (handlers *handlers) authenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
