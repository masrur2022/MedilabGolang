package controllers

import (
	// "medcard/beck/jwtAuthentication"
	"encoding/json"
	"fmt"
	"medcard/beck/bycrypt"
	"medcard/beck/mongoconnection"
	"medcard/beck/structures"
	"net/http"

	jwtauthentication "medcard/beck/jwtAuthentication"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
)

func Doctors_get(c *gin.Context) {
	// """"""""""""""verify does the user has cookie or not""""""""""""""
	// jwtauthentication.Velidation(c)
	// """""""""""""""""""local veriables"""""""""""""""""""
	var userOneArr []structures.DoctorLogScreen
	var globeUserLog structures.DoctorLogScreen

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
		cursorOne.Decode(&globeUserLog)
		userOneArr = append(userOneArr, globeUserLog)
	}
	c.JSON(200, userOneArr)
}
func Clients_get(c *gin.Context) {
	payloadlogin := jwtauthentication.Velidation(c)
	fmt.Println(payloadlogin.Id)
	var accept = Permission(c, payloadlogin.Login)
	if accept != "admin" {
		c.JSON(http.StatusExpectationFailed, gin.H{
			"status": "NO_PERMISSION_TO_DO_THIS",
		})
	} else {
		// """"""""""""""verify does the user has cookie or not""""""""""""""
		// jwtauthentication.Velidation(c)
		// """""""""""""""""""local veriables"""""""""""""""""""
		var userOneArr []structures.ClientLog
		var userOne structures.ClientLog
		// c.ShouldBindJSON()
		// """"""""""""""""""""data base first connection""""""""""""""""""""
		mongoconnection.MongoDB()
		ClientG := mongoconnection.ClientG
		CtxG := mongoconnection.CtxG
		collectionOne := ClientG.Database("MedCard").Collection("clients")
		cur,err := collectionOne.Find(CtxG, bson.M{})
		if err != nil{
			fmt.Fprintf(c.Writer,err.Error())
		}
		defer cur.Close(CtxG)
		for cur.Next(CtxG){
			cur.Decode(&userOne)
			userOneArr = append(userOneArr, userOne)
		}
		//"""""""""""" loop throgth the DB and append the data to userArr""""""""""""

		c.JSON(200, userOneArr)
	}
}
func DoctorProfileChange(c *gin.Context) {
	// """"""""""""""""""get the json request and bind them into var""""""""""""""""""
	var userOne structures.DoctorLog
	c.ShouldBindJSON(&userOne)
	fmt.Println(userOne)
	//""""""" mongoDb connection"""""""
	payloadlogin := jwtauthentication.Velidation(c)
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	collection := ClientG.Database("MedCard").Collection("doctors")
	collection.FindOne(CtxG, bson.M{"login": payloadlogin.Login}).Decode(&globeUserLog)
	// """""""""""""""check the user velidetion"""""""""""""""
	if globeUserLog.Login == "" {
		c.JSON(404, gin.H{
			"STATUS": "CAN_NOT_BE_FOUND",
		})
	} else {
		// """""""""""""""""""""""check are the structure is filled or not"""""""""""""""""""""""
		if userOne.WorkPlace == "" || userOne.Expereance == "" || userOne.Biography == ""|| userOne.Password == "" || userOne.History.Period == "" || userOne.History.ShortInfo == "" || userOne.History.Location == "" {
			c.JSON(404, gin.H{
				"STATUS": "NOT_COMPLETE",
			})
			fmt.Println(userOne)
		} else {
			// """"""get the json request""""""
			fmt.Println(payloadlogin.Id)
			var accept = Permission(c, payloadlogin.Login)
			if accept != "doctors" {
				c.JSON(http.StatusExpectationFailed, gin.H{
					"status": "NO_PERMISSION_TO_DO_THIS",
				})
			} else {
				// """"""""""""hesh the password befor inserting it into databse""""""""""""
				hashedPassword, err := bycrypt.HashPassword(userOne.Password)
				if err != nil {
					fmt.Fprintf(c.Writer, err.Error())
					return
				}
				//""""""""""""""" replace the data"""""""""""""""
				globeUserLog.Password = hashedPassword
				globeUserLog.Login = userOne.Login
				globeUserLog.Name = userOne.Name
				globeUserLog.Sername = userOne.Sername
				globeUserLog.Email = userOne.Email
				globeUserLog.Phone = userOne.Phone
				globeUserLog.WorkPlace = userOne.WorkPlace
				globeUserLog.Expereance = userOne.Expereance
				globeUserLog.Biography = userOne.Biography
				globeUserLog.History.Location = userOne.History.Location
				globeUserLog.History.Period = userOne.History.Period
				globeUserLog.History.ShortInfo = userOne.History.ShortInfo
				//""""""""""""""""" delete the exist data"""""""""""""""""
				collection.DeleteOne(CtxG, bson.M{"login": payloadlogin.Login})
				//"""""""""""""""""" insert the the new data""""""""""""""""""
				collection.InsertOne(CtxG, globeUserLog)

				//"""""""""""""""""" replace the user with new one in `users` DB """"""""""""""""""
				// var usersStruct structures.Users
				collectionUsers := ClientG.Database("MedCard").Collection("users")
				// collection.FindOne(CtxG, bson.M{"login": payloadlogin.Login}).Decode(&globeUserLog)
				collectionUsers.DeleteOne(CtxG,bson.M{"login":payloadlogin.Login})
				// usersStruct.Id = globeUserLog.ID
				// usersStruct.Password = hashedPassword
				// usersStruct.Login = globeUserLog.Login
				// usersStruct.Permissions = globeUserLog.Permissions
				// usersStruct.PrimitiveID = globeUserLog.PrimitiveID
				// collectionUsers.InsertOne(CtxG,usersStruct)
			}
		}
	}
}
func Client_prof_change(c *gin.Context) {
	// """"""""""""""""""get the json request and bind them into var""""""""""""""""""
	var globeClientLog structures.ClientLog
	var globeUserLog structures.ClientLog
	jsonFM := c.Request.FormValue("json")
	fmt.Printf("jsonFM: %v\n", jsonFM)
	files, handler, errIMG := c.Request.FormFile("img")
	if errIMG != nil {
		c.JSON(409, gin.H{
			"sttus": "NOIMGFILEEXIST",
		})
	}

	files.Seek(23, 23)
	fmt.Printf("File Name %s\n", handler.Filename)

	json.Unmarshal([]byte(jsonFM), &globeClientLog)
	fmt.Printf("globeClientLog: %v\n", globeClientLog)
	//""""""" mongoDb connection"""""""
	payloadlogin := jwtauthentication.Velidation(c)
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	collection := ClientG.Database("MedCard").Collection("clients")
	collection.FindOne(CtxG, bson.M{"login": payloadlogin.Login}).Decode(&globeUserLog)
	// """""""""""""""check the user velidetion"""""""""""""""
	if globeUserLog.Login == "" {
		c.JSON(404, gin.H{
			"STATUS": "CAN_NOT_BE_FOUND",
		})
	} else {
		// """""""""""""""""""""""check are the structure is filled or not"""""""""""""""""""""""
		// var isEmailVelid = bycrypt.ValidateMX(globeUserLog.Email) // || isEmailVelid != nil
		if globeClientLog.Phone == "" || errIMG != nil  {
			c.JSON(404, gin.H{
				"STATUS": "NOT_COMPLETE",
			})
		} else {
			// """"""get the json request""""""
			fmt.Println(payloadlogin.Id)
			var accept = Permission(c, payloadlogin.Login)
			if accept != "client" {
				c.JSON(http.StatusExpectationFailed, gin.H{
					"status": "NO_PERMISSION_TO_DO_THIS",
				})
			} else {
				// """"""""""""hesh the password befor inserting it into databse""""""""""""
				hashedPassword, err := bycrypt.HashPassword(globeClientLog.Password)
				compareResult := bycrypt.CompareHashPasswords(hashedPassword, globeClientLog.Password)
				fmt.Printf("compareResult: %v\n", compareResult)

				if err != nil{
					fmt.Fprintf(c.Writer, err.Error())
				}
				fmt.Printf("globeClientLog.Phone: %v\n", globeClientLog.Phone)
				fmt.Printf("globeUserLog.Email: %v\n", globeUserLog.Email)
				//""""""""""""""" replace the data"""""""""""""""
				globeUserLog.Phone = globeClientLog.Phone
				globeUserLog.Email = globeClientLog.Email
				globeUserLog.Password = hashedPassword
				globeUserLog.Login = globeClientLog.Login
				globeUserLog.ProfileImage = bycrypt.ParseFile(c, "static/uploadImg")
				//""""""""""""""""" delete the exist data"""""""""""""""""
				collection.DeleteOne(CtxG, bson.M{"login": payloadlogin.Login})
				//"""""""""""""""""" insert the the new data""""""""""""""""""
				collection.InsertOne(CtxG, globeUserLog)

				//"""""""""""""""""" replace the user with new one in `users` DB """"""""""""""""""
				userCoolection := ClientG.Database("MedCard").Collection("users")
				userCoolection.DeleteOne(CtxG,bson.M{"login":payloadlogin.Login})
				jwtauthentication.Velidation(c)
			}
		}
	}
}
func Accept_decline(c *gin.Context) {
	var accept_decline structures.Accept_Decline
	// """"""get the json request""""""
	c.ShouldBindJSON(&accept_decline)
	// """"""""""mongo connection""""""""""
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	collection := ClientG.Database("MedCard").Collection("views")
	var DbgetUser structures.ViewReq
	collection.FindOne(CtxG, bson.M{"clientid": accept_decline.ClientId}).Decode(&DbgetUser)
	if DbgetUser.ClientId == "" {
		c.JSON(409, gin.H{
			"status": "NO_SUCH_CLIENT_WITH_THAT_CRIDENTIALS",
		})
	} else {
		// """"""get the json request""""""
		payloadlogin := jwtauthentication.Velidation(c)
		fmt.Println(payloadlogin.Id)
		var accept = Permission(c, payloadlogin.Login)
		if accept != "doctors" {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"status": "NO_PERMISSION_TO_DO_THIS",
			})
		} else {
			DbgetUser.Date = accept_decline.Date
			DbgetUser.Status = "checked"

			collection.DeleteOne(CtxG, bson.M{"clientid": accept_decline.ClientId})
			collection.InsertOne(CtxG, DbgetUser)
		}
	}
}
func Signup_cl_view(c *gin.Context) {
	var accept_decline structures.ViewReq
	// """"""get the json request""""""
	c.ShouldBindJSON(&accept_decline)
	// """"""""""mongo connection""""""""""
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	collection := ClientG.Database("MedCard").Collection("views")
	var DbgetUser structures.ViewReq
	collection.FindOne(CtxG, bson.M{"clientid": accept_decline.ClientId, "doctorid": accept_decline.DoctorId}).Decode(&DbgetUser)
	if DbgetUser.Sickness != "" {
		c.JSON(409, gin.H{
			"status": "YOU_POSTED_SUCH_REQUEST_EARLIER",
		})
	} else {
		if  accept_decline.DoctorId == "" || accept_decline.Date == "" || accept_decline.Phone == "" || accept_decline.Sickness == "" {
			c.JSON(409, gin.H{
				"status": "NOT_COMPLETE",
			})
		} else {
			var doctorVerify structures.DoctorLog
			doctorCollectioh := ClientG.Database("MedCard").Collection("doctors")
			doctorCollectioh.FindOne(CtxG,bson.M{"id":accept_decline.DoctorId}).Decode(&doctorVerify)
			if doctorVerify.Login == ""{
				c.JSON(409, gin.H{
					"status": "DOTOR_ID_DOES_NOT_EXIST",
				})
			}else{
				// """"""get the json request""""""
				payloadlogin := jwtauthentication.Velidation(c)
				fmt.Println(payloadlogin.Id)
				var accept = Permission(c, payloadlogin.Login)
				if accept != "client" {
					c.JSON(http.StatusExpectationFailed, gin.H{
						"status": "NO_PERMISSION_TO_DO_THIS",
					})
				} else {
					//"""""""""""""""" conection to client DB to get user info""""""""""""""""
					collectionClient := ClientG.Database("MedCard").Collection("clients")
					var DbgetUserForInfo structures.ClientLog
					collectionClient.FindOne(CtxG, bson.M{"login": payloadlogin.Login}).Decode(&DbgetUserForInfo)
					//"""""""""""" veriables connection""""""""""""
					accept_decline.ClientId = payloadlogin.Id
					accept_decline.ClientName = DbgetUserForInfo.Name
					accept_decline.ClientSername = DbgetUserForInfo.Sername
					accept_decline.Status = "Waiting"
					fmt.Printf("DbgetUserForInfo: %v\n", DbgetUserForInfo)
					collection.InsertOne(CtxG, accept_decline)
				}
			}
		}
	}
}
func Views_get_dr(c *gin.Context) {
	// """"""get the json request""""""
	payloadlogin := jwtauthentication.Velidation(c)
	fmt.Println(payloadlogin.Id)
	var accept = Permission(c, payloadlogin.Login)
	// """"""""""mongo connection""""""""""
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	collection := ClientG.Database("MedCard").Collection("views")
	var DbgetUser structures.ViewReq
	var DbgetUserArr []structures.ViewReq
	cur, err := collection.Find(CtxG, bson.M{"doctorid": payloadlogin.Id})
	if err != nil {
		fmt.Fprintf(c.Writer, err.Error())
	}
	defer cur.Close(CtxG)
	for cur.Next(CtxG) {
		cur.Decode(&DbgetUser)
		DbgetUserArr = append(DbgetUserArr, DbgetUser)
	}
	if payloadlogin.Id == "" {
		c.JSON(409, gin.H{
			"status": "NO_REQUEST_FOUND",
		})
	} else {
		if accept == "doctors" || accept == "client" {
			c.JSON(200, DbgetUserArr)
		} else {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"status": "NO_PERMISSION_TO_DO_THIS",
			})
		}
	}
}
func Views_get_cl(c *gin.Context) {
	// send Id client
	payloadlogin := jwtauthentication.Velidation(c)
	fmt.Println(payloadlogin.Id)
	var accept = Permission(c, payloadlogin.Login)
	// """"""""""mongo connection""""""""""
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	collection := ClientG.Database("MedCard").Collection("views")
	var DbgetUser structures.ViewReq
	var DbgetUserArr []structures.ViewReq
	cur, err := collection.Find(CtxG, bson.M{"doctorid": payloadlogin.Id})
	if err != nil {
		fmt.Fprintf(c.Writer, err.Error())
	}
	defer cur.Close(CtxG)
	for cur.Next(CtxG) {
		cur.Decode(&DbgetUser)
		DbgetUserArr = append(DbgetUserArr, DbgetUser)
	}
	if DbgetUser.Sickness == "" {
		c.JSON(409, gin.H{
			"status": "NO_REQUEST_FOUND",
		})
	} else {
		if accept != "doctors" {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"status": "NO_PERMISSION_TO_DO_THIS",
			})
		} else {
			c.JSON(200, DbgetUserArr)
		}
	}
}
func Emc_get(c *gin.Context) {
	// send Id client
	var ClientReq structures.Dr_get_views
	// """"""get the json request""""""
	c.ShouldBindJSON(&ClientReq)
	// """""""""cookie verificaion"""""""""
	jwtauthentication.Velidation(c)
	payloadlogin := jwtauthentication.Velidation(c)
	fmt.Println(payloadlogin.Id)
	// """"""""""mongo connection""""""""""
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	collection := ClientG.Database("MedCard").Collection("clients")
	var DbgetUser structures.ClientLog
	collection.FindOne(CtxG, bson.M{"id": payloadlogin.Id}).Decode(&DbgetUser)
	if DbgetUser.LastName == "" {
		c.JSON(409, gin.H{
			"status": "NO_REQUEST_FOUND",
		})
	} else {
		var statistc structures.Dr_get_views
		c.ShouldBindJSON(&statistc)
		var accept = Permission(c, string(payloadlogin.Login))
		if accept != "doctors" {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"status": "NO_PERMISSION_TO_DO_THIS",
			})
		} else {
			DbgetUser.Password = ""
			DbgetUser.Login = ""
			DbgetUser.Permissions = ""
			c.JSON(200, DbgetUser)
		}
	}
}
