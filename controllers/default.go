package controllers

import (
	"fmt"
	"io/ioutil"
	"github.com/astaxie/beego"
	"net/http"
	"bytes"
	"strconv"
	"strings"
	"github.com/elgs/gostrgen"
	"time"
	"os"
)

type MainController struct {
	beego.Controller
}

type LaunchController struct {
	beego.Controller
}

type DeleteController struct {
	beego.Controller
}

func (this *MainController) activeContent(view string) {
	this.Layout = "basic-layout.tpl"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "header.tpl"
	this.LayoutSections["Sidebar"] = "sidebar.tpl"
	this.LayoutSections["Footer"] = "footer.tpl"
	this.LayoutSections["Menu"] = "menu.tpl"
	this.TplNames = view + ".tpl"

	sess := this.GetSession("acme")
	if sess != nil {
		this.Data["InSession"] = 1 // for login bar in header.tpl
		m := sess.(map[string]interface{})
		this.Data["First"] = m["first"]
	}
}

func (this *MainController) Get() {
	this.activeContent("index")

}

func (this *MainController) Notice() {
	this.activeContent("notice")

	flash := beego.ReadFromRequest(&this.Controller)
	if n, ok := flash.Data["notice"]; ok {
		this.Data["notice"] = n
	}
}

func (this *LaunchController) activeContent(view string) {
	this.Layout = "basic-layout.tpl"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "header.tpl"
	this.LayoutSections["Sidebar"] = "sidebar.tpl"
	this.LayoutSections["Footer"] = "footer.tpl"
	fileselect := "<select name='filelist' form='form1'>"

	files, err := ioutil.ReadDir("./restcalls")
	if err != nil {

		fmt.Println(err)
	}
    for _, f := range files {
    		fileselect= fileselect + "<option>" + f.Name() + "</option>"
    }

    fileselect = fileselect + "</select>"



	this.TplNames = view + ".tpl"
    this.Data["FileSelect"] = fileselect
	sess := this.GetSession("acme")
	if sess != nil {
		this.Data["InSession"] = 1 // for login bar in header.tpl
		m := sess.(map[string]interface{})
		this.Data["First"] = m["first"]
	}
}

func (this *LaunchController) Get() {
	this.activeContent("appLaunch")
	this.Data["Launching"] = ""
	//******** This page requires login
	sess := this.GetSession("acme")
	if sess == nil {
		this.Redirect("/user/login/home", 302)
		return
	}
	m := sess.(map[string]interface{})
	fmt.Println("username is", m["username"])
	fmt.Println("logged in at", m["timestamp"])
}

func (this *LaunchController) Notice() {
	this.activeContent("appLaunch")

	flash := beego.ReadFromRequest(&this.Controller)
	if n, ok := flash.Data["notice"]; ok {
		this.Data["notice"] = n
	}
}

