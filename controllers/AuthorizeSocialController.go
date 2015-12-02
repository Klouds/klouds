package controllers

import (

	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"github.com/superordinate/klouds/models"
	"io/ioutil"
	"net/url"
	"github.com/antonholmquist/jason"
	"strconv"
	"strings"
)
type AuthorizeSocialController struct {
	AppController
	*render.Render
}


var (
	stateAuth = "kloudstate39"
	facebookCallBack = "http://localhost:1337/user/auth/facebook/callback"
	endpointFBProfile = "https://graph.facebook.com/me?fields=email"
	gitHubCallBack = "http://localhost:1337/user/auth/github/callback"
	endpointGitHubProfile = "https://api.github.com/user"
	googleCallBack = "http://localhost:1337/user/auth/gplus/callback"
	endpointGplusProfile string = "https://www.googleapis.com/oauth2/v2/userinfo"

)
type userAuth struct {
	Id    interface{} `json:"id"`
	Login interface{} `json:"login"`
}

var fbConfig = &oauth2.Config{
	ClientID:     "1620451968207611",
	ClientSecret: "5bb72a6aa6bf30d24346c3e7b71c1f53",
	RedirectURL:  facebookCallBack,
	Scopes:       []string{"email"},
	Endpoint:     facebook.Endpoint,
}
var gitConfig = &oauth2.Config{
	ClientID:     "fc707d20633ddddc6a25",
	ClientSecret: "6bb334feb9a6e111d4c59dabda9c76a98947029b",
	RedirectURL:  gitHubCallBack,
	Scopes:       []string{"user:email", "user"},
	Endpoint:     github.Endpoint,
}
var gplusConfig = &oauth2.Config{
	ClientID:     "1033420881797-3kclq93hc84d1v46rtgki7vs7sdr4l0e.apps.googleusercontent.com",
	ClientSecret: "Ni6buYWXaK7MTG3EYBGnoH6l",
	RedirectURL:  googleCallBack,
	Scopes:       []string{"profile", "email", "openid"},
	Endpoint:     google.Endpoint,
}
func (au AuthorizeSocialController) HandleAuthGithubLogin(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	url := gitConfig.AuthCodeURL(stateAuth, oauth2.AccessTypeOnline)
	http.Redirect(rw, r, url, http.StatusTemporaryRedirect)

}
func (au AuthorizeSocialController) HandleAuthGplusLogin(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	url := gplusConfig.AuthCodeURL(stateAuth, oauth2.AccessTypeOnline)
	http.Redirect(rw, r, url, http.StatusTemporaryRedirect)

}
func (au AuthorizeSocialController) HandleAuthFacebookLogin(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	url := fbConfig.AuthCodeURL(stateAuth, oauth2.AccessTypeOnline)
	http.Redirect(rw, r, url, http.StatusTemporaryRedirect)

}
func (au AuthorizeSocialController) HandleAuthGitHubCallback(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {


	tok, err := gitConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		fmt.Println(err)
	}
	response, err := http.Get(endpointGitHubProfile + "?access_token=" + tok.AccessToken)

	if err != nil {
		fmt.Println(err)
	}
	bits, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	v,err := jason.NewObjectFromBytes(bits)
	id,_:= v.GetString("id")
	email,_:= v.GetString("email")
	username,_:= v.GetString("login")
	pid, _ := strconv.ParseInt(id, 10, 64)

	newUser := &models.User{
		Username: username,
		Provider: "gplus",
		ProviderID:pid,
	}
	if email != ""{
		newUser.Email = email
	}


	if CheckForExistingUsername(newUser) {
		CreateUser(newUser);
		setSession(newUser.Username, rw)
		newUser.Message = "User " + newUser.Username + " successfully created."
		redirect(rw, r, "/user/profile")
		return
	} else {
		setSession(newUser.Username, rw)
		newUser.Message = "Email: " + newUser.Email + " already taken."
		redirect(rw, r, "/user/profile")
		return
	}
	clearSession(rw)
	redirect(rw, r, "/user/login")
	return

}
func (au AuthorizeSocialController) HandleAuthGplusCallback(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	tok, err := gplusConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		fmt.Println(err)
	}
	response, err := http.Get(endpointGplusProfile + "?access_token=" + url.QueryEscape(tok.AccessToken))

	if err != nil {
		fmt.Println(err)
	}
	bits, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	v,err := jason.NewObjectFromBytes(bits)
	id,_:= v.GetString("id")
	email,_:= v.GetString("email")
	pid, _ := strconv.ParseInt(id, 10, 64)

	newUser := &models.User{
		Username: id,
		Provider: "gplus",
		ProviderID:pid,
	}
	if email != ""{
		newUser.Email = email
		e := strings.Split(email,"@")
		newUser.Username = e[0]
	}


	if CheckForExistingEmail(newUser) {
		CreateUser(newUser);
		setSession(newUser.Username, rw)
		newUser.Message = "User " + newUser.Username + " successfully created."
		redirect(rw, r, "/user/profile")
		return
	} else {
		setSession(newUser.Username, rw)
		newUser.Message = "Email: " + newUser.Email + " already taken."
		redirect(rw, r, "/user/profile")
		return
	}
	clearSession(rw)
	redirect(rw, r, "/user/login")
	return
}
func (au AuthorizeSocialController) HandleAuthFacebookCallback(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tok, err := fbConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		fmt.Println(err)
	}
	response, err := http.Get(endpointFBProfile + "&access_token=" + url.QueryEscape(tok.AccessToken))
	if err != nil {
		fmt.Println(err)
	}
	bits, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	v,err := jason.NewObjectFromBytes(bits)
	id,_:= v.GetString("id")
	email,_:= v.GetString("email")
	pid, _ := strconv.ParseInt(id, 10, 64)

	newUser := &models.User{
		Username: id,
		Provider: "facebook",
		ProviderID:pid,
	}
	if email != ""{
		newUser.Email = email
		e := strings.Split(email,"@")
		newUser.Username = e[0]

	}

	clearSession(rw)
	if CheckForExistingUsername(newUser) {
		CreateUser(newUser);
		setSession(newUser.Username, rw)
		newUser.Message = "User " + newUser.Username + " successfully created."
		redirect(rw, r, "/user/profile")
		return
	} else {
		setSession(newUser.Username, rw)
		newUser.Message = "User: " + newUser.Username + " already taken."
		redirect(rw, r, "/user/profile")
		return
	}
	redirect(rw, r, "/user/login")
	return
}

