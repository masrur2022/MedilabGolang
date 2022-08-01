package jwtauthentication

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var myKey = []byte("sekKey")

func GenerateToken(c *gin.Context,email string) string {
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(1 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		fmt.Fprintf(c.Writer, err.Error())
	}
	return tokenString
}

func Velidation(c *gin.Context) {
	// We can obtain the session token from the requests cookies, which come with every request
	cookie, err := c.Request.Cookie("token")
	if err != nil {
		c.JSON(400,gin.H{
			"status":"COOKIE_DOES_NOT_EXIST",
		})
	}
	// Get the JWT string from the cookie
	tknStr := cookie.Value
	// Initialize a new instance of `Claims`
	claims := &Claims{}
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	token, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		fmt.Fprintf(c.Writer, "error 2")
	}
	if !token.Valid {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(c.Writer, "error 3")
	}
}