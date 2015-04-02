package libvirt

import (
    "fmt"
    "encoding/xml"
)

type Listen struct {
    Type    string `xml:"type,attr"`
    Address string `xml:"address,attr"`
}

type Graphics struct {
    Type       string `xml:"type,attr"`
    Port       string `xml:"port,attr"`
    Autoport   string `xml:"autoport,attr"`
    ListenAttr string `xml:"listen,attr"`
    Passwd     string `xml:"passwd,attr"`
    Listen     Listen `xml:"listen"`
}

type Devices struct {
    Graphics Graphics `xml:"graphics"`
}

type Result struct {
    XMLName xml.Name `xml:"domain"`
    Type    string   `xml:"type,attr"`
    Devices Devices  `xml:"devices"`
}

func GetDomainGraphicsFromXML(dom_xml_desc string) (*Graphics, error) {
    var result Result

	err := xml.Unmarshal([]byte(dom_xml_desc), &result)

    return &result.Devices.Graphics, err
}

func GetDomainGraphics(host, port, vm_name string) (*Graphics, error) {
    var graphics *Graphics
    var err error

    virt_conn, err := GetConn(host, port)
    if err != nil {
        return graphics, err
    }

    dom, err := virt_conn.LookupDomainByName(vm_name)
    if err != nil {
        return graphics, err
    }

    dom_str, err := dom.GetXMLDesc(1)
    if err != nil {
        return graphics, err
    }

    dom_is_active, err := dom.IsActive()
    if err != nil {
        return graphics, err
    }

    if !dom_is_active {
        return graphics, fmt.Errorf("Domain %s is not active!\n", vm_name)
    }

    return GetDomainGraphicsFromXML(dom_str)
}
