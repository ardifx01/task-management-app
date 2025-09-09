package utils

import (
	"os"
	"time"


	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(userID primitive.ObjectID)(string,error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID" : userID.Hex(),
		"exp" : time.Now().Add(time.Hour *24).Unix(),
	})

	tokenString,err := token.SignedString(jwtSecret)
	if err !=nil{
		return "",err
	}
	return tokenString,nil
}