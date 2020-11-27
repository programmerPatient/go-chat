package middleware

import (
	"../common/database/redis"
	"../marry"
	"net/http"
)
/**
登录授权
*/
func LoginAuth () marry.HandlerFunc {
	return func(c *marry.Context) {
		redis := redis.Pool.Get()
		defer redis.Release()
		var errs = ""
		cookie := c.GetCookie("token")
		if cookie==nil {
			errs += "请登录"
			http.Redirect(c.W,c.R,"/login?error="+errs,http.StatusFound)
		}
		token := cookie.Value
		account := redis.GetHash("token",token)
		if account == "" {
			errs += "请重新登录"
			http.Redirect(c.W,c.R,"/login?error="+errs,http.StatusFound)
		}
	}
}

