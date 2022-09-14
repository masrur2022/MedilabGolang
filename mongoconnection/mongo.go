package mongoconnection

import (
	"context"
	"fmt"
	"os"
	"time"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CtxG context.Context
var ClientG *mongo.Client

func MongoDB() {
	clientOptions := options.Client().ApplyURI(os.Getenv("DB_URL"))
	fmt.Println(os.Getenv("DB_URL"))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println("Mongo.connect() ERROR: ", err)
		os.Exit(1)
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Minute)
	collection := client.Database("MedCard").Collection("users")
	fmt.Println(collection)

	CtxG = ctx
	ClientG = client
}
