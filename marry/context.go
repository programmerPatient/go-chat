/**
 *请求的数据返回封装
 */
package marry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Context struct{
	W http.ResponseWriter//http请求的返回
	R *http.Request//http请求的
	Method string//请求的方式
	Path string //请求的路由
	Status int //状态
	Params map[string]string//参数
	handlers []HandlerFunc//中間件
	index int
	engine *Engine
	TempParam map[string]interface{}//临时参数

}

func (c *Context) Param(key string) string {
	value , _ := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.R.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.R.URL.Query().Get(key)
}

func nowContext(w http.ResponseWriter,r *http.Request) *Context{
	return &Context{
		W:w,
		R:r,
		Method: r.Method,
		Path: r.URL.Path,
		index: -1,
	}
}

/**
设置返回的数据状态
 */
func (c *Context) SetStatus(code int) {
	c.Status = code
	c.W.WriteHeader(code)
}

/*
设置头部信息
 */
func (c *Context) SetHeader(key string,value string){
	c.W.Header().Set(key,value)
}

/**
获取头部信息
 */
func (c *Context) GetHeader(key string) string {
	return c.R.Header.Get(key)
}

func (c *Context) DelCookie(key string) {
	http.SetCookie(c.W,&http.Cookie{
		Name:       key,
		Value:      "",
		Path:       "",
		Domain:     "",
		Expires:    time.Time{},
		RawExpires: "",
		MaxAge:     -1,
		Secure:     false,
		HttpOnly:   false,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	})
}

/**
设置cookie
 */
func (c *Context) SetCookie(key string,value string){
	http.SetCookie(c.W,&http.Cookie{
		Name:       key,
		Value:      value,
		Path:       "",
		Domain:     "",
		Expires:    time.Time{},
		RawExpires: "",
		MaxAge:     0,
		Secure:     false,
		HttpOnly:   false,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	})
}

/**
获取cookie
 */
func (c *Context) GetCookie(key string) *http.Cookie {
	cookie,err := c.R.Cookie(key)
	if err !=nil {
		return nil
	}
	return cookie
}

/**
code 返回的状态
format 数据格式
 */
func (c *Context) String(code int,format string,values ...interface{}) {
	c.SetStatus(code)
	c.W.Write([]byte(fmt.Sprintf(format,values...)))
}

/**
返回html
 */
func (c *Context) HTML(code int,name string,data interface{}){
	c.SetHeader("Content-Type","text/html")
	c.SetStatus(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.W,name,data); err != nil {
		fmt.Println("未找到模板文件！")
	}
}

func (c *Context) JSON(code int,data interface{}){
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.W)
	if err := encoder.Encode(data); err != nil {
		http.Error(c.W, err.Error(), 500)
	}
}

func (c *Context)  Next(){
	c.index ++
	s := len(c.handlers)
	for ; c.index < s ; c.index++ {
		c.handlers[c.index](c)
	}
}



