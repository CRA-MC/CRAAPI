package api

import (
	"context"
	"craapi/packages/define"
	"craapi/packages/log"
	"craapi/packages/login"
	"craapi/packages/mongodb"
	"craapi/packages/register"

	"github.com/jordan-wright/email"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserReg struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Email    string `json:"Email"`
	AuthCode string `json:"AuthCode"`
}
type Email struct {
	Email string `json:"Email"`
}

var email_ch *chan *email.Email

func Init(Email_ch *chan *email.Email) {
	email_ch = Email_ch
}

func Api_Register(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var user UserReg
	err := jsoniter.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		panic(err)
	}
	err = register.Register(user.Username, user.Password, user.Email, user.AuthCode)
	if viper.GetBool("Debug") {
		if viper.GetBool("DebugShowPass") {
			log.LOGD(
				"Api_Register request: ",
				string(ctx.Request.Body()),
			)
		} else {
			user.Password = "*********"
			s, err := jsoniter.MarshalToString(&user)
			if err != nil {
				panic(err)
			}
			log.LOGD(
				"Api_Register request: ",
				s,
			)
		}
	}
	if err == register.ErrEmailAlreadyExists {
		ctx.SetStatusCode(403)
		ctx.WriteString("{\"error\":\"ErrEmailAlreadyExists\",\"errorMessage\":\"该邮箱已经注册过了\"}")
		return
	}
	if err == register.ErrEmailOverloaded {
		ctx.SetStatusCode(403)
		ctx.WriteString("{\"error\":\"ErrEmailOverloaded\",\"errorMessage\":\"邮箱请求次数过多\"}")
		return
	}
	if err == register.ErrWrongAuthCode {
		ctx.SetStatusCode(403)
		ctx.WriteString("{\"error\":\"ErrWrongAuthCode\",\"errorMessage\":\"验证码错误\"}")
		return
	}
	if err == register.ErrWrongFormat {
		ctx.SetStatusCode(403)
		ctx.WriteString("{\"error\":\"ErrWrongFormat\",\"errorMessage\":\"参数格式错误\"}")
		return
	}
	if err != nil {
		panic(err)
	}
	ctx.SetStatusCode(204)
}
func Api_EmailAuth(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	if viper.GetBool("Debug") {
		log.LOGD(
			"Api_EmailAuth request: ",
			string(ctx.Request.Body()),
		)
	}
	var user Email
	err := jsoniter.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		panic(err)
	}
	err = register.EmailAuth(user.Email)
	if err == register.ErrEmailAlreadyExists {
		ctx.SetStatusCode(403)
		ctx.WriteString("{\"error\":\"ErrEmailAlreadyExists\",\"errorMessage\":\"该邮箱已经注册过了\"}")
		return
	}
	if err == register.ErrEmailOverloaded {
		ctx.SetStatusCode(403)
		ctx.WriteString("{\"error\":\"ErrEmailOverloaded\",\"errorMessage\":\"邮箱请求次数过多\"}")
		return
	}
	if err == register.ErrWrongFormat {
		ctx.SetStatusCode(403)
		ctx.WriteString("{\"error\":\"ErrWrongFormat\",\"errorMessage\":\"参数格式错误\"}")
		return
	}
	if err != nil {
		panic(err)
	}
	ctx.SetStatusCode(204)
}

type UserAuth struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func Api_Login(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var user UserAuth
	err := jsoniter.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		panic(err)
	}

	err, token := login.UserAuth(user.Username, user.Password)
	if viper.GetBool("Debug") {
		if viper.GetBool("DebugShowPass") {
			log.LOGD(
				"Api_Login request: ",
				string(ctx.Request.Body()),
			)
		} else {
			user.Password = "*********"
			s, err := jsoniter.MarshalToString(&user)
			if err != nil {
				panic(err)
			}
			log.LOGD(
				"Api_Login request: ",
				s,
			)
		}
	}
	if err == login.ErrNotHaveThisUser || err == login.ErrWrongPassword {
		ctx.SetStatusCode(403)
		ctx.WriteString("{\"error\":\"ErrWrongPassword\",\"errorMessage\":\"用户名或密码错误\"}")
		return
	}
	if err == login.ErrTooManyRequests {
		ctx.SetStatusCode(403)
		ctx.WriteString("{\"error\":\"ErrOverloaded\",\"errorMessage\":\"请求次数过多\"}")
		return
	}
	if err != nil {
		panic(err)
	}
	ctx.SetStatusCode(200)
	ctx.WriteString("{\"AccessToken\":\"" + token.WebAccess + "\"}")
}

type Token struct {
	Token string `json:"AccessToken"`
}

func Api_GetUserinfo(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var user Token
	err := jsoniter.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		panic(err)
	}
	if viper.GetBool("Debug") {
		log.LOGD(
			"Api_GetUserinfo request: ",
			string(ctx.Request.Body()),
		)
	}
	filter := bson.D{{"WebAccess", user.Token}}
	var res define.DB_Token_webaccess
	err = mongodb.Collection_token.FindOne(context.TODO(), filter).Decode(&res)

	if err == mongo.ErrNoDocuments {
		ctx.SetStatusCode(404)
		ctx.WriteString("{\"error\":\"ErrInvaildToken\",\"errorMessage\":\"Token失效或不存在\"}")
		return
	} else if err != nil {
		panic(err)
	}
	filter = bson.D{{"Username", res.UserID}}
	var dbuser define.DB_User
	err = mongodb.Collection_users.FindOne(context.TODO(), filter).Decode(&dbuser)
	if err == mongo.ErrNoDocuments {
		ctx.SetStatusCode(403)
		ctx.WriteString("{\"error\":\"ErrNoUser\",\"errorMessage\":\"该用户已不存在\"}")
		return
	} else if err != nil {
		panic(err)
	}
	ctx.SetStatusCode(200)
	ctx.WriteString("{\"UserName\":\"" + dbuser.Username + "\",\"UserGroup\":\"" + dbuser.UserGroup + "\"}")
}
