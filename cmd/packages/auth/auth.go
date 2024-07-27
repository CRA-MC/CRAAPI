package auth

import (
	"context"
	"craapi/cmd/packages/define"
	"craapi/cmd/packages/encryption"
	"craapi/cmd/packages/mongodb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection_users, collection_overloaded *mongo.Collection

func Init() {
	collection_users = mongodb.Mongodb_GetCollection("Users", "craapi")
	collection_overloaded = mongodb.Mongodb_GetCollection("Overloaded", "craapi")
}

func Auth(usernameoremail string, password string) (bool, *define.DB_User) {
	filter := bson.D{{"Username", usernameoremail}}
	var res define.DB_User
	err := collection_users.FindOne(context.TODO(), filter).Decode(&res)
	if err == mongo.ErrNoDocuments {
		filter = bson.D{{"Email", usernameoremail}}
		err = collection_users.FindOne(context.TODO(), filter).Decode(&res)
		if err == mongo.ErrNoDocuments {
			return false, nil
		} else if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	filter = bson.D{{"UID", res.ID}}
	err = collection_overloaded.FindOne(context.TODO(), filter).Err()
	if err == nil {
		for err == nil {
			time.Sleep(6 * time.Second)
			err = collection_overloaded.FindOne(context.TODO(), filter).Err()
		}
	}
	if err != mongo.ErrNoDocuments {
		panic(err)
	}

	_, err = collection_overloaded.InsertOne(context.TODO(), define.DB_Overloaded{UserID: res.ID, Time: primitive.NewDateTimeFromTime(time.Now())})
	if err != nil {
		panic(err)
	}

	if encryption.CheckPassowrd(password, res.Salt, res.Iterations, res.PasswordHashed) {
		return true, &res
	}
	return false, nil
}
