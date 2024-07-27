package register

import (
	"bytes"
	"context"
	"craapi/cmd/packages/define"
	"craapi/cmd/packages/encryption"
	"craapi/cmd/packages/mongodb"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"mime"
	"text/template"

	"crypto/rand"

	"github.com/google/uuid"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection_users, collection_profiles, collection_emailoverloaded *mongo.Collection

var ErrEmailAlreadyExists = fmt.Errorf("email already exists")
var ErrEmailOverloaded = fmt.Errorf("email overloaded")

var defaultusergroup string
var email_ch *chan *email.Email

var tmpl *template.Template

func Init(defaultUserGroup string, Email_ch *chan *email.Email) {
	collection_users = mongodb.Mongodb_GetCollection("Users", "craapi")
	collection_emailoverloaded = mongodb.Mongodb_GetCollection("Email Overloaded", "craapi")
	collection_profiles = mongodb.Mongodb_GetCollection("Profiles", "craapi")
	defaultusergroup = defaultUserGroup
	email_ch = Email_ch
	var err error
	tmpl, err = template.ParseFiles("templates/regemail.html")
	if err != nil {
		fmt.Println("create template failed, err", err)
		panic(err)
	}
}

func Register(username string, password string, email string) error {
	filter := bson.D{{"Email", email}}
	err := collection_users.FindOne(context.TODO(), filter).Err()
	if err == mongo.ErrNoDocuments {
		hashedpass := encryption.CreatePassword(password)
		var ans define.DB_User_nonid
		ans.PasswordHashed = hashedpass.Hashedpassword
		ans.Salt = hashedpass.Salt
		ans.Iterations = hashedpass.Iterations
		ans.Email = email
		ans.UserGroup = defaultusergroup
		uuidGen, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		}
		uuidstr := hex.EncodeToString(uuidGen[:])
		filter = bson.D{{"UUID", uuidstr}}
		err = collection_profiles.FindOne(context.TODO(), filter).Err()
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
			err = collection_profiles.FindOne(context.TODO(), filter).Err()
		}
		ans.Profiles = []string{uuidstr}

		result, err := collection_users.InsertOne(context.TODO(), ans)
		if err != nil {
			panic(err)
		}
		fmt.Println("user oid:", result, "Email:", email, "inserted completed")
		return nil
	} else if err == nil {
		//邮箱已经注册
		return ErrEmailAlreadyExists
	}
	panic(err)
}
func EmailAuth(emailaddr string) error {
	filter := bson.D{{"Email", emailaddr}}
	err := collection_users.FindOne(context.TODO(), filter).Err()
	if err == nil {
		return ErrEmailAlreadyExists
	} else if err != mongo.ErrNoDocuments {
		panic(err)
	}
	err = collection_emailoverloaded.FindOne(context.TODO(), filter).Err()
	if err == nil {
		return ErrEmailOverloaded
	} else if err != mongo.ErrNoDocuments {
		panic(err)
	}
	var tpl bytes.Buffer
	salt := make([]byte, 32)
	rand.Read(salt)
	authcode := (base64.URLEncoding.EncodeToString(salt)[0:6])
	err = tmpl.Execute(&tpl, authcode)
	if err != nil {
		panic(err)
	}
	e := email.NewEmail()
	e.From = mime.QEncoding.Encode("UTF-8", viper.GetString("smtp.EmailName")+" <"+viper.GetString("smtp.Address")+">")
	e.To = []string{viper.GetString("smtp.DebugSendAddress")}
	e.Subject = "Craapi Test"
	e.Text = tpl.Bytes()
	*email_ch <- e
	return nil
}
