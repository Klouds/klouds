package controllers

import (

	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
	"os"
	"github.com/smagic39/goth"
	"github.com/smagic39/goth/gothic"
	"github.com/smagic39/goth/providers/facebook"
	"github.com/smagic39/goth/providers/github"
	"github.com/smagic39/goth/providers/gplus"
	"github.com/superordinate/klouds/models"
	"github.com/superordinate/klouds/routers"
)
type AuthorizeSocialController struct {
	AppController
	*render.Render
}
var (
	facebookCallBack = "http://localhost:1337/user/auth/facebook/callback"
	gitHubCallBack = "http://localhost:1337/user/auth/github/callback"
	googleCallBack = "http://localhost:1337/user/auth/gplus/callback"
)

func (au AuthorizeSocialController) HandleAuthLogin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	goth.UseProviders(
		//GITHUB_KEY=fc707d20633ddddc6a25 GITHUB_SECRET=6bb334feb9a6e111d4c59dabda9c76a98947029b go run main.go
		//FACEBOOK_KEY=1620451968207611 FACEBOOK_SECRET=5bb72a6aa6bf30d24346c3e7b71c1f53 go run main.go
		//GPLUS_KEY=1033420881797-3kclq93hc84d1v46rtgki7vs7sdr4l0e.apps.googleusercontent.com GPLUS_SECRET=Ni6buYWXaK7MTG3EYBGnoH6l go run main.go
		facebook.New(os.Getenv("f"), os.Getenv("FACEBOOK_SECRET"), facebookCallBack),
		gplus.New(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"), googleCallBack),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), gitHubCallBack),
	)
	gothic.BeginAuthHandler(w, r, p)
}

func (au AuthorizeSocialController) HandleAuthCallback(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userResponse, err := gothic.CompleteUserAuth(rw, r, p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	newUser := models.User{
		Username:            userResponse.Email,
		Email:                userResponse.Email,
		FirstName:            userResponse.Name,
		Surname:            userResponse.Name,
		Password:            "User123456",
		Role:                "user",
	}
	if CheckForExistingEmail(&newUser) {
		CreateUser(&newUser);
		newUser.Message = "User " + newUser.Username + " successfully created."
		http.HandlerFunc("user/profile")

		return
	} else {
		setSession(newUser.Username, rw)
		newUser.Message = "Email: " + newUser.Email + " already taken."
		http.HandlerFunc("user/profile")
		return
	}
	clearSession(rw)
	http.HandlerFunc("user/login")
	return

}
