package controllers

import (
	"net/http"
	"github.com/superordinate/klouds2.0/models"
	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"
	"strings"
	//"fmt"
)

type UserController struct {
	AppController
	*render.Render
}


//Registration page
func (c *UserController) Register(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	


	if r.Method == "GET" {
		
		var user *models.User

		if getUserName(r) != "" {
			user = GetUserByUsername(getUserName(r))
			c.HTML(rw, http.StatusOK, "user/profile", user)
			
		} else {
			c.HTML(rw, http.StatusOK, "user/register", nil)
		}

	} else if r.Method == "POST" {
		r.ParseForm();
		username := strings.ToLower(r.FormValue("username"))
		email := r.FormValue("email")
		firstname := r.FormValue("firstname")
		surname := r.FormValue("lastname")
		password := r.FormValue("password")
		confirmpassword := r.FormValue("confirmpassword")

		newUser := models.User{
			Username: 			username,
			Email:				email,
			FirstName: 			firstname,
			Surname:			surname,
			Password: 			password,
			ConfirmPassword: 	confirmpassword,
			Role:				"user"}


		newUser.ValidateRegister()

		if (newUser.Message != "") {

			c.HTML(rw, http.StatusOK, "user/register", newUser)
			return;
		}

			
		if CheckForExistingUsername(&newUser) {
			if CheckForExistingEmail(&newUser) {
				CreateUser(&newUser);
				newUser.Message = "User " + newUser.Username + " successfully created."
			} else {
				newUser.Message = "Email: " + newUser.Email + " already taken."
			}
		} else {
			newUser.Message = "User: " + newUser.Username + " already exists."
		}

		c.HTML(rw, http.StatusOK, "user/register", newUser)

	}

	
}

//Login page
func (c *UserController) Login(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if r.Method == "GET" {
		// if logged in, go to profile / else login page
		RedirectToLogin(c,rw,r,p)

	} else if r.Method == "POST" {
		r.ParseForm();

		username := strings.ToLower(r.FormValue("username"))
		password := r.FormValue("password")

		newUser := models.User {
			Username: 	username,
			Password: 	password}

		newUser.ValidateLogin()

		if (newUser.Message != "") {
			newUser.Username = ""
			
			c.HTML(rw, http.StatusOK, "user/login", newUser)
			return;
		}

		if !CheckForExistingUsername(&newUser) {
			//User Exists
			if (CheckForMatchingPassword(&newUser)) {
				//Passwords Match
				//Open Session and forward back to front page
				newUser.Message = "Passwords match."
				setSession(newUser.Username, rw)
				c.HTML(rw, http.StatusOK, "user/loggedIn", newUser)
				return;

			} else {
				newUser.Message = "Password doesn't match record for " + newUser.Username
			}

		} else {
			//User Doesn't exist
			newUser.Message = "User " + newUser.Username + " doesn't exist."
		}

		c.HTML(rw, http.StatusOK, "user/login", newUser)
	}
}

//profile page
func (c *UserController) Profile(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if r.Method == "GET" {
		RedirectToLogin(c,rw,r,p)

	} else if r.Method == "POST" {
		// if logged in, go to profile / else login page
		var user *models.User

		if getUserName(r) != "" {
			user = GetUserByUsername(getUserName(r))
			
		} else {
			c.HTML(rw, http.StatusOK, "user/login", nil)
			return
		}

		r.ParseForm();

		currentPassword := r.FormValue("currentpassword")
		password := r.FormValue("password")
		confirmpassword := r.FormValue("confirmpassword")

		user.Password = currentPassword

		//Check to make sure password matches record
		if (CheckForMatchingPassword(user)) {
			//Passwords Match
			//Validate new password, make sure they match, update user,display success message
			
			user.Password = password
			user.ConfirmPassword = confirmpassword

			user.ValidateNewPassword()

			if user.Message != "" {
				c.HTML(rw, http.StatusOK, "user/profile", user)
				return
			} else {
				user.Message = "Updated user"
				UpdateUser(user)
				c.HTML(rw, http.StatusOK, "user/profile", user)
				return
			}
			

		} else {
			user.Message = "Password doesn't match record for " + user.Username
		}

		c.HTML(rw, http.StatusOK, "user/profile", user)
	}
}

func (c *UserController) Logout(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if r.Method == "POST" {
		clearSession(rw)
		
		var user *models.User

		if getUserName(r) != "" {
			user = GetUserByUsername(getUserName(r))
			c.HTML(rw, http.StatusOK, "user/logout", user)
			
		} else {
			c.HTML(rw, http.StatusOK, "user/logout", nil)
		}

		

	}
}

func (c *UserController) ApplicationList(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		
	if r.Method == "GET"{
		var user *models.User

		if getUserName(r) == "" {
			c.HTML(rw, http.StatusOK, "user/login", nil)
			return
		} 
		
		user = GetUserByUsername(getUserName(r))

		//Get application list for user
		runningapps := []models.RunningApplication{}

		GetRunningApplicationsForUser(&runningapps, user)

		for i:=0; i<len(runningapps);i++ {
			runningapps[i].User = *user
		}

		if len(runningapps) == 0 {
			runningapps = []models.RunningApplication{models.RunningApplication{User: *user}}	
		}
		//pass the application list to the page
		c.HTML(rw, http.StatusOK, "user/apps", runningapps)
	}
}

// if logged in, go to profile / else login page
func RedirectToLogin(c *UserController, rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

		var user *models.User

		if getUserName(r) != "" {
			user = GetUserByUsername(getUserName(r))
			c.HTML(rw, http.StatusOK, "user/profile", user)
			
		} else {
			c.HTML(rw, http.StatusOK, "user/login", nil)
		}
}

func (c *UserController) DeleteApplication(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		
	if r.Method == "GET"{
		var user *models.User

		if getUserName(r) == "" {
			c.HTML(rw, http.StatusOK, "user/login", nil)
			return
		} 
		
		user = GetUserByUsername(getUserName(r))

		//Get application list for user
		runningapp := GetRunningApplicationByName(p.ByName("appID"))

		if DeleteRunningApplication(user.Username, runningapp.Name) {

			runningapp.User = *user
			

			if runningapp.Name == "" {
				runningapp = &models.RunningApplication{User: *user}	
			}
			//pass the application list to the page
			c.HTML(rw, http.StatusOK, "user/deleteapp", runningapp)
		} else {
			runningapp = &models.RunningApplication{User: *user}
			runningapp.Message = "You do not own this application."

			c.HTML(rw, http.StatusOK, "user/login", runningapp)
		}
	}
}

