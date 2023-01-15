package service

import (
	"github.com/gin-gonic/gin"
	"ginchat/models"
	"strconv"
	"fmt"
	"github.com/asaskevich/govalidator"
	"ginchat/utils"
	"math/rand"
	"net/http"
	"github.com/gorilla/websocket"
	"time"
)

// GetUserList
// @Summary 用户列表
// @Tags 首页
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context){
	data := make([] *models.UserBasic,10)
	data = models.GetUserList()
	c.JSON(200,gin.H{
		"message": data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [post]
func CreateUser(c *gin.Context){
	user := models.UserBasic{}
	// user.Name = c.Query("name")
	// password := c.Query("password")
	// repassword := c.Query("repassword")
	user.Name = c.Request.FormValue("name")
	password:=c.Request.FormValue("password")
	repassword := c.Request.FormValue("repassword")	
	// fmt.Println("zz?")
	fmt.Println(user.Name,password,repassword)
	salt := fmt.Sprintf("%06d",rand.Int31())
	if user.Name == "" || password == "" || repassword == "" {
		// fmt.Println("11")
		// fmt.Println("data.Name:")
		// fmt.Println(user.Name)
		// fmt.Println("password:")
		// fmt.Println(password)
		c.JSON(200,gin.H{
			"code":-1,
			"message":"用户名或密码不能为空！",
			"data":user,
		})
		return
	}
	data := models.FindUserByName(user.Name)
	if data.Name!= "" {
		c.JSON(200,gin.H{
			"code":-1,
			"message":"用户名已注册！",
			"data":user,
		})
		return
	}
	if password != repassword {
		c.JSON(200, gin.H{
			"code":-1,
			"message":"两次密码不一致!",
			"data":user,
		})
		return 
	}
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt
	// user.PassWord = password

	models.CreateUser(user)

	c.JSON(200,gin.H{
		"code":0,//0正确，-1失败
		"message": "新增用户成功",
		"data":data,
	})
}

// FindUserByNameAndPwd
// @Summary 登录
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/FindUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context){
	data := models.UserBasic{}
	
	// name := c.Query("name")
	// password := c.Query("password")
	name:=c.Request.FormValue("name")
	password:=c.Request.FormValue("password")
	fmt.Println(name,password)
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200,gin.H{
			"code":-1,//0正确，-1失败
			"message": "该用户不存在",
			"data":data,
		})
	}
	// fmt.Println(user)
	// fmt.Println("user.Salt:",user.Salt)
	// fmt.Println("password:",password)
	// fmt.Println("user.Password:",user.PassWord)
	flag:=utils.ValidPassword(password,user.Salt,user.PassWord)
	if !flag {
		c.JSON(200,gin.H{
			"code":-1,//0正确，-1失败
			"message": "密码不正确",
			"data":data,
		})
	}
	pwd:=utils.MakePassword(password,user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(200,gin.H{
		"code":0,//0正确，-1失败
		"message": "登录成功",
		"data":data,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [post]
func DeleteUser(c *gin.Context){
	user := models.UserBasic{}
	id,_ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200,gin.H{
		"code":0,
		"message": "删除用户成功",
		"data":user,
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "emain"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context){
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Email = c.PostForm("email")
	user.Phone = c.PostForm("phone")
	fmt.Println("update:",user)

	_,err:=govalidator.ValidateStruct(user)
	if err!=nil{
		fmt.Println("err:",err)
		c.JSON(200,gin.H{
			"code":-1,//0正确，-1失败
			"message": "修改参数不匹配",
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200,gin.H{
		"code":0,//0正确，-1失败
		"message": "修改用户成功",
		})
	}
}

//防止跨域站点的伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func SendMsg(c *gin.Context){
	ws,err := upGrade.Upgrade(c.Writer,c.Request,nil)
	if err!=nil{
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn){
		err = ws.Close()
		if err!=nil{
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws,c)
}

func MsgHandler(ws *websocket.Conn,c *gin.Context) {
	for{
		msg,err:=utils.Subscribe(c,utils.PublishKey)
		if err!=nil{
			fmt.Println(err)
		}
		fmt.Println("发送消息:")
		tm:=time.Now().Format("2006-01-02 15:04:05")
		m:=fmt.Sprintf("[ws][%s]:%s",tm,msg)

		err = ws.WriteMessage(1,[]byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}	
}


func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer,c.Request)
}

// SearchFriends
// @Summary 查询好友列表
// @Tags 用户模块
// @param userId formData string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /searchFriends [post]
func SearchFriends(c *gin.Context) {
	id,_:=strconv.Atoi(c.Request.FormValue("userId"))
	users:=models.SearchFriend(uint(id))
	// c.JSON(200,gin.H{
	// 	"code":0,
	// 	"message":"查询好友列表成功",
	// 	"data":users,
	// })
	utils.RespOKList(c.Writer,users,len(users))
}

// func GETmkgirjson(c *gin.Context) {
// 	c.JSON(200,gin.H{
// 		"id":"mkgif",
// 		"size":"middle",
// 		"icon":"icon.gif",
// 		"assets":"1.gif",
// 	})
// }

// func GETemojijson(c *gin.Context) {
// 	c.JSON(200,gin.H{
// 		"id":"emoj",
// 		"size":"small",
// 		"icon":"0.gif",
// 		"assets":"1.gif",
// 	})
// }

// AddFriends
// @Summary 添加好友
// @Tags 用户模块
// @param userId formData string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /searchFriends [post]
func AddFriend(c *gin.Context) {
	userId,_:=strconv.Atoi(c.Request.FormValue("userId"))
	targetId,_:=strconv.Atoi(c.Request.FormValue("targetId"))
	code,msg:=models.AddFriend(uint(userId),uint(targetId))
	// c.JSON(200,gin.H{
	// 	"code":0,
	// 	"message":"查询好友列表成功",
	// 	"data":users,
	// })
	if code ==0 {
		utils.RespOK(c.Writer,code,msg)
	} else {
		utils.RespFail(c.Writer,msg)
	}
}

// CreateCommunity
// @Summary 创建群
// @Tags 用户模块
// @param ownerId formData string false "id"
// @param Name formData string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /searchFriends [post]
func CreateCommunity(c *gin.Context) {
	ownerId,_:=strconv.Atoi(c.Request.FormValue("owner_id"))
	name:=c.Request.FormValue("name")
	community := models.Community{}
	community.OwnerId = uint(ownerId)
	community.Name =name

	code,msg:=models.CreateCommunity(community)
	// c.JSON(200,gin.H{
	// 	"code":0,
	// 	"message":"查询好友列表成功",
	// 	"data":users,
	// })
	if code ==0 {
		utils.RespOK(c.Writer,code,msg)
	} else {
		utils.RespFail(c.Writer,msg)
	}
}

func Loadcommunity(c *gin.Context) {
	owner_id,_ :=strconv.Atoi(c.Request.FormValue("userId"))
	data,msg:=models.Loadcommunity(uint(owner_id))
	if len(data) != 0 {
		utils.RespList(c.Writer,0,data,msg)
	}else {
		utils.RespFail(c.Writer,msg)
	}
}