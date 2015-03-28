package wscontrollers

import (
	"fmt"
	"github.com/astaxie/beego/context"
)

func GetVNCArgs(context *context.Context) (map[string]string, error) {
	return nil, nil
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
	return nil, nil
}

func GetTELNETArgs(context *context.Context) (map[string]string, error) {
	return nil, nil
}
