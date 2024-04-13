package ui

type GetBroadcastParams struct {
	ID string `param:"id"`
}

type BroadcastOut struct {
	ID string `json:"id"`
}

type BroadcastsOut struct {
	Broadcasts []*BroadcastOut `json:"broadcasts"`
}
