package controllers

import (
	"context"
	"log"
	"net/http"
	"restaurant-management/database"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client , "food")

func GetFoods() gin.HandlerFunc{
   return func(c * gin.Context){
     ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	 defer cancel()
	 foods, err := foodCollection.Find(ctx, bson.M{})
	 if err != nil{
		log.Fatal(err);
	 }
	 var allFoods []bson.M
     err = foods.All(ctx, &allFoods)
	 if err != nil{
		log.Fatal(err);
	 }
	 c.JSON(http.StatusOK, allFoods)	 
   }
}

func GetFood() gin.HandlerFunc{
	return func(c * gin.Context){
      var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	  defer cancel()
	  foodId := c.Param("food_id")
	  var food models.Food
	  err := foodCollection.FindOne(ctx,bson.M{"food_id":foodId} ).Decode(&food)
      if err != nil{
		log.Fatal(err)
	  }
	  c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc{
	return func(c * gin.Context){

	}
}

func UpdateFood() gin.HandlerFunc{
	return func(c * gin.Context){

	}
}