package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	// "os"
	"time"

	// "net/smtp"
	"strings"

	// "github.com/dgrijalva/jwt-go"
	"medcard/beck/bycrypt"
	"medcard/beck/converter"
	jwtauthentication "medcard/beck/jwtAuthentication"
	"medcard/beck/mongoconnection"
	"medcard/beck/structures"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

var globeUserLog structures.DoctorLog
var FAQ structures.FrequentlyAskedQuestion
var User structures.Users

var ClientG *mongo.Client
var CtxG context.Context

func Login(c *gin.Context) {
	c.ShouldBindJSON(&globeUserLog)
	fmt.Println(globeUserLog)
	//""""""" mongoDb connection"""""""
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	var DBgetUser structures.DoctorLog
	collection := ClientG.Database("MedCard").Collection("users")
	collection.FindOne(CtxG, bson.M{"login": globeUserLog.Login}).Decode(&DBgetUser)
	// """"""""""""""""""""""compare the password with bycrypt""""""""""""""""""""""
	compareResult := bycrypt.CompareHashPasswords(DBgetUser.Password, globeUserLog.Password)

	if compareResult != true {
		c.Writer.WriteHeader(http.StatusUnauthorized)
	} else {
		// Finally, we set the client cookie for "token" as the JWT we just generated
		// we also set an expiry time which is the same as the token itself
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "token",
			Value:    jwtauthentication.GenerateToken(c, globeUserLog.Email),
			Expires:  time.Now().Add(1 * time.Minute),
			HttpOnly: false,
			Secure:   true,
			Path:     "/",
		})
		// """""""""""""""""""""""""send the login for user"""""""""""""""""""""""""
		c.JSON(200, gin.H{
			"LOGIN": DBgetUser.Login,
		})
		// """""""""""""""""""""connect to online collection to add user and make him online"""""""""""""""""""""
		var online structures.Online
		collectionOnline := ClientG.Database("MedCard").Collection("online")
		online.Login = DBgetUser.Login
		online.Permissions = DBgetUser.Permissions
		// """""""""""""insert data into db"""""""""""""
		collectionOnline.InsertOne(CtxG, online)
	}
}
func SignupDr(c *gin.Context) {
	// jwtauthentication.Velidation(c)
	// """""""""get he json request from client """""""""
	// c.ShouldBindJSON(&globeUserLog)
	converter.Convert("doctors")
	jsonFM := c.Request.FormValue("json")
	files, handler, errIMG := c.Request.FormFile("img")
	// """""""""""""""""""""""check The file on existense"""""""""""""""""""""""
	if errIMG != nil {
		c.JSON(409, gin.H{
			"sttus": "NOIMGFILEEXIST",
		})
	}

	files.Seek(23, 23)
	fmt.Printf("File Name %s\n", handler.Filename)
	// """""""""""""""""""""bind the request data into structure"""""""""""""""""""""
	json.Unmarshal([]byte(jsonFM), &globeUserLog)
	// """"""""""""""""""check the login velidation""""""""""""""""""
	compiledLogin := bycrypt.ChecktheLogin(globeUserLog.Login)
	//""""" check if any field is ampty"""""
	if globeUserLog.Possition == "" || globeUserLog.Login == "" || globeUserLog.Name == "" || globeUserLog.Password == "" || globeUserLog.Sername == "" || errIMG != nil {
		c.JSON(400, gin.H{
			"status": "NOTCOMLETE",
		})
	} else {
		//"""""" connect to data base to get the user and verify it""""""
		mongoconnection.MongoDB()
		CtxG := mongoconnection.CtxG
		ClientG := mongoconnection.ClientG
		collection := ClientG.Database("MedCard").Collection("doctors")
		var DbgetUser structures.DoctorLog
		collection.FindOne(CtxG, bson.M{"login": compiledLogin}).Decode(&DbgetUser)

		collectionUser := ClientG.Database("MedCard").Collection("users")
		var DbgetDoc structures.DoctorLog
		collectionUser.FindOne(CtxG, bson.M{"login": compiledLogin}).Decode(&DbgetDoc)
		// permission check
		var accepted = Permission(c, globeUserLog.RequestLogin)
		if accepted != "admin" {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"status": "NO_PERMISSION_TO_DO_THIS",
			})
		} else {
			//"""""""""""""""""" find out if there any user with such logo""""""""""""""""""
			if DbgetUser.Login != "" || DbgetDoc.Login != ""{
				c.Writer.WriteHeader(http.StatusBadRequest)
				c.JSON(http.StatusExpectationFailed, gin.H{
					"status": "EXIST",
				})
			} else {
				// """"""""""""hesh the password befor inserting it into databse""""""""""""
				hashedPassword, err := bycrypt.HashPassword(globeUserLog.Password)
				if err != nil {
					fmt.Fprintf(c.Writer, err.Error())
					return
				}
				globeUserLog.Password = hashedPassword
				globeUserLog.Login = compiledLogin
				globeUserLog.PrimitiveID = primitive.NewObjectID()
				globeUserLog.ID = globeUserLog.PrimitiveID.Hex()
				// """"""""""""""""""""get the img and bind the url to veriable bellow""""""""""""""""""""
				globeUserLog.ProfileImage = bycrypt.ParseFile(c, "static/uploadImg")
				globeUserLog.Permissions = "doctors"
				globeUserLog.RequestLogin = " "
				fmt.Println(hashedPassword)
				// """""""""""""""Insert the user ino dataBase if its valid"""""""""""""""
				collection.InsertOne(CtxG, globeUserLog)
				converter.Convert("doctors")
			}
		}
	}
}
func SignupCl(c *gin.Context) {
	// jwtauthentication.Velidation(c)
	converter.Convert("clients")
	var globeUserLog structures.ClientLog
	// """""""""get he json request from client """""""""
	c.ShouldBindJSON(&globeUserLog)
	fmt.Println(globeUserLog)
	// """"""""""""""""""check the login velidation""""""""""""""""""
	compiledLogin := bycrypt.ChecktheLogin(globeUserLog.Login)
	//""""" check if any field is ampty"""""
	if globeUserLog.Blood == "" || globeUserLog.Login == "" || globeUserLog.Name == "" || globeUserLog.Password == "" || globeUserLog.Sername == "" || globeUserLog.Gender == "" || globeUserLog.LastName == "" || globeUserLog.RequestLogin == "" {
		c.JSON(400, gin.H{
			"status": "NOT_COMLETE",
		})
	} else {
		//"""""" connect to data base to get the user and verify it""""""
		mongoconnection.MongoDB()
		CtxG := mongoconnection.CtxG
		ClientG := mongoconnection.ClientG
		collection := ClientG.Database("MedCard").Collection("clients")
		var DbgetUser structures.ClientLog
		collection.FindOne(CtxG, bson.M{"login": compiledLogin}).Decode(&DbgetUser)

		collectionUser := ClientG.Database("MedCard").Collection("clients")
		var DbgetDoc structures.ClientLog
		collectionUser.FindOne(CtxG, bson.M{"login": compiledLogin}).Decode(&DbgetDoc)
		// permission check
		var accepted = Permission(c, globeUserLog.RequestLogin)
		if accepted != "admin" {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"status": "NO_PERMISSION_TO_DO_THIS",
			})
		} else {
			//"""""""""""""""""" find out if there any user with such logo""""""""""""""""""
			if DbgetUser.Login != "" {
				c.Writer.WriteHeader(http.StatusBadRequest)
				c.JSON(http.StatusExpectationFailed, gin.H{
					"status": "EXIST",
				})
			} else {
				// """"""""""""hesh the password befor inserting it into databse""""""""""""
				hashedPassword, err := bycrypt.HashPassword(globeUserLog.Password)
				if err != nil {
					fmt.Fprintf(c.Writer, err.Error())
					return
				}
				globeUserLog.Password = hashedPassword
				globeUserLog.Login = compiledLogin
				globeUserLog.PrimitiveID = primitive.NewObjectID()
				globeUserLog.ID = globeUserLog.PrimitiveID.Hex()
				// """"""""""""""""""""get the img and bind the url to veriable bellow""""""""""""""""""""
				globeUserLog.Permissions = "client"
				fmt.Println(hashedPassword)
				// """""""""""""""Insert the user ino dataBase if its valid"""""""""""""""
				collection.InsertOne(CtxG, globeUserLog)
				converter.Convert("clients")
			}
		}
	}
}
func AddQuestion(c *gin.Context) {
	// """"""""""""""verify does the user has cookie or not""""""""""""""
	// jwtauthentication.Velidation(c)
	// """"""""""""""bind the json resived from user to structure""""""""""""""
	c.ShouldBindJSON(&FAQ)
	// """""""""""""""""""""""check the verbles are they filled out """""""""""""""""""""""
	if FAQ.Description == "" || FAQ.Title == "" {
		c.JSON(400, gin.H{
			"status": "NOTCOMLETE",
		})
	} else {
		//""""""" mongoDb connection"""""""
		mongoconnection.MongoDB()
		CtxG := mongoconnection.CtxG
		ClientG := mongoconnection.ClientG
		var DBpushQuestion structures.FrequentlyAskedQuestion
		collection := ClientG.Database("MedCard").Collection("frequently_asked_question")
		collection.FindOne(CtxG, bson.M{"title": FAQ.Title}).Decode(&DBpushQuestion)
		fmt.Println(DBpushQuestion)
		fmt.Println(FAQ)
		//"""""""""""""" check the question is there anything like this""""""""""""""
		if DBpushQuestion.Title != "" {
			c.Writer.WriteHeader(http.StatusBadRequest)
			c.JSON(http.StatusExpectationFailed, gin.H{
				"status": "EXIST",
			})
		} else {
			var accept = Permission(c, FAQ.RequestLogin)
			if accept == "doctors" || accept == "admin" {
				titleArr := strings.Split(FAQ.Description, " ")
				// """""""""""""""check the question response limit"""""""""""""""
				if len(titleArr) < 50 {
					c.JSON(http.StatusExpectationFailed, gin.H{
						"status": "TEXT_IS_TOO_SMALL",
					})
				} else {
					// """""""""""""""insert the data into database"""""""""""""""
					FAQ.PrimitiveID = primitive.NewObjectID()
					FAQ.ID = FAQ.PrimitiveID.Hex()
					FAQ.RequestLogin = ""
					collection.InsertOne(CtxG, FAQ)
				}
			} else {
				c.JSON(http.StatusExpectationFailed, gin.H{
					"status": "NO_PERMISSION_TO_DO_THIS",
				})
			}
		}
	}
}
func RemoveQuestion(c *gin.Context) {
	// """"""""""""""verify does the user has cookie or not""""""""""""""
	// jwtauthentication.Velidation(c)
	// """"""""""""""verify does the user has cookie or not""""""""""""""
	var FAQ structures.Dr_get_views
	c.ShouldBindJSON(&FAQ)
	fmt.Printf("FAQ: %v\n", FAQ)
	if FAQ.Id == "" || FAQ.RequestLogin == "" {
		c.JSON(400, gin.H{
			"status": "NOTCOMLETE",
		})
	} else {
		var accept = Permission(c, FAQ.RequestLogin)
		if accept == "admin" {
			converter.Remove(c, "frequently_asked_question", FAQ.Id)
		} else {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"status": "NO_PERMISSION_TO_DO_THIS",
			})
		}
	}
}
func AdminProfileChange(c *gin.Context) {
	// """"""""""""""verify does the user has cookie or not""""""""""""""
	// jwtauthentication.Velidation(c)
	converter.Convert("admins")
	var globeUserLog structures.AdminLog
	// """""""""get he json request from client """""""""
	// c.ShouldBindJSON(&globeUserLog)
	jsonFM := c.Request.FormValue("json")
	files, handler, errIMG := c.Request.FormFile("img")
	if errIMG != nil {
		c.JSON(409, gin.H{
			"status": "NOIMGFILEEXIST",
		})
	}

	files.Seek(23, 23)
	fmt.Printf("File Name %s\n", handler.Filename)

	json.Unmarshal([]byte(jsonFM), &globeUserLog)
	fmt.Println(globeUserLog)
	compiledLogin := bycrypt.ChecktheLogin(globeUserLog.Login)
	//""""" loop through email and verify it if tthere is @"""""
	if globeUserLog.Email == "" || globeUserLog.Login == "" || globeUserLog.Name == "" || globeUserLog.Password == "" || globeUserLog.Sername == "" || errIMG != nil {
		c.JSON(400, gin.H{
			"status": "NOTCOMLETE",
		})
	} else {
		//"""""" connect to data base to get the user and verify it""""""
		mongoconnection.MongoDB()
		CtxG := mongoconnection.CtxG
		ClientG := mongoconnection.ClientG
		collection := ClientG.Database("MedCard").Collection("admins")
		var DbgetUser structures.DoctorLog
		collection.FindOne(CtxG, bson.M{"login": compiledLogin}).Decode(&DbgetUser)
		if DbgetUser.Login == "" {
			c.Writer.WriteHeader(http.StatusBadRequest)
			c.JSON(http.StatusExpectationFailed, gin.H{
				"status": "DOES_NOT_EXIST",
			})
		} else {
			// """""""""""""""""""""""""check the email on MX velidation"""""""""""""""""""""""""
			var isEmailVelid = bycrypt.ValidateMX(globeUserLog.Email)
			if isEmailVelid != nil {
				c.JSON(http.StatusExpectationFailed, gin.H{
					"status": "EMAIL_DOES_NOT_VALID",
				})
			} else {
				fmt.Printf("globeUserLog.RequestLogin: %v\n", globeUserLog.RequestLogin)
				var accept = Permission(c, globeUserLog.RequestLogin)
				if accept != "admin" {
					c.JSON(http.StatusExpectationFailed, gin.H{
						"status": "NO_PERMISSION_TO_DO_THIS",
					})
				} else {
					// """"""""""""""""""""delete the struc from collections""""""""""""""""""""
					collection.DeleteOne(CtxG, bson.M{"login": compiledLogin})
					ClientG.Database("MedCard").Collection("users").DeleteOne(CtxG, bson.M{"login": compiledLogin})
					// """"""""""""hesh the password befor inserting it into databse""""""""""""
					hashedPassword, err := bycrypt.HashPassword(globeUserLog.Password)
					if err != nil {
						fmt.Fprintf(c.Writer, err.Error())
						return
					}
					globeUserLog.Password = hashedPassword
					globeUserLog.Login = compiledLogin
					globeUserLog.PrimitiveID = primitive.NewObjectID()
					globeUserLog.ID = globeUserLog.PrimitiveID.Hex()
					globeUserLog.ProfileImage = bycrypt.ParseFile(c, "static/uploadImg")
					globeUserLog.Permissions = "admin"
					fmt.Println(hashedPassword)
					// """""""""""""""Insert the user ino dataBase if its valid"""""""""""""""
					collection.InsertOne(CtxG, globeUserLog)
					converter.Convert("admins")
				}
			}
		}
	}
}
func Permission(c *gin.Context, requestLogin string) string {
	var isAccepted string
	// """"""""""mongo connection""""""""""
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	collection := ClientG.Database("MedCard").Collection("users")
	var DbgetUser structures.Users
	collection.FindOne(CtxG, bson.M{"login": requestLogin}).Decode(&DbgetUser)
	isAccepted = DbgetUser.Permissions
	fmt.Println(isAccepted)

	return isAccepted
}
func Statistics(c *gin.Context) {
	var statistc structures.Dr_get_views
	c.ShouldBindJSON(&statistc)
	var accept = Permission(c, statistc.RequestLogin)
	if accept != "admin" {
		c.JSON(http.StatusExpectationFailed, gin.H{
			"status": "NO_PERMISSION_TO_DO_THIS",
		})
	} else {
		// """"""""""""""verify does the user has cookie or not""""""""""""""
		// jwtauthentication.Velidation(c)
		// """""""""""""""""""local veriables"""""""""""""""""""
		var userArr []structures.Users
		var allUsers int

		// """"""""""""""""""""data base connection""""""""""""""""""""
		mongoconnection.MongoDB()
		ClientG := mongoconnection.ClientG
		CtxG := mongoconnection.CtxG
		collection := ClientG.Database("MedCard").Collection("users")
		cursor, err := collection.Find(CtxG, bson.M{})
		// """""""""""""""""""handle the error"""""""""""""""""""
		if err != nil {
			fmt.Fprintf(c.Writer, err.Error())
		}
		//"""""""""""" loop throgth the DB and append the data to userArr""""""""""""
		defer cursor.Close(CtxG)
		for cursor.Next(CtxG) {
			cursor.Decode(&User)
			userArr = append(userArr, User)
		}
		allUsers = len(userArr)
		// """""""""""""""""""local veriables"""""""""""""""""""
		var userOnlineArr []structures.DoctorLog
		var allOnlineUsers int

		// """"""""""""""""""""data base connection""""""""""""""""""""
		mongoconnection.MongoDB()
		ClientO := mongoconnection.ClientG
		CtxO := mongoconnection.CtxG
		collectionO := ClientO.Database("MedCard").Collection("online")
		cursorO, err := collectionO.Find(CtxO, bson.M{})
		// """""""""""""""""""handle the error"""""""""""""""""""
		if err != nil {
			fmt.Fprintf(c.Writer, err.Error())
		}
		//"""""""""""" loop throgth the DB and append the data to userArr""""""""""""
		defer cursorO.Close(CtxG)
		for cursorO.Next(CtxG) {
			cursorO.Decode(&globeUserLog)
			userOnlineArr = append(userOnlineArr, globeUserLog)
		}
		allOnlineUsers = len(userOnlineArr)
		// """""""""""""""""""send the response with amout of users"""""""""""""""""""
		c.JSON(200, gin.H{
			"USERS_AMOUNT": allUsers,
			"USERS_ONLINE": allOnlineUsers,
		})
	}
}
func Users_clients(c *gin.Context) {
	var statistc structures.Dr_get_views
	c.ShouldBindJSON(&statistc)
	var accept = Permission(c, statistc.RequestLogin)
	if accept != "admin" {
		c.JSON(http.StatusExpectationFailed, gin.H{
			"status": "NO_PERMISSION_TO_DO_THIS",
		})
	} else {
		// """"""""""""""verify does the user has cookie or not""""""""""""""
		// jwtauthentication.Velidation(c)
		// """""""""""""""""""local veriables"""""""""""""""""""
		var userOneArr []interface{}
		var globeUserLogDr structures.DoctorLogScreen

		// """"""""""""""""""""data base first connection""""""""""""""""""""
		mongoconnection.MongoDB()
		ClientG := mongoconnection.ClientG
		CtxG := mongoconnection.CtxG
		collectionOne := ClientG.Database("MedCard").Collection("doctors")
		cursorOne, err := collectionOne.Find(CtxG, bson.M{})
		// """""""""""""""""""handle the error"""""""""""""""""""
		if err != nil {
			fmt.Fprintf(c.Writer, err.Error())
		}
		//"""""""""""" loop throgth the DB and append the data to userArr""""""""""""
		defer cursorOne.Close(CtxG)
		for cursorOne.Next(CtxG) {
			cursorOne.Decode(&globeUserLogDr)
			userOneArr = append(userOneArr, globeUserLogDr)
		}
		// '''""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""""
		// """""""""""""""""""local veriables"""""""""""""""""""
		var globeUserLogCl structures.ClientLogScreen
		// """"""""""""""""""""data base first connection""""""""""""""""""""
		collectionTwo := ClientG.Database("MedCard").Collection("clients")
		cursorTwo, err := collectionTwo.Find(CtxG, bson.M{})
		// """""""""""""""""""handle the error"""""""""""""""""""
		if err != nil {
			fmt.Fprintf(c.Writer, err.Error())
		}
		//"""""""""""" loop throgth the DB and append the data to userArr""""""""""""
		defer cursorTwo.Close(CtxG)
		for cursorTwo.Next(CtxG) {
			cursorTwo.Decode(&globeUserLogCl)
			userOneArr = append(userOneArr, globeUserLogCl)
		}
		// """""""""""""""""""send the response with amout of users"""""""""""""""""""
		c.JSON(200, userOneArr)
	}
}
func Questions_get(c *gin.Context) {
	// """"""""""""""verify does the user has cookie or not""""""""""""""
	// jwtauthentication.Velidation(c)
	// """""""""""""""""""local veriables"""""""""""""""""""
	var userOneArr []structures.FrequentlyAskedQuestion

	// """"""""""""""""""""data base first connection""""""""""""""""""""
	mongoconnection.MongoDB()
	ClientG := mongoconnection.ClientG
	CtxG := mongoconnection.CtxG
	collectionOne := ClientG.Database("MedCard").Collection("frequently_asked_question")
	cursorOne, err := collectionOne.Find(CtxG, bson.M{})
	// """""""""""""""""""handle the error"""""""""""""""""""
	if err != nil {
		fmt.Fprintf(c.Writer, err.Error())
	}
	//"""""""""""" loop throgth the DB and append the data to userArr""""""""""""
	defer cursorOne.Close(CtxG)
	for cursorOne.Next(CtxG) {
		cursorOne.Decode(&FAQ)
		userOneArr = append(userOneArr, FAQ)
	}
	c.JSON(200, userOneArr)
}
func Logout(c *gin.Context) {
	// """"""""""""""verify does the user has cookie or not""""""""""""""
	// jwtauthentication.Velidation(c)
	// get the data from the user for logout
	c.ShouldBindJSON(&globeUserLog)
	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    jwtauthentication.GenerateToken(c, globeUserLog.Email),
		Expires:  time.Now().Add(-1 * time.Minute),
		HttpOnly: false,
		Secure:   true,
		Path:     "/",
	})
	// Mongo connection
	mongoconnection.MongoDB()
	ClientG := mongoconnection.ClientG
	CtxG := mongoconnection.CtxG
	collectionOne := ClientG.Database("MedCard").Collection("online")
	collectionOne.DeleteOne(CtxG, bson.M{"login": globeUserLog.Login})
}
func Cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, ResponseType, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
		return
	}
	c.Next()
}
