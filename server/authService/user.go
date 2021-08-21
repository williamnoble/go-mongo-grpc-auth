package authService

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)
var (
	jwtSecret = []byte("KJ34e98saD34AS232£K@£")
)

type User struct {
	ID primitive.ObjectID `bson:"_id"` // ignore
	Username string `bson:"username"`
	Email string `bson:"email"`
	Password string `bson:"password"`
}

func (u User) UserToJWTToken() string{
	js, err := json.Marshal(u)
	if err != nil {
		log.Println("Failed to marshal user")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": string(js),
	})
	// Make sure jwtSecret is a []byte not a string!
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("Failed to generate tokenString for user")
	}
	return tokenString
}

func UserFromJWTToken(token string) User {
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	var user User
	data := []byte(claims["data"].(string))
	json.Unmarshal(data, &user)
	return user
}

