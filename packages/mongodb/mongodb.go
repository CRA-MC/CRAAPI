package mongodb

import (
	"context"
	"craapi/packages/log"
	"fmt"
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongodb_client *mongo.Client
var Mongodb_database *mongo.Database
var Collection_users, Collection_overloaded, Collection_token, Collection_profiles, Collection_serverid, Collection_cookie, Collection_emailauthcode *mongo.Collection

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
	mongo_url += "?authSource="
	if ext != "" {
		mongo_url += ext
	}
	log.LOGI("monogodb url: ", mongo_url)
	var err error
	Mongodb_client, err = mongo.Connect(context.TODO(), options.Client().
		ApplyURI(mongo_url))
	if err != nil {
		panic(err)
	}
	Mongodb_database = Mongodb_client.Database(ext)
}

func Mongodb_GetCollections() {
	Collection_profiles = Mongodb_database.Collection("Profiles")
	Collection_users = Mongodb_database.Collection("Users")
	Collection_overloaded = Mongodb_database.Collection("Overloaded")
	Collection_token = Mongodb_database.Collection("Token")
	Collection_serverid = Mongodb_database.Collection("ServerID")
	Collection_cookie = Mongodb_database.Collection("Cookie")
	Collection_emailauthcode = Mongodb_database.Collection("EmailAuthCode")
}

func Mongodb_GetCollection(CollectionName string) *mongo.Collection {
	return Mongodb_database.Collection(CollectionName)
}
