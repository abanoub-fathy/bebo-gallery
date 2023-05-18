package controllers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/abanoub-fathy/bebo-gallery/model"
	ctx "github.com/abanoub-fathy/bebo-gallery/pkg/context"
	"github.com/gorilla/csrf"
	"golang.org/x/oauth2"
)

type oAuthController struct {
	OAuthConfig *oauth2.Config
	Service     model.OAuthService
}

// NewOAuthController is the constructor for making
func NewOAuthController(oAuthConfig *oauth2.Config, service model.OAuthService) *oAuthController {
	return &oAuthController{
		OAuthConfig: oAuthConfig,
		Service:     service,
	}
}

func (c *oAuthController) Connect(w http.ResponseWriter, r *http.Request) {
	state := csrf.Token(r)

	// create a cookie with the state
	coockie := &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
	}
	// setting the cookie
	http.SetCookie(w, coockie)

	// generate and redirect to authURL
	url := c.OAuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (c *oAuthController) Callback(w http.ResponseWriter, r *http.Request) {
	// get the query params values
	query := r.URL.Query()
	code := query.Get("code")
	state := query.Get("state")

	// get state from request cookie
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if cookie == nil || cookie.Value != state {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	// expire the state cookie
	cookie.Value = ""
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)

	// exchange the code
	token, err := c.OAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get the user from request's context
	user := ctx.UserValue(r.Context())

	// check if the user has token
	existToken, err := c.Service.Find(user.ID.String(), model.OAuthDropboxProvider)
	if err == nil {
		// delete the existToken
		c.Service.Delete(existToken.ID.String())
	} else if err != nil && err != model.ErrNotFound {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	oAuth := &model.OAuth{
		UserID:   user.ID,
		Token:    *token,
		Provider: model.OAuthDropboxProvider,
	}
	if err := c.Service.Create(oAuth); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%+v", token)

	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(token)
}

func (c *oAuthController) Testfunc(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")

	// get the user from ctx
	user := ctx.UserValue(r.Context())

	// get the user OAuth
	oauth, err := c.Service.Find(user.ID.String(), model.OAuthDropboxProvider)
	if err != nil {
		panic(err)
	}

	// create http client
	client := c.OAuthConfig.Client(context.Background(), &oauth.Token)

	req, err := http.NewRequest(http.MethodPost, "https://api.dropboxapi.com/2/files/list_folder", strings.NewReader(fmt.Sprintf(`{"path": "%v"}`, path)))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	w.Header().Add("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}
