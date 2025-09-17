package config
import (
	"sync"
  	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var ( 
	db *gorm.DB
	once sync.Once
)

func Connection() *gorm.DB{
	once.Do(func(){
    dsn := "root:root@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	log.Println("Connected to MySQL successfully!")
	db = d;
	})
	return db;
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("database is not initalized")
	}
	return db
}

func Close() error{
	if db == nil {
		return nil
	}
    sqlDB, err := db.DB();
	if err != nil{
		return err
	}
	return sqlDB.Close()
}