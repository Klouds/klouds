package controllers

import (
	"github.com/superordinate/klouds/models"
	pk "github.com/superordinate/klouds/utilities/pbkdf2"
	"encoding/hex"
	"fmt"
	"github.com/alexcesaro/mail/gomail"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/twinj/uuid"
	"strings"
	"net/http"
	"io/ioutil"
	"time"
	"os" //For getting environment variables
	//"strconv"
	//"bytes"
)

func (this *MainController) Login() {
	this.activeContent("user/login")

	back := strings.Replace(this.Ctx.Input.Param(":back"), ">", "/", -1) // allow for deeper URL such as l1/l2/l3 represented by l1>l2>l3
	fmt.Println("back is", back)
	if this.Ctx.Input.Method() == "POST" {
		flash := beego.NewFlash()
		username:= strings.ToLower(this.GetString("username"))
		
		password := this.GetString("password")
		valid := validation.Validation{}
		
		valid.Required(password, "password")

		if valid.HasErrors() {
			errormap := []string{}
			for _, err := range valid.Errors {
				errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
			}
			this.Data["Errors"] = errormap
			return
		}
		fmt.Println("Authorization is", username, ":", password)

		//******** Read password hash from database
		var x pk.PasswordHash

		x.Hash = make([]byte, 32)
		x.Salt = make([]byte, 16)

		o := orm.NewOrm()
		o.Using("default")
		user := models.AuthUser{Username: username}
		err := o.Read(&user, "username")
		if err == nil {
			if user.Reg_key != "" {
				flash.Error("Account not verified")
				flash.Store(&this.Controller)
				return
			}

			// scan in the password hash/salt
			fmt.Println("Password to scan:", user.Password)
			if x.Hash, err = hex.DecodeString(user.Password[:64]); err != nil {
				fmt.Println("ERROR:", err)
			}
			if x.Salt, err = hex.DecodeString(user.Password[64:]); err != nil {
				fmt.Println("ERROR:", err)
			}
			fmt.Println("decoded password is", x)
		} else {
			flash.Error("No such user")
			flash.Store(&this.Controller)
			return
		}

		//******** Compare submitted password with database
		if !pk.MatchPassword(password, &x) {
			flash.Error("Bad password")
			flash.Store(&this.Controller)
			return
		}

		//******** Create session and go back to previous page
		m := make(map[string]interface{})
		m["first"] = user.First
		m["username"] = username
		m["timestamp"] = time.Now()
		this.SetSession("acme", m)
		this.Redirect("/"+back, 302)
	}
}

func (this *MainController) Logout() {
	this.activeContent("logout")
	this.DelSession("acme")
	this.Redirect("/home", 302)
}

func (this *MainController) Register() {
	this.activeContent("user/register")

	if this.Ctx.Input.Method() == "POST" {
		flash := beego.NewFlash()
		first := this.GetString("first")
		last := this.GetString("last")
		email := this.GetString("email")
		username := strings.ToLower(this.GetString("username"))
		password := this.GetString("password")
		password2 := this.GetString("password2")


		//Check username for special characters
		f := func(r rune) bool {
        	return (r < 'A'  || r > 'Z') && (r < 'a' || r > 'z')  && (r < '0' || r > '9') && r != '-'
    	}

		if strings.IndexFunc(username, f) != -1 {
        	flash.Error(username + " contains special characters. Only A-Z a-z 0-9 and '-' allowed!")
			flash.Store(&this.Controller)
			return
    	}

		valid := validation.Validation{}
		valid.Required(first, "first")
		valid.Email(email, "email")
		valid.MinSize(password, 6, "password")
		valid.Required(password2, "password2")
		if valid.HasErrors() {
			errormap := []string{}
			for _, err := range valid.Errors {
				errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
			}
			this.Data["Errors"] = errormap
			return
		}
		if password != password2 {
			flash.Error("Passwords don't match")
			flash.Store(&this.Controller)
			return
		}
		h := pk.HashPassword(password)

		//******** Save user info to database
		o := orm.NewOrm()
		o.Using("default")

		user := models.AuthUser{First: first, Last: last, Username: username, Email: email}

		// Convert password hash to string
		user.Password = hex.EncodeToString(h.Hash) + hex.EncodeToString(h.Salt)

		// Add user to database with new uuid and send verification email
		u := uuid.NewV4()
		user.Reg_key = u.String()
		_, err := o.Insert(&user)
		if err != nil {
			flash.Error(email + " already registered or " + username + " already taken")
			flash.Store(&this.Controller)
			return
		}

		if !sendVerification(email, u.String()) {
			flash.Error("Unable to send verification email")
			flash.Store(&this.Controller)
			return
		}
		flash.Notice("Your account has been created. You must verify the account in your email.")
		flash.Store(&this.Controller)
		this.Redirect("/notice", 302)
	}
}

