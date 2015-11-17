package routers

import (
	"github.com/superordinate/klouds/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/home", &controllers.MainController{})
	beego.Router("/", &controllers.MainController{})
	beego.Router("/appLaunch", &controllers.LaunchController{})
	beego.Router("/user/login/:back", &controllers.MainController{}, "get,post:Login")
	beego.Router("/user/logout", &controllers.MainController{}, "get:Logout")
	beego.Router("/user/register", &controllers.MainController{}, "get,post:Register")
	beego.Router("/user/profile", &controllers.MainController{}, "get,post:Profile")
	beego.Router("/user/apps", &controllers.MainController{}, "get,post:Apps")
	beego.Router("/deleteApp/:appName", &controllers.DeleteController{})
	beego.Router("/user/verify/:uuid", &controllers.MainController{}, "get:Verify")
	beego.Router("/user/remove", &controllers.MainController{}, "get,post:Remove")
	beego.Router("/notice", &controllers.MainController{}, "get:Notice")
}