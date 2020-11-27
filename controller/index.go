package controller

import (
	"../common/database/redis"
	"../marry"
	"fmt"
	"net/http"
)

func IndexIndex(c *marry.Context){
	redis := redis.Pool.Get()
	defer redis.Release()
	value := redis.HGetAll("room:all")
	fmt.Printf("%v",c.TempParam)
	c.HTML(http.StatusOK,"index.html",value)
}





