package client

import (
	"bufio"
	"fmt"
	"github.com/shelmesky/rconsole/client/inst"
	"github.com/shelmesky/rconsole/utils"
	"net"
	"strings"
	"time"
)

var (
	PROTOCOLS      = []string{"vnc", "rdp", "ssh", "telnet", "spice", "libvirt"}
	PROTOCOL_NAME  = "guacamole"
	BUF_LEN        = 4096
	INST_TERM_BYTE = byte(inst.INST_TERM[0])
)

type Client struct {
	Host           string
	Port           string
	Timeout        time.Duration
	InternalClient net.Conn
	BufReader      *bufio.Reader
	BufWriter      *bufio.Writer
	Connected      bool
	Debug          bool
	ReadBuffer     []byte
}

func ValidProtocol(protocol string) bool {
	for index := range PROTOCOLS {
		if protocol == PROTOCOLS[index] {
			return true
		}
	}
	return false
}

func NewClient(host, port string, timeout time.Duration, debug bool) *Client {
	c := &Client{
		Host:      host,
		Port:      port,
		Timeout:   timeout,
		Connected: false,
		Debug:     debug,
	}
	c.ReadBuffer = make([]byte, BUF_LEN)
	return c
}

func (this *Client) GetClient() net.Conn {
	if this.InternalClient == nil {
		address := net.JoinHostPort(this.Host, this.Port)
		conn, err := net.DialTimeout("tcp", address, this.Timeout)
		if err != nil {
			utils.Println(err)
			return nil
		}

		this.InternalClient = conn
		this.BufReader = bufio.NewReader(conn)
		this.BufWriter = bufio.NewWriter(conn)

		if this.Debug {
			utils.Printf("Client connected with gucad server (%s %s %s)\n",
				this.Host, this.Port, this.Timeout)
		}
		return this.InternalClient
	}

	return this.InternalClient
}

func (this *Client) Close() {
	err := this.GetClient().Close()
	if err != nil {
		utils.Println(err)
	}
	this.InternalClient = nil
	this.Connected = false
	if this.Debug {
		utils.Println("Connection closed.")
	}
}

func (this *Client) BufReceive() string {
	var line []byte
	var err error

	line, err = ReadBytes(this.BufReader, this.ReadBuffer, INST_TERM_BYTE)
	if err != nil {
		return string(line)
	}

	return string(line)
}

func (this *Client) Send(data []byte) error {
	if this.Debug {
		utils.Printf("Sending data: %s\n", string(data))
	}
	client := this.GetClient()
	if client == nil {
		return fmt.Errorf("maybe not connected\n")
	}
	n, err := client.Write(data)
	if err != nil {
		utils.Println("Send data failed:", err)
		this.Close()
		return err
	}
	if n != len(data) {
		utils.Println("Write() has blocked, send: %d, all: %d", n, len(data))
		this.Close()
		return fmt.Errorf("Write() has blocked, send: %d, all: %d", n, len(data))
	}

	return nil
}

func (this *Client) ReadInstruction() *inst.Instruction {
	if this.Debug {
		utils.Println("Read instruction")
	}
	instruction, err := inst.LoadInstruction(this.BufReceive())
	if err != nil {
		utils.Println("Read instruction failed:", err)
	} else {
		return instruction
	}
	return nil
}

func (this *Client) SendInstruction(instruction *inst.Instruction) error {
	if this.Debug {
		utils.Println("Sending instruction:", instruction.String())
	}
	return this.Send([]byte(instruction.Encode()))
}

func (this *Client) HandShake(protocol, width, height, dpi string, audio []string, video []string, KWArgs map[string]string) bool {
	var found_protocol bool
	for index := range PROTOCOLS {
		if protocol == PROTOCOLS[index] {
			found_protocol = true
			break
		}
	}
	if !found_protocol {
		utils.Printf("Invalid protocol: %s\n", protocol)
		return false
	}

	// 1. Send "select" instruction
	utils.Println("Send 'select' instruction")
	err := this.SendInstruction(inst.NewInstruction("select", protocol))
	if err != nil {
		utils.Println("Send 'select' failed:", err)
		return false
	}

	// 2. Receive "args" instruction
	instruction := this.ReadInstruction()
	utils.Printf("Expecting 'args' instruction, received: %s\n", instruction.String())

	if instruction == nil {
		utils.Println("Cannot establish Handshake. Connection Lost!")
		return false
	}

	if instruction.OpCode != "args" {
		utils.Printf("Cannot establish Handshake. Expected opcode 'args', received '%s' instead.\n", instruction.OpCode)
		return false
	}

	// 3. Respond with size, audio & video support
	utils.Printf("Send 'size' instruction (%s %s %s)\n", width, height, dpi)
	size_instruction := inst.NewInstruction("size", width, height, dpi)
	err = this.SendInstruction(size_instruction)
	if err != nil {
		utils.Println("Send 'size' failed:", err)
		return false
	}

	utils.Printf("Send 'audio' instruction (%s)\n", audio)
	err = this.SendInstruction(inst.NewInstruction("audio", audio...))
	if err != nil {
		utils.Println("Send 'audio' failed:", err)
		return false
	}

	utils.Printf("Send 'video' instruction (%s)\n", video)
	err = this.SendInstruction(inst.NewInstruction("video", video...))
	if err != nil {
		utils.Println("Send 'video' failed:", err)
		return false
	}

	var connect_args []string
	for idx := range instruction.Args {
		arg := strings.Replace(instruction.Args[idx], "-", "_", -1)
		if v, ok := KWArgs[arg]; ok {
			connect_args = append(connect_args, v)
		} else {
			connect_args = append(connect_args, "")
		}
	}

	// 4. Send connect arguments
	utils.Printf("Send 'connect' instruction (%s)\n", connect_args)
	err = this.SendInstruction(inst.NewInstruction("connect", connect_args...))
	if err != nil {
		utils.Println("Send 'connect' failed:", err)
		return false
	}

	return true
}

func ReadBytes(reader *bufio.Reader, buf []byte, delim byte) (line []byte, err error) {
	// Use ReadSlice to look for array,
	// accumulating full buffers.
	var frag []byte
	var full [][]byte
	err = nil

	for {
		var e error
		frag, e = reader.ReadSlice(delim)
		if e == nil { // got final fragment
			break
		}
		if e != bufio.ErrBufferFull { // unexpected error
			err = e
			break
		}

		// Make a copy of the buffer.
		buf := make([]byte, len(frag))
		copy(buf, frag)
		full = append(full, buf)
	}

	// Allocate new buffer to hold the full pieces and the fragment.
	n := 0
	for i := range full {
		n += len(full[i])
	}
	n += len(frag)

	// Copy full pieces and fragment in.
	if n > cap(buf) {
		buf = make([]byte, n)
	}

	n = 0
	for i := range full {
		n += copy(buf[n:], full[i])
	}
	copy(buf[n:], frag)
	return buf[:n+len(frag)], err
}

/*
func main() {
	SSH_ARGS := map[string]string{
		"username": "dummy",
		"password": "password",
		"hostname": "127.0.0.1",
		"port":     "22",
	}
	client := NewClient("172.31.31.110", "4822", 3*time.Second, true)
	client.HandShake("ssh", "1024", "768", "96", []string{}, []string{}, SSH_ARGS)
	time.Sleep(7 * time.Second)
}
*/
