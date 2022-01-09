package main

import (

	// "gin-mongo-api/configs"

	"gin-mongo-api/configs"
	"gin-mongo-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	// initialize a new Gin router. Use the Default router because Gin provides some useful middlewares we can use to debug our server.
	router := gin.Default()
	configs.ConnectDB()
	routes.UserRoute(router)
	routes.ProductRoute(router)

	// run our server
	router.Run()
}
