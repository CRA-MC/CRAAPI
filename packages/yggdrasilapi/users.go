package yggdrasilapi

import (
	"context"
	"craapi/packages/define"
	"craapi/packages/encryption"
	"craapi/packages/log"
	"craapi/packages/mongodb"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

type Agent struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

type Auth struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	ClientToken string `json:"clientToken"`
	RequestUser bool   `json:"requestUser,omitempty"`
	Agent       *Agent `json:"agent"`
}
type SignoutAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Properties struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type User struct {
	Username   string        `json:"id"`
	Properties []*Properties `json:"properties"`
}

type AuthVail struct {
	ClientToken string `json:"clientToken,omitempty"`
	AccessToken string `json:"accessToken"`
}

func overload(oid *primitive.ObjectID) bool {

	filter := bson.D{{"UID", *oid}}
	err := mongodb.Collection_overloaded.FindOne(context.TODO(), filter).Err()
	total_time := 0
	if err == nil {
		var overloaded define.DB_Overloaded
		for err == nil {
			time.Sleep(6 * time.Second)
			total_time += 6
			if total_time > 30 {
				return true
			}
			err = mongodb.Collection_overloaded.FindOne(context.TODO(), filter).Decode(&overloaded)
			if time.Since(overloaded.Time.Time()).Seconds() > 5 {
				mongodb.Collection_overloaded.DeleteOne(context.TODO(), filter)
				break
			}
		}
	}
	if err != mongo.ErrNoDocuments && err != nil {
		panic(err)
	}

	_, err = mongodb.Collection_overloaded.InsertOne(context.TODO(), define.DB_Overloaded{UserID: *oid, Time: primitive.NewDateTimeFromTime(time.Now())})
	if err != nil {
		panic(err)
	}
	return false
}

