package controllers

import (
	"context"
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

var orderCollection = database.OpenCollection(database.Client, "order")

func GetOrders()  gin.HandlerFunc{
	return func (c * gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second);
		defer cancel();
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil{
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil{
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		startindex := (page-1)*recordPerPage
        
		matchStage := bson.D{{
			Key:"$match",Value : bson.D{},
		}}
		groupStage := bson.D{{
			Key:"$group", Value: bson.D{
				{ Key : "_id", Value : bson.D{{
				  Key :  "_id",Value : nil, 	
				}}	},

				{
					Key: "totalCount", Value: bson.D{{
                    Key : "$sum", Value: 1,
					}},
				},
				{
					Key : "data", Value: bson.D{
						{Key :"$push", Value : "$$ROOT"},
					},
				},

			},
		}}

		projectStage := bson.D{{
		 Key : "$project", Value :bson.D{
			{
				Key : "_id", Value: 0,
			},
			{
				Key :"total_count",Value:1,
			},
			{
                Key : "orders", Value: bson.D{
					{Key :"$slice", Value :[]interface{}{"$data", startindex, recordPerPage}},
				},
			},

		 },
		}}
		result, err := orderCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})
		if err != nil{
			msg := "err in getting from db"
			c.JSON(http.StatusBadRequest, msg)
			return
		}
        c.JSON(http.StatusOK, result)
	}
}

func GetOrder() gin.HandlerFunc{
	return func ( c * gin.Context){
    ctx , cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	orderId := c.Param("order_id")
	var order = models.Order{}
	err := orderCollection.FindOne(ctx, bson.M{"order_id":orderId}).Decode(&order)
    if err != nil{
	    c.JSON(http.StatusInternalServerError, err.Error())
		return
	}  	
	c.JSON(http.StatusOK, order)
}
}

func CreateOrder() gin.HandlerFunc{
	return func ( c * gin.Context){
     ctx , cancel := context.WithTimeout(context.Background(), 100*time.Second)
	 defer cancel()
	 var order models.Order
	 err := c.BindJSON(&order); if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	 }
	 validateerr := validate.Struct(order)
	 if validateerr != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":validateerr.Error()})
		return
	 }
	 t := time.Now()
	 order.Created_at = &t
	 order.Updated_at = &t
	 order.ID = primitive.NewObjectID()
	 order.Order_id = order.ID.Hex()
	 result, err := orderCollection.InsertOne(ctx, order)
	 if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	 }
	 c.JSON(http.StatusOK, result)          
	}
}

func UpdateOrder() gin.HandlerFunc{
	return func ( c * gin.Context){
     ctx , cancel := context.WithTimeout(context.Background(), 100*time.Second)
	 defer cancel()
	 orderId := c.Param("order_id")
	 var order models.Order
	 err := c.BindJSON(&order); if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	 }
	 var updateobj primitive.D
	 if order.Table_id != nil{
       updateobj = append(updateobj, bson.E{Key :"order_id" , Value : order.Order_id})
	 }
	 t := time.Now()
	 order.Updated_at = &t
	 updateobj = append(updateobj, bson.E{Key : "updated_at" , Value : order.Updated_at})
     upsert := true
	 opt := options.UpdateOptions{
		Upsert: &upsert,
	 }
	 result , err := orderCollection.UpdateOne(ctx, bson.M{"order_id" : orderId} , bson.D{
		{
			Key : "$set", Value: updateobj,
		},
	 }, &opt)
	 
	 if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"err":"error in updating db"})
	    return  
	}
	c.JSON(http.StatusOK, result)
		
	}
}

func CreateOrderbyOrderItem(order models.Order) (string, error) {
	 ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	 defer cancel()
     order.ID = primitive.NewObjectID()
	 order.Order_id = order.ID.Hex()
	 validateerr := validate.Struct(order); if validateerr != nil{
		return "", validateerr
	 }
	 _, err := orderCollection.InsertOne(ctx, order);
	 if err != nil{
		return "", err
	 }
	 return order.Order_id, nil
}
