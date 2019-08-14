package rest

import (
	"encoding/json"
	"github.com/drichtarik/goWebApp/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func getPagesEndpoint(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(handlers.Pages)
}

func BootAllPagesRestApiHandlers(router *mux.Router) {
	router.HandleFunc("/pages/", getPagesEndpoint)
}
