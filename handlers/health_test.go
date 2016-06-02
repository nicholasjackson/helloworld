package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/facebookgo/inject"
	"github.com/nicholasjackson/helloworld/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var healthStatsDMock *mocks.MockStatsD

func healthTestSetup(t *testing.T) {
	// create an injection graph containing the mocked elements we wish to replace

	var g inject.Graph

	healthStatsDMock = &mocks.MockStatsD{}
	HealthDependencies = &HealthDependenciesContainer{}

	err := g.Provide(
		&inject.Object{Value: HealthDependencies},
		&inject.Object{Value: healthStatsDMock, Name: "statsd"},
		&inject.Object{Value: log.New(os.Stdout, "tester", log.Lshortfile)},
	)

	if err != nil {
		fmt.Println(err)
	}

	if err := g.Populate(); err != nil {
		fmt.Println(err)
	}

	healthStatsDMock.Mock.On("Increment", mock.Anything).Return()
}

// Simple test to show how we can use the ResponseRecorder to test our HTTP handlers
func TestHealthHandler(t *testing.T) {
	healthTestSetup(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request

	HealthHandler(&responseRecorder, &request)

	assert.Equal(t, 200, responseRecorder.Code)
}

func TestHealthHandlerSetStats(t *testing.T) {
	healthTestSetup(t)

	var responseRecorder httptest.ResponseRecorder
	var request http.Request

	HealthHandler(&responseRecorder, &request)

	healthStatsDMock.Mock.AssertCalled(t, "Increment", HEALTH_HANDLER+GET+CALLED)
	healthStatsDMock.Mock.AssertCalled(t, "Increment", HEALTH_HANDLER+GET+SUCCESS)
}
