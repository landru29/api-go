package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/landru29/api-go/helpers/config"
	"github.com/landru29/api-go/helpers/mongo"
)


func main() {
	fmt.Printf("Launching API\n")

	config.Load()
	mongo.AutoConnect()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()

}
