package controllers

import (
	"net/http"
	"github.com/superordinate/klouds2.0/models"
	"gopkg.in/unrolled/render.v1"
	"github.com/julienschmidt/httprouter"
	"strings"
	"fmt"
	"io/ioutil"
	"bytes"
	"os"
)

type ApplicationsController struct {
	AppController
	*render.Render
}

func (c *ApplicationsController) AppAdmin(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

		if r.Method == "GET" {
		
		var user *models.User
		

		if getUserName(r) != "" || user != nil {
			user = GetUserByUsername(getUserName(r))
			
			if NotAdministrator(user, c, rw) {
				return
			}


			c.HTML(rw, http.StatusOK, "apps/admin/admin", user)
			
		} else {

			c.HTML(rw, http.StatusOK, "user/login", nil)
		}

	}
		
}

func (c *ApplicationsController) CreateApplication(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	if r.Method == "GET" {
		
		var user *models.User
		//var newapp *models.Application
		if getUserName(r) != "" || user != nil{
			user = GetUserByUsername(getUserName(r))

			if NotAdministrator(user, c, rw) {
				return
			}

			newapp := &models.Application{Username: getUserName(r)}
			c.HTML(rw, http.StatusOK, "apps/admin/newapp", newapp)
			
		} else {
			c.HTML(rw, http.StatusOK, "user/login", nil)
		}

	} else if r.Method == "POST" {
		var user *models.User
		if getUserName(r) != "" || user != nil{

			user = GetUserByUsername(getUserName(r))

			if NotAdministrator(user, c, rw) {
				return
			}

			r.ParseForm();
			//Load Application from form

			appname := strings.ToLower(r.FormValue("appname"))
			image := r.FormValue("dockerimage")
			depends := r.FormValue("dependencies")
			//split them by commas
			dependencies := strings.Split(depends,",")

			//remove white space
			for i:=0; i< len(dependencies); i++ {
				dependencies[i] = stripSpaces(dependencies[i])
				fmt.Println(dependencies[i])
			}

			//Create structure array
			dependecystruct := make ([]models.Dependency, len(dependencies))

			for i:=0; i< len(dependencies); i++ {
				dependecystruct[i] = models.Dependency{Dependency: dependencies[i]}
			}

			envvars := r.FormValue("environmentvariables")
			environmentvariables := strings.Split(envvars,",")

			//remove white space
			for i:=0; i< len(environmentvariables); i++ {
				environmentvariables[i] = stripSpaces(environmentvariables[i])
				fmt.Println(environmentvariables[i])
			}

			//Create structure array
			envvarsstruct := make ([]models.EnvironmentVariable, len(environmentvariables))
			fmt.Println("length! : ",len(environmentvariables))

			for i:=0; i < len(environmentvariables); i++ {
				split := strings.Split(environmentvariables[i], ":")
				if len(envvarsstruct) >= 2 {
					envvarsstruct[i] = models.EnvironmentVariable{	Key: string(split[0]),
																Value: string(split[1])}
				}
			}

			logo := r.FormValue("logo")
			if logo == "" {
				logo ="/images/noimage.png"
			}

			internalport := r.FormValue("internalport")
			protocol := r.FormValue("protocol")
			description := r.FormValue("description")

			newapp := models.Application{
				Name: 					appname,
				DockerImage:			image,
				Dependencies: 			dependecystruct,
				EnvironmentVariables:	envvarsstruct,
				Logo:					logo,
				Description:			description,
				InternalPort: 			internalport,
				Protocol: 				protocol,
				Username:				getUserName(r)}
			
			//Validate fields
			newapp.ValidateApplication()

			//Check against application database for dependencies

			for i:=0; i< len(newapp.Dependencies); i++ {
				
				if !CheckApplicationExists(newapp.Dependencies[i].Dependency) {
					newapp.Message = newapp.Message + "Dependency " + newapp.Dependencies[i].Dependency + " does not exist in the database."
				}
			}

			if (newapp.Message != "") {
				//Validation failed
				c.HTML(rw, http.StatusOK, "apps/admin/newapp", newapp)
				return
			}

			//Create Application

			CreateApplication(&newapp)
			newapp.Message = "Application " + newapp.Name + " created successfully."
			c.HTML(rw, http.StatusOK, "apps/admin/newapp", newapp)
		} else {
			c.HTML(rw, http.StatusOK, "user/login", nil)
		}
	}
		
}

