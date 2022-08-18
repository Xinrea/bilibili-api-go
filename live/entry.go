package live

import (
	"time"
)

type EntryDanmaku struct{}

// From Not implemented yet. Extracts danmaku infos from message. Reference: https://www.wolai.com/nH1frikusirQFR79h4YmA
func (e *EntryDanmaku) From(_ *Message) {
}

type EntryGift struct{}

type EntryGuard struct {
	Sender struct {
		UID  uint64
		Name string
	}
	// 3：舰长 2：提督 1：总督
	Level     int
	Price     int
	StartTime time.Time
	EndTime   time.Time
}

func (e *EntryGuard) From(m *Message) {
	data := m.Data.(map[string]interface{})
	e.Sender.UID = uint64(data["uid"].(float64))
	e.Sender.Name = data["username"].(string)
	e.Level = int(data["guard_level"].(float64))
	e.Price = int(data["price"].(float64))
	e.StartTime = time.Unix(int64(data["start_time"].(float64)), 0)
	e.EndTime = time.Unix(int64(data["end_time"].(float64)), 0)
}
