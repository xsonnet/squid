## Squid框架
#### 特色
- 路由支持正则匹配
- 自带Log
- 支持json配置文件读写
- 支持session
- 支持pongo2模板引擎
#### 安装方法
下载和安装
```
$ go get -u github.com/xsonnet/squid
```
导入到项目代码中
```
import "github.com/xsonnet/squid"
```

#### 示例
```
package main

import "github.com/xsonnet/squid"

func main() {
	app := squid.InitApp([]squid.Router{
	    {"^/$", home},
	})
	app.Run(":8080")
}
func home(ctx squid.Context) {
    ctx.Render("home", squid.Params{"hi": "Hello world."})
}
```