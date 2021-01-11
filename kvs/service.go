package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cloud-native-go/kvs/api"
	"github.com/cloud-native-go/kvs/logger"
	"github.com/gorilla/mux"
)

var transact logger.TransactionLogger

func initializeTransactionLog() error {
	var err error

	transact, err = logger.NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := transact.ReadEvents()
	ok, e := true, logger.Event{}

	for ok && err != nil {
		select {
		case err, _ = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case logger.EventDelete:
				err = api.Delete(e.Key)
			case logger.EventPut:
				err = api.Put(e.Key, e.Value)
			}
		}
	}

	transact.Run()

	return err

}

// keyValuePutHandler expects to be called with a PUT request for
// the "/v1/key/{key}"
func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	err = api.Put(key, string(value))

	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// keyValueGetHandler expects to be called with a PUT request for
// the "/v1/key/{key}"
func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := api.Get(key)

	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	transact.WritePut(key, value)

	w.Write([]byte(value))
}

func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := api.Delete(key)

	if err != nil {
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	transact.WriteDelete(key, value)

	w.WriteHeader(http.StatusOK)
}

func main() {

	initializeTransactionLog()

	r := mux.NewRouter()

	// Register keyValuePutHandler as the handler function for PUT
	// requests matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")

	// Register keyValueGetHandler as the handler function for GET
	// requests matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")

	// Register keyValueGetHandler as the handler function for DELETE
	// requests matching "/v1/{key}"
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
