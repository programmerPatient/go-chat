package middleware

import (
	"../common/database/redis"
	"../common/lib"
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
		if cookie == nil {
			errs += "请登录"
			http.Redirect(c.W,c.R,"/login?error="+errs,http.StatusFound)
		}
		token := cookie.Value
		_,b := lib.JwtDecode(token)
		if !b {
			c.DelCookie("token")
			errs += "令牌过期或无效，请重新登陆！"
			http.Redirect(c.W,c.R,"/login?error="+errs,http.StatusFound)
		}

	}
}