func (this *LaunchController) Post() {

	sess := this.GetSession("acme")
	if sess == nil {
		this.Redirect("/user/login/home", 302)
		return
	}
	m := sess.(map[string]interface{})
	username, _ := m["username"].(string)
	username = strings.ToLower(username)

	this.activeContent("appLaunch")
	fmt.Println("app launch started")
	numbertolaunch, _ := strconv.Atoi(this.GetString("launchNumber"))

	file := this.GetString("filelist")

	for i:=0; i<numbertolaunch;i++ {

		//Read File for launching
		podfile, err := ioutil.ReadFile("./restcalls/" + file)

		if err != nil {
			panic(err)
		}

		chararactersToGenerate := 8
		set := gostrgen.Lower | gostrgen.Digit
		include := "" // optionally include some additional letters
		exclude := "Ol"     //exclude big 'O' and small 'l' to avoid confusion with zero and one.

		randomstring, err := gostrgen.RandGen(chararactersToGenerate, set, include, exclude)
		if err != nil {
			fmt.Println(err)
		}

		//append username into JSON FILE

		var newstring string =""
		var usernamestring = username + "-" + randomstring
		
		JSONSPLIT := strings.Split(string(podfile), "#USERNAME")
		newstring = newstring + JSONSPLIT[0] + usernamestring + JSONSPLIT[1] + usernamestring + JSONSPLIT[2]

		fmt.Println(newstring)

		//Create the request
		url := "http://" + os.Getenv("MARATHON_ENDPOINT") + "/v2/apps/"
		bytestring := []byte(newstring)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytestring))

		//Make the request
		res, err := http.DefaultClient.Do(req)

		if err != nil {
	    	panic(err) //Something is wrong while sending request
	 	}

		if res.StatusCode != 201 {
			fmt.Printf("Success expected: %d", res.StatusCode) //Uh-oh this means our test failed
		}


	  	this.Data["Launching"] = "Launching successful! Redirecting in 2s ..."
	  	time.Sleep(5 * time.Second)
	  	this.Redirect("/user/apps", 302)

	}
	flash := beego.ReadFromRequest(&this.Controller)
	if n, ok := flash.Data["notice"]; ok {
		this.Data["notice"] = n
	}
}


func (this *DeleteController) activeContent(view string) {
	this.Layout = "basic-layout.tpl"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "header.tpl"
	this.LayoutSections["Sidebar"] = "sidebar.tpl"
	this.LayoutSections["Footer"] = "footer.tpl"

	appName := this.Ctx.Input.Param(":appName")

	this.TplNames = view + ".tpl"
    this.Data["ApplicationName"] = appName

	sess := this.GetSession("acme")
	if sess != nil {
		this.Data["InSession"] = 1 // for login bar in header.tpl
		m := sess.(map[string]interface{})
		this.Data["First"] = m["first"]
	}
}

func (this *DeleteController) Get() {
	this.activeContent("deleteApp")

	//******** This page requires login
	sess := this.GetSession("acme")
	if sess == nil {
		this.Redirect("/user/login/home", 302)
		return
	}

	appName  := this.Ctx.Input.Param(":appName")
	this.Data["YesOrNo"]="<h1>Really delete this application? (this process is irreversible):</h1><br> <form action='/deleteApp/" + appName + 
	"' method='post' id='form1'><button type='submit' " +
	"form='form1' value='YES'>Yes, do it!</button> </form><br><form action='/user/profile' method='get' id='form2'>" +
	" <button type='submit' form='form2' value='NO'>I changed my mind! SAVE MY APP!</button></form>"

/*
	time.Sleep(5 * time.Second)
	this.Redirect("/user/profile", 302)
*/
	m := sess.(map[string]interface{})
	fmt.Println("username is", m["username"])
	fmt.Println("logged in at", m["timestamp"])
}

func (this *DeleteController) Post() {
	this.activeContent("deleteApp")
	this.Data["YesOrNo"] = "<h1>Deleting application. Will be redirected in 2 seconds.</h1>"
		//******** This page requires login
	sess := this.GetSession("acme")
	if sess == nil {
		this.Redirect("/user/login/home", 302)
		return
	}

	appName  := this.Ctx.Input.Param(":appName")

	//Create the request
	url := "http://" + os.Getenv("MARATHON_ENDPOINT")+ "/v2/apps/" + string(appName)

	req, err := http.NewRequest("DELETE", url, nil)

	//Make the request
	res, err := http.DefaultClient.Do(req)

	if err != nil {
    	panic(err) //Something is wrong while sending request
 	}

	if res.StatusCode != 201 {
		fmt.Printf("Success expected: %d", res.StatusCode) //Uh-oh this means our test failed
	}

	fmt.Println(res.StatusCode)

	time.Sleep(2 * time.Second)
	this.Redirect("/user/apps", 302)
	m := sess.(map[string]interface{})
	fmt.Println("username is", m["username"])
	fmt.Println("logged in at", m["timestamp"])
}