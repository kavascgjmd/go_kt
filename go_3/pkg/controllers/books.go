package controllers

import (
	"books/pkg/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllBooks(w http.ResponseWriter,  r * http.Request){
	 w.Header().Set("Content-Type","application/json")
     books := models.GetAllBooks()
	 json.NewEncoder(w).Encode(books)
}

func GetBookbyId(w http.ResponseWriter, r * http.Request){
    vars := mux.Vars(r);
    id := vars["bookid"]
	ID, _ := strconv.ParseInt(id, 0 , 0)
	book := models.GetBookbyId(ID)
	w.Header().Set("Content-Type","application/json")
    json.NewEncoder(w).Encode(book)

}

func CreateBook(w http.ResponseWriter, r * http.Request){
	w.Header().Set("Content-Type","application/json")
    var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book); if err != nil {
		http.Error(w, "failed to create", http.StatusInternalServerError)
		return
	}
	book.CreateBook()
	json.NewEncoder(w).Encode(book)
}


