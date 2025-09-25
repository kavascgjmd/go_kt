package controllers

import (
	"context"
	"net/http"
	"restaurant-management/database"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tableCollection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc{
	return func( c * gin.Context){
     ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	 defer cancel()
	 result, err :=  tableCollection.Find(ctx, bson.M{}) 
	 if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"err":err.Error()})
		return
	 }
	 var allresult [] bson.M
	 err = result.All(ctx, &allresult)
	 if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"err":err.Error()})
		return
	 }
	 c.JSON(http.StatusOK, allresult)
	}
}
func GetTable() gin.HandlerFunc{
	return func( c * gin.Context){
	  ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
      defer cancel()
	  tableId := c.Param("table_id")
	  var table models.Table
	  err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
	  if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()});
		return
	  }
	  c.JSON(http.StatusOK, table)
	}
}
func CreateTable() gin.HandlerFunc{
	return func( c * gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var table models.Table
        err := c.BindJSON(&table); if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
			return
		}
		validaterr := validate.Struct(table)
		if validaterr != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"err":validaterr.Error()})
			return
		}
		t := time.Now()
		table.Created_at = &t;
		table.Updated_at = &t
		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()
		result , err := tableCollection.InsertOne(ctx, table)
        if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"err":err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
		
	}
}
func UpdateTable() gin.HandlerFunc{
	return func( c * gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		tableId := c.Param("table_id")
		var table models.Table
		err := c.BindJSON(&table); if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return
		}
		var updateobj primitive.D
		if table.Table_number != nil{
			updateobj = append(updateobj, bson.E{Key:"table_number" , Value: table.Table_number})
		}
		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result , err := tableCollection.UpdateOne(ctx, bson.M{"table_id":tableId}, bson.D{
			{
				Key: "$set",
				Value: updateobj,
			},
		}, &opt)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
