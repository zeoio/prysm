package a

type PeersResponse struct {
	Peers []*Peer
}

type Peer struct {
	Inbound bool
}
