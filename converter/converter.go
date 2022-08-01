package converter

import (
	"context"
	"fmt"	
	"medcard/beck/mongoconnection"
	"medcard/beck/structures"


	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gin-gonic/gin"

) 
var client *mongo.Client
var ctx context.Context
var userAuthentication structures.Users
var FAQ          structures.FrequentlyAskedQuestion


func Convert(collectionName string){
	// """""""""""connection to doctors collection"""""""""""
	mongoconnection.MongoDB()
	client := mongoconnection.ClientG
	ctx := mongoconnection.CtxG

	//"""""""""""""""" veriables for declering the date into """"""""""""""""
	var collectionFile structures.Users
	var collection []structures.Users
	//"""""""""""""""" collection specification""""""""""""""""
	DRcollection := client.Database("MedCard").Collection(collectionName)
	// """"""""""""""""""""get all users from doctors collection""""""""""""""""""""
	cur ,err := DRcollection.Find(ctx,bson.M{})

	if err != nil{
		fmt.Println("rer")
	}
	// """"""""""""""""decoding the data to veriables spesified earlier""""""""""""""""
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		cur.Decode(&collectionFile)
		collection = append(collection, collectionFile)
	}
	// """""""""user connection"""""""""
	//"""""""""""""""" veriables for declering the date into """"""""""""""""
	var userFile structures.Users
	var user []structures.Users
	//"""""""""""""""" collection specification""""""""""""""""
	UScollection := client.Database("MedCard").Collection("users")
	// """"""""""""""""""""get all users from doctors collection""""""""""""""""""""
	UScur ,err := UScollection.Find(ctx,bson.M{})

	if err != nil{
		fmt.Println("rer")
	}
	// """"""""""""""""decoding the data to veriables spesified earlier""""""""""""""""
	defer UScur.Close(ctx)
	for UScur.Next(ctx) {
		UScur.Decode(&userFile)
		user = append(user, userFile)
	}
	if len(user) == 0 {
		for i := 0; i < len(collection) ; i++{
			UScollection.InsertOne(ctx,collection[i])
		}
	}else{
		// devide the users arrey into two part 1 users with permission we got from another collecton  2 deos not needed
		var newUsersPermission []structures.Users
		var newUserFile structures.Users
		// loop thrugh the arrey / get users needed/ append them into newUsersPermission
		for u := 0;u < len(user);u++{
			if user[u].Permissions == collectionFile.Permissions{
				newUserFile = user[u]
				newUsersPermission = append(newUsersPermission, newUserFile)
			}
		}
		//"""""""""""""""""""" transfering the data from doctors to users""""""""""""""""""""
		diference := len(collection) - len(newUsersPermission)
		for i := len(collection) - diference ; i < len(collection) ; i++{
				UScollection.InsertOne(ctx,collection[i])
		}
	}
}
//""""""""""""""""""""""""""""""""""""""" remove data from dataBase"""""""""""""""""""""""""""""""""""""""
func Remove(c *gin.Context,collectionName string ,id string){
		//""""""" mongoDb connection""""""" 
		mongoconnection.MongoDB()
		CtxG := mongoconnection.CtxG
		ClientG := mongoconnection.ClientG
		var DBdeleteQuestion structures.FrequentlyAskedQuestion
		collection := ClientG.Database("MedCard").Collection(collectionName)
		collection.FindOne(CtxG,bson.M{"id":id}).Decode(&DBdeleteQuestion)
		if DBdeleteQuestion.Title == ""{
			c.JSON(400,gin.H{
				"status":"CANNOT_FIND_THE_SPECIFYED_TITLE",
			})
		}else{
			collection.DeleteOne(CtxG,bson.M{"id":id})
		}
}