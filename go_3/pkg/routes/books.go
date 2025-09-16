package routes

import (
	"github.com/gorilla/mux"
	"books/pkg/controllers"
)

var CreateBookRoutes = func(r * mux.Router){
  r.HandleFunc("/book", controllers.GetAllBooks).Methods("GET")
  r.HandleFunc("/book/{bookid}", controllers.GetBookbyId).Methods("GET")
  r.HandleFunc("/book", controllers.CreateBook).Methods("PUT")
}