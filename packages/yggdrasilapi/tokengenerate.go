package yggdrasilapi

import (
	"context"
	"craapi/packages/define"
	"craapi/packages/mongodb"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/jxskiss/base62"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func tokengen(clientToken string, userID primitive.ObjectID, UUID string) define.DB_Token {
	var token define.DB_Token
	salt := make([]byte, 128)
	rand.Read(salt)
	tmp := sha512.Sum512(salt)
	token.AccessToken = base62.EncodeToString(tmp[:])
	filter := bson.D{{"AccessToken", token.AccessToken}}
	err := mongodb.Collection_token.FindOne(context.TODO(), filter).Err()
	for err != mongo.ErrNoDocuments {
		salt = make([]byte, 128)
		rand.Read(salt)
		tmp := sha512.Sum512(salt)
		token.AccessToken = base62.EncodeToString(tmp[:])
		err = mongodb.Collection_token.FindOne(context.TODO(), filter).Err()
		if err != mongo.ErrNoDocuments && err != nil {
			panic(err)
		}
	}
	token.UserID = userID
	token.Validate = true
	if uuid.Validate(clientToken) != nil {
		uuidGen, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		}
		clientToken = hex.EncodeToString(uuidGen[:])
	}
	token.ClientToken = clientToken
	token.Time = primitive.NewDateTimeFromTime(time.Now())
	token.UUID = UUID
	return token
}
