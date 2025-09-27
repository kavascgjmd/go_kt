package middleware

import (
	"net/http"
	"restaurant-management/helpers"

	"github.com/gin-gonic/gin"
)

func Authenticaiton() gin.HandlerFunc{
	return func(c *gin.Context){
		clientToken := c.Request.Header.Get("token")
		if clientToken == ""{
			c.JSON(http.StatusBadRequest, gin.H{"error" : "invalid token"})
			return
		}
		claims, err := helpers.ValidateToken(clientToken)
		if err != ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("user_id", claims.Uid)
		c.Next()
	}
}