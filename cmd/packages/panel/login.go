package panel

import (
	"context"
	"craapi/cmd/packages/define"
	"craapi/cmd/packages/encryption"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"time"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CookieShow(key []byte, value []byte) {
	fmt.Println("Cookie:", string(key), string(value))

}

func Login_page_get(ctx *fasthttp.RequestCtx) {
	// fmt.Println(ctx.Request.Header.String())
	cookieid := string(ctx.Request.Header.Cookie("id"))
	filter := bson.D{{"Cookie", cookieid}}
	var cookiehandle define.DB_Cookie
	err := collection_cookie.FindOne(context.TODO(), filter).Decode(&cookiehandle)
	if err == nil {
		// filter = bson.D{{"_id", cookiehandle.UserID}}
		// var res define.DB_User
		// err = collection_users.FindOne(context.TODO(), filter).Decode(&res)
		// ctx.WriteString("欢迎回来，" + res.Username)
		ctx.Redirect("/", 307)
		ctx.SetStatusCode(fasthttp.StatusTemporaryRedirect)
		return
	}
	ctx.SetContentType("text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		fmt.Println("create template failed, err", err)
		ctx.Redirect("/static/403.html", 307)
		ctx.SetStatusCode(fasthttp.StatusTemporaryRedirect)
		return
	}
	Wrongpass := false
	tmpl.Execute(ctx, Wrongpass)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func Login_page_post(ctx *fasthttp.RequestCtx) {
	usernameoremail := string(ctx.FormValue("username"))
	password := string(ctx.FormValue("password"))
	filter := bson.D{{"Username", usernameoremail}}
	var res define.DB_User
	err := collection_users.FindOne(context.TODO(), filter).Decode(&res)
	if err == mongo.ErrNoDocuments {
		filter = bson.D{{"Email", usernameoremail}}
		err = collection_users.FindOne(context.TODO(), filter).Decode(&res)
		if err == mongo.ErrNoDocuments {
			ctx.SetContentType("text/html; charset=utf-8")
			tmpl, err := template.ParseFiles("templates/login.html")
			if err != nil {
				fmt.Println("create template failed, err", err)
				ctx.Redirect("/static/403.html", 307)
				ctx.SetStatusCode(fasthttp.StatusTemporaryRedirect)
				return
			}
			Wrongpass := true
			tmpl.Execute(ctx, Wrongpass)
			ctx.SetStatusCode(fasthttp.StatusOK)
			return
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

		salt := make([]byte, 32)
		rand.Read(salt)
		tmp := sha256.Sum256([]byte(res.Email + string(salt)))
		cookieid := base64.URLEncoding.EncodeToString(tmp[:])
		var cookie fasthttp.Cookie
		cookie.SetMaxAge(30 * 86400)
		cookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
		cookie.SetValue(cookieid)
		cookie.SetKey("id")
		fmt.Println(cookie.Key())
		ctx.Response.Header.SetCookie(&cookie)

		_, err = collection_cookie.InsertOne(context.TODO(), define.DB_Cookie{UserID: res.ID, Time: primitive.NewDateTimeFromTime(time.Now()), Cookie: cookieid})
		if err != nil {
			panic(err)
		}
		fmt.Println("登陆成功！ 账户：", usernameoremail, "CookieID:", cookieid)
		ctx.Redirect("/", 303)
		ctx.SetStatusCode(fasthttp.StatusSeeOther)
		return
	}

	ctx.SetContentType("text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		fmt.Println("create template failed, err", err)
		ctx.Redirect("/static/403.html", 307)
		ctx.SetStatusCode(fasthttp.StatusTemporaryRedirect)
		return
	}
	Wrongpass := true
	tmpl.Execute(ctx, Wrongpass)
	ctx.SetStatusCode(fasthttp.StatusOK)

}
