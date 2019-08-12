package main

import (
	"github.com/drichtarik/goWebApp/handlers"
	"github.com/drichtarik/goWebApp/rest"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	muxRouter := mux.NewRouter().StrictSlash(true)
	handlers.BootAllPageHandlers(muxRouter)
	rest.BootAllRestApiHandlers(muxRouter)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
