package controllers

import (

	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
	"os"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
)
type AuthorizeSocialController struct {
	AppController
	*render.Render
}

func (au AuthorizeSocialController) HandleGitHubLogin(w http.ResponseWriter, r *http.Request,p httprouter.Params) {
	goth.UseProviders(
		//GITHUB_KEY=fc707d20633ddddc6a25 GITHUB_SECRET=6bb334feb9a6e111d4c59dabda9c76a98947029b go run main.go
		//FACEBOOK_KEY=1620451968207611 FACEBOOK_SECRET=5bb72a6aa6bf30d24346c3e7b71c1f53 go run main.go
		//GPLUS_KEY=1033420881797-3kclq93hc84d1v46rtgki7vs7sdr4l0e.apps.googleusercontent.com GPLUS_SECRET=Ni6buYWXaK7MTG3EYBGnoH6l go run main.go

		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:1337/user/github/callback"),
	)
	gothic.BeginAuthHandler(w,r,p)
}
func (au AuthorizeSocialController) HandleFaceBookLogin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	goth.UseProviders(
		//GITHUB_KEY=fc707d20633ddddc6a25 GITHUB_SECRET=6bb334feb9a6e111d4c59dabda9c76a98947029b go run main.go
		//FACEBOOK_KEY=1620451968207611 FACEBOOK_SECRET=5bb72a6aa6bf30d24346c3e7b71c1f53 go run main.go
		//GPLUS_KEY=1033420881797-3kclq93hc84d1v46rtgki7vs7sdr4l0e.apps.googleusercontent.com GPLUS_SECRET=Ni6buYWXaK7MTG3EYBGnoH6l go run main.go

		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost:1337/user/facebook/callback"),
	)
	gothic.BeginAuthHandler(w,r,p)

}
func (au AuthorizeSocialController) HandleGoogleLogin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	goth.UseProviders(
		//GITHUB_KEY=fc707d20633ddddc6a25 GITHUB_SECRET=6bb334feb9a6e111d4c59dabda9c76a98947029b go run main.go
		//FACEBOOK_KEY=1620451968207611 FACEBOOK_SECRET=5bb72a6aa6bf30d24346c3e7b71c1f53 go run main.go
		//GPLUS_KEY=1033420881797-3kclq93hc84d1v46rtgki7vs7sdr4l0e.apps.googleusercontent.com GPLUS_SECRET=Ni6buYWXaK7MTG3EYBGnoH6l go run main.go

		gplus.New(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"), "http://localhost:1337/user/gplus/callback"),
	)
}

func (au AuthorizeSocialController) HandleGitHubCallback(w http.ResponseWriter, r *http.Request,p httprouter.Params) {

	gothic.GetState = func(req *http.Request) string {
		// Get the state string from the query parameters.
		return req.URL.Query().Get("state")

	}
	user, err := gothic.CompleteUserAuth(w, r,p)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprintln(w, user)
	return
}
func (au AuthorizeSocialController) HandleFacebookCallback(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {


}
func (au AuthorizeSocialController) HandleGoogleCallback(w http.ResponseWriter, r *http.Request,_ httprouter.Params) {

}
