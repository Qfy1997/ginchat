package router

import (
	"github.com/gin-gonic/gin"
	"ginchat/service"
	"ginchat/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
) 

func Router() *gin.Engine {
	r := gin.Default()
	//swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any",ginSwagger.WrapHandler(swaggerfiles.Handler))

	//静态资源
	r.Static("/asset","asset/")
	r.LoadHTMLGlob("views/**/**")
	//首页
	r.GET("/",service.GetIndex)
	r.GET("/index",service.GetIndex)
	r.GET("/toRegister",service.ToRegister)
	r.GET("/toChat",service.ToChat)
	r.GET("/chat",service.Chat)
	r.GET("/tocreateCommunity",service.TocreateCommunity)
	
	
	//加载表情包
	// r.POST("/asset/plugins/doutu/mkgif/info.json",service.GETmkgirjson)
	// r.POST("/asset/plugins/doutu/emoj/info.json",service.GETemojijson)

	r.POST("/searchFriends",service.SearchFriends)
	//用户模块
	r.POST("/user/getUserList",service.GetUserList)
	r.POST("/user/createUser",service.CreateUser)
	r.POST("/user/deleteUser",service.DeleteUser)
	r.POST("/user/updateUser",service.UpdateUser)
	r.POST("/user/FindUserByNameAndPwd",service.FindUserByNameAndPwd)

	//发送消息
	r.GET("/user/sendMsg",service.SendMsg)
	r.GET("/user/sendUserMsg",service.SendUserMsg)
	//上传文件
	r.POST("/attach/upload",service.Upload)
	//添加好友
	r.POST("/contact/addfriend",service.AddFriend)
	//创建群
	r.POST("/contact/createcommunity",service.CreateCommunity)
	//加载群
	r.POST("contact/loadcommunity",service.Loadcommunity)
	return r
}