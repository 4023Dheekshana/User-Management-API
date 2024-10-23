package main

import (
	"log"
	"userapi/database"
	"userapi/router"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	//gin.SetMode(gin.ReleaseMode) optional not to get warning lines
	route := gin.Default()
	database.ConnectDataBase()
	router.RegisterRoutes(route)
	err := route.Run(":8080")
	if err != nil {
		log.Fatalf("Error runing")
		panic(err)
	}

}
