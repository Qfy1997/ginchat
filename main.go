package main

import (
	// "github.com/gin-gonic/gin"
	 "ginchat/router"
	 "ginchat/utils"
)
func main() {
	utils.InitConfig()
	utils.InitMySQL()
	utils.InitRedis()
	r := router.Router()
	r.Run(":8080")
}