func (c *ApplicationsController) ApplicationList(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	if r.Method == "GET" {
		applicationList := []models.Application{}

		GetApplications(&applicationList)

		for i:=0;i<len(applicationList);i++ {
			applicationList[i].Username = getUserName(r)
		}

		if len(applicationList) == 0 {
			applicationList = []models.Application{models.Application{Username: getUserName(r)}}	
		}

		//Display Application list page
		c.HTML(rw, http.StatusOK, "apps/list", applicationList)

	} else if r.Method == "POST" {

		//Don't think this is needed but who knows :D
	}
}	

func (c *ApplicationsController) Application(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	if r.Method == "GET" {
		application := GetApplicationByName(p.ByName("appID"))

		application.Username = getUserName(r)

		//Display Application list page
		c.HTML(rw, http.StatusOK, "apps/application", application)

	} else if r.Method == "POST" {

		//Don't think this is needed but who knows :D
	}
}	

func (c *ApplicationsController) Launch(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

 	if r.Method == "POST" {
 		//If logged in
 		var user *models.User

 		if getUserName(r) != "" {

 			application := GetApplicationByName(p.ByName("appID"))
			user = GetUserByUsername(getUserName(r))

			if NotAdministrator(user, c, rw) {
				return
			}

			//Read the JSON template
			podfile, err := ioutil.ReadFile("public/json/template.json")

			if err != nil {
				panic(err)
			}

			SplitAtID := strings.Split(string(podfile), "#APPLICATIONNAME")
			SplitAtImage := strings.Split(string(SplitAtID[1]), "#DOCKERIMAGE")
			SplitAtPort := strings.Split(string(SplitAtImage[1]), "#INTERNALPORT")
			SplitAtProtocol := strings.Split(string(SplitAtPort[1]), "#PROTOCOL")
			SplitAtVariables := strings.Split(string(SplitAtProtocol[1]), "#ENVIRONMENTVARIABLES")
			SplitAtHTTP := strings.Split(string(SplitAtVariables[1]), "#ISITHTTP")
			SplitAtRoutingName := strings.Split(string(SplitAtHTTP[1]), "#ROUTING")

			//Inject values into JSON
			//make some custom things
			envvariables := ""

			for i:=0; i< len(application.EnvironmentVariables); i++ {
				envvariables = envvariables + `"` + application.EnvironmentVariables[i].Key + `":"` +
								application.EnvironmentVariables[i].Value + `"`

				if i != len(application.EnvironmentVariables) - 1 {
					envvariables = envvariables + `,`
				}
			}

			protocol := ""

			ishttp := "false";

			if strings.ToLower(application.Protocol) == "http" {
				ishttp = "true"
				protocol ="tcp"
			} else {
				protocol = strings.ToLower(application.Protocol)
			}

			application.Name = user.Username + "-" + strings.ToLower(RandString(8)) + "-" + application.Name


			//merge the new string
			newstring := SplitAtID[0] + application.Name + SplitAtImage[0] + application.DockerImage +
							SplitAtPort[0] + application.InternalPort + SplitAtProtocol[0] + 
							protocol + SplitAtVariables[0] + envvariables +
							SplitAtHTTP[0] + ishttp + SplitAtRoutingName[0] + application.Name + 
							SplitAtRoutingName[1]	
			
			fmt.Println(newstring)
			//Launch against marathon
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


			//Add launched application to DB
			runningapp := models.RunningApplication{
					Name:	application.Name,
					ApplicationID:	application.Id,
					Owner:	user.Id,
					AccessUrl:	application.Name + "." + os.Getenv("KLOUDS_DOMAIN"),
			}

			AddRunningApplication(&runningapp)
			//Display new application
			application.Username = getUserName(r)

			c.HTML(rw, http.StatusOK, "apps/launch", application)
			
		} else {
			c.HTML(rw, http.StatusOK, "user/login", nil)
		}
	}
}	