package live

type APIResponse[T interface{}] struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type DataRoomInit struct {
	RoomId      uint64 `json:"room_id"`
	ShortId     uint64 `json:"short_id"`
	Uid         uint64 `json:"uid"`
	NeedP2P     int    `json:"need_p2p"`
	IsHidden    bool   `json:"is_hidden"`
	IsLocked    bool   `json:"is_locked"`
	IsPortrait  bool   `json:"is_portrait"`
	LiveStatus  int    `json:"live_status"`
	HiddenTill  int    `json:"hidden_till"`
	LockTill    int    `json:"lock_till"`
	Encrypted   bool   `json:"encrypted"`
	PwdVerified bool   `json:"pwd_verified"`
	LiveTime    int64  `json:"live_time"`
	RoomShield  int    `json:"room_shield"`
	IsSp        int    `json:"is_sp"`
	SpecialType int    `json:"special_type"`
}
