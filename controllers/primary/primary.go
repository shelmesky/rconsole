package controllers

import (
	"github.com/astaxie/beego"
	"github.com/shelmesky/rconsole/client"
	"github.com/shelmesky/rconsole/controllers/websocket"
	"github.com/shelmesky/rconsole/controllers/spice"
	"github.com/shelmesky/rconsole/utils"
	"strings"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	var found_protocol bool
	var url_args map[string]string
	var connect_args string
	var err error
	var args_list []string

	Type := this.GetString("type")
	for index := range client.PROTOCOLS {
		if Type == client.PROTOCOLS[index] {
			found_protocol = true
			break
		}
	}

	if !found_protocol {
		utils.Printf("Invalid protocol: %s\n", Type)
		this.Abort("500")
		return
	}

	if Type == "vnc" {
		url_args, err = wscontrollers.GetVNCArgs(this.Ctx)
	}

	if Type == "rdp" {
		url_args, err = wscontrollers.GetRDPArgs(this.Ctx)
	}

	if Type == "ssh" {
		url_args, err = wscontrollers.GetSSHArgs(this.Ctx)
	}

	if Type == "telnet" {
		url_args, err = wscontrollers.GetTELNETArgs(this.Ctx)
	}

    if Type == "spice" {
		url_args, err = spicecontrollers.GetSPICEArgs(this.Ctx)
    }

	if err != nil {
		utils.Println(err)
		this.Abort("500")
		return
	}

	if url_args == nil {
		utils.Printf("empty url_args for protcol: %s\n", Type)
		this.Abort("500")
		return
	}

	for key := range url_args {
		args_list = append(args_list, key+"="+url_args[key])
	}

	connect_args = strings.Join(args_list, "&")

	this.Data["CONNECT_ARGS"] = connect_args

    if Type == "spice" {
        this.Data["PASSWORD"] = url_args["password"]
    }

    if Type != "spice" {
	    this.TplNames = "index.html"
    } else {
        this.TplNames = "spice-old.html"
    }
	this.Render()
}
