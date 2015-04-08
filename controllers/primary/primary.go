package controllers

import (
	"github.com/astaxie/beego"
	"github.com/shelmesky/rconsole/client"
	"github.com/shelmesky/rconsole/controllers/libvirt"
	"github.com/shelmesky/rconsole/controllers/spice"
	"github.com/shelmesky/rconsole/controllers/websocket"
	"github.com/shelmesky/rconsole/libvirt"
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
	var kvm_console_type string

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

	if Type == "libvirt" {
		url_args, err = libvirtcontrollers.GetLIBVIRTArgs(this.Ctx)
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

	// if type is spice, maybe spice server need password
	// so we try to get it from url args
	if Type == "spice" {
		this.Data["PASSWORD"] = url_args["password"]
	}

	// if type is libvirt(with KVM), the console type is vnc or spice
	// if need password when connect, get it use libvirt API
	if Type == "libvirt" {
		host := url_args["hostname"]
		port := url_args["port"]
		vm_name := url_args["vm"]

		graphics, err := libvirt.GetDomainGraphics(host, port, vm_name)
		if err != nil {
			utils.Println("call GetDomainGraphics failed:", err)
			this.Abort("500")
			return
		}

		this.Data["PASSWORD"] = graphics.Passwd
		kvm_console_type = graphics.Type
	}

	if Type == "spice" {
		this.TplNames = "spice.html"
	} else if Type == "libvirt" && kvm_console_type == "vnc" {
		if url_args["shared"] == "yes" {
			this.Data["VNC_SHARED"] = "yes"
		} else {
			this.Data["VNC_SHARED"] = "no"
		}
		this.TplNames = "libvirt_vnc.html"
	} else if Type == "libvirt" && kvm_console_type == "spice" {
		this.TplNames = "libvirt_spice.html"
	} else {
		this.TplNames = "index.html"
	}
	this.Render()
}
