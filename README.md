# go-chat
## 框架集成了路由、中间件、控制器等功能

### (1) 框架运行入口：
> 执行目录下的 main.go 文件 可以修改 文件内的main函数来修改框架监听的端口和地址
> 格式为：ip:port 默认ip为127.0.0.1 可以直接修改为 :port 代表直接监听127.0.0.1:port

### (2) 路由文件：
> 在router文件下的 app.go文件参考其中的格式书写

### (3) 中间件文件：

> 在middleware文件下，中间件函数命名为 文件名+函数名

### (4) 控制器文件：

> 在controller文件下，控制器函数命名为 文件名+函数名
> 所有的路由处理函数对应控制器下的函数

### (5) 静态资源文件：

> 静态资源存放地为static文件夹下，所有的静态资源的路由为 /asset/ 映射到static文件夹下
> 所以静态资源的访问方式为： /asset/+static下对应的问价路径