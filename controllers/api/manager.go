package managercontrollers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/shelmesky/rconsole/client"
	"github.com/shelmesky/rconsole/lib"
	"github.com/shelmesky/rconsole/mongo"
	"github.com/shelmesky/rconsole/utils"
)

type ResponseMessage struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

type ConnectionManagerController struct {
	beego.Controller
}

func (this *ConnectionManagerController) ListConnection() {
}

func (this *ConnectionManagerController) CreateConnection() {
	var resp_message ResponseMessage
	var found_protocol bool
	var decode_failed bool
	var decode_failed_reason string
	var insert_failed bool
	var insert_failed_reason string
	var insert_id string
	var valid_failed bool
	var valid_failed_reason string

	conn_type := this.Ctx.Input.Param(":conn_type")

	found_protocol = client.ValidProtocol(conn_type)

	if found_protocol {

		if conn_type == "vnc" {
			args, err := this.DecodeVNCArgs()
			if err != nil {
				decode_failed = true
				decode_failed_reason = fmt.Sprintf("decode vnc args failed: %s", err)
			} else {
				err = this.Valid(args)
				if err != nil {
					valid_failed = true
					valid_failed_reason = fmt.Sprintf("valid error: %s", err)
				} else {
					args.UUID = utils.MakeRandomID()
					insert_id = args.UUID
					err = mongo.InsertOne("connection", *args)
					if err != nil {
						insert_failed = true
						insert_failed_reason = fmt.Sprintf("save vnc args failed: %s, args: %s", err, *args)
					}
				}
			}
		}

		if conn_type == "rdp" {
			args, err := this.DecodeRDPArgs()
			if err != nil {
				decode_failed = true
				decode_failed_reason = fmt.Sprintf("decode rdp args failed: %s", err)
			} else {
				err = this.Valid(args)
				if err != nil {
					valid_failed = true
					valid_failed_reason = fmt.Sprintf("valid error: %s", err)
				} else {
					args.UUID = utils.MakeRandomID()
					insert_id = args.UUID
					err = mongo.InsertOne("connection", *args)
					if err != nil {
						insert_failed = true
						insert_failed_reason = fmt.Sprintf("save rdp args failed: %s, args: %s", err, *args)
					}
				}
			}
		}

		if conn_type == "ssh" {
			args, err := this.DecodeSSHArgs()
			if err != nil {
				decode_failed = true
				decode_failed_reason = fmt.Sprintf("decode ssh args failed: %s", err)
			} else {
				err = this.Valid(args)
				if err != nil {
					valid_failed = true
					valid_failed_reason = fmt.Sprintf("valid error: %s", err)
				} else {
					args.UUID = utils.MakeRandomID()
					insert_id = args.UUID
					err = mongo.InsertOne("connection", *args)
					if err != nil {
						insert_failed = true
						insert_failed_reason = fmt.Sprintf("save ssh args failed: %s, args: %s", err, *args)
					}
				}
			}
		}

		if conn_type == "telnet" {
			args, err := this.DecodeTELNETArgs()
			if err != nil {
				decode_failed = true
				decode_failed_reason = fmt.Sprintf("decode telnet args failed: %s", err)
			} else {
				err = this.Valid(args)
				if err != nil {
					valid_failed = true
					valid_failed_reason = fmt.Sprintf("valid error: %s", err)
				} else {
					args.UUID = utils.MakeRandomID()
					insert_id = args.UUID
					err = mongo.InsertOne("connection", *args)
					if err != nil {
						insert_failed = true
						insert_failed_reason = fmt.Sprintf("save telnet args failed: %s, args: %s", err, *args)
					}
				}
			}
		}

		if conn_type == "spice" {
			args, err := this.DecodeSPICEArgs()
			if err != nil {
				decode_failed = true
				decode_failed_reason = fmt.Sprintf("decode spice args failed: %s", err)
			} else {
				err = this.Valid(args)
				if err != nil {
					valid_failed = true
					valid_failed_reason = fmt.Sprintf("valid error: %s", err)
				} else {
					args.UUID = utils.MakeRandomID()
					insert_id = args.UUID
					err = mongo.InsertOne("connection", *args)
					if err != nil {
						insert_failed = true
						insert_failed_reason = fmt.Sprintf("save spice args failed: %s, args: %s", err, *args)
					}
				}
			}
		}

		if conn_type == "libvirt" {
			args, err := this.DecodeLIBVIRTArgs()
			if err != nil {
				decode_failed = true
				decode_failed_reason = fmt.Sprintf("decode libvirt args failed: %s", err)
			} else {
				err = this.Valid(args)
				if err != nil {
					valid_failed = true
					valid_failed_reason = fmt.Sprintf("valid error: %s", err)
				} else {
					args.UUID = utils.MakeRandomID()
					insert_id = args.UUID
					err = mongo.InsertOne("connection", *args)
					if err != nil {
						insert_failed = true
						insert_failed_reason = fmt.Sprintf("save libvirt args failed: %s, args: %s", err, *args)
					}
				}
			}
		}

		if decode_failed {
			utils.Println(decode_failed_reason)
			resp_message.Code = 1
			resp_message.Message = decode_failed_reason
			this.Ctx.Output.Status = 500
		}

		if insert_failed {
			utils.Println(insert_failed_reason)
			resp_message.Code = 2
			resp_message.Message = insert_failed_reason
			this.Ctx.Output.Status = 500
		}

		if valid_failed {
			utils.Println(valid_failed_reason)
			resp_message.Code = 3
			resp_message.Message = valid_failed_reason
			this.Ctx.Output.Status = 500
		}

	} else {
		resp_message.Code = 4
		resp_message.Message = "wrong protocol type"
		this.Ctx.Output.Status = 400
	}

	if resp_message.Code == 0 {
		resp_message.Code = 0
		resp_message.Message = insert_id
	}

	this.Data["json"] = resp_message
	this.ServeJson()
}

