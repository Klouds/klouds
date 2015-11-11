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
	Mux    *httprouter.Router

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
	r.Mux.POST("/user/register", u.Register)
	r.Mux.GET("/user/register", u.Register)
	r.Mux.POST("/user/logout", u.Logout)
	r.Mux.POST("/user/login", u.Login)
	r.Mux.GET("/user/login", u.Login)
	r.Mux.GET("/user/profile", u.Profile)
	r.Mux.POST("/user/profile", u.Profile)

	//Application Pages
	r.Mux.GET("/apps/admin", a.AppAdmin)
	r.Mux.GET("/apps/admin/newapp", a.CreateApplication)
	r.Mux.POST("/apps/admin/newapp", a.CreateApplication)

	r.Mux.NotFound = http.FileServer(http.Dir("public"))
}