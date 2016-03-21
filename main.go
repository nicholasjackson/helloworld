package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nicholasjackson/helloworld/global"
	"github.com/nicholasjackson/helloworld/handlers"

	"github.com/alexcesaro/statsd"
	"github.com/facebookgo/inject"
)

func main() {
	config := os.Args[1]
	rootfolder := os.Args[2]

	global.LoadConfig(config, rootfolder)

	setupInjection()
	setupHandlers()
}

func setupHandlers() {
	http.Handle("/", handlers.GetRouter())

	fmt.Println("Listening for connections on port", 8001)
	http.ListenAndServe(fmt.Sprintf(":%v", 8001), nil)
}

func setupInjection() {
	var g inject.Graph

	var err error

	statsdClient, err := statsd.New(statsd.Address(global.Config.StatsDServerIP)) // reference to a statsd client
	if err != nil {
		panic(fmt.Sprintln("Unable to create StatsD Client: ", err))
	}

	err = g.Provide(
		&inject.Object{Value: handlers.RouterDependencies},
		&inject.Object{Value: handlers.HealthDependencies},
		&inject.Object{Value: handlers.EchoDependencies},
		&inject.Object{Value: statsdClient, Name: "statsd"},
	)

	if err != nil {
		fmt.Println(err)
	}

	// Here the Populate call is creating instances of NameAPI &
	// PlanetAPI, and setting the HTTPTransport on both to the
	// http.DefaultTransport provided above:
	if err := g.Populate(); err != nil {
		fmt.Println(err)
	}
}
