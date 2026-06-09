package main

import (
	"sespima_api/config"
	"sespima_api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	config.ConnectDB()



	r := gin.Default()

	routes.Register(r, config.DB)

	r.Run(":8000")
}