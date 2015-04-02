package libvirtcontrollers

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/shelmesky/rconsole/libvirt"
	"github.com/shelmesky/rconsole/utils"
	"net"
	"net/http"
	"time"
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
	var console_host string
	var console_port string
	ws, err := upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	if err != nil {
		utils.Println("websocket upgrade failed:", err)
		return
	}

	defer func() {
		utils.Println("websocket client disconnected:", ws.RemoteAddr())
		ws.Close()
	}()

	libvirt_args, err := GetLIBVIRTArgs(this.Ctx)
	if err != nil {
		utils.Println(err)
		return
	}

	if libvirt_args == nil {
		utils.Println("empty args for LIBVIRT")
		return
	}

	libvirt_host := libvirt_args["hostname"]
	libvirt_port := libvirt_args["port"]
	vm_name := libvirt_args["vm"]

	graphics, err := libvirt.GetDomainGraphics(libvirt_host, libvirt_port, vm_name)
	if err != nil {
		utils.Println(err)
		return
	}

	if graphics.Listen.Address == "0.0.0.0" {
		console_host = libvirt_host
	} else {
		console_host = graphics.Listen.Address
	}

	console_port = graphics.Port

	console_conn, err := net.DialTimeout("tcp", console_host+":"+console_port, 3*time.Second)
	if err != nil {
		utils.Println("Can not connect to console host:", err)
		return
	}

	writer_buf := make([]byte, 4096)

	go func() {
		for {
			n, err := console_conn.Read(writer_buf)
			if err != nil {
				utils.Println("Error read from console host:", err)
				return
			}

			err = ws.WriteMessage(websocket.BinaryMessage, writer_buf[:n])
			if err != nil {
				utils.Println("websocket write failed:", err)
				return
			}
		}
	}()

	// TODO: use ws.NextReader to reuse memory allocation

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			utils.Println("websocket readmessage failed:", err)
			return
		}

		console_conn.Write(data)
	}
}
