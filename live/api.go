package live

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const roomInitUrl = "https://api.live.bilibili.com/room/v1/Room/room_init"

type LiveAPI struct{}

func (*LiveAPI) RoomInit(c *http.Client, roomID uint64) (*DataRoomInit, error) {
	requestUrl, err := url.Parse(roomInitUrl)
	if err != nil {
		panic(err)
	}
	queries := requestUrl.Query()
	queries.Set("id", strconv.Itoa(int(roomID)))
	requestUrl.RawQuery = queries.Encode()
	resp, err := c.Get(requestUrl.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrUrlIncorrect
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response APIResponse[DataRoomInit]
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(requestUrl.String(), string(body))
		return nil, err
	}
	return &response.Data, nil
}
