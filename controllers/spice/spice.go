package spicecontrollers

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/shelmesky/rconsole/utils"
    "net"
	"net/http"
)

type SpiceController struct {
	beego.Controller
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin:     func(req *http.Request) bool { return true },
	}
)

func (this *SpiceController) Get() {
    ws, err := upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	if err != nil {
		utils.Println("websocket upgrade failed:", err)
		return
	}

	defer func() {
		utils.Println("websocket client disconnected:", ws.RemoteAddr())
		ws.Close()
	}()

	Type := this.GetString("type")

    if Type != "spice" {
		utils.Printf("invalid protocol: %s\n", Type)
		return
    }

    spice_args, err := GetSpiceArgs(this.Ctx)
    if err != nil {
        utils.Prinln("get args for SPICE failed:", err)
        return
    }

    if spice_args == nil {
        utils.Prinln("empty args for SPICE")
        return
    }

    hostname := spice_args["hostname"]
    port := spice_args["port"]
    password := spice_args["password"]
}
