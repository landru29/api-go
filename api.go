package main

import (
	"fmt"

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
		router := routes.DefineRoutes()
		router.Run(":" + viper.GetString("api_port"))
	},
}

func init() {
	flags := mainCommand.Flags()
	flags.String("api-host", "your-api-host", "API host")
	flags.String("api-port", "3000", "API port")
	viper.BindPFlag("api_host", flags.Lookup("api-host"))
	viper.BindPFlag("api_port", flags.Lookup("api-port"))
}

func main() {
	mainCommand.Execute()

}

/*
   mongo.AutoConnect()

   r := gin.Default()

   r.GET("/ping", func(c *gin.Context) {
       c.JSON(200, gin.H{
           "message": "pong",
       })
   })
   r.Run()*/
