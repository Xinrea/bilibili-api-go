package live

import (
	"net/http"
	"testing"
	"time"
)

func TestLiveAPI_RoomInit(t *testing.T) {
	api := new(LiveAPI)
	data, err := api.RoomInit(http.DefaultClient, 21484828)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", data)
}

func TestNew(t *testing.T) {
	room, err := New(704808)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(room.id)
	room.Register(GUARD_HANDLER, func(g EntryGuard) {
		t.Log(g.Sender.Name)
		t.Log(g.Price)
		t.Log(g.Level)
	})
	room.Connect("")
	time.Sleep(time.Second * 120)
	room.Disconnect()
}
