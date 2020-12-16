package tx

import (
	"encoding/json"
	"github.com/ToolPackage/fse/common/utils"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestNewChannel(t *testing.T) {
	pipe12 := &DualStream{buf: make([]byte, 0)}
	pipe21 := &DualStream{buf: make([]byte, 0)}
	chan1 := NewChannel(pipe21, pipe12)
	chan2 := NewChannel(pipe12, pipe21)
	chan1.NewPacket("test").
		StatusCode(http.StatusAccepted).
		Header("name", "asd").
		Body([]string{"hello"}).
		Emit()
	packet := chan2.RecvPacket()
	assert.NotNil(t, packet)
	assert.Equal(t, "test", packet.Action)
	assert.Equal(t, http.StatusAccepted, packet.StatusCode)
	assert.Equal(t, 9, packet.ContentLength)
	assert.Equal(t, "asd", packet.Headers["name"])
	assert.Equal(t, "[\"hello\"]", string(packet.Content))

	data := make([]string, 0)
	assert.Nil(t, json.Unmarshal(packet.Content, &data))
	assert.Equal(t, 1, len(data))
	assert.Equal(t, "hello", data[0])
}

type DualStream struct {
	buf        []byte
	readOffset int
}

func (s *DualStream) Write(data []byte) (int, error) {
	s.buf = append(s.buf, data...)
	return len(data), nil
}

func (s *DualStream) Read(buf []byte) (int, error) {
	if s.readOffset >= len(s.buf) {
		return 0, io.EOF
	}
	n := utils.Min(len(buf), len(s.buf)-s.readOffset)
	for i := 0; i < n; i++ {
		buf[i] = s.buf[s.readOffset+i]
	}
	s.readOffset += n
	return n, nil
}
