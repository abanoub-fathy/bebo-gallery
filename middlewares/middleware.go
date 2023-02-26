package middlewares

import (
	"fmt"
	"net/http"

	"github.com/abanoub-fathy/bebo-gallery/model"
)

type RequireUser struct {
	Service *model.Service
}

func (mw *RequireUser) ApplyFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the user token
		token, err := r.Cookie("token")
		if err != nil {
			fmt.Println("error while getting cookie", err)
			http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
			return
		}

		// get the user from token
		_, err = mw.Service.UserService.FindUserByRememberToken(token.Value)
		if err != nil {
			fmt.Println("error while getting user from cookie", err)
			http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
			return
		}

		// call the next handler func
		next(w, r)
	}
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFunc(next.ServeHTTP)
}
