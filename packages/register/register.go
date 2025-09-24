package register

import (
	"bytes"
	"context"
	"craapi/packages/define"
	"craapi/packages/encryption"
	"craapi/packages/genrand"
	"craapi/packages/mongodb"
	"craapi/packages/overloaded"
	"encoding/hex"
	"fmt"
	"mime"
	"text/template"
	"time"

	"github.com/dlclark/regexp2"

	"github.com/google/uuid"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrEmailAlreadyExists = fmt.Errorf("email already exists")
var ErrEmailOverloaded = fmt.Errorf("email overloaded")
var ErrWrongAuthCode = fmt.Errorf("email auth code wrong or not exist")
var ErrWrongFormat = fmt.Errorf("email, username or password format wrong")

var defaultusergroup string
var email_ch *chan *email.Email

var tmpl *template.Template

func isValidEmail(email string) bool {
	if len(email) < 2 {
		return false
	}
	regexPattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp2.MustCompile(regexPattern, 0)
	match, err := re.MatchString(email)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	return match
}
func isValidpassword(pass string) bool {
	regexPattern := `^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z_]{8,16}`
	re := regexp2.MustCompile(regexPattern, 0)
	match, err := re.MatchString(pass)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	return match
}
func isValidusername(username string) bool {
	if len(username) < 4 || len(username) > 20 {
		return false
	}
	return true
}

func Init(defaultUserGroup string, Email_ch *chan *email.Email) {
	defaultusergroup = "user"
	email_ch = Email_ch
	var err error
	tmpl, err = template.ParseFiles("templates/regemail.html")
	if err != nil {
		fmt.Println("create template failed, err", err)
		panic(err)
	}
}

func Register(username string, password string, email string, authcode string) error {
	if (!isValidEmail(email)) || (!isValidpassword(password)) || (!isValidusername(username)) {
		return ErrWrongFormat
	}
	filter := bson.D{{"Email", email}}
	err := mongodb.Collection_users.FindOne(context.TODO(), filter).Err()
	if err == mongo.ErrNoDocuments {
		if overloaded.Overload(filter) {
			return ErrEmailOverloaded
		} else {
			var document define.DB_EmailOverloaded
			document.Email = email
			document.Time = primitive.NewDateTimeFromTime(time.Now())
			_, err = mongodb.Collection_overloaded.InsertOne(context.TODO(), document)
			if err != nil {
				panic(err)
			}
		}
		filter = bson.D{{"Email", email}, {"AuthCode", authcode}}
		err := mongodb.Collection_emailauthcode.FindOne(context.TODO(), filter).Err()
		if err == mongo.ErrNoDocuments {
			return ErrWrongAuthCode
		}
		if err != nil {
			panic(err)
		}
		hashedpass := encryption.CreatePassword(password)
		var ans define.DB_User_nonid
		ans.Username = username
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
		newprofile.Name = username
		newprofile.UUID = uuidstr
		newprofile.Textures = ""
		newprofile.SkinUploadable = false
		newprofile.CapeUploadable = false
		result, err = mongodb.Collection_profiles.InsertOne(context.TODO(), newprofile)
		if err != nil {
			panic(err)
		}
		fmt.Println("Profile oid:", result, "Name:", username, "inserted completed")
		return nil
	} else if err == nil {
		//邮箱已经注册
		return ErrEmailAlreadyExists
	}
	panic(err)
}

func EmailAuth(emailaddr string) error {
	if !isValidEmail(emailaddr) {
		return ErrWrongFormat
	}
	// fmt.Print("sending to ", emailaddr)
	filter := bson.D{{"Email", emailaddr}}
	err := mongodb.Collection_users.FindOne(context.TODO(), filter).Err()
	if err == nil {
		return ErrEmailAlreadyExists
	} else if err != mongo.ErrNoDocuments {
		panic(err)
	}
	filter = bson.D{{"Email", emailaddr}, {"Time", bson.D{{"$gt", primitive.NewDateTimeFromTime(time.Now().Add(-time.Minute))}}}}
	err = mongodb.Collection_emailauthcode.FindOne(context.TODO(), filter).Err()
	if err == nil {
		return ErrEmailOverloaded
	} else if err != mongo.ErrNoDocuments {
		panic(err)
	}
	var tpl bytes.Buffer
	authcode := genrand.GenerateRandomString(8)
	err = tmpl.Execute(&tpl, authcode)
	if err != nil {
		panic(err)
	}
	e := email.NewEmail()
	e.From = mime.QEncoding.Encode("UTF-8", viper.GetString("smtp.EmailName")+" <"+viper.GetString("smtp.Address")+">")
	e.To = []string{emailaddr}
	e.Subject = "CRA API 验证码"
	e.HTML = tpl.Bytes()
	*email_ch <- e
	var putin define.DB_auth_code
	putin.Email = emailaddr
	putin.AuthCode = authcode
	putin.Time = primitive.NewDateTimeFromTime(time.Now())
	a, err := mongodb.Collection_emailauthcode.InsertOne(context.TODO(), putin)
	if err != nil {
		panic(err)
	}
	fmt.Println(a.InsertedID)
	return nil
}
