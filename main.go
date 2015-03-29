package main

import (
	"github.com/astaxie/beego"
	"github.com/shelmesky/rconsole/controllers/primary"
	"github.com/shelmesky/rconsole/controllers/spice"
	"github.com/shelmesky/rconsole/controllers/websocket"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	beego.SetStaticPath("/static", "./static")
	beego.Router("/connect", &controllers.MainController{})
	beego.Router("/ws", &wscontrollers.WebSocketController{})
	beego.Router("/ws/spice", &spicecontrollers.SpiceController{})
	beego.Run()
}
