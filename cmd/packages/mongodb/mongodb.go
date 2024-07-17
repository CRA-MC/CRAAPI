package mongodb

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongodb_client *mongo.Client

func Mongodb_INIT(username string, password string, host string, port int, ext string) {
	//init mongodb
	mongo_url := "mongodb://"
	if username != "" {
		mongo_url += username
		if password != "" {
			mongo_url += ":" + password
		}
		mongo_url += "@"
	} else {
		if password != "" {
			fmt.Println("Error: mongoDB have password but not user name")
			os.Exit(1)
		}
	}
	if host != "" {
		mongo_url += host + ":" + strconv.Itoa(port)
	} else {
		fmt.Println("Error: mongoDB have no url")
		os.Exit(1)
	}
	mongo_url += "/"
	if ext != "" {
		mongo_url += ext
	}
	fmt.Println("monogodb connection: ", mongo_url)
	var err error
	mongodb_client, err = mongo.Connect(context.TODO(), options.Client().
		ApplyURI(mongo_url))
	if err != nil {
		panic(err)
	}
}

func Mongodb_GetCollection(CollectionName string, Database string) *mongo.Collection {
	return mongodb_client.Database(Database).Collection(CollectionName)
}
