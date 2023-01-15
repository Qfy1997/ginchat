package models

import (
	"gorm.io/gorm"
	"fmt"
	"ginchat/utils"
)
type Community struct {
	gorm.Model
	Name 	string
	OwnerId uint
	Img 	string
	Desc 	string
}

func CreateCommunity(community Community)(int, string){
	if len(community.Name) == 0 {
		return -1,"群名称不能为空"
	}
	if community.OwnerId == 0 {
		return -1,"请先登录"
	}
	if err:= utils.DB.Create(&community).Error; err!=nil {
		fmt.Println(err)
		return -1,"建群失败"
	}
	return 0,"建群成功"
}


func Loadcommunity(ownerId uint) ([]*Community,string) {
	fmt.Println(ownerId)
	data := make([]*Community,10)
	utils.DB.Where("owner_id = ?",ownerId).Find(&data)
	for _,v:=range data {
		fmt.Println(v)
	}
	
	// utils.DB.Where()
	return data,"查询成功"
}