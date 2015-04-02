package libvirt

import (
	"gopkg.in/alexzorin/libvirt-go.v2"
	"sync"
)

type GlobalConn struct {
	conn_map map[string]*libvirt.VirConnection
	lock     sync.RWMutex
}

var (
	global_conn *GlobalConn
)

func init() {
	global_conn = new(GlobalConn)
	global_conn.conn_map = make(map[string]*libvirt.VirConnection, 0)
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

	new_conn, err := libvirt.NewVirConnection("qemu+tcp://" + host + ":" + port + "/system")
	if err != nil {
		return &new_conn, err
	}

	global_conn.lock.Lock()
	global_conn.conn_map[host] = &new_conn
	global_conn.lock.Unlock()

	return &new_conn, nil
}
