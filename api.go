package main

import (
	"fmt"

	"github.com/landru29/api-go/mongo"
	"github.com/landru29/api-go/routes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
		mongo.ConnectMongo(viper.GetString("mongo_host"), viper.GetString("mongo_port"), viper.GetString("mongo_user"), viper.GetString("mongo_password"), viper.GetString("mongo_db_name"))
		router := routes.DefineRoutes()
		router.Run(":" + viper.GetString("api_port"))
	},
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

	flags.String("default-pagination-limit", "10", "Default pagination limit")

	viper.BindPFlag("api_host", flags.Lookup("api-host"))
	viper.BindPFlag("api_port", flags.Lookup("api-port"))

	viper.BindPFlag("mongo_host", flags.Lookup("mongo-host"))
	viper.BindPFlag("mongo_port", flags.Lookup("mongo-port"))
	viper.BindPFlag("mongo_user", flags.Lookup("mongo-user"))
	viper.BindPFlag("mongo_password", flags.Lookup("mongo-password"))
	viper.BindPFlag("mongo_db_name", flags.Lookup("mongo-db-name"))

	viper.BindPFlag("default_pagination_limit", flags.Lookup("default-pagination-limit"))
}

func main() {
	mainCommand.Execute()
}
