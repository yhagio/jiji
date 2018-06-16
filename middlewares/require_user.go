package middlewares

import (
	"jiji/models"
	"net/http"
	"strings"
)

type RequireUser struct {
	models.UserService
}

func (mw *RequireUser) ApplyFunc(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is requesting a static asset or image
		// we will not need to lookup the current user so we skip doing that.
		path := r.URL.Path
		if strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/images/") {
			next(w, r)
			return
		}

		user := LookUpUserFromContext(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFunc(next.ServeHTTP)
}
