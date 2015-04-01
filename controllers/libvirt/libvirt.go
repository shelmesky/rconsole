package libvirtcontrollers

import (
    "net/http"
	"github.com/astaxie/beego"
	"github.com/shelmesky/rconsole/libvirt"
	"github.com/shelmesky/rconsole/utils"
	"github.com/gorilla/websocket"
)

type LibvirtController struct {
	beego.Controller
}

var (
	SUB_PROTOCOLS = []string{"binary"}
	upgrader      = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin:     func(req *http.Request) bool { return true },
		Subprotocols:    SUB_PROTOCOLS,
	}
)

func (this *LibvirtController) Get() {
	ws, err := upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	if err != nil {
		utils.Println("websocket upgrade failed:", err)
		return
	}

	defer func() {
		utils.Println("websocket client disconnected:", ws.RemoteAddr())
		ws.Close()
	}()

    virt_conn, err := libvirt.GetConn("127.0.0.1", "16509")
    if err != nil {
        utils.Println(err)
    }

    dom, err := virt_conn.LookupDomainByName("win7")
    if err != nil {
        utils.Println(err)
    }

    dom_str, err := dom.GetXMLDesc(0)
    if err != nil {
        utils.Println(err)
    }

    utils.Println(dom_str)
}
