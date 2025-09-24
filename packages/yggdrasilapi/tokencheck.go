package yggdrasilapi

import (
	"context"
	"craapi/packages/define"
	"craapi/packages/mongodb"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func tokencheck(accessToken string, clientToken string) bool {
	filter := bson.D{{"AccessToken", accessToken}}
	if uuid.Validate(clientToken) == nil {
		filter = bson.D{{"AccessToken", accessToken}, {"ClientToken", clientToken}}
	}
	var Token define.DB_Token
	err := mongodb.Collection_token.FindOne(context.TODO(), filter).Decode(&Token)
	if err == mongo.ErrNoDocuments {
		return false
	}
	if err != nil {
		panic(err)
	}
	if time.Since(Token.Time.Time()).Abs().Hours() >= 12 {
		return false
	}
	return Token.Validate
}

func tokencheckexist(accessToken string, clientToken string) bool {

	filter := bson.D{{"AccessToken", accessToken}}
	if uuid.Validate(clientToken) == nil {
		filter = bson.D{{"AccessToken", accessToken}, {"ClientToken", clientToken}}
	}
	err := mongodb.Collection_token.FindOne(context.TODO(), filter).Err()
	if err == mongo.ErrNoDocuments {
		return false
	}
	if err != nil {
		panic(err)
	}
	return true
}
