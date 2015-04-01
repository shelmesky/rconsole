package libvirt

import (
    "gopkg.in/alexzorin/libvirt-go.v2"
)

var (
    global_conn map[string]*libvirt.VirConnection
)

func init() {
    global_conn = make(map[string]*libvirt.VirConnection, 0)
}

func GetConn(host, port string) (*libvirt.VirConnection, error) {
    var conn *libvirt.VirConnection
    var err error
    var ok bool

    if conn, ok = global_conn[host]; ok {
        if ok, err = conn.IsAlive(); ok {
            return conn, nil
        }
    }

    new_conn, err := libvirt.NewVirConnection("qemu+tcp://"+host+":"+port+"/system")
    if err != nil {
        return &new_conn, err
    }

    global_conn[host] = &new_conn

    return &new_conn, nil
}

