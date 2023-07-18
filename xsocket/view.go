package xsocket

import (
	"encoding/json"
)

type CompatMsg struct {
	Event string      `json:"event,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type CompatMsgDataItem struct {
	Type string      `json:"type,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func (c CompatMsg) JSON() []byte {
	data, _ := json.Marshal(c)
	return data
}

func (c CompatMsgDataItem) JSON() []byte {
	data, _ := json.Marshal(c)
	return data
}
