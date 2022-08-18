package live

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/recws-org/recws"
	"log"
	"net/http"
	"strings"
	"time"
)

type LiveRoom struct {
	id      uint64
	shortID uint64

	guardHandler   func(e EntryGuard)
	giftHandler    func(e EntryGift)
	danmakuHandler func(e EntryDanmaku)

	ctx    context.Context
	ws     *recws.RecConn
	cancel context.CancelFunc
}

const (
	GUARD_HANDLER = iota
	GIFT_HANDLER
	DANMAKU_HANDLER
)

// New creates live room with target room id(short id or normal id).
func New(id uint64) (*LiveRoom, error) {
	resp, err := new(LiveAPI).RoomInit(http.DefaultClient, id)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &LiveRoom{
		resp.RoomId,
		resp.ShortId,
		nil,
		nil,
		nil,
		ctx,
		nil,
		cancel,
	}, nil
}

func (l *LiveRoom) Register(handlerType int, handler interface{}) {
	switch handlerType {
	case GUARD_HANDLER:
		l.guardHandler = handler.(func(EntryGuard))
	case DANMAKU_HANDLER:
		l.danmakuHandler = handler.(func(EntryDanmaku))
	case GIFT_HANDLER:
		l.giftHandler = handler.(func(EntryGift))
	}
}

// Connect to live room with sessData. sessData can be empty but connection number might be limited by websocket server.
func (l *LiveRoom) Connect(sessData string) {
	// Check if already connected
	if l.ws != nil {
		return
	}
	l.ws = &recws.RecConn{KeepAliveTimeout: 0}
	requestHeader := http.Header{}
	if len(sessData) > 0 {
		requestHeader.Set("Cookie", "SESSDATA="+sessData)
	} else {
		requestHeader = nil
	}
	l.ws.Dial("wss://broadcastlv.chat.bilibili.com:443/sub", requestHeader)
	// Start processing messages and sending heartbeats.
	go l.heartbeat()
	go l.run()
}

func (l *LiveRoom) Disconnect() {
	l.cancel()
	l.ws.Close()
}

func (l *LiveRoom) join() {
	if l.ws == nil {
		return
	}
	loginJson, _ := json.Marshal(LoginRequest{
		UID:      0,
		RoomID:   int(l.id),
		Platform: "web",
		Protover: WsBodyProtocolVersionZlib,
		Type:     2,
		Key:      "",
	})
	err := l.ws.WriteMessage(websocket.BinaryMessage, encodeMessage(loginJson, WsBodyProtocolVersionNormal, WsOpUserAuthentication))
	if err != nil {
		log.Println("Join:", err)
	}
}

func (l *LiveRoom) run() {
	for {
		select {
		case <-l.ctx.Done():
			return
		default:
		}
		l.join()
		l.messageDispatcher()
	}
}

func (l *LiveRoom) messageDispatcher() {
	for {
		select {
		case <-l.ctx.Done():
			return
		default:
			_, rawBytes, err := l.ws.ReadMessage()
			if err != nil {
				time.Sleep(time.Second)
				return
			}
			// A message may contain several entries.
			entries := decodeMessage(rawBytes)
			for _, e := range entries {
				// Sometimes danmaku entry may have Cmd like "DANMU_MSG_0:11" for special danmakus.
				switch {
				case strings.HasPrefix(e.Cmd, "DANMU_MSG"):
					if l.danmakuHandler != nil {
						var danmakuEntry EntryDanmaku
						danmakuEntry.From(&e)
						l.danmakuHandler(danmakuEntry)
					}
				case e.Cmd == "GUARD_BUY" || e.Cmd == "USER_TOAST_MSG":
					if l.guardHandler != nil {
						var guardEntry EntryGuard
						guardEntry.From(&e)
						l.guardHandler(guardEntry)
					}
				default:
				}
			}

		}
	}
}

// Strange const content, seems to be a private key. Same as https://github.com/nodejs/help/issues/3000
const HeartbeatContent = "5b6f626a656374204f626a6563745d"

func (l *LiveRoom) heartbeat() {
	if l.ws == nil {
		return
	}
	content, _ := hex.DecodeString(HeartbeatContent)
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			_ = l.ws.WriteMessage(websocket.BinaryMessage, encodeMessage(content, WsBodyProtocolVersionNormal, WsOpHeartbeat))
		case <-l.ctx.Done():
			return
		}
	}
}
