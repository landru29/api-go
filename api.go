package main

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"

	"github.com/landru29/api-go/routes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var _currentMongoSession *mgo.Session

var mainCommand = &cobra.Command{
	Use:   "api-go",
	Short: "API by noopy",
	Long:  "Full API by noopy",
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetEnvPrefix("noopy")
		viper.AutomaticEnv()
		viper.SetConfigType("json")
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println(err.Error())
		}

		// Application statup here
		ConnectMongo(viper.GetString("mongo_host"), viper.GetString("mongo_port"), viper.GetString("mongo_user"), viper.GetString("mongo_password"), viper.GetString("mongo_db_name"))
		router := routes.DefineRoutes()
		router.Run(":" + viper.GetString("api_port"))
	},
}

// GetMongoSession get the current mongo session
func GetMongoSession() *mgo.Session {
	return _currentMongoSession
}

// ConnectMongo connect to the mongo database
func ConnectMongo(host, port, user, password, name string) (session *mgo.Session, err error) {
	mongooseConnectionChain := "mongodb://" +
		map[bool]string{true: user, false: ""}[len(user) > 0] +
		map[bool]string{true: ":" + password, false: ""}[len(password) > 0] +
		map[bool]string{true: "@", false: ""}[len(user) > 0] +
		host +
		map[bool]string{true: ":" + port, false: ""}[len(port) > 0] +
		map[bool]string{true: "/" + name, false: ""}[len(name) > 0]

	_currentMongoSession, err = mgo.Dial(mongooseConnectionChain)

	if err != nil {
		fmt.Println(mongooseConnectionChain)
		panic(err.Error())
	}

	return
}

func init() {
	flags := mainCommand.Flags()
	flags.String("api-host", "your-api-host", "API host")
	flags.String("api-port", "3000", "API port")

	flags.String("mongo-host", "127.0.0.1", "MongoDb host")
	flags.String("mongo-port", "27017", "MongoDb port")
	flags.String("mongo-user", "your-mongodb-user", "MongoDb user")
	flags.String("mongo-password", "your-mongodb-password", "MongoDb password")
	flags.String("mongo-db-name", "your-mongodb-dbname", "MongoDb base name")

	viper.BindPFlag("api_host", flags.Lookup("api-host"))
	viper.BindPFlag("api_port", flags.Lookup("api-port"))

	viper.BindPFlag("mongo_host", flags.Lookup("mongo-host"))
	viper.BindPFlag("mongo_port", flags.Lookup("mongo-port"))
	viper.BindPFlag("mongo_user", flags.Lookup("mongo-user"))
	viper.BindPFlag("mongo_password", flags.Lookup("mongo-password"))
	viper.BindPFlag("mongo_db_name", flags.Lookup("mongo-db-name"))
}

func main() {
	mainCommand.Execute()

}

/*

   r := gin.Default()

   r.GET("/ping", func(c *gin.Context) {
       c.JSON(200, gin.H{
           "message": "pong",
       })
   })
   r.Run()*/
