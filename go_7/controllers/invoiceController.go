package controllers

import (
	"context"
	"net/http"
	"restaurant-management/database"
	"restaurant-management/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
)

var invoiceCollection = database.OpenCollection(database.Client, "invoice")

func GetInvoices() gin.HandlerFunc{
	return func(c * gin.Context){
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second )
	defer cancel()
    result , err := invoiceCollection.Find(ctx, bson.M{})
	var allresult []bson.M
	err = result.All(ctx, allresult)
	if err != nil{
	   c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
       return
	}
	c.JSON(http.StatusOK , allresult)
	}
}

func GetInvoice() gin.HandlerFunc{
	return func(c * gin.Context){
    ctx, cancel:= context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	invoiceId := c.Param("invoice_id")
	var invoice models.Invoice
	err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id":invoiceId}).Decode(&invoice)
    if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()});
		return
	}
	c.JSON(http.StatusOK, invoice);
	}
}

func CreateInvoice() gin.HandlerFunc{
	return func(c * gin.Context){
    ctx, cancel:= context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var invoice models.Invoice
	err := c.BindJSON(&invoice) ; if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}
	t := time.Now()
	invoice.Created_at = &t
	invoice.Updated_at = &t
	invoice.ID = primitive.NewObjectID()
	invoice.Invoice_id = invoice.ID.Hex()
	result , err := invoiceCollection.InsertOne(ctx, invoice);
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured which creating invoice"})
	    return 
	}
	c.JSON(http.StatusOK, result)
	}
}

func UpdateInvoice() gin.HandlerFunc{
	return func(c * gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoice models.Invoice
		invoiceId := c.Param("invoice_id")
		err := c.BindJSON(&invoice) ;if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
			return
		}
		var updateobj primitive.D
		if invoice.Payment_Method != nil{
           updateobj = append(updateobj, bson.E{Key :"payment_method", Value: invoice.Payment_Method})
		}
		if invoice.Payment_Status != nil{
           updateobj = append(updateobj, bson.E{Key :"payment_status", Value: invoice.Payment_Status})
		}
		if invoice.Payment_due_date != nil{
           updateobj = append(updateobj, bson.E{Key :"payment_due_date", Value: invoice.Payment_due_date})
		}
		upsert:= true
        opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result , err := invoiceCollection.UpdateOne(ctx, bson.M{"invoice_id":invoiceId}, bson.D{
			{
				Key: "$set", Value: updateobj,
			},
		}, &opt)
		
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"err":"error in updating invice"});
			return
		}
        c.JSON(http.StatusOK, result)
	}
}