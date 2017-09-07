package main

// how to add middlewares for specific routes using hhtprouter router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

// handler function for profile endpoint
func profileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the content of the profile controller\n")
	log.Println("executing profile controller")
}

// handler function for /hello/:name endpoint
func helloHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Hello, %s!\n", ps.ByName("name"))
	log.Println("executing hello controller")
}

func loginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "This is the content of the login controller\n")
	log.Println("executing login controller")
}

// auth middleware
func auth(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// do some stuff before
	log.Println("auth middleware -> before executing controller")

	// call endpoint handler
	next(rw, r)

	// do some stuff after
	log.Println("auth middleware -> after the controller was executed")
}

// getUrlParams function is extracting URL parameters
func getUrlParams(router *httprouter.Router, req *http.Request) httprouter.Params {

	_, params, _ := router.Lookup(req.Method, req.URL.Path)

	return params
}

// callwithParams function is helping us to call controller from middleware having access to URL params
func callwithParams(router *httprouter.Router, handler func(w http.ResponseWriter, r *http.Request, ps httprouter.Params)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		params := getUrlParams(router, r)
		handler(w, r, params)
	}
}

func main() {
	router := httprouter.New()
	router.POST("/login", loginHandler)

	// add middleware for a specific route
	nProfile := negroni.New()
	nProfile.Use(negroni.HandlerFunc(auth))
	nProfile.UseHandlerFunc(profileHandler)
	router.Handler("GET", "/", nProfile)

	// add middleware for a specific route and get params from route
	nHello := negroni.New()
	nHello.Use(negroni.HandlerFunc(auth))
	nHello.UseHandlerFunc(callwithParams(router, helloHandler))
	router.Handler("GET", "/hello/:name", nHello)

	// Includes some default middlewares to all routes
	n := negroni.Classic()
	n.UseHandler(router)

	log.Fatal(http.ListenAndServe(":8080", n))
}
