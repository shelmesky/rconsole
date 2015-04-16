package libvirt

import (
	"github.com/shelmesky/rconsole/utils"
	"gopkg.in/alexzorin/libvirt-go.v2"
	"sync"
)

type GlobalConn struct {
	conn_map map[string]*libvirt.VirConnection
	lock     sync.RWMutex
}

var (
	global_conn *GlobalConn
	event_map   map[string]bool
)

func init() {
	global_conn = new(GlobalConn)
	global_conn.conn_map = make(map[string]*libvirt.VirConnection, 0)
	event_map = make(map[string]bool, 0)
}

func GetConn(host, port string) (*libvirt.VirConnection, error) {
	var conn *libvirt.VirConnection
	var err error
	var ok bool

	global_conn.lock.RLock()
	if conn, ok = global_conn.conn_map[host]; ok {
		if ok, err = conn.IsAlive(); ok {
			global_conn.lock.RUnlock()
			return conn, nil
		}
	}
	global_conn.lock.RUnlock()

	address := host + ":" + port
	new_conn, err := libvirt.NewVirConnection("qemu+tcp://" + address + "/system")
	if err != nil {
		return &new_conn, err
	}

	if v, ok := event_map[address]; !ok {
		if v == false {
			utils.Println("register event...")

			var Callback libvirt.DomainEventCallback
			Callback = EventCallback

			test := func() {}
			//dom := libvirt.VirDomain{}
			dom, _ := new_conn.LookupDomainByName("i-2-161-VM")

			ret := new_conn.DomainEventRegister(dom, libvirt.VIR_DOMAIN_EVENT_ID_LIFECYCLE, &Callback, test)
			utils.Println("LIFECYCLE ret:", ret)

			//ret = new_conn.DomainEventRegister(dom, libvirt.VIR_DOMAIN_EVENT_ID_REBOOT, &Callback, test)
			//utils.Println("ret:", ret)

			//ret = new_conn.DomainEventRegister(dom, libvirt.VIR_DOMAIN_EVENT_ID_BLOCK_JOB, &Callback, test)
			//utils.Println("ret:", ret)

			go func() {
				for {
					ret := libvirt.EventRunDefaultImpl()
					utils.Println("EventRunDefaultImpl ret:", ret)
				}
			}()

			event_map[address] = true
		}
	}

	global_conn.lock.Lock()
	global_conn.conn_map[host] = &new_conn
	global_conn.lock.Unlock()

	return &new_conn, nil
}
