package service

import (
	"github.com/gin-gonic/gin"
	"text/template"
	"ginchat/models"
	"strconv"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} Helloworld
// @Router /index [get]
func GetIndex(c *gin.Context){
	ind,err:=template.ParseFiles("index.html","views/chat/head.html")
	if err!=nil{
		panic(err)
	}
	ind.Execute(c.Writer,"index")
	// c.JSON(200,gin.H{
	// 	"message":"welcom!!",
	// })
}
// TocreateCommunity
// @Tags 创建群
// @Success 200 {string} Helloworld
// @Router /tocreateCommunity [get]
func TocreateCommunity(c *gin.Context){
	ind,err:=template.ParseFiles("views/chat/createcom.html","views/chat/head.html")
	if err!=nil{
		panic(err)
	}
	userId,_:=strconv.Atoi(c.Query("userId"))
	token:=c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	ind.Execute(c.Writer,"toCreatecommunity")
}

// ToRegister
// @Tags 注册
// @Success 200 {string} Helloworld
// @Router /toRegister [get]
func ToRegister(c *gin.Context) {
	ind,err:=template.ParseFiles("views/user/register.html")
	if err!=nil{
		panic(err)
	}
	ind.Execute(c.Writer,"register")
}

// ToChat
// @Tags 进入聊天室
// @Success 200 {string} Helloworld
// @Router /toChat [get]
func ToChat(c *gin.Context) {
	ind,err:=template.ParseFiles("views/chat/index.html",
	"views/chat/head.html",
	"views/chat/foot.html",
	"views/chat/tabmenu.html",
	"views/chat/concat.html",
	"views/chat/group.html",
	"views/chat/profile.html",
	"views/chat/main.html",)
	if err!=nil{
		panic(err)
	}
	userId,_:=strconv.Atoi(c.Query("userId"))
	token:=c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	ind.Execute(c.Writer,"register")
}

func Chat(c *gin.Context) {
	models.Chat(c.Writer,c.Request)
}