package panel

import (
	"context"
	"craapi/packages/define"
	"craapi/packages/encryption"
	"craapi/packages/mongodb"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register_page_get(ctx *fasthttp.RequestCtx) {
	staticfileget(ctx, "register.tpl")
}

func Register_page_post(ctx *fasthttp.RequestCtx) {
	verificationCode := string(ctx.FormValue("verificationCode"))
	if verificationCode != "hMF8bHDeFrtWBs9d2k0G" {
		ctx.WriteString("验证码错误，请刷新重试")
		return
	}
	account := string(ctx.FormValue("account"))
	email := string(ctx.FormValue("Email"))
	password := string(ctx.FormValue("password"))
	filter := bson.D{{"Username", account}}
	err := mongodb.Collection_users.FindOne(context.TODO(), filter).Err()
	if err == nil {
		ctx.WriteString("该用户名已注册")
		return
	}
	if err != mongo.ErrNoDocuments {
		panic(err)
	}
	filter = bson.D{{"Email", email}}
	err = mongodb.Collection_users.FindOne(context.TODO(), filter).Err()
	if err == mongo.ErrNoDocuments {
		hashedpass := encryption.CreatePassword(password)
		var ans define.DB_User_nonid
		ans.PasswordHashed = hashedpass.Hashedpassword
		ans.Salt = hashedpass.Salt
		ans.Iterations = hashedpass.Iterations
		ans.Email = email
		ans.UserGroup = defaultconfig.DefaultUserGroup
		ans.Username = account
		uuidGen, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		}
		uuidstr := hex.EncodeToString(uuidGen[:])
		filter = bson.D{{"UUID", uuidstr}}
		err = mongodb.Collection_profiles.FindOne(context.TODO(), filter).Err()
		for err != mongo.ErrNoDocuments {
			if err != nil {
				panic(err)
			}
			uuidGen, err = uuid.NewRandom()
			if err != nil {
				panic(err)
			}
			uuidstr = hex.EncodeToString(uuidGen[:])
			filter = bson.D{{"UUID", uuidstr}}
			err = mongodb.Collection_profiles.FindOne(context.TODO(), filter).Err()
		}
		ans.Profiles = []string{uuidstr}
		result, err := mongodb.Collection_users.InsertOne(context.TODO(), ans)
		if err != nil {
			panic(err)
		}
		fmt.Println("user oid:", result, "Email:", email, "inserted completed")
		var newprofile define.DB_profiles_profile
		newprofile.Name = account
		newprofile.UUID = uuidstr
		result, err = mongodb.Collection_profiles.InsertOne(context.TODO(), newprofile)
		if err != nil {
			panic(err)
		}
		fmt.Println("Profile oid:", result, "Name:", account, "inserted completed")
		ctx.WriteString("用户：" + account + "注册成功，UUID：" + uuidstr)

		return
	} else if err == nil {
		//邮箱已经注册
		return
	}
	panic(err)
}
