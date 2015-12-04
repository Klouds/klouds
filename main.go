package main

import (
	"net/http"
	"github.com/superordinate/klouds/routers"
)


//	Action defines a standard function signature for us to use when 
// 	creating controller actions. A controller action is basically just a method attached
//	to a controller.


func main() {
	var newmux routers.Routing
	newmux.Init()
	http.ListenAndServe("0.0.0.0:1337", newmux.Mux)
}