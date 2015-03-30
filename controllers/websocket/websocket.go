package wscontrollers

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	rclient "github.com/shelmesky/rconsole/client"
	"github.com/shelmesky/rconsole/utils"
	"net/http"
	"time"
)

type WebSocketController struct {
	beego.Controller
}

var (
	SUB_PROTOCOLS = []string{rclient.PROTOCOL_NAME}
	upgrader      = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin:     func(req *http.Request) bool { return true },
		Subprotocols:    SUB_PROTOCOLS,
	}
)

func (this *WebSocketController) Get() {
	var url_args map[string]string
	var err error
	var found_protocol bool

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
	for index := range rclient.PROTOCOLS {
		if Type == rclient.PROTOCOLS[index] {
			found_protocol = true
			break
		}
	}

	if !found_protocol {
		utils.Printf("invalid protocol: %s\n", Type)
		return
	}

	if Type == "vnc" {
		url_args, err = GetVNCArgs(this.Ctx)
	}

	if Type == "rdp" {
		url_args, err = GetRDPArgs(this.Ctx)
	}

	if Type == "ssh" {
		url_args, err = GetSSHArgs(this.Ctx)
	}

	if Type == "telnet" {
		url_args, err = GetTELNETArgs(this.Ctx)
	}

	if err != nil {
		utils.Println("get url args failed:", err)
		return
	}

	if url_args == nil {
		utils.Printf("empty url_args for protcol: %s\n", Type)
		return
	}

	protocol_type := url_args["type"]
	width := url_args["width"]
	height := url_args["height"]
	dpi := url_args["dpi"]

	client := rclient.NewClient("172.31.31.110", "4822", 3*time.Second, false)

	ret := client.HandShake(protocol_type, width, height, dpi, []string{}, []string{}, url_args)
	if !ret {
		utils.Println("handshake failed!")
		return
	}

	go func() {
		for {
			instruction := client.BufReceive()
			err := ws.WriteMessage(websocket.TextMessage, []byte(instruction))
			if err != nil {
				utils.Println("websocket writemessage failed:", err)
				return
			}
		}
	}()

	// TODO: use ws.NextReader to reuse memory allocation
	for {
		message_type, message, err := ws.ReadMessage()
		if err != nil {
			utils.Println("websocket readmessage failed:", err)
			return
		}
		if message_type != websocket.TextMessage {
			utils.Println("invalid message type:", message_type)
			return
		}

		client.Send(message)
	}
}
