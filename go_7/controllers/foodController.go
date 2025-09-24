package controllers

import (
	"context"
	"log"
	"net/http"
	"restaurant-management/database"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
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
      
      recordPerPage, err := strconv.ParseInt(c.Query("recordPerPage"), 10, 64)
	  if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	  }
      
	  page , err := strconv.ParseInt(c.Query("startIndex"), 10, 64)
	  if err != nil || page < 1{
		page  = 1
	  }

	  startIndex := (page - 1) * recordPerPage
      
	  matchStage := bson.D{{Key: "$match",   Value:  bson.D{}}}
      groupStage := bson.D{
		{Key : "$group", Value : bson.D{
			{Key : "_id", Value : nil},
			{Key: "total_count",Value : bson.D{{Key :"$sum", Value : 1}}},
			{Key : "data", Value : bson.D{{Key : "$push", Value :"$$ROOT"}}},
		}},
	  }
	  projectStage := bson.D{
		{
			Key  : "$project",Value : bson.D{
				{Key : "_id",Value:  0},
				{Key : "total_count", Value :1},
				{Key : "food_items", Value : bson.D{
					{Key  : "$slice", Value:  []interface{}{"$data", int(startIndex), int(recordPerPage)}},
				}},
			},
		},
	  }
	  result , err := foodCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, groupStage, projectStage,
	  })
      if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in getting from db"})
	  }
	  var allFoods []bson.M
	  if err = result.All(ctx, &allFoods); err != nil{
		log.Fatal(err)
	  }
	  c.JSON(http.StatusOK, allFoods);
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