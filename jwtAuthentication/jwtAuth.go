package jwtauthentication

import (
	"encoding/json"
	"fmt"
	"medcard/beck/mongoconnection"
	"medcard/beck/structures"
	"net/http"
	"medcard/beck/converter"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type ClaimsTok struct {
	Login string `json:"login"`
	Id string `json:"id"`
	jwt.StandardClaims
}

var myKey = []byte("sekKey")

func GenerateToken(c *gin.Context,login string) string {
	// explore he db tofind user id
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	collection := ClientG.Database("MedCard").Collection("users")
	var DbgetUser structures.ClientLog
	collection.FindOne(CtxG, bson.M{"login": login}).Decode(&DbgetUser)
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(30 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &ClaimsTok{
		Login: login,
		Id: DbgetUser.ID,
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

func Velidation(c *gin.Context) ClaimsTok {
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
	claims := &ClaimsTok{}
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

	payloadLogin := GetAccessDetails(tknStr)
	jsonStr := string(payloadLogin)
	fmt.Println(jsonStr)
	fmt.Println("first")
	// loginString := strings.Split(strings.Split(jsonStr, ":")[2], `"`)[1]
	userId := strings.Split(strings.Split(jsonStr, ":")[2], `"`)[1]
	loginString := strings.Split(strings.Split(jsonStr, ":")[3], `"`)[1]

	fmt.Println(userId)
	fmt.Println(strings.Split(strings.Split(jsonStr, ":")[3], `"`)[1])
	fmt.Printf("1%v\n",loginString)

	var ClaimsObj ClaimsTok
	ClaimsObj.Id = userId
	ClaimsObj.Login = loginString

	//"""""""""""""" Online check if the user is online or not""""""""""""""
	onlineCkeck(claims.Login)

	return ClaimsObj
}
func GetAccessDetails(tokenStr string) ([]byte) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		 // check token signing method etc
		 return []byte("sekKey"), nil
	})

	if err != nil {
		return nil
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid{
		fmt.Println(claims)
		ClaimsTok, err := json.Marshal(claims)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}
		return ClaimsTok
	}else{
		return nil
	}
}
func onlineCkeck(login string){
	//  """""""""DB connection to get the user """""""""
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	//""""""""""""""""" connection to online DB"""""""""""""""""
	collection := ClientG.Database("MedCard").Collection("users")
	var Dbgetonline structures.Users
	collection.FindOne(CtxG, bson.M{"login": login}).Decode(&Dbgetonline)
	
	//"""""""""""""""" compare the result arreys and add them into the `online` DB""""""""""""""""
	date := strings.Split(time.Now().String(), " ")[0]
	hour := strings.Split(strings.Split(time.Now().String(), " ")[1], ":")[0]
	minutes := strings.Split(strings.Split(time.Now().String(), " ")[1], ":")[1]
	holeDate := date +":"+hour+":"+minutes
	//"""""""""""""""""""""""" converter function recolling""""""""""""""""""""""""
	converter.Convert()
	// """""""""defining the difference"""""""""
	spliteDate := strings.Split(Dbgetonline.LastActive, ":")
	hourInt, err := strconv.Atoi(hour)
	minutesInt, err := strconv.Atoi(minutes)
	hourIntDb, err := strconv.Atoi(spliteDate[1])
	fmt.Printf("hourIntDb: %v\n", hourIntDb)
	minutesIntDb, err := strconv.Atoi(spliteDate[2])
	if err != nil{
		fmt.Println("Its not a string")
	}
	diference := ((hourInt * 60)+minutesInt) - ((hourIntDb * 60)+minutesIntDb)
	// """""""""check online process"""""""""
	if Dbgetonline.Login == ""{
		fmt.Println("not found")
	}else{
		if spliteDate[0] != date{
			fmt.Println("1")
			Dbgetonline.LastActive = holeDate
			Dbgetonline.UserStatus = "online"
			collection.DeleteOne(CtxG,bson.M{"login":login})
			collection.InsertOne(CtxG,Dbgetonline)
			go Timer(login)
		}else{
				if diference < 2{
					fmt.Println("2")
					Dbgetonline.LastActive = holeDate
					Dbgetonline.UserStatus = "online"
					collection.DeleteOne(CtxG,bson.M{"login":login})
					collection.InsertOne(CtxG,Dbgetonline)
					go Timer(login)
				}else{
					fmt.Println("3")
					Dbgetonline.LastActive = holeDate
					Dbgetonline.UserStatus = "online"
					collection.DeleteOne(CtxG,bson.M{"login":login})
					collection.InsertOne(CtxG,Dbgetonline)
					go Timer(login)
				}
		}
	}
}
func Timer(login string){	
	//  """""""""DB connection to get the user """""""""
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG

	collection := ClientG.Database("MedCard").Collection("users")
	var Dbgetonline structures.Users
	collection.FindOne(CtxG, bson.M{"login": login}).Decode(&Dbgetonline)

	timerOne := time.NewTimer(10 * time.Minute)
	<-timerOne.C
	//"""""""""""""""" compare the result arreys and add them into the `online` DB""""""""""""""""
	date := strings.Split(time.Now().String(), " ")[0]
	hour := strings.Split(strings.Split(time.Now().String(), " ")[1], ":")[0]
	minutes := strings.Split(strings.Split(time.Now().String(), " ")[1], ":")[1]
	holeDate := date +":"+hour+":"+minutes
	//""""""""""""""" inserting data intodatabase"""""""""""""""
	Dbgetonline.UserStatus = "offline"
	Dbgetonline.LastActive = holeDate
	collection.DeleteOne(CtxG,bson.M{"login":login})
	collection.InsertOne(CtxG,Dbgetonline)
}