package main

import (
	"encoding/json"
	"flag"
	"github.com/astaxie/beego"
	"github.com/shelmesky/rconsole/client"
	"github.com/shelmesky/rconsole/controllers/api"
	"github.com/shelmesky/rconsole/controllers/libvirt"
	"github.com/shelmesky/rconsole/controllers/primary"
	"github.com/shelmesky/rconsole/controllers/spice"
	"github.com/shelmesky/rconsole/controllers/websocket"
	"github.com/shelmesky/rconsole/mongo"
	"github.com/shelmesky/rconsole/utils"
	"gopkg.in/alexzorin/libvirt-go.v2"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// 全局配置
type GlobalConfig struct {
	LogFile        string `json:"log_file"`
	ListenAddress  string `json:"listen_address"`
	ProfileAddress string `json:"profile_address"`
	MongoURL       string `json:"mongo_url"`
	MongoTimeout   int    `json:"mongo_timeout"`

	ServerProfile bool
	ShowCheckConn bool           `json:"show_check_conn"`
	ServerDebug   bool           `json:"debug"`
	GuacdServer   []client.Guacd `json:"guacd"`
}

var (
	SignalChan chan os.Signal
	Config     GlobalConfig

	// 命令行配置
	ConfigFile = flag.String("config", "./config.json", "")

	LogFile        = flag.String("log_file", "", "logging file, default: /var/log/rconsole.log")
	ProfileAddress = flag.String("profile_address", "", "profile server listen on, default: 0.0.0.0:9998")
	ListenAddress  = flag.String("listen", "", "server listen on, default: 0.0.0.0:9999")
	MongoURL       = flag.String("mongo_url", "", "MongoDB URL, default: mongo://127.0.0.1/rconsole")
	MongoTimeout   = flag.Int("mongo_timeout", 0, "Wait n seconds for connect to MongoDB, default: 1")

	ShowCheckConn = flag.Bool("show_check_conn", false, "Show check MongoDB server informatioin")

	ServerProfile = flag.Bool("profile", true, "Start web profile interface")
	ServerDebug   = flag.Bool("debug", false, "Print debug information when server is running.")
)

// 判断文件或目录是否存在
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

/*
判断命令行参数和配置文件
如果配置文件不存在，则会使用命令行的参数
否则直接使用配置文件中的参数
*/
func Init() {
	if *ServerProfile == true {
		Config.ServerProfile = true
	}

	if !Exist(*ConfigFile) {
		if *LogFile == "" {
			*LogFile = "/var/log/rconsole.log"
		}
		Config.LogFile = *LogFile

		if *ProfileAddress == "" {
			*ProfileAddress = "0.0.0.0:9998"
		}
		Config.ProfileAddress = *ProfileAddress

		if *ListenAddress == "" {
			*ListenAddress = "0.0.0.0:9999"
		}
		Config.ListenAddress = *ListenAddress

		if *MongoURL == "" {
			*MongoURL = "mongo://127.0.0.1/rconsole"
		}
		Config.MongoURL = *MongoURL

		if *MongoTimeout == 0 {
			*MongoTimeout = 1
		}
		Config.MongoTimeout = *MongoTimeout

		if *ShowCheckConn == true {
			*ShowCheckConn = true
		}
		Config.ShowCheckConn = *ShowCheckConn

		if *ServerDebug == true {
			*ServerDebug = true
		}
		Config.ServerDebug = *ServerDebug

	} else {
		data, err := ioutil.ReadFile(*ConfigFile)
		if err != nil {
			utils.Println(err)
			return
		}

		err = json.Unmarshal(data, &Config)
		if err != nil {
			utils.Println(err)
		}

		if *LogFile != "" {
			Config.LogFile = *LogFile
		}

		if *ProfileAddress != "" {
			Config.ProfileAddress = *ProfileAddress
		}

		if *ListenAddress != "" {
			Config.ListenAddress = *ListenAddress
		}

		if *MongoURL != "" {
			Config.MongoURL = *MongoURL
		}

		if *MongoTimeout != 0 {
			Config.MongoTimeout = *MongoTimeout
		}

		if *ShowCheckConn != false {
			Config.ShowCheckConn = *ShowCheckConn
		}

		if *ServerDebug != false {
			Config.ServerDebug = *ServerDebug
		}
	}
}

// POSIX信号发生时的回调
func SignalCallback() {
	for s := range SignalChan {
		sig := s.String()
		utils.Println("Got Signal: ", sig)

		if s == syscall.SIGINT {
			utils.Println("RConsole server exit...")
			os.Exit(0)
		}
	}
}

func main() {
	ret := libvirt.EventRegisterDefaultImpl()
	utils.Println("EventRegisterDefaultImpl ret:", ret)

	runtime.GOMAXPROCS(runtime.NumCPU())

	defer func() {
		if err := recover(); err != nil {
			utils.Println("Server Error: ", err.(string))
		}
	}()

	// 解析命令行参数
	flag.Parse()

	// 自定义初始化
	Init()

	// 设置guacd客户端连接时的调试信息
	client.ClientDebug = Config.ServerDebug

	// 初始化guacd服务连接池
	client.Pool.Init(Config.GuacdServer)

	utils.Println("Try connect to:", Config.MongoURL)

	// 初始化MongoDB的连接
	timeout := time.Duration(Config.MongoTimeout) * time.Second
	err := mongo.InitMongoDB(Config.MongoURL, timeout)
	if err != nil {
		utils.Println(err)
		os.Exit(1)
	}

	// 启动独立的线程测试MongoDB连接
	go mongo.TouchMongoDB(Config.ShowCheckConn)

	utils.Println("Connect to MongoDB OK")

	// 捕捉信号并设置回调函数
	SignalChan = make(chan os.Signal, 1)
	signal.Notify(SignalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGPIPE,
		syscall.SIGALRM,
		syscall.SIGPIPE)

	go SignalCallback()

	// 启动性能调试接口
	if Config.ServerProfile == true {
		go func() {
			http.ListenAndServe(Config.ProfileAddress, nil)
		}()
	}

	utils.Printf("http profile server Running on %s\n", Config.ProfileAddress)

	addr_splited := strings.Split(Config.ListenAddress, ":")
	if len(addr_splited) < 2 {
		utils.Println("the format of listen address is invalid!")
		os.Exit(1)
	}

	beego.HttpAddr = addr_splited[0]
	beego.HttpPort, err = strconv.Atoi(addr_splited[1])

	beego.SetStaticPath("/static", "./static")
	beego.Router("/connect", &controllers.MainController{})
	beego.Router("/ws", &wscontrollers.WebSocketController{})
	beego.Router("/ws/spice", &spicecontrollers.SpiceController{})
	beego.Router("/ws/libvirt", &libvirtcontrollers.LibvirtController{})

	beego.Router("/api/conn/list", &managercontrollers.ConnectionManagerController{}, "get:ListConnection")
	beego.Router("/api/conn/create/:conn_type([a-z]+)", &managercontrollers.ConnectionManagerController{}, "post:CreateConnection")
	beego.Router("/api/conn/update/:conn_type([a-z]+)", &managercontrollers.ConnectionManagerController{}, "put:UpdateConnection")
	beego.Router("/api/conn/delete", &managercontrollers.ConnectionManagerController{}, "delete:DeleteConnection")

	beego.Router("/api/libvirt/host", &libvirtcontrollers.HostController{})

	beego.Run()
}