func (this *ConnectionManagerController) UpdateConnection() {
}

func (this *ConnectionManagerController) DeleteConnection() {
}

func (this *ConnectionManagerController) DecodeVNCArgs() (*lib.VNCArgs, error) {
	var vnc_args lib.VNCArgs
	vnc_args.Type = "vnc"
	err := this.ParseForm(&vnc_args)
	return &vnc_args, err
}

func (this *ConnectionManagerController) DecodeRDPArgs() (*lib.RDPArgs, error) {
	var rdp_args lib.RDPArgs
	rdp_args.Type = "rdp"
	err := this.ParseForm(&rdp_args)
	return &rdp_args, err
}

func (this *ConnectionManagerController) DecodeSSHArgs() (*lib.SSHArgs, error) {
	var ssh_args lib.SSHArgs
	ssh_args.Type = "ssh"
	err := this.ParseForm(&ssh_args)
	return &ssh_args, err
}

func (this *ConnectionManagerController) DecodeTELNETArgs() (*lib.TELNETArgs, error) {
	var telnet_args lib.TELNETArgs
	telnet_args.Type = "telnet"
	err := this.ParseForm(&telnet_args)
	return &telnet_args, err
}

func (this *ConnectionManagerController) DecodeSPICEArgs() (*lib.SPICEArgs, error) {
	var spice_args lib.SPICEArgs
	spice_args.Type = "spice"
	err := this.ParseForm(&spice_args)
	return &spice_args, err
}

func (this *ConnectionManagerController) DecodeLIBVIRTArgs() (*lib.LIBVIRTArgs, error) {
	var libvirt_args lib.LIBVIRTArgs
	libvirt_args.Type = "libvirt"
	err := this.ParseForm(&libvirt_args)
	return &libvirt_args, err
}

func (this *ConnectionManagerController) Valid(value interface{}) error {
	valid := validation.Validation{}
	b, err := valid.Valid(value)
	if err != nil {
		return err
	}
	if !b {
		for _, err := range valid.Errors {
			return fmt.Errorf("%s %s", err.Key, err.Message)
		}
	}
	return nil
}
