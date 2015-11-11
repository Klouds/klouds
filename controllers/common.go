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

    dbm, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/klouds?charset=utf8&parseTime=True")

    if(err != nil){
        panic("Unable to connect to the database")
    } else {
    	fmt.Println("Database connection established.")
    }

    db = &dbm
    dbm.DB().Ping()
    dbm.DB().SetMaxIdleConns(10)
    dbm.DB().SetMaxOpenConns(100)
    db.LogMode(false)
 
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