package models

import (
	"gorm.io/gorm"
	"ginchat/utils"
	"fmt"
	"time"
)
type UserBasic struct {
	gorm.Model
	Name	 		string
	PassWord 		string
	Phone 	 		string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email 	 		string `valid:"email"`
	Identity 		string
	ClientIp 		string
	ClientPort  	string
	Salt			string
	LoginTime 		uint64
	HeartbeatTime 	uint64
	LoginOutTime 	uint64
	IsLogout 		bool
	DeviceInfo 		string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([] *UserBasic,10)
	utils.DB.Find(&data)
	for _,v := range data {
		fmt.Println(v)
	}
	return data
}

func FindUserByNameAndPwd(name string,password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and pass_word=?",name,password).First(&user)
	//token加密
	str:=fmt.Sprintf("%d",time.Now().Unix())
	temp:=utils.MD5Encode(str)
	utils.DB.Model(&user).Where("id = ?",user.ID).Update("identity",temp)
	return user
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?",name).First(&user)
	return user
}

func FindUserByPhone(phone string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("Phone = ?",phone).First(&user)
	return user
}

func FindUserByEmail(email string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("email = ?", email).First(&user)
	return user
}

func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}
func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{Name:user.Name,PassWord:user.PassWord,Email:user.Email,Phone:user.Phone})
}

func FindById(id uint) UserBasic {
	user:=UserBasic{}
	utils.DB.Where("id = ?", id).First(&user)
	return user
}