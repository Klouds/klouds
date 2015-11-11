package main

import (
	"net/http"

	"github.com/superordinate/klouds2.0/routers"
)

type Book struct {
	Title string	`json:"title"`
	Author string	`json:"author"`
}

//	Action defines a standard function signature for us to use when 
// 	creating controller actions. A controller action is basically just a method attached
//	to a controller.




func main() {

	var newmux routers.Routing
	newmux.Init()
	http.ListenAndServe("0.0.0.0:8080", newmux.Mux)
}