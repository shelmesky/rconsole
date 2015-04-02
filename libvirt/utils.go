package libvirt

import (
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

func GetDomainGraphics(dom_xml_desc string) (*Graphics, error) {
    var result Result

	err := xml.Unmarshal([]byte(dom_xml_desc), &result)

    return &result.Devices.Graphics, err
}