func sendVerification(email, u string) bool {
	link := "http://klouds.org/user/verify/" + u
	host := "smtp.mailgun.org"
	port := 2525

	postmaster := os.Getenv("POSTMASTER")
	postmaster_key := os.Getenv("POSTMASTER_KEY")

	fmt.Println(postmaster_key + " " + postmaster)
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", "p.a.mauviel@gmail.com", "Klouds.io")
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Account Verification for Klouds.io")
	msg.SetBody("text/html", "To verify your account, please click on the link: <a href=\""+link+
		"\">"+link+"</a><br><br>Best Regards,<br>Klouds.io team.")
	m := gomail.NewMailer(host, postmaster, postmaster_key, port)
	if err := m.Send(msg); err != nil {
		return false
	}
	return true
}

func (this *MainController) Verify() {
	this.activeContent("user/verify")

	u := this.Ctx.Input.Param(":uuid")
	o := orm.NewOrm()
	o.Using("default")
	user := models.AuthUser{Reg_key: u}
	err := o.Read(&user, "Reg_key")
	if err == nil {
		this.Data["Verified"] = 1
		user.Reg_key = ""
		if _, err := o.Update(&user); err != nil {
			delete(this.Data, "Verified")
		}
	}
}

func (this *MainController) Profile() {
	this.activeContent("user/profile")

	//******** This page requires login
	sess := this.GetSession("acme")
	if sess == nil {
		this.Redirect("/user/login/home", 302)
		return
	}
	m := sess.(map[string]interface{})

	flash := beego.NewFlash()

	//******** Read password hash from database
	var x pk.PasswordHash

	x.Hash = make([]byte, 32)
	x.Salt = make([]byte, 16)

	o := orm.NewOrm()
	o.Using("default")
	user := models.AuthUser{Username: m["username"].(string)}
	err := o.Read(&user, "Username")
	if err == nil {
		// scan in the password hash/salt
		if x.Hash, err = hex.DecodeString(user.Password[:64]); err != nil {
			fmt.Println("ERROR:", err)
		}
		if x.Salt, err = hex.DecodeString(user.Password[64:]); err != nil {
			fmt.Println("ERROR:", err)
		}
	} else {
		flash.Error("Internal error")
		flash.Store(&this.Controller)
		return
	}

	// this deferred function ensures that the correct fields from the database are displayed
	defer func(this *MainController, user *models.AuthUser) {
		this.Data["First"] = user.First
		this.Data["Last"] = user.Last
		this.Data["Username"] = user.Username
		this.Data["Email"] = user.Email
	}(this, &user)

	if this.Ctx.Input.Method() == "POST" {
		first := this.GetString("first")
		last := this.GetString("last")
		username := strings.ToLower(this.GetString("username"))
		email := this.GetString("email")
		current := this.GetString("current")
		password := this.GetString("password")
		password2 := this.GetString("password2")
		valid := validation.Validation{}

		valid.Required(first, "first")
		valid.Email(email, "email")
		valid.Required(current, "current")

		if valid.HasErrors() {
			errormap := []string{}
			for _, err := range valid.Errors {
				errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
			}
			this.Data["Errors"] = errormap
			return
		}

		if password != "" {
			valid.MinSize(password, 6, "password")
			valid.Required(password2, "password2")
			if valid.HasErrors() {
				errormap := []string{}
				for _, err := range valid.Errors {
					errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
				}
				this.Data["Errors"] = errormap
				return
			}

			if password != password2 {
				flash.Error("Passwords don't match")
				flash.Store(&this.Controller)
				return
			}
			h := pk.HashPassword(password)

			// Convert password hash to string
			user.Password = hex.EncodeToString(h.Hash) + hex.EncodeToString(h.Salt)
		}

		//******** Compare submitted password with database
		if !pk.MatchPassword(current, &x) {
			flash.Error("Bad current password")
			flash.Store(&this.Controller)
			return
		}

		//******** Save user info to database
		user.First = first
		user.Last = last
		user.Username = username
		user.Email = email

		_, err := o.Update(&user)
		if err == nil {
			flash.Notice("Profile updated")
			flash.Store(&this.Controller)
			m["username"] = username
		} else {
			flash.Error("Internal error")
			flash.Store(&this.Controller)
			return
		}
	}
}

