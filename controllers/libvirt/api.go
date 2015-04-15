package libvirtcontrollers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/shelmesky/rconsole/mongo"
	"github.com/shelmesky/rconsole/utils"
)

type LibvirtHost struct {
	UUID     string `form:"-" bson:"uuid"`
	Hostname string `valid:"Required" form:"hostname"`
	Port     string `valid:"Required; Numeric" form:"port"`
	Username string `form:"username"`
	Password string `form:"password"`
}

type ResponseMessage struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

type HostController struct {
	beego.Controller
}

func (this *HostController) Get() {
}

func (this *HostController) Post() {
	var resp_message ResponseMessage
	var libvirt_host LibvirtHost
	var insert_id string

	err := this.ParseForm(&libvirt_host)
	if err != nil {
		err_msg := fmt.Sprintf("parseform failed: %s", err)
		utils.Println(err_msg)
		resp_message.Code = 1
		resp_message.Message = err_msg
		this.Ctx.Output.Status = 500
		goto end
	}

	err = this.Valid(libvirt_host)
	if err != nil {
		err_msg := fmt.Sprintf("validation failed: %s", err)
		utils.Println(err_msg)
		resp_message.Code = 2
		resp_message.Message = err_msg
		this.Ctx.Output.Status = 500
		goto end
	}

	libvirt_host.UUID = utils.MakeRandomID()
	insert_id = libvirt_host.UUID
	err = mongo.InsertOne("libvirt_host", libvirt_host)

	if resp_message.Code == 0 {
		resp_message.Code = 0
		resp_message.Message = insert_id
	}

end:
	this.Data["json"] = resp_message
	this.ServeJson()
}

func (this *HostController) Put() {
}

func (this *HostController) Delete() {
}

func (this *HostController) Valid(value interface{}) error {
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
