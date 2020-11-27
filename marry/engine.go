/*
 *封装http请求的包
 */
package marry

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
)


type Engine struct{
	*RouterGroup
	htmlTemplates *template.Template
	funcMap template.FuncMap
	router *Router
	groups []*RouterGroup // store all groups
}
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}

func (g *RouterGroup) addRouter(method string,path string,handler HandlerFunc)  {
	path = g.prefix + "/" +path
	g.engine.router.addRouter(method,path,handler)
}

func (g *RouterGroup) GET(path string,handler HandlerFunc){

	g.addRouter("GET",path,handler)
}

func (g *RouterGroup) POST(path string,handler HandlerFunc){
	g.addRouter("POST",path,handler)
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares,middlewares...)
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		prefix:      prefix,
		middlewares: g.middlewares,
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}






func (engine *Engine) Run(addr string)(err error){
	return http.ListenAndServe(addr,engine)
}

func (engine *Engine) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := relativePath
	fileServer := http.StripPrefix(absolutePath,http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		fmt.Println(file)
		if _,err := fs.Open(file);err != nil{
			c.Status = http.StatusNotFound
			return
		}
		fileServer.ServeHTTP(c.W,c.R)
	}
}
/**
访问静态资源
*/
func (engine *Engine) Static(relativePath string,root string) {
	handler := engine.createStaticHandler(relativePath,http.Dir(root))
	urlPattern := path.Join(relativePath,"/*filepath")
	engine.GET(urlPattern,handler)
}


func (engine *Engine) SetFuncMap(funcMap template.FuncMap){
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string){
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}



func New() *Engine{
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}


func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request){
	var middleware []HandlerFunc
	for _,group := range engine.groups {
		if strings.HasPrefix(req.URL.Path[1:],group.prefix) {
			middleware = append(middleware,group.middlewares...)
		}
	}
	context := nowContext(w,req)
	context.handlers = middleware
	context.engine = engine
	engine.router.handle(context)
}
