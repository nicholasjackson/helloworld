package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/context"
	"github.com/nicholasjackson/helloworld/logging"
)

func requestValidationHandler(mainHandlerRef string, t reflect.Type, statsD logging.StatsD, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		request := reflect.New(t).Interface()

		defer r.Body.Close()
		data, _ := ioutil.ReadAll(r.Body)

		err := json.Unmarshal(data, &request)
		if err != nil {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			statsD.Increment(mainHandlerRef + BAD_REQUEST)
			return
		}

		_, err = govalidator.ValidateStruct(request)
		if err != nil {
			fmt.Println("Validation Error:", err)
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			statsD.Increment(mainHandlerRef + INVALID_REQUEST)
			return
		}

		context.Set(r, "request", request)
		statsD.Increment(mainHandlerRef + VALID_REQUEST)
		next.ServeHTTP(w, r)
	})
}
