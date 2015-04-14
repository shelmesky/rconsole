package wscontrollers

import (
	"bytes"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	rclient "github.com/shelmesky/rconsole/client"
	"github.com/shelmesky/rconsole/utils"
	"net/http"
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

	found_protocol = rclient.ValidProtocol(Type)

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

	client, err := rclient.Pool.Get()
	if err != nil {
		utils.Println("can not find any client")
		return
	}

	defer client.Close()

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

	read_buffer := make([]byte, 0, 4096)
	buf := bytes.NewBuffer(read_buffer)

	for {
		message_type, reader, err := ws.NextReader()
		if err != nil {
			utils.Printf("failed when call NextReader: %s, message_type: %d\n", err, message_type)
			return
		}

		buf.Truncate(0)
		n, err := buf.ReadFrom(reader)
		if err != nil {
			utils.Println(n, err)
		}

		client.Send(buf.Bytes())
	}
}
