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
	"go.mongodb.org/mongo-driver/mongo"
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
		order_id := c.Param("order_id")
		allOrderItems , err := ItemsByOrder(order_id)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		    return
		}
		c.JSON(http.StatusOK, allOrderItems)
	}
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error){
	 ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	 defer cancel();

	 matchStage := bson.D{{Key : "$match", Value : bson.D{{Key : "order_id" ,Value: id}}}}
	 lookupStage := bson.D{{Key : "$lookup", Value : bson.D{{Key : "from", Value : "food"} , {Key : "localField", Value : "food_id"}, {Key : "foreignField", Value:  "food_id"}, {Key: "as",Value:  "food"}}}}
     unwindStage := bson.D{{Key : "$unwind", Value : bson.D{{Key :"path", Value: "$food"}, {Key : "preserveNullAndEmptyArrays", Value : true}}}}

	 lookUpOrderStage := bson.D{{Key : "$lookup", Value: bson.D{{Key : "from", Value : "order"}, {Key : "localField", Value : "order_id" }, {Key : "foreignField", Value: "order_id"}, {Key : "as", Value : "order"}}}}
     unwindOrderStage := bson.D{{Key : "$unwind", Value : bson.D{{Key : "path", Value : "$order"}, {Key : "preserveNullAndEmptyArrays", Value: true}}}}

	 lookUpTableStage := bson.D{{Key : "$lookup", Value: bson.D{{Key : "from", Value: "table"}, {Key : "localField", Value :"order.table_id"}, {Key :"foreignField", Value: "table_id"}, {Key : "as", Value: "table"}}}}
     unwindTableStage := bson.D{{Key : "$unwind",Value:  bson.D{{Key : "path",Value :  "$table"}, {Key  : "preserveNullAndEmptyArrays",Value : true}}}}
     
	 projectStage := bson.D{
		{Key : "$project", Value : bson.D{
			{Key : "id",Value: 0},
			{Key :"amount", Value : "$food.price"},
			{Key : "food_name", Value :"$food.name"},
			{Key : "food_image",Value : "$food.food_image"},
			{Key : "table_number", Value :"$table.table_number"},
			{Key : "table_id", Value: "$table.table_id"},
			{Key : "order_id", Value: "$order.order_id"},
			{Key : "price", Value : "$food.price"},
		},
		},
	 }

     groupStage := bson.D{{Key : "$group", Value: bson.D{{Key : "_id", Value : bson.D{{Key:  "order_id", Value:  "$order_id"}, { Key : "table_id",Value:  "$table_id"}, {Key : "table_number", Value :"$table_number"}}}, {Key : "payment_due",Value : bson.D{{Key : "$sum", Value:  "$amount"}}}, {Key : "orderitems", Value :  bson.D{{Key : "$push", Value :"$$ROOT"}}}  }}}
	 projectStage2 := bson.D{
		{Key : "$project", Value:  bson.D{
			{Key : "id", Value: 0	},
			{Key : "payment_due", Value: 1},
			{Key : "table_number", Value :"$_id.table_number"},
			{Key : "order_items", Value: 1 },
		}},
	 }
     result, err :=  orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookUpOrderStage,
		unwindOrderStage,
		lookUpTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	 })

	 if err != nil{
		return OrderItems, err
	 }

	 if err = result.All(ctx, &OrderItems); err != nil{
		return OrderItems, err
	 }
     return OrderItems, err
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



