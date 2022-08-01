package controllers

import (
	// "medcard/beck/jwtAuthentication"
	"encoding/json"
	"fmt"
	"medcard/beck/bycrypt"
	"medcard/beck/mongoconnection"
	"medcard/beck/structures"
	"net/http"

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
		var userOneArr structures.ClientLog
		// c.ShouldBindJSON()
		// """"""""""""""""""""data base first connection""""""""""""""""""""
		mongoconnection.MongoDB()
		ClientG := mongoconnection.ClientG
		CtxG := mongoconnection.CtxG
		collectionOne := ClientG.Database("MedCard").Collection("doctors")
		collectionOne.Find(CtxG, bson.M{})
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
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	collection := ClientG.Database("MedCard").Collection("doctors")
	collection.FindOne(CtxG, bson.M{"login": userOne.Login}).Decode(&globeUserLog)
	// """""""""""""""check the user velidetion"""""""""""""""
	if globeUserLog.Login == "" {
		c.JSON(404, gin.H{
			"STATUS": "CAN_NOT_BE_FOUND",
		})
	} else {
		// """""""""""""""""""""""check are the structure is filled or not"""""""""""""""""""""""
		if userOne.WorkPlace == "" || userOne.Expereance == "" || userOne.Biography == "" || userOne.History.Period == "" || userOne.History.ShortInfo == "" || userOne.History.Location == "" {
			c.JSON(404, gin.H{
				"STATUS": "NOT_COMPLETE",
			})
		} else {
			var accept = Permission(c, globeUserLog.Login)
			if accept != "doctors" {
				c.JSON(http.StatusExpectationFailed, gin.H{
					"status": "NO_PERMISSION_TO_DO_THIS",
				})
			} else {
				//""""""""""""""" replace the data"""""""""""""""
				userOne.PrimitiveID = globeUserLog.PrimitiveID
				userOne.ID = globeUserLog.ID
				userOne.Password = globeUserLog.Password
				userOne.Login = globeUserLog.Login
				userOne.Name = globeUserLog.Name
				userOne.Sername = globeUserLog.Sername
				userOne.Possition = globeUserLog.Possition
				userOne.ProfileImage = globeUserLog.ProfileImage
				userOne.Permissions = globeUserLog.Permissions
				//""""""""""""""""" delete the exist data"""""""""""""""""
				collection.DeleteOne(CtxG, bson.M{"login": userOne.Login})
				//"""""""""""""""""" insert the the new data""""""""""""""""""
				collection.InsertOne(CtxG, userOne)
			}
		}
	}
}
func Client_prof_change(c *gin.Context) {
	// """"""""""""""""""get the json request and bind them into var""""""""""""""""""
	var globeClientLog structures.ClientLog
	var globeUserLog structures.ClientLog
	jsonFM := c.Request.FormValue("json")
	files, handler, errIMG := c.Request.FormFile("img")
	if errIMG != nil {
		c.JSON(409, gin.H{
			"sttus": "NOIMGFILEEXIST",
		})
	}

	files.Seek(23, 23)
	fmt.Printf("File Name %s\n", handler.Filename)

	json.Unmarshal([]byte(jsonFM), &globeClientLog)
	fmt.Println(globeClientLog)
	fmt.Println(jsonFM)
	//""""""" mongoDb connection"""""""
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	collection := ClientG.Database("MedCard").Collection("clients")
	collection.FindOne(CtxG, bson.M{"login": globeClientLog.Login}).Decode(&globeUserLog)
	// """""""""""""""check the user velidetion"""""""""""""""
	if globeUserLog.Login == "" {
		c.JSON(404, gin.H{
			"STATUS": "CAN_NOT_BE_FOUND",
		})
	} else {
		// """""""""""""""""""""""check are the structure is filled or not"""""""""""""""""""""""
		var isEmailVelid = bycrypt.ValidateMX(globeUserLog.Email)
		if globeClientLog.Phone == "" || errIMG != nil || isEmailVelid != nil {
			c.JSON(404, gin.H{
				"STATUS": "NOT_COMPLETE",
			})
		} else {
			var statistc structures.Dr_get_views
			c.ShouldBindJSON(&statistc)
			var accept = Permission(c, statistc.RequestLogin)
			if accept != "clients" {
				c.JSON(http.StatusExpectationFailed, gin.H{
					"status": "NO_PERMISSION_TO_DO_THIS",
				})
			} else {
				//""""""""""""""" replace the data"""""""""""""""
				globeClientLog.PrimitiveID = globeUserLog.PrimitiveID
				globeClientLog.ID = globeUserLog.ID
				globeClientLog.Phone = globeUserLog.Phone
				globeClientLog.Password = globeUserLog.Password
				globeClientLog.Login = globeUserLog.Login
				globeClientLog.Name = globeUserLog.Name
				globeClientLog.Sername = globeUserLog.Sername
				globeClientLog.LastName = globeUserLog.LastName
				globeClientLog.Gender = globeUserLog.Gender
				globeClientLog.Blood = globeUserLog.Blood
				globeClientLog.Permissions = globeUserLog.Permissions
				globeClientLog.ProfileImage = bycrypt.ParseFile(c, "static/uploadImg")
				//""""""""""""""""" delete the exist data"""""""""""""""""
				collection.DeleteOne(CtxG, bson.M{"login": globeUserLog.Login})
				//"""""""""""""""""" insert the the new data""""""""""""""""""
				collection.InsertOne(CtxG, globeClientLog)
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
		var statistc structures.Dr_get_views
		c.ShouldBindJSON(&statistc)
		var accept = Permission(c, statistc.RequestLogin)
		if accept != "doctors" {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"status": "NO_PERMISSION_TO_DO_THIS",
			})
		} else {
			DbgetUser.Date = accept_decline.Date
			DbgetUser.Status = accept_decline.Status

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
		if accept_decline.ClientId == "" || accept_decline.DoctorId == "" || accept_decline.ClientName == "" || accept_decline.ClientSername == "" || accept_decline.Phone == "" || accept_decline.Sickness == "" {
			c.JSON(409, gin.H{
				"status": "NOT_COMPLETE",
			})
		} else {
			var statistc structures.Dr_get_views
			c.ShouldBindJSON(&statistc)
			var accept = Permission(c, statistc.RequestLogin)
			if accept != "client" {
				c.JSON(http.StatusExpectationFailed, gin.H{
					"status": "NO_PERMISSION_TO_DO_THIS",
				})
			} else {
				collection.InsertOne(CtxG, accept_decline)
			}
		}
	}
}
func Views_get_dr(c *gin.Context) {
	var ViewReq structures.Dr_get_views
	// """"""get the json request""""""
	c.ShouldBindJSON(&ViewReq)
	fmt.Println(ViewReq)
	// """"""""""mongo connection""""""""""
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	collection := ClientG.Database("MedCard").Collection("views")
	var DbgetUser structures.ViewReq
	var DbgetUserArr []structures.ViewReq
	cur, err := collection.Find(CtxG, bson.M{"doctorid": ViewReq.Id})
	if err != nil {
		fmt.Fprintf(c.Writer, err.Error())
	}
	defer cur.Close(CtxG)
	for cur.Next(CtxG) {
		cur.Decode(&DbgetUser)
		DbgetUserArr = append(DbgetUserArr, DbgetUser)
	}
	if ViewReq.Id == "" {
		c.JSON(409, gin.H{
			"status": "NO_REQUEST_FOUND",
		})
	} else {
		var statistc structures.Dr_get_views
		c.ShouldBindJSON(&statistc)
		var accept = Permission(c, statistc.RequestLogin)
		if accept != "doctors" {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"status": "NO_PERMISSION_TO_DO_THIS",
			})
		} else {
			c.JSON(200, DbgetUserArr)
		}
	}
}
func Views_get_cl(c *gin.Context) {
	// send Id client
	var ViewReq structures.Dr_get_views
	// """"""get the json request""""""
	c.ShouldBindJSON(&ViewReq)
	// """"""""""mongo connection""""""""""
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	collection := ClientG.Database("MedCard").Collection("views")
	var DbgetUser structures.ViewReq
	var DbgetUserArr []structures.ViewReq
	cur, err := collection.Find(CtxG, bson.M{"doctorid": ViewReq.Id})
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
		var statistc structures.Dr_get_views
		c.ShouldBindJSON(&statistc)
		var accept = Permission(c, statistc.RequestLogin)
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
	// """"""""""mongo connection""""""""""
	mongoconnection.MongoDB()
	CtxG := mongoconnection.CtxG
	ClientG := mongoconnection.ClientG
	collection := ClientG.Database("MedCard").Collection("clients")
	var DbgetUser structures.ClientLog
	collection.FindOne(CtxG, bson.M{"id": ClientReq.Id}).Decode(&DbgetUser)
	if DbgetUser.LastName == "" {
		c.JSON(409, gin.H{
			"status": "NO_REQUEST_FOUND",
		})
	} else {
		var statistc structures.Dr_get_views
		c.ShouldBindJSON(&statistc)
		var accept = Permission(c, statistc.RequestLogin)
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
