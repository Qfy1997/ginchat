package service

import (
	// "net/http"
	"github.com/gin-gonic/gin"
	"ginchat/utils"
	"time"
	"fmt"
	"strings"
	"math/rand"
	"os"
	"io"
)

func Upload(c *gin.Context) {
	 w := c.Writer
	req := c.Request
	srcFile,head,err := req.FormFile("file")
	if err!=nil{
		utils.RespFail(w, err.Error())
	}
	suffix := ".png"
	ofilName := head.Filename
	tem := strings.Split(ofilName,".")
	if len(tem) >1 {
		suffix = "." + tem[len(tem)-1]
	}
	fileName := fmt.Sprintf("%d%04d%s",time.Now().Unix(),rand.Int31(),suffix)
	dstFile,err:=os.Create("./asset/upload/"+fileName)
	if err!= nil {
		utils.RespFail(w,err.Error())
	}
	_,err = io.Copy(dstFile,srcFile)
	if err != nil {
		utils.RespFail(w,err.Error())
	}
	url :="./asset/upload/"+fileName
	utils.RespOK(w,url,"发送图片成功")
}