package define

import "go.mongodb.org/mongo-driver/bson/primitive"

type DB_User struct {
	ID               primitive.ObjectID `bson:"_id"`
	PerferedLanguage string             `bson:"PerferedLanguage"`
	Username         string             `bson:"Username"`
	Email            string             `bson:"Email"`
	PasswordHashed   string             `bson:"PasswordHashed"`
	UserGroup        string             `bson:"UserGroup"`
	Profiles         []string           `bson:"Profiles"`
	Iterations       int                `bson:"Iterations"`
	Salt             string             `bson:"Salt"`
}
type DB_EmailOverloaded struct {
	Email string             `bson:"UID"`
	Time  primitive.DateTime `bson:"Time"`
}
type DB_Overloaded struct {
	UserID primitive.ObjectID `bson:"UID"`
	Time   primitive.DateTime `bson:"Time"`
}
type DB_Cookie struct {
	UserID primitive.ObjectID `bson:"UID"`
	Cookie string             `bson:"Cookie"`
	Time   primitive.DateTime `bson:"Time"`
}
type DB_User_nonid struct {
	PerferedLanguage string   `bson:"PerferedLanguage"`
	Username         string   `bson:"Username"`
	Email            string   `bson:"Email"`
	PasswordHashed   string   `bson:"PasswordHashed"`
	UserGroup        string   `bson:"UserGroup"`
	Profiles         []string `bson:"Profiles"`
	Iterations       int      `bson:"Iterations"`
	Salt             string   `bson:"Salt"`
}
type DB_auth_code struct {
	Email    string             `bson:"Email"`
	AuthCode string             `bson:"AuthCode"`
	Time     primitive.DateTime `bson:"Time"`
}
type DB_profiles_profile struct {
	UUID           string `bson:"UUID"`
	Name           string `bson:"Name"`
	Textures       string `bson:"Textures"`
	SkinUploadable bool   `bson:"SkinUploadable"`
	CapeUploadable bool   `bson:"CapeUploadable"`
}
type DB_Token struct {
	AccessToken string             `bson:"AccessToken"`
	ClientToken string             `bson:"ClientToken"`
	UUID        string             `bson:"UUID"`
	Validate    bool               `bson:"Validate"`
	Time        primitive.DateTime `bson:"Time"`
	UserID      primitive.ObjectID `bson:"UID"`
}
type DB_Token_webaccess struct {
	WebAccess string             `bson:"WebAccess"`
	Time      primitive.DateTime `bson:"Time"`
	UserID    string             `bson:"UserID"`
}
type Properties struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature"`
}
type Profile struct {
	UUID       string        `json:"id"`
	Name       string        `json:"name"`
	Properties []*Properties `json:"properties,omitempty"`
}
type DB_ServerID struct {
	ServerID    string             `bson:"ServerID"`
	AccessToken string             `bson:"AccessToken"`
	UUID        string             `bson:"UUID"`
	Time        primitive.DateTime `bson:"Time"`
}
type DB_FilesMOD struct {
	FileName string             `bson:"FileName"`
	Time     primitive.DateTime `bson:"Time"`
	SHA256   string             `bson:"AccessToken"`
	Name     string             `bson:"Name"`
	Version  string             `bson:"Version"`
}
type DB_FilesTexturesSKIN struct {
	FileName string             `bson:"FileName"`
	Time     primitive.DateTime `bson:"Time"`
	SHA256   string             `bson:"AccessToken"`
	Name     string             `bson:"Name"`
	Version  string             `bson:"Version"`
}
