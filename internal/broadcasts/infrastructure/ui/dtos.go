package ui

type GetBroadcastParams struct {
	ID string `param:"id"`
}

type BroadcastIn struct {
	Title string `json:"title" validate:"required"`
}

type BroadcastOut struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type CandidateIn struct {
	Candidate any `json:"candidate" validate:"required"`
}

type OfferIn struct {
	SDP string `json:"sdp" validate:"required"`
}

type BroadcastCreateOut struct {
	SDP string `json:"sdp"`
}

type BroadcastsOut struct {
	Broadcasts []*BroadcastOut `json:"broadcasts"`
}

type ViewerIn struct {
	BroadcastID string `json:"broadcast_id" validate:"required"`
	SDP         string `json:"sdp" validate:"required"`
}

type ViewerOut struct {
	SDP string `json:"sdp"`
}

type WsEvent struct {
	Data  any    `json:"data"`
	Event string `json:"event"`
}
