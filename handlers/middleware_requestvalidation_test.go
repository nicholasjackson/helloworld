package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/nicholasjackson/helloworld/mocks"
)

type mockType struct {
	FirstName string `json:"first_name" valid:"alphanum,stringlength(1|255),required"`
}

var mockHandler *mocks.MockHandler
var mockRequestStatsD *mocks.MockStatsD

func setupRequestValidationTests(t *testing.T) {
	mockHandler = &mocks.MockHandler{}
	mockRequestStatsD = &mocks.MockStatsD{}

	mockRequestStatsD.Mock.On("Increment", mock.Anything)
	mockHandler.Mock.On("ServeHTTP", mock.Anything, mock.Anything)
}

func TestCallsNextOnSuccessfulValidation(t *testing.T) {
	setupRequestValidationTests(t)
	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`{"first_name": "Nic"}`))

	handlerFunc := requestValidationHandler(HEALTH_HANDLER, reflect.TypeOf(mockType{}), mockRequestStatsD, mockHandler)

	handlerFunc.ServeHTTP(&responseRecorder, &request)

	mockHandler.Mock.AssertCalled(t, "ServeHTTP", mock.Anything, mock.Anything)
	mockRequestStatsD.Mock.AssertCalled(t, "Increment", HEALTH_HANDLER+VALID_REQUEST)
}

func TestSetsContextSuccessfully(t *testing.T) {
	setupRequestValidationTests(t)
	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`{"first_name": "Nic"}`))

	handlerFunc := requestValidationHandler(HEALTH_HANDLER, reflect.TypeOf(mockType{}), mockRequestStatsD, mockHandler)

	handlerFunc.ServeHTTP(&responseRecorder, &request)
	requestObj := context.Get(&request, "request").(*mockType)
	assert.Equal(t, "Nic", requestObj.FirstName)
}

func TestReturnsBadRequestWhenNoObject(t *testing.T) {
	setupRequestValidationTests(t)
	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(``))

	handlerFunc := requestValidationHandler(HEALTH_HANDLER, reflect.TypeOf(mockType{}), mockRequestStatsD, mockHandler)

	handlerFunc.ServeHTTP(&responseRecorder, &request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	mockRequestStatsD.Mock.AssertCalled(t, "Increment", HEALTH_HANDLER+BAD_REQUEST)
}

func TestReturnsBadRequestWhenRequestInvalid(t *testing.T) {
	setupRequestValidationTests(t)
	var responseRecorder httptest.ResponseRecorder
	var request http.Request
	request.Body = ioutil.NopCloser(bytes.NewBufferString(`{"first_name": ""}`))

	handlerFunc := requestValidationHandler(HEALTH_HANDLER, reflect.TypeOf(mockType{}), mockRequestStatsD, mockHandler)

	handlerFunc.ServeHTTP(&responseRecorder, &request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	mockRequestStatsD.Mock.AssertCalled(t, "Increment", HEALTH_HANDLER+INVALID_REQUEST)
}
