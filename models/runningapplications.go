package models

type RunningApplication struct {
  Id       					int64
  Name 						string					`sql:"size:255; not null; unique;"`
  ApplicationID				int64 					
  Owner				    	int64			    	`sql:"not null;"`
  AccessUrl          		string         			`sql:"size:255; not null;"` 
  IsRunning					bool
  HostIP					string
  HostPort					int
  ServicePort				int 
    Logo            		string          		`sql:"-"`
  	User					User					`sql:"-"`
  	Message					string 					`sql:"-"`	//These dont get put in the database
}
