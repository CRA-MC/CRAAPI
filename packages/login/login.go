package login

import (
	"context"
	"craapi/packages/define"
	"craapi/packages/encryption"
	"craapi/packages/log"
	"craapi/packages/mongodb"
	"craapi/packages/overloaded"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"time"

	"github.com/jxskiss/base62"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrWrongPassword = fmt.Errorf("wrong Password")
var ErrTooManyRequests = fmt.Errorf("to many requests")
var ErrNotHaveThisUser = fmt.Errorf("this email or username have no user in database")

func tokengen(userID string) define.DB_Token_webaccess {
	var token define.DB_Token_webaccess
	salt := make([]byte, 128)
	rand.Read(salt)
	tmp := sha512.Sum512(salt)
	token.WebAccess = base62.EncodeToString(tmp[:])
	filter := bson.D{{"AccessToken", token.WebAccess}}
	err := mongodb.Collection_token.FindOne(context.TODO(), filter).Err()
	for err != mongo.ErrNoDocuments {
		salt = make([]byte, 128)
		rand.Read(salt)
		tmp := sha512.Sum512(salt)
		token.WebAccess = base62.EncodeToString(tmp[:])
		err = mongodb.Collection_token.FindOne(context.TODO(), filter).Err()
		if err != mongo.ErrNoDocuments && err != nil {
			panic(err)
		}
	}
	token.UserID = userID
	token.Time = primitive.NewDateTimeFromTime(time.Now())
	return token
}

func UserAuth(username string, password string) (error, define.DB_Token_webaccess) {
	filter := bson.D{{"Username", username}}
	var res define.DB_User
	err := mongodb.Collection_users.FindOne(context.TODO(), filter).Decode(&res)
	if err == mongo.ErrNoDocuments {
		filter = bson.D{{"Email", username}}
		err = mongodb.Collection_users.FindOne(context.TODO(), filter).Decode(&res)
		if err == mongo.ErrNoDocuments {
			return ErrNotHaveThisUser, define.DB_Token_webaccess{}
		} else if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	if overloaded.Overload_id(&res.ID) {
		return ErrTooManyRequests, define.DB_Token_webaccess{}
	}

	if encryption.CheckPassowrd(password, res.Salt, res.Iterations, res.PasswordHashed) {
		token := tokengen(res.Username)
		filter = bson.D{{"UserID", res.Username}}
		sortfilter := bson.D{{"Time", 1}}
		var tokencount int64
		tokencount, err = mongodb.Collection_token.CountDocuments(context.TODO(), filter)
		if err != nil && err != mongo.ErrNoDocuments {
			panic(err)
		}
		if tokencount >= 10 {
			deletedtokencount := int64(0)
			var deletedtoken define.DB_Token
			for i := int64(0); i < (tokencount - 9); i++ {
				err = mongodb.Collection_token.FindOneAndDelete(context.TODO(), filter, new(options.FindOneAndDeleteOptions).SetSort(sortfilter)).Decode(&deletedtoken)
				if err != nil {
					if err != mongo.ErrNoDocuments {
						panic(err)
					}
				} else {
					deletedtokencount++
				}
				if viper.GetBool("Debug") {
					// log.LOGD("CRA API LOGD [", time.Now().Format("2006-01-02 15:04:05"), "] deleted Token for username: ", res.Username, " Token date: ", deletedtoken.Time.Time().Format("2006-01-02 15:04:05"))
				}
			}
			log.LOGD(
				"deleted ", deletedtokencount, " Token(s) for username: ", res.Username,
			)
		}
		_, err = mongodb.Collection_token.InsertOne(context.TODO(), token)
		if err != nil {
			panic(err)
		}
		return nil, token
	}
	return ErrWrongPassword, define.DB_Token_webaccess{}
}
