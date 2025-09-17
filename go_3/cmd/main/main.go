package main

import (
	"books/pkg/routes"
	"log"
	"net/http"
    "books/pkg/config"
	"books/pkg/models"
	"github.com/gorilla/mux"
)

func main(){
	config.Connection()
	models.SetDB()
	r := mux.NewRouter();
	routes.CreateBookRoutes(r)
	log.Fatal(http.ListenAndServe(":8080",r))
}

