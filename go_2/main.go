package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type Movie struct {
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r * http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func createMovies( w http.ResponseWriter, r * http.Request){
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(2)
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movies)
}

func main(){
	r := mux.NewRouter()
	director := Director{Firstname: "he", Lastname : "she"};
	movies = append(movies, Movie{ID: "1", Isbn:"4", Title: "ok", Director: &director})
	r.HandleFunc("/movies", getMovies).Methods("GET");
	r.HandleFunc("/movies", createMovies).Methods("POST");
	log.Fatal(http.ListenAndServe(":8080", r));
}