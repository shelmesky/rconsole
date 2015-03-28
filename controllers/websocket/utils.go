package wscontrollers

import (
	"fmt"
	"github.com/astaxie/beego/context"
)

func GetVNCArgs(context *context.Context) (map[string]string, error) {
	vnc_args := make(map[string]string, 0)

	Type := context.Input.Query("type")
	if Type == "" && Type != "vnc" {
		return vnc_args, fmt.Errorf("got wrong protocol: %s\n", Type)
	}

	vnc_args["type"] = Type

	vnc_args["hostname"] = context.Input.Query("hostname")
	vnc_args["port"] = context.Input.Query("port")
	vnc_args["password"] = context.Input.Query("password")
	vnc_args["width"] = context.Input.Query("width")
	vnc_args["height"] = context.Input.Query("height")
	vnc_args["dpi"] = context.Input.Query("dpi")

	if vnc_args["hostname"] == "" || vnc_args["port"] == "" {
		goto get_args_failed
	}

	if vnc_args["password"] == "" || vnc_args["width"] == "" {
		goto get_args_failed
	}

	if vnc_args["height"] == "" || vnc_args["dpi"] == "" {
		goto get_args_failed
	}

	return vnc_args, nil

get_args_failed:
	return vnc_args, fmt.Errorf("get args for VNC protocol failed!\n")

}

func GetRDPArgs(context *context.Context) (map[string]string, error) {
	rdp_args := make(map[string]string, 0)

	Type := context.Input.Query("type")
	if Type == "" && Type != "rdp" {
		return rdp_args, fmt.Errorf("got wrong protocol: %s\n", Type)
	}

	rdp_args["type"] = Type

	rdp_args["hostname"] = context.Input.Query("hostname")
	rdp_args["port"] = context.Input.Query("port")
	rdp_args["username"] = context.Input.Query("username")
	rdp_args["width"] = context.Input.Query("width")
	rdp_args["height"] = context.Input.Query("height")
	rdp_args["dpi"] = context.Input.Query("dpi")

	if rdp_args["hostname"] == "" || rdp_args["port"] == "" {
		goto get_args_failed
	}

	if rdp_args["username"] == "" || rdp_args["width"] == "" {
		goto get_args_failed
	}

	if rdp_args["height"] == "" || rdp_args["dpi"] == "" {
		goto get_args_failed
	}

	return rdp_args, nil

get_args_failed:
	return rdp_args, fmt.Errorf("get args for RDP protocol failed!\n")
}

func GetSSHArgs(context *context.Context) (map[string]string, error) {
	ssh_args := make(map[string]string, 0)

	Type := context.Input.Query("type")
	if Type == "" && Type != "ssh" {
		return ssh_args, fmt.Errorf("got wrong protocol: %s\n", Type)
	}

	ssh_args["type"] = Type

	ssh_args["hostname"] = context.Input.Query("hostname")
	ssh_args["port"] = context.Input.Query("port")
	ssh_args["username"] = context.Input.Query("username")
	ssh_args["width"] = context.Input.Query("width")
	ssh_args["height"] = context.Input.Query("height")
	ssh_args["dpi"] = context.Input.Query("dpi")

	if ssh_args["hostname"] == "" || ssh_args["port"] == "" {
		goto get_args_failed
	}

	if ssh_args["username"] == "" || ssh_args["width"] == "" {
		goto get_args_failed
	}

	if ssh_args["height"] == "" {
		goto get_args_failed
	}

	return ssh_args, nil

get_args_failed:
	return ssh_args, fmt.Errorf("get args for SSH protocol failed!\n")

}

func GetTELNETArgs(context *context.Context) (map[string]string, error) {
	telnet_args := make(map[string]string, 0)

	Type := context.Input.Query("type")
	if Type == "" && Type != "telnet" {
		return telnet_args, fmt.Errorf("got wrong protocol: %s\n", Type)
	}

	telnet_args["type"] = Type

	telnet_args["hostname"] = context.Input.Query("hostname")
	telnet_args["port"] = context.Input.Query("port")
	telnet_args["username"] = context.Input.Query("username")
	telnet_args["width"] = context.Input.Query("width")
	telnet_args["height"] = context.Input.Query("height")
	telnet_args["dpi"] = context.Input.Query("dpi")

	if telnet_args["hostname"] == "" || telnet_args["port"] == "" {
		goto get_args_failed
	}

	if telnet_args["username"] == "" || telnet_args["width"] == "" {
		goto get_args_failed
	}

	if telnet_args["height"] == "" {
		goto get_args_failed
	}

	return telnet_args, nil

get_args_failed:
	return telnet_args, fmt.Errorf("get args for TELNET protocol failed!\n")

}