func (this *MainController) Remove() {
	this.activeContent("user/remove")

	//******** This page requires login
	sess := this.GetSession("acme")
	if sess == nil {
		this.Redirect("/user/login/home", 302)
		return
	}
	m := sess.(map[string]interface{})

	if this.Ctx.Input.Method() == "POST" {
		current := this.GetString("current")
		valid := validation.Validation{}
		valid.Required(current, "current")
		if valid.HasErrors() {
			errormap := []string{}
			for _, err := range valid.Errors {
				errormap = append(errormap, "Validation failed on "+err.Key+": "+err.Message+"\n")
			}
			this.Data["Errors"] = errormap
			return
		}

		flash := beego.NewFlash()

		//******** Read password hash from database
		var x pk.PasswordHash

		x.Hash = make([]byte, 32)
		x.Salt = make([]byte, 16)

		o := orm.NewOrm()
		o.Using("default")
		user := models.AuthUser{Username: m["username"].(string)}
		err := o.Read(&user, "Username")
		if err == nil {
			// scan in the password hash/salt
			if x.Hash, err = hex.DecodeString(user.Password[:64]); err != nil {
				fmt.Println("ERROR:", err)
			}
			if x.Salt, err = hex.DecodeString(user.Password[64:]); err != nil {
				fmt.Println("ERROR:", err)
			}
		} else {
			flash.Error("Internal error")
			flash.Store(&this.Controller)
			return
		}

		//******** Compare submitted password with database
		if !pk.MatchPassword(current, &x) {
			flash.Error("Bad current password")
			flash.Store(&this.Controller)
			return
		}

		//******** Delete user record
		_, err = o.Delete(&user)
		if err == nil {
			flash.Notice("Your account is deleted.")
			flash.Store(&this.Controller)
			this.DelSession("acme")
			this.Redirect("/notice", 302)
		} else {
			flash.Error("Internal error")
			flash.Store(&this.Controller)
			return
		}
	}
}

