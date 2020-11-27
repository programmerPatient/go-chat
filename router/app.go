package router

import (
	"../controller"
	"../marry"
	"../middleware"
)

/**
路由控制，
 */
func Run(engine *marry.Engine){
	/**
	不分组路由
	 */
	engine.GET("/login",controller.LoginIndex)
	engine.POST("/login",controller.LoginCheck)
	engine.GET("/register",controller.RegisterIndex)
	engine.POST("/register",controller.RegisterCheck)

	//分组路由 user代表路有前缀
	g := engine.RouterGroup.Group("user")
	//为分组路由添加中间件
	g.Use(middleware. LoginAuth())
	{
		g.GET("/room",controller.RoomIndex)
		g.GET("/index",controller.IndexIndex)
		g.POST("/room/add",controller.RoomAdd)
	}
}
