package main

import (
	"bytes"
	"fmt"

	"github.com/tidwall/resp"
)

const (
	CommandSET    = "SET"
	CommandGET    = "GET"
	CommandHELLO  = "HELLO"
	CommandClient = "CLIENT"
)

// Command represents a generic command interface
type Command interface{}

// SetCommand represents a SET command
type SetCommand struct {
	key, val []byte
}

// GetCommand represents a GET command
type GetCommand struct {
	key []byte
}

// HelloCommand represents a HELLO command
type HelloCommand struct {
	value string
}

// ClientCommand represents a CLIENT command
type ClientCommand struct {
	value string
}

// respWriteMap serializes a map into RESP format
func respWriteMap(m map[string]string) []byte {
	buf := &bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("%%%d\r\n", len(m)))
	rw := resp.NewWriter(buf)
	for k, v := range m {
		rw.WriteString(k)
		rw.WriteString(":" + v)
	}
	return buf.Bytes()
}
