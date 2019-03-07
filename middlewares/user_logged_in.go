package middlewares

import (
	"jiji/models"
	"net/http"
)

type User struct {
	models.UserService
}

func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFunc(next.ServeHTTP)
}

func (mw *User) ApplyFunc(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("authToken")
		if err != nil {
			next(w, r)
			return
		}

		user, err := mw.UserService.GetByToken(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}

		ctx := r.Context()
		ctx = AssignUserToContext(ctx, user)
		r = r.WithContext(ctx)

		next(w, r)
	})
}
