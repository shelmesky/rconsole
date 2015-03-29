package spicecontrollers

import (
	"fmt"
	"github.com/astaxie/beego/context"
)

func GetSPICEArgs(context *context.Context) (map[string]string, error) {
	spice_args := make(map[string]string, 0)

	Type := context.Input.Query("type")
	if Type == "" || Type != "spice" {
		return spice_args, fmt.Errorf("got wrong protocol: %s\n", Type)
	}

	spice_args["type"] = Type

	spice_args["hostname"] = context.Input.Query("hostname")
	spice_args["port"] = context.Input.Query("port")
	spice_args["password"] = context.Input.Query("password")

	if spice_args["hostname"] == "" || spice_args["port"] == "" {
		goto get_args_failed
	}

	return spice_args, nil

get_args_failed:
	return spice_args, fmt.Errorf("get args for SPICE protocol failed!\n")
}
