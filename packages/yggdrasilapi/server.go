package yggdrasilapi

import (
	"context"
	"craapi/packages/define"
	"craapi/packages/mongodb"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthJoin struct {
	SelectedProfile string `json:"selectedProfile"`
	AccessToken     string `json:"accessToken"`
	ServerId        string `json:"serverId"`
}

func GetIP(ctx *fasthttp.RequestCtx) (string, error) {
	var ip string
	ip = string(ctx.Request.Header.Peek("X-Real-IP"))
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = string(ctx.Request.Header.Peek("X-Forward-For"))

	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip = string(ctx.RemoteIP())

	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	return "", errors.New("no valid ip found")
}

func ServerJoin(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var autho AuthJoin
	err := jsoniter.Unmarshal(ctx.Request.Body(), &autho)
	if err != nil {
		panic(err)
	}

	filter := bson.D{{"AccessToken", autho.AccessToken}}
	var Token define.DB_Token
	err = mongodb.Collection_token.FindOne(context.TODO(), filter).Decode(&Token)
	if err != nil && err != mongo.ErrNoDocuments {
		panic(err)
	}
	if err == mongo.ErrNoDocuments || autho.SelectedProfile != Token.UUID {
		ctx.Response.SetStatusCode(403)
		fmt.Fprint(ctx, "{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"角色错误\"}")
		return
	}
	ctx.Response.SetStatusCode(204)
	_, err = mongodb.Collection_serverid.InsertOne(context.TODO(), define.DB_ServerID{
		ServerID:    autho.ServerId,
		AccessToken: autho.AccessToken,
		UUID:        autho.SelectedProfile,
		Time:        primitive.NewDateTimeFromTime(time.Now()),
	})
	if err != nil {
		panic(err)
	}
}

func ServerhasJoined(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	values := ctx.QueryArgs()
	username := string(values.Peek("username"))
	serverId := string(values.Peek("serverId"))
	filter := bson.D{{"ServerID", serverId}}
	if values.Has("ip") && !viper.GetBool("IgnoreClientIP") {
		ip := string(values.Peek("ip"))
		reip, err := GetIP(ctx)
		if err != nil {
			panic(err)
		}
		if ip != reip {
			ctx.Response.SetStatusCode(403)
			fmt.Fprint(ctx, "{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"访问IP错误\"}")
			return
		}
	}
	var id define.DB_ServerID
	err := mongodb.Collection_serverid.FindOne(context.TODO(), filter).Decode(&id)
	if err == mongo.ErrNoDocuments {
		ctx.Response.SetStatusCode(403)
		fmt.Fprint(ctx, "{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"serverID不存在\"}")
		return
	}
	if err != nil {
		panic(err)
	}
	filter = bson.D{{"UUID", id.UUID}, {"Name", username}}
	var profile define.DB_profiles_profile
	err = mongodb.Collection_profiles.FindOne(context.TODO(), filter).Decode(&profile)
	if err == mongo.ErrNoDocuments {
		ctx.Response.SetStatusCode(403)
		fmt.Fprint(ctx, "{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"无此角色\"}")
		return
	}
	if err != nil {
		panic(err)
	}
	ctx.Response.SetStatusCode(200)

	var profilec define.Profile
	profilec.UUID = profile.UUID
	profilec.Name = profile.Name
	encoded, error := jsoniter.Marshal(profilec)
	if error != nil {
		panic(error)
	}
	ctx.Write(encoded)
}
