package ui

type GetBroadcastParams struct {
	ID string `param:"id"`
}

type BroadcastIn struct {
	SDP string `json:"sdp" validate:"required"`
}

type BroadcastOut struct {
	ID string `json:"id"`
}

type BroadcastsOut struct {
	Broadcasts []*BroadcastOut `json:"broadcasts"`
}
