package define

import "go.mongodb.org/mongo-driver/bson/primitive"

type DB_User struct {
	ID                 primitive.ObjectID `bson:"_id"`
	PerferedLanguage   string
	Username           string
	Email              string
	PasswordHashed     string
	UserGroup          string
	Profiles           []string
	CertificationCodes []string
	Iterations         int
	Salt               string
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
	PerferedLanguage   string
	Username           string
	Email              string
	PasswordHashed     string
	UserGroup          string
	Profiles           []string
	CertificationCodes []string
	Iterations         int
	Salt               string
}
