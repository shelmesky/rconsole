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
	ws, err := upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	if err != nil {
		utils.Println("websocket upgrade failed:", err)
		return
	}

	defer func() {
		utils.Println("websocket client disconnected:", ws.RemoteAddr())
		ws.Close()
	}()

	connect_args := map[string]string{
		"username": "roy",
		"password": "password",
		"hostname": "172.31.31.110",
		"port":     "22",
	}

	client := rclient.NewClient("172.31.31.110", "4822", 3*time.Second, false)
	ret := client.HandShake("ssh", "1024", "650", "96", []string{}, []string{}, connect_args)
	if !ret {
		utils.Println("handshake failed!")
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

	for {
		message_type, message, err := ws.ReadMessage()
		if err != nil {
			utils.Println("websocket readmessage failed:", err)
			return
		}
		if message_type != websocket.TextMessage {
			utils.Println("invalid message type:", message_type)
		}

		client.Send(message)
	}
}
