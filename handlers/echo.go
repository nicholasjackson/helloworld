package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/nicholasjackson/helloworld/logging"
)

type EchoDependenciesContainer struct {
	StatsD logging.StatsD `inject:"statsd"`
	Log    *log.Logger    `inject:""`
}

var EchoDependencies *EchoDependenciesContainer = &EchoDependenciesContainer{}

const EHTAGNAME = "EchoHandler: "

// use the validation middleware to automatically validate input
// github.com/asaskevich/govalidator
type Echo struct {
	Echo string `json:"echo" valid:"stringlength(1|255),required"`
}

func EchoHandler(rw http.ResponseWriter, r *http.Request) {
	EchoDependencies.StatsD.Increment(ECHO_HANDLER + POST + CALLED)
	EchoDependencies.Log.Printf("%v Called GET\n", EHTAGNAME)

	// request is set into the context from the middleware
	request := context.Get(r, "request").(*Echo)
	fmt.Println("r: ", request)

	encoder := json.NewEncoder(rw)
	encoder.Encode(request)

	EchoDependencies.StatsD.Increment(ECHO_HANDLER + POST + SUCCESS)
}
