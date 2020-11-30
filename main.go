package main

import (
	"./common/database/redis"
	"./config"
	"./marry"
	"./router"
	"./server"
	"fmt"
	"html/template"
	"time"
)

func FormatAsData(t time.Time) string{
	year,month,day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d",year,month,day)
}

func main() {
	redis.Pool = redis.New()
	go func() {
		ws := server.NewWs()
		ws.Start()
	}()
	engine := marry.New()
	//静态访问路径
	engine.Static("/assets","./static")
	engine.SetFuncMap(template.FuncMap{
		"FormatAsData":FormatAsData,
	})
	engine.LoadHTMLGlob("./templates/*")
	router.Run(engine)
	addr := config.ListenIp+":"+config.Port
	err := engine.Run(addr)
	if err !=  nil {
		fmt.Println(err)
	}
}
