package controller

import (
	"../common/database/redis"
	"../marry"
	"net/http"
)

/**
注册
*/
func RegisterIndex(c *marry.Context) {
	cookie := c.GetCookie("token")
	if cookie!=nil {
		http.Redirect(c.W,c.R,"/user/index",http.StatusFound)
	}
	c.HTML(http.StatusOK,"register.html",nil)
}

func RegisterCheck(c *marry.Context) {
	redis := redis.Pool.Get()
	defer redis.Release()
}

