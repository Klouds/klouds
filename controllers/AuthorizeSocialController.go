package controllers

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/superordinate/klouds/models"
)
type AuthorizeSocialController struct {
	AppController
	*render.Render
}
var (
	oauthConf = &oauth2.Config{
		ClientID:     "fc707d20633ddddc6a25",
		ClientSecret: "6bb334feb9a6e111d4c59dabda9c76a98947029b",
		Scopes:       []string{"user:email", "repo"},
		Endpoint:     githuboauth.Endpoint,
	}
	oauthStateString = "kloudgithublogin39"
)
func (au AuthorizeSocialController)  HandleGitHubLogin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
func (au AuthorizeSocialController) HandleGitHubCallback(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {

	state := r.FormValue("state")
	if state != oauthStateString {
		http.Redirect(w, r,  "/user/login", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Redirect(w, r,  "/user/login", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
		return
	}
	newUser := models.User{
		Username: 			*user.Login,
		Email:				*user.Email,
		FirstName: 			*user.Name,
		Surname:			*user.Name,
		Password: 			"user123456789",
		ConfirmPassword: 	"user123456789",
		Role:				"user",
	}
	if CheckForExistingUsername(&newUser) {
		if CheckForExistingEmail(&newUser) {
			CreateUser(&newUser);
			newUser.Message = "User " + newUser.Username + " successfully created."
			setSession(newUser.Username, w)
			http.Redirect(w, r, "/user/profile", http.StatusTemporaryRedirect)

		} else {
			newUser.Message = "Email: " + newUser.Email + " already taken."
		}
	} else {
		newUser.Message = "User: " + newUser.Username + " already exists."
	}
	http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)

}