package live

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
)

const (
	WsOpHeartbeat               = 2  //心跳
	WsOpUserAuthentication      = 7  //用户进入房间
	WsPackageHeaderTotalLength  = 16 //头部字节大小
	WsBodyProtocolVersionJson   = 0  //JSON消息
	WsBodyProtocolVersionNormal = 1  //普通消息
	WsBodyProtocolVersionZlib   = 2  //压缩消息
)

type LoginRequest struct {
	UID      int    `json:"uid"`
	RoomID   int    `json:"roomid"`
	Protover int    `json:"protover"`
	Platform string `json:"platform"`
	Type     int    `json:"type"`
	Key      string `json:"key"`
}

type WebsocketHeader struct {
	TotalLength  uint32
	HeaderLength uint16
	Version      uint16
	Operation    int32
	Fixed        uint32
}

type Message struct {
	Cmd  string        `json:"cmd"`
	Data interface{}   `json:"data"`
	Info []interface{} `json:"info"`
}

func encodeMessage(content []byte, version uint16, operation int32) []byte {
	var totalLength = WsPackageHeaderTotalLength + uint32(len(content))
	msg := new(bytes.Buffer)
	header := WebsocketHeader{
		TotalLength:  totalLength,
		HeaderLength: WsPackageHeaderTotalLength,
		Version:      version,
		Operation:    operation,
		Fixed:        1,
	}
	_ = binary.Write(msg, binary.BigEndian, header)
	_ = binary.Write(msg, binary.BigEndian, content)
	return msg.Bytes()
}

func decodeMessage(raw []byte) []Message {
	var result []Message
	for len(raw) > 0 {
		var header WebsocketHeader
		buffer := bytes.NewBuffer(raw)
		_ = binary.Read(buffer, binary.BigEndian, &header)
		if header.Operation == 5 {
			result = append(result, processMessage(header, raw[header.HeaderLength:header.TotalLength])...)
		}
		raw = raw[header.TotalLength:]
	}
	return result
}

func processMessage(header WebsocketHeader, raw []byte) []Message {
	switch header.Version {
	case WsBodyProtocolVersionJson:
		var msg Message
		_ = json.Unmarshal(raw, &msg)
		return []Message{msg}
	case WsBodyProtocolVersionZlib:
		rawDecompressed := doZlibUnCompress(raw)
		return decodeMessage(rawDecompressed)
	default:
		log.Println("Unsupported protocol version: ", header.Version)
		return nil
	}
}

func doZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	_, err := io.Copy(&out, r)
	if err != nil {
		log.Println(err)
	}
	return out.Bytes()
}
