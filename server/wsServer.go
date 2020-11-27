package server

import (
	"./roomLogic"
	"./serverConfig"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)



type WsServer struct{
	*BaseServer
}

func NewWs() *WsServer{

	w := &WsServer{
	}
	w.BaseServer = &BaseServer{
		host: serverConfig.WsHost,
		port: serverConfig.WsPort,
		addr: serverConfig.WsHost + ":" + serverConfig.WsPort,
	}
	return w
}


func wsHandler(ws *websocket.Conn) {
	var err error
	for {
		var reply string
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			roomLogic.Clear(ws,roomLogic.WsAcc[ws])
			break
		}else{
			//消息转发
			nowData := &roomLogic.RoomData{
				MessType: 0,
				Data:     nil,
			}
			_ = json.Unmarshal([]byte(reply),nowData)
			roomLogic.OnReceive(ws,nowData)
		}
	}

}


func (w *WsServer) Start()  {
	http.Handle("/",websocket.Handler(wsHandler))
	err := http.ListenAndServe(":"+w.port, nil)
	if err != nil {
		fmt.Println(err)
	}

}
