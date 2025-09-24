package cmd

import (
	"context"
	"craapi/packages/mongodb"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initalize() {
	fmt.Println("config file:" + cfgFile)
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("toml")
	} else {
		viper.SetConfigFile("./craapi.config")
		viper.SetConfigType("toml")
	}
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
	mongodb.Mongodb_INIT(viper.GetString("mongoDB.User"), viper.GetString("mongoDB.Password"), viper.GetString("mongoDB.Host"), viper.GetInt("mongoDB.Port"), viper.GetString("mongoDB.Database"))
	mongodb.Mongodb_database = mongodb.Mongodb_client.Database(viper.GetString("mongoDB.Database"))
	err = mongodb.Mongodb_database.CreateCollection(context.TODO(), "Users")
	if err != nil {
		fmt.Println("[error] couldn't create collection Users", err)
		os.Exit(1)
	}
	err = mongodb.Mongodb_database.CreateCollection(context.TODO(), "Profiles")
	if err != nil {
		fmt.Println("[error] couldn't create collection Profiles", err)
		os.Exit(1)
	}
	err = mongodb.Mongodb_database.CreateCollection(context.TODO(), "Overloaded")
	if err != nil {
		fmt.Println("[error] couldn't create collection Overloaded", err)
		os.Exit(1)
	}
	err = mongodb.Mongodb_database.CreateCollection(context.TODO(), "Token")
	if err != nil {
		fmt.Println("[error] couldn't create collection Token", err)
		os.Exit(1)
	}
	err = mongodb.Mongodb_database.CreateCollection(context.TODO(), "ServerID")
	if err != nil {
		fmt.Println("[error] couldn't create collection ServerID", err)
		os.Exit(1)
	}
	err = mongodb.Mongodb_database.CreateCollection(context.TODO(), "Cookie")
	if err != nil {
		fmt.Println("[error] couldn't create collection Cookie", err)
		os.Exit(1)
	}
	err = mongodb.Mongodb_database.CreateCollection(context.TODO(), "EmailAuthCode")
	if err != nil {
		fmt.Println("[error] couldn't create collection Cookie", err)
		os.Exit(1)
	}
	coll := mongodb.Mongodb_database.Collection("Overloaded")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"Time", 1}},
		Options: options.Index().SetExpireAfterSeconds(6),
	}
	name, err := coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		fmt.Println("[error] couldn't create index of Overloaded", err)
	} else {
		fmt.Println("Name of Index Created: " + name)
	}
	coll = mongodb.Mongodb_database.Collection("Token")
	indexModel = mongo.IndexModel{
		Keys:    bson.D{{"Time", 1}},
		Options: options.Index().SetExpireAfterSeconds(604800),
	}
	name, err = coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		fmt.Println("[error] couldn't create index of Token", err)
	} else {
		fmt.Println("Name of Index Created: " + name)
	}
	coll = mongodb.Mongodb_database.Collection("Cookie")
	name, err = coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		fmt.Println("[error] couldn't create index of Cookie", err)
	} else {
		fmt.Println("Name of Index Created: " + name)
	}
	coll = mongodb.Mongodb_database.Collection("ServerID")
	indexModel = mongo.IndexModel{
		Keys:    bson.D{{"Time", 1}},
		Options: options.Index().SetExpireAfterSeconds(30),
	}
	name, err = coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		fmt.Println("[error] couldn't create index of ServerID", err)
	} else {
		fmt.Println("Name of Index Created: " + name)
	}
	coll = mongodb.Mongodb_database.Collection("EmailAuthCode")
	indexModel = mongo.IndexModel{
		Keys:    bson.D{{"Time", 1}},
		Options: options.Index().SetExpireAfterSeconds(300),
	}
	name, err = coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		fmt.Println("[error] couldn't create index of EmailAuthCode", err)
	} else {
		fmt.Println("Name of Index Created: " + name)
	}

}

var initalizeCmd = &cobra.Command{
	Use:   "init",
	Short: "initalize CRA API",
	Long:  `All software has versions.`,
	Run: func(cmd *cobra.Command, args []string) {
		initalize()
	},
}

func init() {
	initalizeCmd.Flags().StringVarP(&cfgFile, "config", "c", "./craapi.config", "config file")
	rootCmd.AddCommand(initalizeCmd)
}

/*


	Collection_profiles = mongodb_database.Collection("Profiles")
	Collection_users = mongodb_database.Collection("Users")
	Collection_overloaded = mongodb_database.Collection("Overloaded")
	Collection_token = mongodb_database.Collection("Token")
	Collection_serverid = mongodb_database.Collection("ServerID")
	Collection_cookie = mongodb_database.Collection("Cookie")

*/
