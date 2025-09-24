package panel

import (
	"context"
	"craapi/packages/mongodb"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserName struct {
	Username string
}

func NewUserNameCheckPost(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var user UserName
	err := jsoniter.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		panic(err)
	}
	filter := bson.D{{"Username", user.Username}}
	err = mongodb.Collection_users.FindOne(context.TODO(), filter).Err()
	if err == nil {
		ctx.SetStatusCode(fasthttp.StatusAccepted)
		ctx.WriteString("UsernamePassed!")
	} else if err != mongo.ErrNoDocuments {
		panic(err)
	}
	ctx.SetStatusCode(200)
	ctx.WriteString("{\"Username\":\"Exist\"}")
}

type Email struct {
	Email string
}

func NewUserEmailCheckPost(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var user Email
	err := jsoniter.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		panic(err)
	}
	filter := bson.D{{"Email", user.Email}}
	err = mongodb.Collection_users.FindOne(context.TODO(), filter).Err()
	if err == nil {
		ctx.SetStatusCode(fasthttp.StatusAccepted)
		ctx.WriteString("EmailPassed!")
	} else if err != mongo.ErrNoDocuments {
		panic(err)
	}
	ctx.SetStatusCode(200)
	ctx.WriteString("{\"Email\":\"Exist\"}")
}

func Api_Login(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var user Email
	err := jsoniter.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		panic(err)
	}
	filter := bson.D{{"Email", user.Email}}
	err = mongodb.Collection_users.FindOne(context.TODO(), filter).Err()
	if err == nil {
		ctx.SetStatusCode(fasthttp.StatusAccepted)
		ctx.WriteString("EmailPassed!")
	} else if err != mongo.ErrNoDocuments {
		panic(err)
	}
	ctx.SetStatusCode(200)
	ctx.WriteString("{\"Email\":\"Exist\"}")
}

type Emailauth struct {
	Email string
}

func Api_Email_Auth(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var user Email
	err := jsoniter.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		panic(err)
	}
	filter := bson.D{{"Email", user.Email}}
	err = mongodb.Collection_users.FindOne(context.TODO(), filter).Err()
	if err == nil {
		ctx.SetStatusCode(403)
		return
	} else if err != mongo.ErrNoDocuments {

	}
	panic(err)
}
