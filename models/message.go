package models

import (
	"gorm.io/gorm"
	"github.com/gorilla/websocket"
	"sync"
	"net/http"
	"gopkg.in/fatih/set.v0"
	"fmt"
	"strconv"
	// "net"
	"encoding/json"
	"context"
	"github.com/go-redis/redis/v8"
	// "context"
	"github.com/spf13/viper"
	// "ginchat/utils"
)
type Message struct {
	gorm.Model
	FormId 		int64   `json:"userid"`//发送着
	TargetId 	int64 	`json:"dstid"`//接受者 
	Type 		int		`json:"cmd"`//消息类型 10.私聊 11.群聊 0.广播
	Media 		int		`json:"media"`//消息类型  1.文字 2.表情包 3.图片 4.音频
	Content 	string 	`json:"content"`//消息内容
	Pic 		string
	Url 		string
	Desc 		string
	Amount 		int		//其他数字统计
}



func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}
//映射关系
var clientMap map[int64] *Node = make(map[int64] *Node,0)
//读写锁
var rwLocker sync.RWMutex


//需要：发送者ID，接受者ID，消息类型1，发送的内容context，发送类型1， token校验
func Chat(writer http.ResponseWriter, request *http.Request){
	ctx := request.Context()
	//1.检验token 等合法性
	//token := query.Get("token")
	query := request.URL.Query()
	//获取对应的websocket连接号，很重要！！！！websocket init
	Id := query.Get("userId")
	fmt.Println("ID:")
	fmt.Println(Id)
	userId,_ := strconv.ParseInt(Id,10,64)
	ctx = context.WithValue(ctx,"userId",userId)
	// msgType :=query.Get("cmd")
	// fmt.Println("msgType:")
	// fmt.Println(msgType)
	// targetId := query.Get("targetId")
	// context := query.Get("context")
	subcli:= redis.NewClient(&redis.Options{
		Addr:viper.GetString("redis.addr"),
		Password:viper.GetString("redis.password"),
		DB:viper.GetInt("redis.DB"),
		PoolSize:viper.GetInt("redis.poolSize"),
		MinIdleConns:viper.GetInt("redis.minIdleConn"),
	})
	topic:=fmt.Sprintf("%v", userId)
	isvalida :=true //checkToken()  待。。。
	conn,err := (&websocket.Upgrader{
		//token 校验
		CheckOrigin :func (r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer,request,nil)
	if err!=nil{
		fmt.Println(err)
	}
	//2.获取链接conn
	node := &Node{
		Conn : conn,
		DataQueue: make(chan []byte,50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//3.用户关系
	//4.userid跟node绑定并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//完成接收逻辑，接收websocket连接(前端)发来的数据，并对此消息进行后端逻辑处理
	go recvProc(ctx,node)
	//建立redis客户端并订阅响应频道，接收消息报文
	sub:=subcli.Subscribe(context.TODO(),topic)
	// msg,err := sub.ReceiveMessage(ctx)
	rdmsg,zz := sub.ReceiveMessage(context.TODO())
	if zz!=nil{
		fmt.Println("redis sub err:",zz)
		return 
	}
	fmt.Println("subcli message:",rdmsg.Payload)
	// dispatch(ctx,[]byte(rdmsg.Payload))
	//完成发送逻辑，将redis中订阅的消息写入websocket，即：前端显示
	go sendProc(ctx,node,[]byte(rdmsg.Payload))
	// sendMsg(ctx,userId,[]byte("欢迎进入聊天室"))
}

func sendProc(ctx context.Context,node *Node,rdmsg []byte){
	err := node.Conn.WriteMessage(websocket.TextMessage,rdmsg)
	if err != nil {
		fmt.Println(err)
		return 
	}
	// userId:=ctx.Value("userId")
	// for {
	// 	select {
	// 	case data:= <- node.DataQueue:
	// 		fmt.Println(userId,"'s [ws] sendMsg >>>>>>> msg",string(data))
	// 		err := node.Conn.WriteMessage(websocket.TextMessage,rdmsg)

	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return 
	// 		}
	// 	}
	// }
	// fmt.Println(userId,"'s [ws] sendMsg >>>>>>> msg",string(rdmsg))
}

func recvProc(ctx context.Context,node *Node){
	for {
		userId:=ctx.Value("userId")
		// topic:=fmt.Sprintf("%v", userId)
		_,data,err := node.Conn.ReadMessage()
		// sub:=subcli.Subscribe(context.TODO(),topic)
		// msg,err := sub.ReceiveMessage(ctx)
		// rdmsg,zz := sub.ReceiveMessage(context.TODO())
		if err != nil {
			fmt.Println(err)
			return 
		}
		// if zz!=nil{
		// 	fmt.Println("redis sub err:",zz)
		// 	return 
		// }
		dispatch(ctx,data)
		// broadMsg(data)
		fmt.Println(userId,"'s [ws] recvProc <<<<< ",string(data))
	}
}

//后端调度逻辑,将消息发布到对应频道
func dispatch(ctx context.Context,data []byte){
	// fmt.Println(data)
	userId:=ctx.Value("userId")
	msg := Message{}
	// var p map[string]interface{}
	err := json.Unmarshal(data,&msg)
	// b,_ := json.Marshal(p)
	fmt.Println(msg)
	if err != nil {
		fmt.Println(err)
		return 
	}
	// fmt.Println("dipatch msg.Type:")
	// fmt.Println(msg.Type)
	switch msg.Type {
	case 10: //私信
		fmt.Println(userId,"'s [ws] dispatch data:",string(data))
		// sendMsg(ctx,msg.TargetId, data)
		topic:=strconv.FormatInt(msg.TargetId,10)
		pubcli:= redis.NewClient(&redis.Options{
			Addr:viper.GetString("redis.addr"),
			Password:viper.GetString("redis.password"),
			DB:viper.GetInt("redis.DB"),
			PoolSize:viper.GetInt("redis.poolSize"),
			MinIdleConns:viper.GetInt("redis.minIdleConn"),
		})
		err:=pubcli.Publish(context.TODO(),topic,string(data))
		if err!=nil{
			fmt.Println(err)
		}
		// err:=utils.Red.Publish(ctx,topic,string(data)).Err()
		// if err!=nil{
		// 	fmt.Println(err)
		// }
		// case 2:  //群发
		// sendGroupMsg()
		// case 3:  //广播
		// sendAllMsg()
		// case 4:
	}
}

// var udpsendChan chan []byte = make(chan []byte, 1024)
// func broadMsg(data []byte){
// 	udpsendChan <- data 
// }

// func init(){
// 	// go udpSendProc()
// 	// go udpRecvProc()
// 	fmt.Println("init goroutine...")
// }

// //完成udp数据发送协程
// func udpSendProc() {
// 	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
// 		IP:net.IPv4(192,168,18,1),
// 		Port:3000,
// 	})
// 	defer con.Close()
// 	if err!=nil{
// 		fmt.Println(err)
// 	}
// 	for {
// 		select {
// 		case data := <- udpsendChan:
// 			fmt.Println("udpSendProc data:",string(data))
// 			_, err := con.Write(data)
// 			if err != nil{
// 				fmt.Println(err)
// 				return 
// 			}
// 		}
// 	}
// }

// //完成udp数据接收协程
// func udpRecvProc() {
// 	con,err :=net.ListenUDP("udp",&net.UDPAddr{
// 		IP:net.IPv4zero,
// 		Port:3000,
// 	})
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer con.Close()
// 	for{
// 		var buf [512] byte
// 		n,err:=con.Read(buf[0:])
// 		if err != nil {
// 			fmt.Println(err)
// 			return 
// 		}
// 		fmt.Println("udpRecvProc  data:",string(buf[0:n]))
// 		// dispatch(buf[0:n])
// 	}
// }



// func sendMsg(ctx context.Context,targetId int64,msg []byte){
// 	userId:=ctx.Value("userId")
// 	fmt.Println(userId,"'s [ws] sendMsg >>> targetID:",targetId," msg:",string(msg))
// 	rwLocker.RLock()
// 	node,ok := clientMap[targetId]
// 	rwLocker.RUnlock()
// 	if ok{
// 		node.DataQueue <- msg
// 	}
// }
