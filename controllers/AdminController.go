package controllers

import (
	"net/http"
	"github.com/superordinate/klouds/models"
	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"
)

type AdminController struct {
	AppController
	*render.Render
}

func (c *AdminController) AppAdmin(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

		if r.Method == "GET" {
		
		var user *models.User
		


		if getUserName(r) != "" || user != nil {
			user = GetUserByUsername(getUserName(r))
			
			if NotAdministrator(user, c, rw) {
				return
			}

					//Get application list for user
			runningapps := []models.RunningApplication{}

			GetRunningApplications(&runningapps)

			for i:=0; i<len(runningapps);i++ {
				runningapps[i].User = *user
			}

			if len(runningapps) == 0 {
				runningapps = []models.RunningApplication{models.RunningApplication{User: *user}}	
			}

			c.HTML(rw, http.StatusOK, "apps/admin/admin", runningapps)
			
		} else {

			c.HTML(rw, http.StatusOK, "user/login", nil)
		}

	}
		
}


func (c *AdminController) UserAdmin(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

		if r.Method == "GET" {
		
		var user *models.User
		


		if getUserName(r) != "" || user != nil {
			user = GetUserByUsername(getUserName(r))
			
			if NotAdministrator(user, c, rw) {
				return
			}

					//Get application list for user
			users := []models.User{}

			GetUsers(&users)

			for i:=0; i<len(users);i++ {
				if i < 1 {
					users[i] = *user
				} else {
					break;
				}
				
			}

			if len(users) == 0 {
				users = []models.User{models.User{}}	
			}

			c.HTML(rw, http.StatusOK, "apps/admin/users", users)
			
		} else {

			c.HTML(rw, http.StatusOK, "user/login", nil)
		}

	} else if r.Method == "POST" {
		var user *models.User
		


		if getUserName(r) != "" || user != nil {
			user = GetUserByUsername(getUserName(r))
			
			if NotAdministrator(user, c, rw) {
				return
			}

					//Get application list for user
			runningapps := []models.RunningApplication{}

			GetRunningApplications(&runningapps)

			for i:=0; i<len(runningapps);i++ {
				runningapps[i].User = *user
			}

			if len(runningapps) == 0 {
				runningapps = []models.RunningApplication{models.RunningApplication{User: *user}}	
			}

			c.HTML(rw, http.StatusOK, "apps/admin/users", runningapps)
			
		} else {

			c.HTML(rw, http.StatusOK, "user/login", nil)
		}
	}
		
}

func (c *AdminController) UserPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	if r.Method == "GET" {
		var user *models.User
		
		if getUserName(r) != "" || user != nil {
			user = GetUserByUsername(getUserName(r))
			userapps := GetUserByUsername(p.ByName("username"))

			if NotAdministrator(user, c, rw) {
				return
			}

			//Get application list for user
			runningapps := []models.RunningApplication{}

			GetRunningApplicationsForUser(&runningapps, userapps)

			for i:=0; i<len(runningapps);i++ {
				runningapps[i].User = *user
			}

			if len(runningapps) == 0 {
				runningapps = []models.RunningApplication{models.RunningApplication{User: *user}}	
			}

			//pass the application list to the page
			c.HTML(rw, http.StatusOK, "apps/admin/user", runningapps)
		} else {

			c.HTML(rw, http.StatusOK, "user/login", nil)
		}
		

	} else if r.Method == "POST" {
		var user *models.User
		
		if getUserName(r) != "" || user != nil {
			user := GetUserByUsername(getUserName(r))
			userapps := GetUserByUsername(p.ByName("username"))

			if NotAdministrator(user, c, rw) {
				return
			}

			//Get application list for user
			runningapps := []models.RunningApplication{}

			GetRunningApplicationsForUser(&runningapps, userapps)

			for i:=0; i<len(runningapps);i++ {
				runningapps[i].User = *user
			}

			if len(runningapps) == 0 {
				runningapps = []models.RunningApplication{models.RunningApplication{User: *user}}	
			}

			//pass the application list to the page
			c.HTML(rw, http.StatusOK, "apps/admin/user", runningapps)
		} else {

			c.HTML(rw, http.StatusOK, "user/login", nil)
		}
	}
		
}



//Checks if user is admin and redirects to login page if s/he isn't
func NotAdministrator(user *models.User, c *AdminController, rw http.ResponseWriter) bool {

	if user.Role != "admin"{
		user.Message = "You are not an administrator."
		
		c.HTML(rw, http.StatusOK, "user/login", user)
		return true
	} else {
		return false
	}
}