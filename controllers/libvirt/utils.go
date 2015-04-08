package libvirtcontrollers

import (
	"fmt"
	"github.com/astaxie/beego/context"
)

func GetLIBVIRTArgs(context *context.Context) (map[string]string, error) {
	libvirt_args := make(map[string]string, 0)

	Type := context.Input.Query("type")
	if Type == "" || Type != "libvirt" {
		return libvirt_args, fmt.Errorf("got wrong protocol: %s\n", Type)
	}

	libvirt_args["type"] = Type

	libvirt_args["hostname"] = context.Input.Query("hostname")
	libvirt_args["port"] = context.Input.Query("port")
	libvirt_args["vm"] = context.Input.Query("vm")
	libvirt_args["shared"] = context.Input.Query("shared")

	if libvirt_args["hostname"] == "" || libvirt_args["port"] == "" {
		goto get_args_failed
	}

	if libvirt_args["vm"] == "" {
		goto get_args_failed
	}

	return libvirt_args, nil

get_args_failed:
	return libvirt_args, fmt.Errorf("get args for LIBVIRT protocol failed!\n")
}
