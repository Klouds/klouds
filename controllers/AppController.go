package controllers

import (
	"net/http"
)
type Action func(rw http.ResponseWriter, r *http.Request) error

//This is our base controller
type AppController struct{}

//The Action functions helps with error handling in a controller
func (c *AppController) Action(a Action) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := a(rw, r); err != nil {
			http.Error(rw, err.Error(), 500)
		}

	})

}
func redirect(rw http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(rw, r, url, http.StatusMovedPermanently)
}
func readHttpBody(response *http.Response) string {


	bodyBuffer := make([]byte, 5000)
	var str string

	count, err := response.Body.Read(bodyBuffer)

	for ; count > 0; count, err = response.Body.Read(bodyBuffer) {

		if err != nil {

		}

		str += string(bodyBuffer[:count])
	}

	return str

}
