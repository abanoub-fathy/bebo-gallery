package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/abanoub-fathy/bebo-gallery/pkg/context"

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
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// get the user from token
		user, err := mw.Service.UserService.FindUserByRememberToken(token.Value)
		if err != nil {
			fmt.Println("error while getting user from cookie", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// get ctx from request
		ctx := r.Context()

		// create context with user
		ctx = context.WithUser(ctx, user)

		// set the new ctx to request
		r = r.WithContext(ctx)

		// call the next handler func
		next(w, r)
	}
}

func (mw *RequireUser) Apply(next http.Handler) http.Handler {
	return mw.ApplyFunc(next.ServeHTTP)
}

type UserMiddleware struct {
	Service *model.Service
}

func (userMW *UserMiddleware) UserInCtxApply(next http.Handler) http.Handler {
	return userMW.UserInCtxApplyFn(next.ServeHTTP)
}

func (userMW *UserMiddleware) UserInCtxApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if the path for getting public assets
		// we don't need to set user in ctx so we will
		// call next and return
		path := r.URL.Path
		if strings.HasPrefix(path, "/images/") || strings.HasPrefix(path, "/assets/") {
			next(w, r)
			return
		}

		// get the user token
		token, err := r.Cookie("token")
		if err != nil {
			next(w, r)
			return
		}

		// get the user from token
		user, err := userMW.Service.UserService.FindUserByRememberToken(token.Value)
		if err != nil {
			next(w, r)
			return
		}

		// create context with user
		ctx := context.WithUser(r.Context(), user)

		// set the new ctx to request
		r = r.WithContext(ctx)

		// call the next handler func
		next(w, r)
	}
}
