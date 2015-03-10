package main

import (
	"github.com/astaxie/beego"
	"github.com/shelmesky/rconsole/controllers/primary"
	"github.com/shelmesky/rconsole/controllers/websocket"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	beego.SetStaticPath("/static", "./static")
	beego.Router("/", &controllers.MainController{})
	beego.Router("/ws/", &wscontrollers.WebSocketController{})
	beego.Run()
}
