package handlers

import (
	"net/http"
	"reflect"

	"github.com/gorilla/pat"
	"github.com/nicholasjackson/helloworld/logging"
)

type RouterDependenciesContainer struct {
	StatsD logging.StatsD `inject:"statsd"`
}

var RouterDependencies *RouterDependenciesContainer = &RouterDependenciesContainer{}

func GetRouter() *pat.Router {
	r := pat.New()

	r.Get("/v1/health", HealthHandler)

	r.Add("POST", "/v1/echo", requestValidationHandler(
		ECHO_HANDLER+POST,
		reflect.TypeOf(Echo{}),
		RouterDependencies.StatsD,
		http.HandlerFunc(EchoHandler),
	))

	//Add routing for static routes
	s := http.StripPrefix("/swagger/", http.FileServer(http.Dir("/swagger")))
	r.PathPrefix("/swagger/").Handler(s)

	return r
}
