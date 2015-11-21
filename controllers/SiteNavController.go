package controllers

import (
	"net/http"
	"github.com/superordinate/klouds/models"
	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"

)

type SiteNavController struct {
	AppController
	*render.Render
}

func (c *SiteNavController) ServeFull(rw http.ResponseWriter, r *http.Request, p httprouter.Params, pageName string) {

	
}

func (c *SiteNavController) Index(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var user *models.User

		if getUserName(r) != "" {
			user = GetUserByUsername(getUserName(r))
			c.HTML(rw, http.StatusOK, "index", user)
			
		} else {
			c.HTML(rw, http.StatusOK, "index", nil)
		}
}


func (c *SiteNavController) About(rw http.ResponseWriter,r *http.Request,p httprouter.Params){
	c.HTML(rw,http.StatusOK,"about",nil)

}
func (c *SiteNavController) Customer(rw http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	c.HTML(rw,http.StatusOK,"customer",nil)

}
func (c *SiteNavController) Pricing(rw http.ResponseWriter,r *http.Request,p httprouter.Params){
	c.HTML(rw,http.StatusOK,"pricing",nil)

}
func (c *SiteNavController) Contact(rw http.ResponseWriter,r *http.Request,p httprouter.Params)  {

	c.HTML(rw,http.StatusOK,"contact",nil)

}
func  (c *SiteNavController) Supporting(rw http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	c.HTML(rw,http.StatusOK,"supporting",nil)

}
func (c *SiteNavController) Blog(rw http.ResponseWriter,r *http.Request,p httprouter.Params){
	c.HTML(rw,http.StatusOK,"blog",nil)
}