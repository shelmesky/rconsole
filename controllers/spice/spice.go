package spicecontrollers

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/shelmesky/rconsole/utils"
	"net"
	"net/http"
	"time"
)

type SpiceController struct {
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

	spice_args, err := GetSPICEArgs(this.Ctx)
	if err != nil {
		utils.Println("get args for SPICE failed:", err)
		return
	}

	if spice_args == nil {
		utils.Println("empty args for SPICE")
		return
	}

	hostname := spice_args["hostname"]
	port := spice_args["port"]

	spice_conn, err := net.DialTimeout("tcp", hostname+":"+port, 3*time.Second)
	if err != nil {
		utils.Println("Can not connect to spice server:", err)
		return
	}

	writer_buf := make([]byte, 4096)

	go func() {
		for {
			n, err := spice_conn.Read(writer_buf)
			if err != nil {
				utils.Println("Error read from spice server:", err)
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

		spice_conn.Write(data)
	}
}
