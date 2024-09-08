package ui

type GetBroadcastParams struct {
	ID string `param:"id"`
}

type BroadcastIn struct {
	SDP   string `json:"sdp" validate:"required"`
	Title string `json:"title" validate:"required"`
}

type BroadcastOut struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type BroadcastCreateOut struct {
	SDP string `json:"sdp"`
}

type BroadcastsOut struct {
	Broadcasts []*BroadcastOut `json:"broadcasts"`
}

type BroadcastConnectParams struct {
	BroadcastID string `param:"broadcastID" validate:"required"`
	SDP         string `json:"sdp" validate:"required"`
}

type BroadcastConnectOut struct {
	SDP string `json:"sdp"`
}

type WsEvent struct {
	Data  any    `json:"data"`
	Event string `json:"event"`
}
