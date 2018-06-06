package middlewares

import (
	"jiji/models"
	"net/http"
)

type RequireUser struct {
	models.UserService
}

func (mw *RequireUser) ApplyFunc(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("authToken")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user, err := mw.UserService.GetByToken(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Get context from request
		ctx := r.Context()
		// Creates a context with user attached
		ctx = AssignUserToContext(ctx, user)
		// Overwrites request with our context attached
		r = r.WithContext(ctx)

		next(w, r)
	})
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFunc(next.ServeHTTP)
}
