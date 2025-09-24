package yggdrasilapi

import (
	"context"
	"craapi/packages/define"
	"craapi/packages/log"
	"craapi/packages/mongodb"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProfileSearch(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	uuid := ctx.UserValue("uuid").(string)
	filter := bson.D{{"UUID", uuid}}
	var profile define.DB_profiles_profile
	err := mongodb.Collection_profiles.FindOne(context.TODO(), filter).Decode(&profile)
	if err == mongo.ErrNoDocuments {
		ctx.Response.SetStatusCode(204)
		if viper.GetBool("Debug") {
			log.LOGD(
				"ProfileSearch request: ",
				string(ctx.Request.Body()),
				"无此角色",
			)
		}
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
	if viper.GetBool("Debug") {
		log.LOGD(
			"ProfileSearch uuid: ",
			uuid,
			" respon: ",
			string(encoded),
		)
	}
}
func ProfileMutiSearch(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var profiles []string
	names := jsoniter.Get(ctx.Request.Body())
	names.ToVal(&profiles)
	var filter bson.D
	var profile define.DB_profiles_profile
	respon := "["
	visited := false
	var profilec define.Profile
	for i, value := range profiles {
		if i > viper.GetInt("yggdrasilapi.SearchProfileMAX") {
			break
		}
		filter = bson.D{{"Name", value}}
		err := mongodb.Collection_profiles.FindOne(context.TODO(), filter).Decode(&profile)
		if err != nil && err != mongo.ErrNoDocuments {
			panic(err)
		}
		if err == nil {
			profilec.UUID = profile.UUID
			profilec.Name = profile.Name
			encoded, error := jsoniter.Marshal(profilec)
			if error != nil {
				panic(error)
			}
			if visited {
				respon += ","
			} else {
				visited = true
			}
			respon += string(encoded)
		}
	}
	respon += "]"
	ctx.Response.SetStatusCode(200)
	ctx.WriteString(respon)
	if viper.GetBool("Debug") {
		log.LOGD(
			"ProfileMutiSearch request: ",
			string(ctx.Request.Body()),
			" respon: ",
			respon,
		)
	}
}
