package controller

import (
	"../common/database/redis"
	"../marry"
	"fmt"
	"github.com/chilts/sid"
	"net/http"
)
func RoomIndex(c *marry.Context) {
	redis := redis.Pool.Get()
	defer redis.Release()
	roomId := c.Query("room_id")
	roomName := c.Query("room_name")
	cookie := c.GetCookie("token")
	account := redis.GetHash("token",cookie.Value)
	val := make(map[interface{}]interface{})
	val["roomId"] = roomId
	val["roomName"] = roomName
	val["account"] = account
	val["accountName"] = redis.GetHash("user:"+account,"name")
	c.HTML(http.StatusOK,"room.html",val)
}

func RoomAdd(c *marry.Context){
	roomName := c.PostForm("name")
	roomId := sid.Id()
	redis := redis.Pool.Get()
	defer redis.Release()
	res := redis.SetSAdd("room:name",roomName)
	if !res {
		fmt.Print("房间名已经存在")
	}
	redis.SetHash("room:all",roomId,roomName)
	var data map[string]string
	data = make(map[string]string)
	data["room_id"] = roomId
	data["room_name"] = roomName
	c.JSON(http.StatusOK,data)
}

