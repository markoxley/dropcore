package dropcore

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	MSG_START = "DRIP"
	MSG_END   = "PIRD"
)

type Message struct {
	Source      string
	Destination string
	Data        map[string]string
}

func NewMessage(source string, destination string) *Message {
	return &Message{
		Source:      source,
		Destination: destination,
		Data:        make(map[string]string),
	}
}

func Parse(source string) (*Message, error) {
	if len(source) < len(MSG_START)+len(MSG_END)+5 {
		return nil, errors.New("message too short")
	}
	chksum := byte(source[len(source)-1])
	source = source[:len(source)-1]
	if !strings.HasPrefix(source, MSG_START) || !strings.HasSuffix(source, MSG_END) {
		return nil, errors.New("message incorrect format")
	}
	if createChecksum(source) != chksum {
		return nil, errors.New("invalid checksum")
	}
	source = source[len(MSG_START) : len(source)-len(MSG_END)]
	msgLength, err := strconv.ParseInt(source[:4], 16, 64)
	if err != nil {
		return nil, err
	}
	if int(msgLength) != len(source)-4 {
		return nil, errors.New("message length invalid")
	}
	source = source[4:]
	parts := strings.Split(source, ",")
	msg := NewMessage("", "")
	for i, p := range parts {
		switch i {
		case 0:
			msg.Source = p
		case 1:
			msg.Destination = p
		default:
			kv := strings.Split(p, "=")
			if len(kv) != 2 {
				return nil, errors.New("invalid data format")
			}
			msg.Data[kv[0]] = restoreFromSafe(kv[1])
		}
	}
	return msg, nil
}

func (m *Message) AddString(name string, value string) *Message {
	m.Data[name] = value
	return m
}

func (m *Message) AddInt(name string, value int) *Message {
	m.Data[name] = fmt.Sprintf("%v", value)
	return m
}

func (m *Message) AddFloat(name string, value float64) *Message {
	m.Data[name] = fmt.Sprintf("%f", value)
	return m
}

func (m *Message) AddBool(name string, value bool) *Message {
	if value {
		m.Data[name] = "1"
	} else {
		m.Data[name] = "0"
	}
	return m
}

func (m *Message) ToString() string {
	data := fmt.Sprintf("%s,%s", m.Source, m.Destination)
	for k, v := range m.Data {
		data += fmt.Sprintf(",%s=%s", makeSafe(k), makeSafe(v))
	}
	msg := fmt.Sprintf("%s%04x%s%s", MSG_START, len(data), data, MSG_END)
	chksum := createChecksum(msg)
	b := []byte(msg)
	b = append(b, chksum)
	msg = string(b)

	return msg
}

func makeSafe(msg string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(msg, "_", "_0"), ",", "_1"), "=", "_2")
}

func restoreFromSafe(msg string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(msg, "_2", "="), "_1", ","), "_0", "_")
}

func createChecksum(text string) byte {
	v := uint64(0)
	for _, c := range text {
		v += uint64(c)
	}
	return byte(v & 0xff)
}
