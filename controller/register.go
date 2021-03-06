package controller

import (
	"../common/database/redis"
	"../common/lib"
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
	account := c.PostForm("account")
	name := c.PostForm("name")
	password := c.PostForm("password")
	confimPassword := c.PostForm("confim_password")
	var errs string
	if password != confimPassword {
		errs += "密码不一致"
		http.Redirect(c.W,c.R,"/register?error="+errs,http.StatusFound)//重定向
		return
	}
	res := redis.HGetAll("user:"+account)
	if  len(res) > 0  {
		errs += "账号已经存在"
		http.Redirect(c.W,c.R,"/register?error="+errs,http.StatusFound)//重定向
		return
	}
	ress := redis.GetHash("user:name",name)
	if ress != "" {
		errs += "用户名已经存在"
		http.Redirect(c.W,c.R,"/register?error="+errs,http.StatusFound)//重定向
		return
	}
	redis.SetHash("user:name",name,account)
	redis.SetHash("user:"+account,"name",name)
	redis.SetHash("user:"+account,"account",account)
	has := lib.MD5(password)
	redis.SetHash("user:"+account,"password",has)
	resss := lib.JwtEncode(account)
	c.SetCookie("token",resss)
	http.Redirect(c.W,c.R,"/user/index",http.StatusFound)
}

