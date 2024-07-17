package panel

import (
	"context"
	"craapi/cmd/packages/define"
	"craapi/cmd/packages/encryption"
	"fmt"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register_page_get(ctx *fasthttp.RequestCtx) {
	staticfileget(ctx, "register.tpl")
}

func Register_page_post(ctx *fasthttp.RequestCtx) {
	username := string(ctx.FormValue("username"))
	email := string(ctx.FormValue("email"))
	password := string(ctx.FormValue("password"))
	filter := bson.D{{"Username", username}}
	err := collection_users.FindOne(context.TODO(), filter).Err()
	if err == mongo.ErrNoDocuments {
		filter = bson.D{{"Email", email}}
		err := collection_users.FindOne(context.TODO(), filter).Err()
		if err == mongo.ErrNoDocuments {
			hashedpass := encryption.CreatePassword(password)
			var ans define.DB_User_nonid
			ans.PasswordHashed = hashedpass.Hashedpassword
			ans.Salt = hashedpass.Salt
			ans.Iterations = hashedpass.Iterations
			ans.Email = email
			ans.UserGroup = defaultconfig.DefaultUserGroup
			result, err := collection_users.InsertOne(context.TODO(), ans)
			if err != nil {
				panic(err)

			}
			fmt.Println("user oid:", result, "name:", username, "inserted completed")
			return
		} else if err == nil {
			//邮箱已经注册
			return
		}
		panic(err)
	} else if err == nil {
		// 用户名已经注册
		return
	}
	panic(err)
}
