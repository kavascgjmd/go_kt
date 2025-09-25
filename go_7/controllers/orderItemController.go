package controllers

import (
	"context"
	"net/http"
	"time"

	"restaurant-management/database"
	"restaurant-management/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderitemPack struct{
	Table_id * string
	OrderItemList []models.OrderItem
}

var orderItemCollection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc{
	return func( c * gin.Context){
       ctx, cancel := context.WithTimeout(context.Background(),100*time.Second )
	   defer cancel()
	   result, err := orderItemCollection.Find(ctx, bson.M{})
	   if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	   }
	   var allOrderItems []bson.M
	   err = result.All(ctx, &allOrderItems)
	   if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return		
	   }
	   c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItem() gin.HandlerFunc{
	return func( c * gin.Context){
     ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	 defer cancel()
	 orderItemId := c.Param("order_item_id")
	 var orderItem models.OrderItem
	 err :=  orderCollection.FindOne(ctx, bson.M{"order_item_id":orderItemId}).Decode(&orderItem) 
     if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	 }
	 c.JSON(http.StatusOK, orderItem)
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc{
	return func(c * gin.Context){
		
	}
}

func CreateOrderItem() gin.HandlerFunc{
	return func( c * gin.Context){
       ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	   defer cancel()
	   var orderItemPack OrderitemPack
	   var order models.Order
	   err := c.BindJSON(&orderItemPack); if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
		return
	   }
	   var t = time.Now()
	   order.Created_at = &t
	   order.Updated_at = &t
	   order.Table_id = orderItemPack.Table_id
	   order.Order_id , err = CreateOrderbyOrderItem(order)
	   if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	   }
	   for _, orderItem := range orderItemPack.OrderItemList{
          orderItem.Created_at = &t
		  orderItem.Updated_at = &t
		  orderItem.Order_id = order.Order_id
		  orderItem.ID = primitive.NewObjectID()
		  val := orderItem.ID.Hex()
		  orderItem.Order_item_id = &val
		  _, err := orderItemCollection.InsertOne(ctx, orderItem)
		  if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		    return
		  }
	   }
	   c.JSON(http.StatusOK, "Successful");
	}
}


func UpdateOrderItem() gin.HandlerFunc{
	return func( c * gin.Context){
        ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	    defer cancel()
	    var orderItem models.OrderItem
	   	err := c.BindJSON(&orderItem); if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		}
	    orderitemid := c.Param("orderitem_id")
	    val := orderItemCollection.FindOne(ctx, bson.M{"orderitem_id":orderitemid});
	    if val == nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"cant find orderitem"})
		return
		}
	   var updateobj primitive.D
	   if(orderItem.Food_id != nil){
          updateobj = append(updateobj, bson.E{Key:"food_id" , Value: orderItem.Food_id})
	   }
	   t := time.Now()
	   orderItem.Updated_at = &t;
	   updateobj = append(updateobj, bson.E{Key : "updated_at", Value: orderItem.Updated_at})
	   result, err := orderItemCollection.UpdateOne(ctx, bson.M{"order_item_id" : orderitemid}, bson.D{
		{
			Key:"$set",
			Value: updateobj,
		},
	   })
	   if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	   }
        c.JSON(http.StatusOK, result)       
	}
}



