package c

import (
	a "github.com/prysmaticlabs/prysm/tools/type-migration/a"
	b "github.com/prysmaticlabs/prysm/tools/type-migration/b"
)


// AToBPeersResponse --
func AToBPeersResponse(src *a.PeersResponse) *b.PeersResponse {
	if src == nil {
		return &b.PeersResponse{}
	}
	return &b.PeersResponse{
	
		Peers: src.Peers,
	}
}

// BToAPeersResponse --
func BToAPeersResponse(src *b.PeersResponse) *a.PeersResponse {
	if src == nil {
		return &a.PeersResponse{}
	}
	return &a.PeersResponse{
	
		Peers: src.Peers,
	}
}

