package routers

import (
	"net/http"
	"html/template"
	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"
	"github.com/superordinate/klouds2.0/controllers"
	"fmt"
)

type Routing struct {

	Render *render.Render
	Mux *httprouter.Router

}

func (r *Routing) Init() {

	controllers.Init()
	r.Render = render.New(render.Options{Directory: "views",
		Funcs: []template.FuncMap{
        {

            "str2html": func(raw string) template.HTML {
            	fmt.Println(raw)
                return template.HTML(raw)
            },
            "add": func(x,y int) int {
                return x + y
            },
            "mod": func(x,y int) int {
                return x % y
            },
        },
    },
    })
	r.Mux = httprouter.New()

	c := &controllers.SiteNavController{Render: r.Render}
	u := &controllers.UserController{Render: r.Render}
	a := &controllers.ApplicationsController{Render: r.Render}
	
	r.Mux.GET("/", c.Index)

	//User Pages
	r.Mux.GET("/user", u.Login)
	r.Mux.GET("/user/apps", u.ApplicationList)
	r.Mux.GET("/user/apps/deleteapp/:appID", u.DeleteApplication)
	r.Mux.POST("/user/register", u.Register)
	r.Mux.GET("/user/register", u.Register)
	r.Mux.POST("/user/logout", u.Logout)
	r.Mux.POST("/user/login", u.Login)
	r.Mux.GET("/user/login", u.Login)
	r.Mux.GET("/user/profile", u.Profile)
	r.Mux.POST("/user/profile", u.Profile)

	//Application Pages
	r.Mux.GET("/apps/list", a.ApplicationList)
	r.Mux.GET("/apps/", a.ApplicationList)
	r.Mux.GET("/apps/app/:appID", a.Application)
	r.Mux.POST("/apps/app/:appID/launch", a.Launch)
	r.Mux.GET("/admin", a.AppAdmin)
	r.Mux.GET("/apps/app/:appID/delete", a.DeleteApplication)
	r.Mux.POST("/apps/app/:appID/delete", a.DeleteApplication)
	r.Mux.GET("/apps/app/:appID/edit", a.EditApplication)
	r.Mux.POST("/apps/app/:appID/edit", a.EditApplication)
	r.Mux.GET("/admin/newapp", a.CreateApplication)
	r.Mux.POST("/admin/newapp", a.CreateApplication)
	
	r.Mux.NotFound = http.FileServer(http.Dir("public"))
}