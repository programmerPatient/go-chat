package controller

import (
	"../common/database/redis"
	"../common/lib"
	"../marry"
	"fmt"
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
		http.RedirectHandler("/register?error="+errs,http.StatusFound)//重定向
	}
	fmt.Printf(account,name,password,confimPassword)
	res := redis.HGetAll("user:"+account)
	fmt.Printf("%v",res)
	if res != nil {
		errs += "账号已经存在"
		http.RedirectHandler("/register?error="+errs,http.StatusFound)//重定向
	}
	res = redis.GetHash("user:name",name)
	fmt.Printf("%v",res)
	if res != nil {
		errs += "用户名已经存在"
		http.RedirectHandler("/register?error="+errs,http.StatusFound)//重定向
	}
	redis.SetHash("user:name",name,account)
	redis.SetHash("user:"+account,"name",name)
	redis.SetHash("user:"+account,"account",account)
	has := lib.MD5(password)
	redis.SetHash("user:"+account,"password",has)

	http.Redirect(c.W,c.R,"/user/index",http.StatusFound)
}

