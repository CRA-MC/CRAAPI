package overloaded

import (
	"context"
	"craapi/packages/define"
	"craapi/packages/mongodb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Overload(filter bson.D) bool {
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
	return false
}

func Overload_id(oid *primitive.ObjectID) bool {

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
