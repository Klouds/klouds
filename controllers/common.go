package controllers

import (
	"net/http"
 	_ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    "github.com/superordinate/klouds2.0/models"
    "github.com/gorilla/securecookie"
    "time"
    "strings"
    "unicode"
    "fmt"
    "math/rand"
    "os"
)

type ErrorMessage struct {
	Message	string

}

var (
	db *gorm.DB
	cookieHandler *securecookie.SecureCookie 
)

//Initializes supporting functions
func Init() {
	InitCookieHandler()
	InitDB()
}


/*Session Management */
//Initialize the cookie handler
func InitCookieHandler() {
	cookieHandler = securecookie.New(
    	securecookie.GenerateRandomKey(64),
     	securecookie.GenerateRandomKey(32))

}

//Open a new session
func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}

	if encoded, err := cookieHandler.Encode("kloudsSession", value); err == nil {
	 	cookie := &http.Cookie {
		    Name:  "kloudsSession",
		    Value: encoded,
		    Path:  "/",
		    HttpOnly: true,

		}
		
		http.SetCookie(response, cookie)
	}
}

//Gets the logged in username
func getUserName(request *http.Request) (userName string) {
    if cookie, err := request.Cookie("kloudsSession"); err == nil {
       	cookieValue := make(map[string]string)

       	if err = cookieHandler.Decode("kloudsSession", cookie.Value, &cookieValue); err == nil {
           userName = cookieValue["name"]
       	}
    }

   	return userName
}

//clears the active session
func clearSession(response http.ResponseWriter) {
	loc, _ := time.LoadLocation("UTC")

   	cookie := &http.Cookie{
    	Name:   "kloudsSession",
        Value:  "",
        Path:   "/",
        Expires: time.Date(1970, 1, 1,1,1,1,0,loc),
        MaxAge: 0,
    }
    fmt.Println("clearing session")
    http.SetCookie(response, cookie)
}

/* DATABASE FUNCTIONALITY */
// connect to the db
func InitDB() {

	fmt.Println("Initializing Database connection.")

	mysqlhost := os.Getenv("MYSQL_HOST")
	mysqluser := os.Getenv("MYSQL_USER")
	mysqlpass := os.Getenv("MYSQL_PASSWORD")

    dbm, err := gorm.Open("mysql", mysqluser+ ":" + mysqlpass + 
    		"@(" + mysqlhost + ")/klouds?charset=utf8&parseTime=True")

    if(err != nil){
        panic("Unable to connect to the database")
    } else {
    	fmt.Println("Database connection established.")
    }

    db = &dbm
    dbm.DB().Ping()
    dbm.DB().SetMaxIdleConns(10)
    dbm.DB().SetMaxOpenConns(100)
    db.LogMode(true)
 
    if !dbm.HasTable(&models.User{}){
        dbm.CreateTable(&models.User{})
    }
    if !dbm.HasTable(&models.Application{}){
        dbm.CreateTable(&models.Application{})
    }
    if !dbm.HasTable(&models.Dependency{}){
        dbm.CreateTable(&models.Dependency{})
    }
    if !dbm.HasTable(&models.EnvironmentVariable{}){
        dbm.CreateTable(&models.EnvironmentVariable{})
    }
    if !dbm.HasTable(&models.RunningApplication{}){
        dbm.CreateTable(&models.RunningApplication{})
    }   
}


/* HELPER FUNCTIONS */

//Checks if user is admin and redirects to login page if s/he isn't
func NotAdministrator(user *models.User, c *ApplicationsController, rw http.ResponseWriter) bool {

	if user.Role != "admin"{
		user.Message = "You are not an administrator."
		c.HTML(rw, http.StatusOK, "user/login", user)
		return true
	} else {
		return false
	}
}

//strips all whitespace out of a string
func stripSpaces(str string) string {
    return strings.Map(func(r rune) rune {
        if unicode.IsSpace(r) {
            // if the character is a space, drop it
            return -1
        }
        // else keep it in the string
        return r
    }, str)
}

/* USER DATABASE CALLS */
//Create a new user in the database
func CreateUser(u *models.User) {
	fmt.Println("Creating user: " + u.Username)

	db.Create(&u)
}

//Check if Username exists, returns false if username not taken
func CheckForExistingUsername(u *models.User) bool {
	newUser := &models.User{}

	db.Where(&models.User{Username: u.Username}).First(&newUser)

	return newUser.Id == 0
} 

//Check if Email exists, returns false if email not taken
func CheckForExistingEmail(u *models.User) bool {
	newUser := &models.User{}

	db.Where(&models.User{Email: u.Email}).First(&newUser)

	return newUser.Id == 0
}

//Check if passwords match
func CheckForMatchingPassword(u *models.User) bool {
	newUser := &models.User{}

	db.Where(&models.User{Email: u.Email}).First(&newUser)

	return newUser.Password == u.Password
} 

//Get User
func GetUserByUsername(username string) *models.User {
	newUser := &models.User{}

	db.Where(&models.User{Username: username}).First(&newUser)

	return newUser
}

//UpdateUser
func UpdateUser(u *models.User) {
	db.Save(&u)
} 


/* APPLICATION DATABASE THINGS */

//Check if application exists in database, returns true if app exists
func CheckApplicationExists(appname string) bool {
	if appname == "" {
		return true
	}

	newapp := &models.Application{}

	db.Where(&models.Application{Name: appname}).First(&newapp)

	return newapp.Id != 0
}

//Create application in DB
func CreateApplication(a *models.Application) {
	fmt.Println("Creating Application: " + a.Name)

	db.Create(&a)

}

//Get Application List
func GetApplications(a *[]models.Application) {
	fmt.Println("Getting list of all applications")

	//Returns a list of all applications 
	applicationList := []models.Application{}

	db.Find(&applicationList)
	//makes the list externally available
	*a = applicationList
}

//Get application by name
func GetApplicationByName(appname string) *models.Application {

	newapp := &models.Application{}
	dependencies := []models.Dependency{}
	envvariables := []models.EnvironmentVariable{}

	//Get an application -- this doesnt grab associated dbs
	db.Where(&models.Application{Name: appname}).First(&newapp)

	db.Model(&newapp).Related(&dependencies)
	db.Model(&newapp).Related(&envvariables)

	fmt.Println(newapp)
	fmt.Println(dependencies)
	fmt.Println(envvariables)

	newapp.Dependencies = dependencies
	newapp.EnvironmentVariables = envvariables
	//Get dependencies

	return newapp
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

/* USER APPLICATION THINGS */

func AddRunningApplication(a *models.RunningApplication) {
	fmt.Println("Adding running Application: " + a.Name)

	db.Create(&a)
}


func GetRunningApplications(a *[]models.RunningApplication) {

	fmt.Println("Getting list of all running applications")

	//Returns a list of all applications 
	applicationList := []models.RunningApplication{}

	db.Find(&applicationList)
	//makes the list externally available
	*a = applicationList

}

func GetRunningApplicationsForUser(a *[]models.RunningApplication, u *models.User) {

	runningapps := []models.RunningApplication{}

	db.Where("owner = ?", u.Id).Find(&runningapps)

	LoadLogoForRunningApplications(&runningapps)
	*a = runningapps
}

func LoadLogoForRunningApplications(a *[]models.RunningApplication) {

	for index := range *a {
		application := models.Application{}

		db.Where("id = ?", (*a)[index].ApplicationID).First(&application)

		(*a)[index].Logo = application.Logo
	}
}