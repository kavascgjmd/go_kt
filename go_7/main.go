package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"restaurant-management/routes"
	"restaurant-management/middleware"
	"restaurant-management/database"
)

func main(){
	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)

    router.Use(middleware.Authentication())
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	router.Run(":" + port);
	
	
}