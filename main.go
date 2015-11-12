package main

import (
	"net/http"
	"github.com/superordinate/klouds2.0/routers"
	"github.com/joho/godotenv"
	"log"
)

type Book struct {
	Title string	`json:"title"`
	Author string	`json:"author"`
}

//	Action defines a standard function signature for us to use when 
// 	creating controller actions. A controller action is basically just a method attached
//	to a controller.




func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}


	var newmux routers.Routing
	newmux.Init()
	http.ListenAndServe("0.0.0.0:1337", newmux.Mux)
}