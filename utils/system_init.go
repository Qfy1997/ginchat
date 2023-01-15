package utils

import (
	"github.com/spf13/viper"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"gorm.io/gorm/logger"
	"time"
	"context"
	// "ginchat/models"
	"github.com/go-redis/redis/v8"
)

var (
	DB *gorm.DB
	Red *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println("config app init")
	// fmt.Println("config app:",viper.Get("app"))
	// fmt.Println("config mysql:",viper.Get("mysql"))
	// fmt.Println("zz:",viper.Get("mysql.dns"))
}


func InitMySQL()  {
	//自定义日志模版，打印SQL语句
	newLogger := logger.New(
		log.New(os.Stdout,"\r\n",log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢SQL阀值
			LogLevel:logger.Info, //级别
			Colorful:true,//彩色
		},
	)
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")),&gorm.Config{Logger: newLogger})	
	fmt.Println("Mysql inited...")
	// user := models.UserBasic{}
	// db.Find(&user)
	// fmt.Println(user)	
}


func InitRedis()  {
	Red= redis.NewClient(&redis.Options{
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
}

const (
	PublishKey = "websocket"
)

//Publish 发布消息到redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("Publish...",msg)
	err = Red.Publish(ctx,channel,msg).Err()
	if err!=nil{
		fmt.Println(err)
	}
	return err
}

//Subscribe 订阅redis消息
func Subscribe(ctx context.Context, channel string) (string,error) {
	sub:=Red.Subscribe(ctx,channel)
	fmt.Println("Subscribe1...",ctx)
	fmt.Println("zz")
	msg,err := sub.ReceiveMessage(ctx)
	if err!=nil{
		fmt.Println("报错了。。")
		fmt.Println(err)
		return "",err
	}
	fmt.Println("Subscribe...",msg.Payload)
	return msg.Payload,err
}

