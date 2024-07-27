package api

import (
	"craapi/cmd/packages/auth"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"github.com/jordan-wright/email"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type UserLogin struct {
	username string
	password string
}
type UserReg struct {
	username string
	password string
	email    string
}
type Email struct {
	email string
}

var email_ch *chan *email.Email

func Init(Email_ch *chan *email.Email) {
	email_ch = Email_ch
}

func Api_Login(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var user UserLogin
	err := jsoniter.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		panic(err)
	}
	cor, res := auth.Auth(user.username, user.password)
	if cor {
		salt := make([]byte, 32)
		rand.Read(salt)
		tmp := sha256.Sum256([]byte(res.Email + string(salt)))
		token := base64.URLEncoding.EncodeToString(tmp[:])
		ctx.SetStatusCode(200)
		ctx.WriteString("{\"accessToken\":\"" + token + "\"}")
	} else {
		ctx.SetStatusCode(403)
		ctx.WriteString("{\"error\":\"ForbiddenOperationException\",\"errorMessage\":\"Invalid credentials. Invalid username or password.\"}")
	}
}
func Api_Register(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var user UserReg
	err := jsoniter.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		panic(err)
	}
}
func Api_EmailAuth(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	var user Email
	err := jsoniter.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		panic(err)
	}

}
