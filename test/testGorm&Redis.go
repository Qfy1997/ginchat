package main

import (
	// "gorm.io/driver/mysql"
	// "gorm.io/gorm"
	// "ginchat/models"
	"fmt"
	"github.com/go-redis/redis/v8"
	"context"
	"github.com/spf13/viper"
)

const (
	PublishKey = "websocket"
)

func main() {
	// db, err := gorm.Open(mysql.Open("root:Qfy123123@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"),&gorm.Config{})
	// if err!= nil {
	// 	panic("failed to connect database")
	// }
	//迁移 schema
	//db.AutoMigrate(&models.UserBasic{})
	// db.AutoMigrate(&models.Message{})
	//db.AutoMigrate(&models.GroupBasic{})
	//db.AutoMigrate(&models.Contact{})
	// db.AutoMigrate(&models.Community{})
	//Create
	// user := &models.UserBasic{}
	// user.Name = "zz"
	// db.Create(user)
	// contact :=&models.Contact{}
	// contact.OwnerId=8
	// contact.TargetId=5
	// contact.Type=1
	// db.Create(contact)
	// contact := &models.Community{}

	//Read
	// fmt.Println(db.First(user,1))  //根据整形主键查找
	
	// db.Model(user).Update("PassWord","1234")

	Red:= redis.NewClient(&redis.Options{
		Addr:viper.GetString("redis.addr"),
		Password:viper.GetString("redis.password"),
		DB:viper.GetInt("redis.DB"),
		PoolSize:viper.GetInt("redis.poolSize"),
		MinIdleConns:viper.GetInt("redis.minIdleConn"),
	})
	
	// ctx:=context.Context
	_,err := Red.Ping(context.TODO()).Result()
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("Redis init..")
	// for {
	// 	fmt.Println("监听中。。。")
	// 	sub:=Red.Subscribe(context.TODO(),PublishKey)
	// 	msg,err := sub.ReceiveMessage(context.TODO())
	// 	if err!=nil{
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println(msg)
	// }
	err = Red.Publish(context.TODO(),PublishKey,PublishKey).Err()
	if err!=nil{
		fmt.Println(err)
	}
	
	
}