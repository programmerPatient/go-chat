package roomLogic

type RoomData struct {
	MessType int `json:"messType"`
	Data map[string]string `json:"data"`
}
