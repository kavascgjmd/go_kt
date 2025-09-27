package controllers

import (
	"context"
	"log"
	"net/http"
	"restaurant-management/database"
	"restaurant-management/models"
	"restaurant-management/helpers"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

var userCollection = database.OpenCollection(database.Client, "user")


func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context){
     ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	 defer cancel()
	 recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
	 if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"errror" : err.Error()})
		return
	 }
     page, err := strconv.Atoi(c.Query("page"))
	 if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"errror" : err.Error()})
		return
	 }
	 startIndex := (page - 1)*recordPerPage;
	 matchStage := bson.D{{Key : "$match", Value : bson.D{}}}
	 projectStage := bson.D{
		{Key : "$project"},
		{Key : "_id",Value :  0},
		{Key : "user_items", Value:  bson.D{{Key :"$slice" , Value : []interface{}{"$data", startIndex, recordPerPage}}}},
	 }

     result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage, projectStage,
	 })
	 if err != nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() })
		return
	 }
	 c.JSON(http.StatusOK, result)
	}
}


func GetUser() gin.HandlerFunc{
	return func(c * gin.Context){
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	userId := c.Param("user_id")
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"user_id":userId}).Decode(&user)
    if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
	}
}

func SignUp() gin.HandlerFunc{
	return func(c * gin.Context){
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel() 
	var user models.User
	err := c.BindJSON(&user); if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	err = validate.Struct(user); if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	count , err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil || count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id already exist"})
		return
	}
	hashpassword := HashPassword(user.HashPassword)
    user.HashPassword = hashpassword
	t := time.Now()
	user.Created_at = &t
	user.Updated_at = &t
	_, refreshToken , err := helpers.GenerateAllToken(user.Email, user.FirstName, user.LastName, user.User_id )
    if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.RefreshToken = &refreshToken
	result, err := userCollection.InsertOne(ctx, user) ; if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
	}
}

func Login() gin.HandlerFunc{
	return func(c * gin.Context){
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel() 
	var user models.User
	var foundUser models.User
	err := c.BindJSON(&user); if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userId := c.Param("user_id")
	err = userCollection.FindOne(ctx, bson.M{"user_id":userId}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"No such user"})
		return
	}
    
	passwordIsValid , msg := VerfiyPassword(foundUser.HashPassword, user.HashPassword)
	if !passwordIsValid {
		c.JSON(http.StatusBadRequest, gin.H{"error" : msg})
		return
	}
	token , refreshToken, _ := helpers.GenerateAllToken(foundUser.Email , foundUser.FirstName, foundUser.LastName, foundUser.User_id)
    helpers.UpdateAlltokens(token, refreshToken, userId)
	c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func HashPassword(password string) string{
    bytes , err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil{
		log.Panic(err)
	}
	return string(bytes)
}

func VerfiyPassword(userPassword string, providedPassword string)( bool, string){
      err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword));
	  check := true 
	  msg := ""
	  if err != nil{
         msg = "login or password is incorrect "
		 check = false
	  }
	  return check , msg
}