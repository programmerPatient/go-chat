package marry

import (
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type Router struct {
	roots map[string]*node
	handler map[string]HandlerFunc
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern,"/")

	parts := make([]string,0)
	for _,item := range vs {
		if item != "" {
			parts = append(parts,item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

/**
添加路由
 */
func (router *Router) addRouter(method string,path string,handler HandlerFunc){
	key := method + "-" + path
	parts := parsePattern(path)

	_,ok := router.roots[method]
	if !ok {
		router.roots[method] = &node{}
	}
	router.roots[method].insert(path,parts,0)
	router.handler[key] = handler
}

func (router *Router) getRoute(method string,path string) (*node,map[string]string){
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root,ok := router.roots[method]
	if !ok {
		return nil,nil
	}
	n := root.search(searchParts,0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index,part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:],"/")
			}
		}
		return n,params
	}
	return nil,nil
}

func (router *Router) handle(c *Context) {
	n, params := router.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers,router.handler[key])
	} else {
		c.handlers = append(c.handlers,func(c *Context){
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}

func newRouter() *Router{
	return &Router{
		roots: make(map[string]*node),
		handler: make(map[string]HandlerFunc),
	}
}
