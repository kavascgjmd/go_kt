package controllers

import (
	"context"
	"log"
	"net/http"
	"restaurant-management/database"
	"restaurant-management/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
      ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second);
	  defer cancel()
	  var food models.Food
	  err := c.BindJSON(&food); if err != nil{
		c.JSON(http.StatusBadRequest, err.Error())
		return
	  }
	  t := time.Now()
	  food.Created_at = &t;
      food.Updated_at = &t;
	  food.ID = primitive.NewObjectID()
	  food.Food_id = food.ID.Hex()

	  result, err := foodCollection.InsertOne(ctx, food);
	  if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	  }
      c.JSON(http.StatusOK, result);
	}
}

func UpdateFood() gin.HandlerFunc{
	return func(c * gin.Context){
      var food models.Food;
	   err := c.BindJSON(&food);if err  != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error" : err.Error()})
		return 
	  }
      ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	  defer cancel()
	  foodId := c.Query("food_id")
	  var updateobj primitive.D
	  if food.Food_image != nil{
          updateobj = append(updateobj, bson.E{Key : "food_image", Value: food.Food_image })
	  }
      if food.Menu_id != nil{
          updateobj= append(updateobj, bson.E{Key : "menu_id" , Value: food.Menu_id})
	  }
	  if food.Food_image != nil{
          updateobj = append(updateobj, bson.E{Key: "food_image", Value: food.Food_image})
	  }
	  if food.Price != nil{
          updateobj = append(updateobj, bson.E{Key : "price", Value: food.Price})
	  }
	  t := time.Now();
	  food.Updated_at = &t;
	  updateobj = append(updateobj, bson.E{Key: "updated_at", Value: food.Updated_at})
      upsert := true 
	  opt := options.UpdateOptions{
		Upsert: &upsert,
	  }

	  result, err := foodCollection.UpdateOne(
		ctx, bson.M{"food_id" : foodId} ,bson.D{
			{
				Key:"$set",Value:  updateobj,
		}} , &opt,
	  )

	  if err != nil{
		msg := "Food upation failed"
		c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
	  }
	  c.JSON(http.StatusOK, result);
	}
}