package helpers

import (
	"context"
	"log"
	"restaurant-management/database"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetials struct{
	Email string
	FirstName string
	LastName string
	Uid string
	jwt.RegisteredClaims
}

var userCollection = database.OpenCollection(database.Client, "user");

var SECRET_KEY string = "kavascg"

func GenerateAllToken(email string, firstName string, lastName string, uid string)(string, string, error){
	claims := &SignedDetials{
		Email: email,
		FirstName:  firstName,
		LastName: lastName,
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
           ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour)),
		},
	}
	refreshClaims := &SignedDetials{
		RegisteredClaims: jwt.RegisteredClaims{
           ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour)),},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString([]byte(SECRET_KEY))
    if err != nil{
		log.Panic(err)
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims).SignedString([]byte(SECRET_KEY))
    if err != nil{
		log.Panic(err)
	}
	return token, refreshToken, err
}

func UpdateAlltokens(singedToken string, signedRefreshToken string, userId string){
	 ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	 defer cancel()
	 var updateObj primitive.D
	 updateObj = append(updateObj, bson.E{Key : "refresh_token", Value: signedRefreshToken})
	 upsert := true
	 opt := options.UpdateOptions{
		Upsert: &upsert,
	 }
	 filter := bson.M{
		"user_id":userId,
	 }
	 _, err := userCollection.UpdateOne(ctx, filter, bson.D{{Key : "$set", Value: updateObj},}, &opt)
     if err != nil{
		log.Panic(err)
		return
	 }
	}


func ValidateToken (singedToken string)(claims *SignedDetials, msg string){
     token , err := jwt.ParseWithClaims(
		singedToken,
		&SignedDetials{},
		func(token *jwt.Token)(interface{}, error){
             return []byte(SECRET_KEY), nil
		},
	 )

	 claims , ok := token.Claims.(*SignedDetials)
	 if !ok {
         msg = "the token is invalid"
		 msg = err.Error()
		 return 
	 }

	 if claims.ExpiresAt.Time.Before(time.Now()) {
		msg = "token is expired"
		msg = err.Error()
		return
	 }
	 return claims , msg

}