func UserAuth(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var autho Auth
	err := jsoniter.Unmarshal(ctx.Request.Body(), &autho)
	if err != nil {
		panic(err)
	}
	filter := bson.D{{"Username", autho.Username}}
	var res define.DB_User
	err = mongodb.Collection_users.FindOne(context.TODO(), filter).Decode(&res)
	if err == mongo.ErrNoDocuments {
		filter = bson.D{{"Email", autho.Username}}
		err = mongodb.Collection_users.FindOne(context.TODO(), filter).Decode(&res)
		if err == mongo.ErrNoDocuments {
			ctx.Response.SetStatusCode(403)
			fmt.Fprint(ctx, "{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"密码不匹配或者用户名不存在\"}")
			if viper.GetBool("Debug") {
				if viper.GetBool("DebugShowPass") {
					log.LOGD(
						"UserAuth request: ",
						string(ctx.Request.Body()),
						" respon: ",
						"{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"密码不匹配或者用户名不存在\"}",
					)
				} else {
					autho.Password = "*********"
					s, err := jsoniter.MarshalToString(&autho)
					if err != nil {
						panic(err)
					}
					log.LOGD(
						"UserAuth request: ",
						s,
						" respon: ",
						"{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"密码不匹配或者用户名不存在\"}",
					)
				}
			}
			return
		} else if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	if overload(&res.ID) {
		ctx.Response.SetStatusCode(429)
		fmt.Fprint(ctx, "{\"error\":\"Too Many Requests\",\"errorMessage\":\"请求超时\"}")
		if viper.GetBool("Debug") {
			if viper.GetBool("DebugShowPass") {
				log.LOGD(
					"UserAuth request: ",
					string(ctx.Request.Body()),
					" respon: ",
					"{\"error\":\"Too Many Requests\",\"errorMessage\":\"请求超时\"} User: ", res.Username,
				)
			} else {
				autho.Password = "*********"
				s, err := jsoniter.MarshalToString(&autho)
				if err != nil {
					panic(err)
				}
				log.LOGD(
					"UserAuth request: ",
					s,
					" respon: ",
					"{\"error\":\"Too Many Requests\",\"errorMessage\":\"请求超时\"} User: ", res.Username,
				)
			}
		}
		return
	}

	if encryption.CheckPassowrd(autho.Password, res.Salt, res.Iterations, res.PasswordHashed) {
		token := tokengen(autho.ClientToken, res.ID, res.Profiles[0])
		filter = bson.D{{"UID", res.ID}}
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
		respon := "{\"accessToken\":\"" + token.AccessToken + "\",\"clientToken\":\"" + token.ClientToken + "\",\"availableProfiles\":["
		var profile define.DB_profiles_profile
		var profilec define.Profile
		var selected string
		for i, uuid := range res.Profiles {
			filter = bson.D{{"UUID", uuid}}
			err = mongodb.Collection_profiles.FindOne(context.TODO(), filter).Decode(&profile)
			if err != nil {
				panic(err)
			}
			profilec.UUID = uuid
			profilec.Name = profile.Name
			encoded, error := jsoniter.Marshal(profilec)
			if error != nil {
				panic(error)
			}
			if i == 0 {
				selected = string(encoded)
			} else {
				respon += ","
			}
			respon += string(encoded)
		}
		respon += "],\"selectedProfile\":" + selected
		if autho.RequestUser {
			user := User{
				Username:   res.ID.Hex(),
				Properties: []*Properties{{"preferredLanguage", "zh_CN"}},
			}
			tmp, err := jsoniter.Marshal(user)
			if err != nil {
				panic(err)
			}
			respon += ",\"user\":" + string(tmp)
		}
		respon += "}"
		ctx.Response.SetStatusCode(200)
		fmt.Fprint(ctx, respon)
		if viper.GetBool("Debug") {
			if viper.GetBool("DebugShowPass") {
				log.LOGD(
					"UserAuth request: ",
					string(ctx.Request.Body()),
					" respon: ",
					respon,
				)
			} else {
				autho.Password = "*********"
				s, err := jsoniter.MarshalToString(&autho)
				if err != nil {
					panic(err)
				}
				log.LOGD(
					"UserAuth request: ",
					s,
					" respon: ",
					respon,
				)
			}
		}
		return
	}
	ctx.Response.SetStatusCode(403)
	fmt.Fprint(ctx, "{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"密码不匹配或者用户名不存在\"}")
	if viper.GetBool("Debug") {
		if viper.GetBool("DebugShowPass") {
			log.LOGD(
				"UserAuth request: ",
				string(ctx.Request.Body()),
				" respon: ",
				"{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"密码不匹配或者用户名不存在\"}",
			)
		} else {
			autho.Password = "*********"
			s, err := jsoniter.MarshalToString(&autho)
			if err != nil {
				panic(err)
			}
			log.LOGD(
				"UserAuth request: ",
				s,
				" respon: ",
				"{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"密码不匹配或者用户名不存在\"}",
			)
		}
	}
}

func TokenRefresh(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	accessTokenany := jsoniter.Get(ctx.Request.Body(), "accessToken")
	if accessTokenany.LastError() != nil {
		ctx.Response.SetStatusCode(403)
		fmt.Fprint(ctx, "{\"error\":\"WrongRequest\",\"errorMessage\":\"请求有误\"}")
		if viper.GetBool("Debug") {
			log.LOGD(
				"TokenRefresh request: ",
				string(ctx.Request.Body()),
				" respon: ",
				"{\"error\":\"WrongRequest\",\"errorMessage\":\"请求有误\"}",
			)
		}
		return
	}
	accessToken := accessTokenany.ToString()
	clientTokenany := jsoniter.Get(ctx.Request.Body(), "clientToken")
	clientToken := ""
	if clientTokenany.LastError() == nil {
		clientToken = clientTokenany.ToString()
	}
	requestUserany := jsoniter.Get(ctx.Request.Body(), "requestUser")
	if clientTokenany.LastError() == nil {
		clientToken = clientTokenany.ToString()
	}
	requestUser := requestUserany.ToBool()
	if tokencheckexist(accessToken, clientToken) {
		filter := bson.D{{"AccessToken", accessToken}}
		var res define.DB_Token
		err := mongodb.Collection_token.FindOne(context.TODO(), filter).Decode(&res)
		if err != nil {
			panic(err)
		}
		var user define.DB_User
		filter = bson.D{{"_id", res.UserID}}
		err = mongodb.Collection_users.FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			panic(err)
		}
		selectedprofile := user.Profiles[0]
		selectedProfileany := jsoniter.Get(ctx.Request.Body(), "selectedProfile")
		if selectedProfileany.LastError() == nil {
			var profile define.Profile
			selectedProfileany.ToVal(profile)
			if IsContain(user.Profiles, profile.UUID) {
				selectedprofile = profile.UUID
			}
		}
		token := tokengen(clientToken, res.UserID, selectedprofile)
		_, err = mongodb.Collection_token.InsertOne(context.TODO(), token)
		if err != nil {
			panic(err)
		}
		filter = bson.D{{"AccessToken", res.AccessToken}}
		_, err = mongodb.Collection_token.DeleteOne(context.TODO(), filter)
		if err != nil {
			panic(err)
		}
		respon := "{\"accessToken\":\"" + token.AccessToken + "\",\"clientToken\":\"" + token.ClientToken + "\",\"selectedProfile\":"
		var profile define.DB_profiles_profile
		var profilec define.Profile
		filter = bson.D{{"UUID", token.UUID}}
		err = mongodb.Collection_profiles.FindOne(context.TODO(), filter).Decode(&profile)
		if err != nil {
			panic(err)
		}
		profilec.UUID = profile.UUID
		profilec.Name = profile.Name
		encoded, error := jsoniter.Marshal(profilec)
		if error != nil {
			panic(error)
		}
		respon += string(encoded)
		if requestUser {
			respon += ",\"user\":{\"id\":\"" + user.ID.Hex() + "\",\"properties\":[{\"name\":\"preferredLanguage\",\"value\":\"zh_CN\"}]}"
		}
		respon += "}"
		ctx.Response.SetStatusCode(200)
		fmt.Fprint(ctx, respon)
		if viper.GetBool("Debug") {
			log.LOGD(
				"TokenRefresh request: ",
				string(ctx.Request.Body()),
				" respon: ",
				respon,
			)
		}
	} else {
		ctx.Response.SetStatusCode(403)
		fmt.Fprint(ctx, "{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"无此令牌或者令牌失效\"}")
		if viper.GetBool("Debug") {
			log.LOGD(
				"TokenRefresh request: ",
				string(ctx.Request.Body()),
				" respon: ",
				"{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"无此令牌或者令牌失效\"}",
			)
		}
		return
	}
}
func TokenVail(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var autho AuthVail
	err := jsoniter.Unmarshal(ctx.Request.Body(), &autho)
	if err != nil {
		panic(err)
	}
	if tokencheck(autho.AccessToken, autho.ClientToken) {
		ctx.Response.SetStatusCode(204)
		if viper.GetBool("Debug") {
			log.LOGD(
				"TokenVail request: ",
				string(ctx.Request.Body()),
				"经过验证 该令牌有效",
			)
		}
		return
	}
	ctx.Response.SetStatusCode(403)
	fmt.Fprint(ctx, "{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"无此令牌或者令牌失效\"}")
	if viper.GetBool("Debug") {
		log.LOGD(
			"TokenVail request: ",
			string(ctx.Request.Body()),
			" respon: ",
			"{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"无此令牌或者令牌失效\"}",
		)
	}
}
func TokeninValid(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var autho AuthVail
	err := jsoniter.Unmarshal(ctx.Request.Body(), &autho)
	if err != nil {
		panic(err)
	}
	filter := bson.D{{"AccessToken", autho.AccessToken}}
	_, err = mongodb.Collection_token.DeleteOne(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	ctx.Response.SetStatusCode(204)
	if viper.GetBool("Debug") {
		log.LOGD(
			"TokeninValid request: ",
			string(ctx.Request.Body()),
			" 令牌: ",
			autho.AccessToken,
			"删除成功",
		)
	}
}
func Signout(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var autho Auth
	err := jsoniter.Unmarshal(ctx.Request.Body(), &autho)
	if err != nil {
		panic(err)
	}
	filter := bson.D{{"Username", autho.Username}}
	var res define.DB_User
	err = mongodb.Collection_users.FindOne(context.TODO(), filter).Decode(&res)
	if err == mongo.ErrNoDocuments {
		filter = bson.D{{"Email", autho.Username}}
		err = mongodb.Collection_users.FindOne(context.TODO(), filter).Decode(&res)
		if err == mongo.ErrNoDocuments {
			ctx.Response.SetStatusCode(403)
			fmt.Fprint(ctx, "{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"密码不匹配或者用户名不存在\"}")
			if viper.GetBool("Debug") {
				if viper.GetBool("DebugShowPass") {
					log.LOGD(
						"Signout request: ",
						string(ctx.Request.Body()),
						" respon: ",
						"{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"密码不匹配或者用户名不存在\"}",
					)
				} else {
					autho.Password = "*********"
					s, err := jsoniter.MarshalToString(&autho)
					if err != nil {
						panic(err)
					}
					log.LOGD(
						"Signout request: ",
						s,
						" respon: ",
						"{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"密码不匹配或者用户名不存在\"}",
					)
				}
			}
			return
		} else if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	if overload(&res.ID) {
		ctx.Response.SetStatusCode(429)
		fmt.Fprint(ctx, "{\"error\":\"Too Many Requests\",\"errorMessage\":\"请求超时\"}")
		if viper.GetBool("Debug") {
			if viper.GetBool("DebugShowPass") {
				log.LOGD(
					"Signout request: ",
					string(ctx.Request.Body()),
					" respon: ",
					"{\"error\":\"Too Many Requests\",\"errorMessage\":\"请求超时\"} User: ", res.Username,
				)
			} else {
				autho.Password = "*********"
				s, err := jsoniter.MarshalToString(&autho)
				if err != nil {
					panic(err)
				}
				log.LOGD(
					"Signout request: ",
					s,
					" respon: ",
					"{\"error\":\"Too Many Requests\",\"errorMessage\":\"请求超时\"} User: ", res.Username,
				)
			}
		}
		return
	}

	if encryption.CheckPassowrd(autho.Password, res.Salt, res.Iterations, res.PasswordHashed) {
		filter = bson.D{{"UID", res.ID}}
		resu, err := mongodb.Collection_token.DeleteMany(context.TODO(), filter)
		if err != nil {
			panic(err)
		}
		log.LOGD(
			"Signout Deleted ",
			resu.DeletedCount,
			" Token(s)",
		)
		ctx.Response.SetStatusCode(204)
		return
	}
	ctx.Response.SetStatusCode(403)
	fmt.Fprint(ctx, "{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"密码不匹配或者用户名不存在\"}")
	if viper.GetBool("Debug") {
		if viper.GetBool("DebugShowPass") {
			log.LOGD(
				"Signout request: ",
				string(ctx.Request.Body()),
				" respon: ",
				"{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"密码不匹配或者用户名不存在\"}",
			)
		} else {
			autho.Password = "*********"
			s, err := jsoniter.MarshalToString(&autho)
			if err != nil {
				panic(err)
			}
			log.LOGD(
				"Signout request: ",
				s,
				" respon: ",
				"{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"密码不匹配或者用户名不存在\"}",
			)
		}
	}
}
