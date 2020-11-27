package controller

import (
	"../common/database/redis"
	"../common/lib"
	"../marry"
	"net/http"
)

func LoginIndex(c *marry.Context) {
	cookie := c.GetCookie("token")
	if cookie != nil {
		http.Redirect(c.W,c.R,"/user/index",http.StatusFound)
	}
	c.HTML(http.StatusOK,"login.html",nil)
}

/**
登录验证
*/
func LoginCheck(c *marry.Context) {
	redis := redis.Pool.Get()
	defer redis.Release()
	account := c.PostForm("account")
	password := c.PostForm("password")
	var errs string
	res := redis.GetHash("user:"+account,"password")
	if res == "" {
		errs += "账号不存在"
		http.Redirect(c.W,c.R,"/login?error="+errs,http.StatusFound)
		return
	}
	if lib.MD5(password) != res {
		errs += "密码错误"
		http.Redirect(c.W,c.R,"/login?error="+errs,http.StatusFound)
		return
	}else{
		jwt := lib.JwtNew(account)
		res,err := jwt.JwtEncode()
		if !err {
			errs += "内部服务错误请重试"
			http.Redirect(c.W,c.R,"/login?error="+errs,http.StatusFound)
			return
		}
		redis.SetHash("token",res,account)
		c.SetCookie("token",res)
	}
	http.Redirect(c.W,c.R,"/user/index",http.StatusFound)
}
