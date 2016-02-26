package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facebookgo/inject"
	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/nicholasjackson/helloworld/mocks"
)

var echoStatsDMock *mocks.MockStatsD

func echoTestSetup(t *testing.T) {
	// create an injection graph containing the mocked elements we wish to replace

	var g inject.Graph

	echoStatsDMock = &mocks.MockStatsD{}
	EchoDependencies = &EchoDependenciesContainer{}

	err := g.Provide(
		&inject.Object{Value: EchoDependencies},
		&inject.Object{Value: echoStatsDMock, Name: "statsd"},
	)

	if err != nil {
		fmt.Println(err)
	}

	if err := g.Populate(); err != nil {
		fmt.Println(err)
	}

	echoStatsDMock.Mock.On("Increment", mock.Anything).Return()
}

func TestEchoHandlerSetStats(t *testing.T) {
	echoTestSetup(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request

	echo := Echo{Echo: "Hello World"}
	context.Set(&request, "request", &echo)

	EchoHandler(&responseRecorder, &request)

	echoStatsDMock.Mock.AssertCalled(t, "Increment", ECHO_HANDLER+POST+CALLED)
	echoStatsDMock.Mock.AssertCalled(t, "Increment", ECHO_HANDLER+POST+SUCCESS)
}

func TestEchoHandlerCorrectlyEchosResponse(t *testing.T) {
	echoTestSetup(t)

	var responseRecorder *httptest.ResponseRecorder
	var request http.Request

	responseRecorder = httptest.NewRecorder()

	echo := Echo{Echo: "Hello World"}
	context.Set(&request, "request", &echo)

	EchoHandler(responseRecorder, &request)

	body := responseRecorder.Body.Bytes()
	response := Echo{}
	json.Unmarshal(body, &response)

	assert.Equal(t, 200, responseRecorder.Code)
	assert.Equal(t, response.Echo, "Hello World")
}
