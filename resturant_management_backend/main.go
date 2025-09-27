package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"restaurant-management/routes"
	"restaurant-management/middleware"
)

func main(){
	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)

    router.Use(middleware.Authenticaiton())
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	router.Run(":" + port);
	
	
}