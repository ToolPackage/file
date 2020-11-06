package tx

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Packet struct {
	Action        string
	StatusCode    int
	ContentLength int
	Headers       map[string]string
	Content       []byte
}

type ActionHandler func(c *Channel, p *Packet)

var handlerMap = make(map[string]ActionHandler)

func Register(action string, handler ActionHandler) {
	_, ok := handlerMap[action]
	if ok {
		log.Fatal("duplicate action " + action)
	}
	handlerMap[action] = handler
}

type Channel struct {
	in  *bufio.Reader
	out *bufio.Writer
}

func NewChannel(in io.Reader, out io.Writer) *Channel {
	return &Channel{in: bufio.NewReader(in), out: bufio.NewWriter(out)}
}

// main loop
func (c *Channel) Process() {
	for {
		packet := c.RecvPacket()
		handler, ok := handlerMap[packet.Action]
		if ok {
			handler(c, packet)
		} else {
			c.ForceClose(errors.New("invalid action " + packet.Action))
		}
	}
}

func (c *Channel) ForceClose(err error) {
	panic(err)
}

// read next packet
func (c *Channel) RecvPacket() *Packet {
	action := c.readString(' ')
	statusCode := c.readInt(' ')
	contentLength := c.readInt('\n')
	headers := make(map[string]string)
	for {
		header := c.readLine()
		if len(header) == 0 {
			break
		}
		idx := strings.IndexRune(header, ':')
		if idx == -1 {
			c.ForceClose(errors.New("failed to parse packet headers"))
		}
		headers[strings.TrimSpace(header[:idx])] = strings.TrimSpace(header[idx+1:])
	}
	// read content
	buf := make([]byte, contentLength)
	n, err := c.in.Read(buf)
	if err != nil || n != contentLength {
		if n != contentLength {
			err = errors.New("expect to read " + strconv.Itoa(contentLength-n) + " more bytes of content")
		}
		c.ForceClose(err)
	}

	return &Packet{
		Action:        action,
		StatusCode:    statusCode,
		ContentLength: contentLength,
		Headers:       headers,
		Content:       buf,
	}
}

func (c *Channel) readInt(delim byte) int {
	str, err := c.in.ReadString(delim)
	if err != nil {
		c.ForceClose(err)
	}
	v, err := strconv.Atoi(str[:len(str)-1])

	return v
}

func (c *Channel) readString(delim byte) string {
	str, err := c.in.ReadString(delim)
	if err != nil {
		c.ForceClose(err)
	}
	return str[:len(str)-1]
}

func (c *Channel) readLine() string {
	line, err := c.in.ReadString('\n')
	if err != nil {
		c.ForceClose(err)
	}
	return line[:len(line)-1]
}

func (c *Channel) NewPacket(action string) *PacketBuilder {
	if !isValidActionNaming(action) {
		log.Fatal("invalid action naming: " + action)
	}
	return newPacketBuilder(c, action)
}

func isValidActionNaming(action string) bool {
	return strings.IndexRune(action, ' ') == -1
}

func (c *Channel) sendPacket(p *Packet) {
	var data strings.Builder
	data.WriteString(p.Action)
	data.WriteRune(' ')
	data.WriteString(strconv.Itoa(p.StatusCode))
	data.WriteRune(' ')
	data.WriteString(strconv.Itoa(p.ContentLength))
	data.WriteRune('\n')
	for key, value := range p.Headers {
		data.WriteString(key)
		data.WriteRune(':')
		data.WriteString(value)
		data.WriteRune('\n')
	}
	data.WriteRune('\n')
	c.sendString(data.String())
	// send content
	if p.ContentLength > 0 {
		_, err := c.out.Write(p.Content)
		if err != nil {
			c.ForceClose(err)
		}
	}
	_ = c.out.Flush()
}

func (c *Channel) sendString(s string) {
	_, err := c.out.WriteString(s)
	if err != nil {
		c.ForceClose(err)
	}
}

type PacketBuilder struct {
	channel *Channel
	packet  *Packet
}

func newPacketBuilder(c *Channel, action string) *PacketBuilder {
	packet := &Packet{
		Action:     action,
		StatusCode: http.StatusOK,
		Headers:    make(map[string]string),
	}
	return &PacketBuilder{channel: c, packet: packet}
}

func (pb *PacketBuilder) StatusCode(code int) *PacketBuilder {
	pb.packet.StatusCode = code
	return pb
}

func (pb *PacketBuilder) Body(content interface{}) *PacketBuilder {
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		log.Fatal(err)
	}
	pb.packet.Content = jsonBytes
	pb.packet.ContentLength = len(jsonBytes)
	return pb
}

func (pb *PacketBuilder) header(name string, value string) *PacketBuilder {
	name = strings.TrimSpace(name)
	value = strings.TrimSpace(value)
	pb.packet.Headers[name] = value
	return pb
}

func (pb *PacketBuilder) Emit() {
	pb.channel.sendPacket(pb.packet)
}
