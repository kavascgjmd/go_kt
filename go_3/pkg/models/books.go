package models

import (
	"books/pkg/config"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name string `json:"name"`
	Author string `json:"author"`
	Publication string `json:"publication"`
}

var db *gorm.DB

func SetDB(){
	db = config.GetDB()
	db.AutoMigrate(&Book{});
}

func (b * Book) CreateBook() * Book{
    db.Create(b);
    return b
}

func GetAllBooks() [] Book {
	 var Books [] Book
	 db.Find(Books)
	 return Books
}

func GetBookbyId(Id int64) * Book{
	 book := &Book{}
     db.Where("ID == ?",Id).Find(book);
	 return book
}
