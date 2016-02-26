package handlers

import (
	"encoding/json"
	"net/http"
	
	"github.com/nicholasjackson/helloworld/logging"
)

// This is not particularlly a real world example it mearly shows how a builder or a factory could be injected
// into the HealthHandler
type HealthResponseBuilder struct {
	statusMessage string
}

func (b *HealthResponseBuilder) SetStatusMessage(message string) *HealthResponseBuilder {
	b.statusMessage = message
	return b
}

func (b *HealthResponseBuilder) Build() HealthResponse {
	var hr HealthResponse
	hr.StatusMessage = b.statusMessage
	return hr
}

type HealthDependenciesContainer struct {
	// if not specified will create singleton
	SingletonBuilder *HealthResponseBuilder `inject:""`

	// statsD interface must use a name type as injection cannot infer ducktypes
	Stats logging.StatsD `inject:"statsd"`

	// if not specified in the graph will automatically create private instance
	PrivateBuilder *HealthResponseBuilder `inject:"private"`
}

type HealthResponse struct {
	StatusMessage string `json:"status_message"`
}

var HealthDependencies *HealthDependenciesContainer = &HealthDependenciesContainer{}

func HealthHandler(rw http.ResponseWriter, r *http.Request) {
	// all HealthHandlerDependencies are automatically created by injection process
	HealthDependencies.Stats.Increment(HEALTH_HANDLER + GET + CALLED)

	response := HealthDependencies.SingletonBuilder.SetStatusMessage("OK").Build()

	encoder := json.NewEncoder(rw)
	encoder.Encode(&response)

	HealthDependencies.Stats.Increment(HEALTH_HANDLER + GET + SUCCESS)
}
