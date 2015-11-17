package models

type Application struct {
  Id       					int64
  Name 						string					`sql:"size:255; not null; unique;"`
  DockerImage				string 					`sql:"size:255; not null; unique;"` 
  Dependencies				[]Dependency			`sql:"size:255; not null;"`
  EnvironmentVariables 		[]EnvironmentVariable
  Logo 						string 					`sql:"size:255; not null;"`   
  Description 				string 					`sql:"size:255; not null;"`
  InternalPort 				string 					`sql:"size:30"` 
  Protocol					string 					`sql:"size:30"`
  	User					User					`sql:"-"`
  	Message					string 					`sql:"-"`	//These dont get put in the database
}

type Dependency struct {
	ID 				int64	
	ApplicationID	int64 	`sql:index`
	Dependency		string
}

type EnvironmentVariable struct {
	ID 				int64
	ApplicationID	int64 	`sql:index`
	Key				string
	Value			string
}

//Validates the sql application
func (a *Application) ValidateApplication() {
	    a.Message = ""

		//Check for valid appplication name (letters, numbers and -)
	    if !VerifyApplicationName() {
	    	a.Message = a.Message + "Not a valid application name -- "
	    }

		//Check for valid docker image (letters,numbers and - / letters,numbers and -)
	    if !VerifyImage() {
	    	a.Message = a.Message + "Not a valid image name -- "
	    }

		//Check for valid url
	    if !VerifyURL() {
	    	a.Message = a.Message + "Logo url is invalid -- "
	    }

		//Check description for some nasty
		if !VerifyDescription() {
			a.Message = a.Message + "There's some nastiness in your description -- "
		}

		//Check that port is valid
		if !VerifyPort() {
			a.Message = a.Message + "The port provided is invalid -- "
		}

		//Protocol must be TCP, UDP or HTTP
		if !VerifyProtocol() {
			a.Message = a.Message + "Protocol must be HTTP, TCP or UDP -- "
		}

}

func VerifyApplicationName() bool {

	return true
}

func VerifyImage() bool {
	return true
}

func VerifyURL() bool {
	return true
}

func VerifyDescription() bool {
	return true
}

func VerifyPort() bool {
	return true
}

func VerifyProtocol() bool {
	return true;
}