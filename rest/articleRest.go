package rest

import (
	"encoding/json"
	"fmt"
	"github.com/drichtarik/goWebApp/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(handlers.Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["name"]
	for _, article := range handlers.Articles {
		if article.Title == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article handlers.Article
	json.Unmarshal(reqBody, &article)
	handlers.Articles = append(handlers.Articles, article)
	json.NewEncoder(w).Encode(article)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["name"]
	for i, article := range handlers.Articles {
		if article.Title == key {
			fmt.Println("delete 2")
			fmt.Println(article)
			handlers.Articles = append(handlers.Articles[:i], handlers.Articles[i+1:]...)
		}
	}
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["name"]
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article handlers.Article
	json.Unmarshal(reqBody, &article)
	for i, a := range handlers.Articles {
		if a.Title == key {
			fmt.Println("Found!")
			handlers.Articles[i] = article
			fmt.Println(handlers.Articles)
		}
	}
	json.NewEncoder(w).Encode(handlers.Articles)
}

func BootAllArticlesRestApiHandlers(router *mux.Router) {
	router.HandleFunc("/articles/", returnAllArticles)
	router.HandleFunc("/article/{name}", deleteArticle).Methods("DELETE")
	router.HandleFunc("/article/", createNewArticle).Methods("POST")
	router.HandleFunc("/article/{name}", updateArticle).Methods("PUT")
	router.HandleFunc("/article/{name}", returnSingleArticle)
}
