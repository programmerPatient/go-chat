package roomLogic

import (
	"../../common/database/redis"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"strings"
)


var rooms = make(map[string]*Room)

var WsAcc = make(map[*websocket.Conn]string)

 var accWs = make(map[string]*websocket.Conn)

var userToRoom = make(map[string]string)

//处理客户端消息
func OnReceive(ws *websocket.Conn,newData *RoomData){
	switch newData.MessType {
	//心跳检测
	case 100:
		return
	//初始化房间信息
	case 102:
		initRoom(ws ,newData.Data)
	//聊天
	case 104:
		chat(newData.Data)
	default:
	}
}

/**
初始化房间信息
 */
func initRoom(ws *websocket.Conn,data map[string]string){
	redis := redis.Pool.Get()
	defer redis.Release()
	roomId := strings.Replace(data["roomId"], " ", "", -1)
	account := strings.Replace(data["account"], " ", "", -1)
	/**
	关闭老连接，异地登录
	 */
	oldWs := accWs[account]
	accWs[account] = ws
	delete(WsAcc,oldWs)
	WsAcc[ws] = account
	if oldWs != nil {
		oldWs.Close()
	}
	if rooms[roomId] == nil {
		rooms[roomId] = NewRoom()
	}
	if userToRoom[account] != "" && userToRoom[account] != roomId {
		rooms[userToRoom[account]].Pop(account)
		go offine(account,userToRoom[account])
	}
	userToRoom[account] = roomId
	rooms[roomId].Push(account)
	re := redis.GetHash("room:all",roomId)
	nickname := redis.GetHash("user:"+account,"name")
	if re == "" {
		return
	}
	wss := rooms[roomId].GetPeople()
	da := make(map[string]interface{})
	datas := make(map[string]interface{})
	newDatas,newWss := getWsAndNameByacc(wss)
	datas["roomUser"] = newDatas
	da = dataPack(400,"欢迎"+nickname+"进入房间",datas)
	sendMessageToClient(newWss,da)
}


/**
根据账号来获取连接和名称
*/
func getWsAndNameByacc(wss []string) (map[string]string,map[string]*websocket.Conn)  {
	redis := redis.Pool.Get()
	defer redis.Release()
	newDatas := make(map[string]string)
	newWss := make(map[string]*websocket.Conn)
	for _,t := range wss {
		name := redis.GetHash("user:"+t,"name")
		newDatas[t] = name
		newWss[t] = accWs[t]
	}
	return newDatas,newWss
}

/**
数据封装
*/
func dataPack(code int,message string,data interface{}) map[string]interface{} {
	da := make(map[string]interface{})
	da["code"] = code
	da["message"] = message
	da["data"] = data
	return da
}

/**
发送消息给客户端
*/
func sendMessageToClient(ws map[string]*websocket.Conn,data interface{}){
	go func() {
		datas,_ := json.Marshal(data)
		for _,w := range ws {
			if err := websocket.Message.Send(w,string(datas)); err != nil {
				fmt.Println(err)
				continue
			}
		}
	}()
}



/**
聊天
*/
func chat(data map[string]string){
	roomId := strings.Replace(data["roomId"], " ", "", -1)
	account := strings.Replace(data["account"]," ","",-1)
	toAccount := strings.Replace(data["toUid"]," ","",-1)
	message := data["msg"]
	redis := redis.Pool.Get()
	defer redis.Release()
	newWss := make(map[string]*websocket.Conn)
	if toAccount == "-1" {
		wss := rooms[roomId].GetPeople()
		_,newWss = getWsAndNameByacc(wss)
	}else{
		newWss[toAccount] = accWs[toAccount]
		newWss[account] = accWs[account]
	}
	da := make(map[string]string)
	da["from"] = redis.GetHash("user:"+account,"name")
	da["msg"] = message
	datas := dataPack(400,"",da)
	sendMessageToClient(newWss,datas)
}

/**
下线
*/
func offine(account string,roomId string){
	redis := redis.Pool.Get()
	defer redis.Release()
	nickname := redis.GetHash("user:"+account,"name")
	message := "用户"+nickname+"下线了！"
	wss := rooms[roomId].GetPeople()
	newDatas,newWss := getWsAndNameByacc(wss)
	datas := make(map[string]interface{})
	datas["roomUser"] = newDatas
	newData := dataPack(400,message,datas)
	sendMessageToClient(newWss,newData)
}



/**
清理信息
*/
func Clear(ws *websocket.Conn,account string) {
	if WsAcc[ws] != "" {
		roomId := userToRoom[account]
		room := rooms[roomId]
		room.Pop(account)
		offine(account,roomId)
	}
}

