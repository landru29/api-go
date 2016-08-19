package main

import (
    "fmt"

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
        fmt.Println(viper.GetString("api_host"))
    },
}

func init() {
    flags := mainCommand.Flags()
    flags.String("api-host", "your-api-host", "API host")
    viper.BindPFlag("api_host", flags.Lookup("api-host"))
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