func (this *MainController) Apps() {
	this.activeContent("user/apps")

	//******** This page requires login
	sess := this.GetSession("acme")
	if sess == nil {
		this.Redirect("/user/login/home", 302)
		return
	}
	m := sess.(map[string]interface{})

	flash := beego.NewFlash()

	//******** Read password hash from database
	var x pk.PasswordHash

	x.Hash = make([]byte, 32)
	x.Salt = make([]byte, 16)

	o := orm.NewOrm()
	o.Using("default")
	user := models.AuthUser{Username: m["username"].(string)}
	err := o.Read(&user, "Username")
	if err == nil {
		// scan in the password hash/salt
		if x.Hash, err = hex.DecodeString(user.Password[:64]); err != nil {
			fmt.Println("ERROR:", err)
		}
		if x.Salt, err = hex.DecodeString(user.Password[64:]); err != nil {
			fmt.Println("ERROR:", err)
		}
	} else {
		flash.Error("Internal error")
		flash.Store(&this.Controller)
		return
	}

	// Create an applist

	//Get all running apps
	url := "http://" + os.Getenv("MARATHON_ENDPOINT") + "/v2/apps"
	//bytestring := []byte(newstring)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}

	//Make the request
	res, err := http.DefaultClient.Do(req)

	if err != nil {
    	panic(err) //Something is wrong while sending request
 	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}





	JSONSPLIT := strings.Split(string(body), strings.ToLower(user.Username))
	numberOfApps := len(JSONSPLIT) - 1
	//this.Data["AppList"] = string(body)
	appNames := make([]string, numberOfApps)
	appType := make([]string, numberOfApps)

	fmt.Println("Number of apps for user: ", len(JSONSPLIT) - 1)

	for i:=1; i<len(JSONSPLIT); i++ {
		JSONBYTES := []byte(JSONSPLIT[i])

		appNames[i-1] = strings.ToLower(user.Username)




		serviceType := strings.Split(string(JSONSPLIT[i]), "udp")

		if (len(serviceType) == 1) { 
			serviceType = strings.Split(string(JSONSPLIT[i]), "tcp")
		} else {
			appType[i-1] = "udp"

			for j:=0; j<len(JSONBYTES);j++ {
				if JSONBYTES[j] == '"' {
					break;
				}
				appNames[i-1] = appNames[i-1] + string(JSONBYTES[j])
			}
			fmt.Println(appType[i-1])
			fmt.Println(appNames[i-1])			
			continue;
		}
		if (len(serviceType) == 1 ) {
			serviceType = strings.Split(string(JSONSPLIT[i]), "http")
		} else {
			appType[i-1] = "tcp"
			for j:=0; j<len(JSONBYTES);j++ {
				if JSONBYTES[j] == '"' {
					break;
				}
				appNames[i-1] = appNames[i-1] + string(JSONBYTES[j])
			}
			fmt.Println(appType[i-1])
			fmt.Println(appNames[i-1])
			continue;
		}

		if (len(serviceType) != 1) {
			appType[i-1] = "http"
			for j:=0; j<len(JSONBYTES);j++ {
				if JSONBYTES[j] == '"' {
					break;
				}
				appNames[i-1] = appNames[i-1] + string(JSONBYTES[j])
			}
			fmt.Println(appType[i-1])
			fmt.Println(appNames[i-1])			
			continue;
		}
	}

	formstring :=""

	for i:=0;i<len(appNames);i+=2 {
		formstring = formstring + "<tr><td>" + appNames[i] + "</td>"

		/* LINK TO CONTAINER PART */
		if (appType[i] == "tcp") {
			http := false;
			checkForHttp := strings.Split(string(JSONSPLIT[i+1]), "HAPROXY_HTTP")
			if (len(checkForHttp) > 1) {
				httpBytes := []byte(checkForHttp[1])

				for j:=0; j < len(httpBytes);j++ {
					if (httpBytes[j] == 't') {
						http =true;
						break
					} else {
						continue
					}
				}
			}

			if (http) {
				formstring = formstring + "<td><a href='http://" + appNames[i] + ".klouds.org' target='_blank'> GO TO SITE</a></td>"
			} else {
				lookForServicePort := strings.Split(string(JSONSPLIT[i+1]), "servicePort")
				portString := ""

				if (len(lookForServicePort) > 1) {
					servicePortBytes := []byte(lookForServicePort[1])
					
					

					for j:=0; j < len(servicePortBytes);j++ {
						if ((servicePortBytes[j] == '"' || servicePortBytes[j] == ':') && j < 5) {
							continue
						} else if ((servicePortBytes[j] == '"' || servicePortBytes[j] == ',') && j >= 5){
							break
						} else {
							portString = portString + string(servicePortBytes[j])
						}
					}
				}
				formstring = formstring + "<td><input width='200' value='" + appNames[i] + ".klouds.org:" + portString + "' type='text' readonly/></td>"
				}
		} else if (appType[i] == "udp") {

			//REST CALL TO APPNAME
			url := "http://" + os.Getenv("MARATHON_ENDPOINT") + "/v2/apps/" + appNames[i] 
			//bytestring := []byte(newstring)
			req, err := http.NewRequest("GET", url, nil)

			if err != nil {
				panic(err)
			}

			//Make the request
			res, err := http.DefaultClient.Do(req)

			if err != nil {
		    	panic(err) //Something is wrong while sending request
		 	}

			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				panic(err)
			}

			splitbyHost := strings.Split(string(body), "host")	//[3] is the one we want 

			if (len(splitbyHost) < 3) {
				panic ("Fucked something up scraping host")
			}
			splitbyPort := strings.Split(string(splitbyHost[2]), "ports"); // [2] is the one we want 

			if (len (splitbyPort) < 2) {
				panic ("Fucked something up scraping port: " + string(splitbyHost[3]))
			}
			//SCRAPE HOST OUT

			hostBytes := []byte(splitbyHost[2])
			hostString := ""

			for j:=0; j < len(hostBytes);j++ {
				if ((hostBytes[j] == '"' || hostBytes[j] == ':') && j < 5) {
					continue
				} else if ((hostBytes[j] == '"' || hostBytes[j] == ',') && j >= 5){
					break
				} else {
					hostString = hostString + string(hostBytes[j])
				}
			}
						

			//SCRAPE PORT OUT
			hostPort := ""
			portBytes := []byte(splitbyPort[1])
			for j:=0; j < len(portBytes);j++ {
				if ((portBytes[j] == '"' || portBytes[j] == ':') && j < 3) {
					continue
				} else if ((portBytes[j] == '"' || portBytes[j] == ',') && j >= 3){
					break
				} else if (portBytes[j] != '[' && portBytes[j] != ']'){
					hostPort = hostPort + string(portBytes[j])
				}
			}
			//PUT INTO STRING
			formstring = formstring + "<td><input width='200' value='" + hostString + ":" + hostPort + "' type='text' readonly/></td>"
		}
		
		/* PUT CODE TO DO FANCY STUFF HERE */


		formstring = formstring + "<td><a href='../deleteApp/"+ appNames[i] + "'> DELETE APP </a></td>"
	}
	formstring = formstring + "</tr></table>"
	this.Data["AppList"] = formstring
}
