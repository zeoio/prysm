package migration

import (
	v1Alpha1 "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	v2 "github.com/prysmaticlabs/prysm/proto/prysm/v2"
)

// V1Alpha1ToV2Genesis --
func V1Alpha1ToV2Genesis(src *v1Alpha1.Genesis) *v2.Genesis {
	if src == nil {
		return &v2.Genesis{}
	}
	return &v2.Genesis{

		GenesisTime:            src.GenesisTime,
		DepositContractAddress: src.DepositContractAddress,
		GenesisValidatorsRoot:  src.GenesisValidatorsRoot,
	}
}

// V2ToV1Alpha1Genesis --
func V2ToV1Alpha1Genesis(src *v2.Genesis) *v1Alpha1.Genesis {
	if src == nil {
		return &v1Alpha1.Genesis{}
	}
	return &v1Alpha1.Genesis{

		GenesisTime:            src.GenesisTime,
		DepositContractAddress: src.DepositContractAddress,
		GenesisValidatorsRoot:  src.GenesisValidatorsRoot,
	}
}

// V1Alpha1ToV2HostData --
func V1Alpha1ToV2HostData(src *v1Alpha1.HostData) *v2.HostData {
	if src == nil {
		return &v2.HostData{}
	}
	return &v2.HostData{

		Addresses: src.Addresses,
		PeerId:    src.PeerId,
		Enr:       src.Enr,
	}
}

// V2ToV1Alpha1HostData --
func V2ToV1Alpha1HostData(src *v2.HostData) *v1Alpha1.HostData {
	if src == nil {
		return &v1Alpha1.HostData{}
	}
	return &v1Alpha1.HostData{

		Addresses: src.Addresses,
		PeerId:    src.PeerId,
		Enr:       src.Enr,
	}
}

// V1Alpha1ToV2Peer --
func V1Alpha1ToV2Peer(src *v1Alpha1.Peer) *v2.Peer {
	if src == nil {
		return &v2.Peer{}
	}
	return &v2.Peer{

		Address:         src.Address,
		Direction:       v2.PeerDirection(src.Direction),
		ConnectionState: v2.ConnectionState(src.ConnectionState),
		PeerId:          src.PeerId,
		Enr:             src.Enr,
	}
}

// V2ToV1Alpha1Peer --
func V2ToV1Alpha1Peer(src *v2.Peer) *v1Alpha1.Peer {
	if src == nil {
		return &v1Alpha1.Peer{}
	}
	return &v1Alpha1.Peer{

		Address:         src.Address,
		Direction:       v1Alpha1.PeerDirection(src.Direction),
		ConnectionState: v1Alpha1.ConnectionState(src.ConnectionState),
		PeerId:          src.PeerId,
		Enr:             src.Enr,
	}
}

// V1Alpha1ToV2Peers --
func V1Alpha1ToV2Peers(src *v1Alpha1.Peers) *v2.Peers {
	if src == nil {
		return &v2.Peers{}
	}
	peers := make([]*v2.Peer, 0, len(src.Peers))
	for _, p := range src.Peers {
		peers = append(peers, V1Alpha1ToV2Peer(p))
	}
	return &v2.Peers{
		Peers: peers,
	}
}

// V2ToV1Alpha1Peers --
func V2ToV1Alpha1Peers(src *v2.Peers) *v1Alpha1.Peers {
	if src == nil {
		return &v1Alpha1.Peers{}
	}
	peers := make([]*v1Alpha1.Peer, 0, len(src.Peers))
	for _, p := range src.Peers {
		peers = append(peers, V2ToV1Alpha1Peer(p))
	}
	return &v1Alpha1.Peers{
		Peers: peers,
	}
}

// V1Alpha1ToV2SyncStatus --
func V1Alpha1ToV2SyncStatus(src *v1Alpha1.SyncStatus) *v2.SyncStatus {
	if src == nil {
		return &v2.SyncStatus{}
	}
	return &v2.SyncStatus{

		Syncing: src.Syncing,
	}
}

// V2ToV1Alpha1SyncStatus --
func V2ToV1Alpha1SyncStatus(src *v2.SyncStatus) *v1Alpha1.SyncStatus {
	if src == nil {
		return &v1Alpha1.SyncStatus{}
	}
	return &v1Alpha1.SyncStatus{

		Syncing: src.Syncing,
	}
}

// V1Alpha1ToV2Version --
func V1Alpha1ToV2Version(src *v1Alpha1.Version) *v2.Version {
	if src == nil {
		return &v2.Version{}
	}
	return &v2.Version{

		Version:  src.Version,
		Metadata: src.Metadata,
	}
}

// V2ToV1Alpha1Version --
func V2ToV1Alpha1Version(src *v2.Version) *v1Alpha1.Version {
	if src == nil {
		return &v1Alpha1.Version{}
	}
	return &v1Alpha1.Version{

		Version:  src.Version,
		Metadata: src.Metadata,
	}
}
