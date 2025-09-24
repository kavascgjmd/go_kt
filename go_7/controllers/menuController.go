package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaurant-management/database"
	"restaurant-management/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var validate = validator.New()
func GetMenus() gin.HandlerFunc{
	return func(c * gin.Context){
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel();
	result , err := menuCollection.Find(ctx, bson.M{}) //bsom.M is map[string]interface{} how mongo store data, here it passed for no filters 
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to find menues in mongodb"})
	}
	var allMenus []bson.M 
	if err = result.All(ctx, &allMenus); err != nil{ // result is mongo.Cursor so it need to be converted into something that is go, to be converted to json
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenu() gin.HandlerFunc{
	return func(c * gin.Context){
       ctx, cancel := context.WithTimeout(context.Background() , 100*time.Second)
	   defer cancel()
	   menuId := c.Param("menu_id")
	   var menu = &models.Menu{} 
	   err := menuCollection.FindOne(ctx, bson.M{"menu_id":menuId}).Decode(menu)
	   if err != nil{
		log.Fatal(err)
	   }
	   c.JSON(http.StatusOK,menu)
	}
}

func CreateMenu() gin.HandlerFunc{
	return func(c * gin.Context){
      ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	  defer cancel();
	  var menu models.Menu
	  err := c.BindJSON(&menu);
	  if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return 
	  }

	  validationErr := validate.Struct(menu)
	  if validationErr != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error" : validationErr.Error()})
		return 
	  }
      t := time.Now()
	  menu.Created_at = &t
	  menu.Updated_at = &t
	  menu.ID = primitive.NewObjectID()
      menu.Menu_id = menu.ID.Hex()

	  result, insertErr := menuCollection.InsertOne(ctx, menu)
	  if insertErr != nil{
		msg := fmt.Sprintf("Menu item is not created")
		c.JSON(http.StatusInternalServerError, msg)
        return
	}
	c.JSON(http.StatusOK , result)
	}
}

func UpdateMenu() gin.HandlerFunc{
	return func(c * gin.Context){

		var ctx , cancel = context.WithTimeout(context.Background(), 100*time.Second)
        var menu models.Menu
		err := c.BindJSON(&menu); if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
        menuId := c.Param("menu_id")
		filter := bson.M{"menu_id":menuId}

		var updateObj primitive.D // this is mongo_document type primitive.D // that is a slice of bson.E []bson.E

		if menu.Start_Date != nil && menu.End_Date != nil{
			if !inTimeSpan(menu.Start_Date, menu.End_Date, time.Now()){
				msg := "kindly retype the time"
				c.JSON(http.StatusBadRequest, gin.H{"error":msg})
				defer cancel()
				return
			}

			updateObj = append(updateObj, bson.E{Key : "start_date",Value  : menu.Start_Date}) // bson.E is a type of struct of {key:string, value:interface{}}
	        updateObj = append(updateObj, bson.E{Key : "end_date", Value : menu.End_Date})
		}
		if menu.Name != ""{
			updateObj = append(updateObj, bson.E{ Key : "name",Value :menu.Name})
		}
		if menu.Category != "" {
			updateObj = append(updateObj, bson.E{Key : "category", Value : menu.Category})
		}
        t := time.Now();
		menu.Updated_at = &t;
        updateObj = append(updateObj, bson.E{Key : "updated_at",Value : menu.Updated_at});
		upsert := true // this is for if feild doesn't exist then create
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := menuCollection.UpdateOne(
			ctx, 
			filter,
			bson.D{ // this is also a slice of bson.E
				{Key : "$set", Value : updateObj}, // if i just pass updateObj than non-existing fields will be deleted 
			},
			&opt,

		)

		if err != nil {
			msg := "Menu updation failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
		}
        c.JSON(http.StatusOK, result)
	}
}

func inTimeSpan(start *time.Time, end *time.Time, check time.Time) bool{
	return start.After(check) && end.After(*start)
}
