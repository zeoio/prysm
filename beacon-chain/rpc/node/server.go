// Package node defines a gRPC node service implementation, providing
// useful endpoints for checking a node's sync status, peer info,
// genesis data, and version information.
package node

import (
	"context"
	"errors"
	"fmt"

	ptypes "github.com/gogo/protobuf/types"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1"
	"github.com/prysmaticlabs/prysm/beacon-chain/db"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p"
	"github.com/prysmaticlabs/prysm/beacon-chain/sync"
	"github.com/prysmaticlabs/prysm/shared/version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server defines a server implementation of the gRPC Node service,
// providing RPC endpoints for verifying a beacon node's sync status, genesis and
// version information, and services the node implements and runs.
type Server struct {
	SyncChecker  sync.Checker
	Server       *grpc.Server
	BeaconDB     db.ReadOnlyDatabase
	PeersFetcher p2p.PeersProvider
	PeerManager  p2p.PeerManager
}

// GetIdentity return data about the node's network presence.
func (ns *Server) GetIdentity(ctx context.Context, _ *ptypes.Empty) (*ethpb.IdentityResponse, error) {
	hostAddrs := ns.PeerManager.Host().Addrs()
	stringAddrs := make([]string, len(hostAddrs))
	for i, addr := range hostAddrs {
		stringAddrs[i] = addr.String()
	}
	record := ns.PeerManager.ENR()
	enr := ""
	err := error(nil)
	if record != nil {
		enr, err = p2p.SerializeENR(record)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Unable to serialize enr: %v", err)
		}
	}

	peerID := ns.PeerManager.PeerID()
	metadata, err := ns.PeersFetcher.Peers().Metadata(peerID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unable to get metadata: %v", err)
	}
	identity := &ethpb.Identity{
		PeerId:       peerID.String(),
		Enr:          enr,
		P2PAddresses: stringAddrs,
		Metadata: &ethpb.Metadata{
			Attnets:   metadata.Attnets,
			SeqNumber: metadata.SeqNumber,
		},
	}
	return &ethpb.IdentityResponse{
		Data: identity,
	}, nil
}

// GetPeer returns the data known about the peer defined by the provided peer id.
func (ns *Server) GetPeer(ctx context.Context, peerReq *ethpb.PeerRequest) (*ethpb.PeerResponse, error) {
	pid, err := peer.Decode(peerReq.PeerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Unable to parse provided peer id: %v", err)
	}
	addr, err := ns.PeersFetcher.Peers().Address(pid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Requested peer does not exist: %v", err)
	}
	dir, err := ns.PeersFetcher.Peers().Direction(pid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Requested peer does not exist: %v", err)
	}
	pbDirection := ethpb.PeerDirection_UNKNOWN
	switch dir {
	case network.DirInbound:
		pbDirection = ethpb.PeerDirection_INBOUND
	case network.DirOutbound:
		pbDirection = ethpb.PeerDirection_OUTBOUND
	}
	connState, err := ns.PeersFetcher.Peers().ConnectionState(pid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Requested peer does not exist: %v", err)
	}
	record, err := ns.PeersFetcher.Peers().ENR(pid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Requested peer does not exist: %v", err)
	}
	enr := ""
	if record != nil {
		enr, err = p2p.SerializeENR(record)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Unable to serialize enr: %v", err)
		}
	}
	peerInfo := &ethpb.Peer{
		PeerId:    peerReq.PeerId,
		Enr:       enr,
		Address:   addr.String(),
		State:     ethpb.ConnectionState(connState),
		Direction: pbDirection,
	}
	return &ethpb.PeerResponse{
		Data: peerInfo,
	}, nil
}

// ListPeers lists the info of all peers connected to this node.
func (ns *Server) ListPeers(ctx context.Context, _ *ptypes.Empty) (*ethpb.PeersResponse, error) {
	peers := ns.PeersFetcher.Peers().Connected()
	res := make([]*ethpb.Peer, 0, len(peers))
	for _, pid := range peers {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		multiaddr, err := ns.PeersFetcher.Peers().Address(pid)
		if err != nil {
			continue
		}
		direction, err := ns.PeersFetcher.Peers().Direction(pid)
		if err != nil {
			continue
		}
		record, err := ns.PeersFetcher.Peers().ENR(pid)
		if err != nil {
			continue
		}
		enr := ""
		if record != nil {
			enr, err = p2p.SerializeENR(record)
			if err != nil {
				continue
			}
		}

		address := fmt.Sprintf("%s/p2p/%s", multiaddr.String(), pid.Pretty())
		pbDirection := ethpb.PeerDirection_UNKNOWN
		switch direction {
		case network.DirInbound:
			pbDirection = ethpb.PeerDirection_INBOUND
		case network.DirOutbound:
			pbDirection = ethpb.PeerDirection_OUTBOUND
		}
		res = append(res, &ethpb.Peer{
			PeerId:    pid.String(),
			Enr:       enr,
			Address:   address,
			State:     ethpb.ConnectionState_CONNECTED,
			Direction: pbDirection,
		})
	}

	return &ethpb.PeersResponse{
		Data: res,
	}, nil
}

// GetVersion checks the version information of the beacon node.
func (ns *Server) GetVersion(ctx context.Context, _ *ptypes.Empty) (*ethpb.VersionResponse, error) {
	versionInfo := &ethpb.Version{
		Version: version.GetVersion(),
	}
	return &ethpb.VersionResponse{
		Data: versionInfo,
	}, nil
}

// GetSyncStatus requests the beacon node to describe if it's currently syncing or not, and
// if it is, what block it is up to.
func (ns *Server) GetSyncStatus(ctx context.Context, _ *ptypes.Empty) (*ethpb.SyncingResponse, error) {
	return nil, errors.New("unimplemented")
}

// GetHealth returns node health status in http status codes. Useful for load balancers.
// Response Usage:
//    "200":
//      description: Node is ready
//    "206":
//      description: Node is syncing but can serve incomplete data
//    "503":
//      description: Node not initialized or having issues
func (ns *Server) GetHealth(ctx context.Context, _ *ptypes.Empty) (*ptypes.Empty, error) {
	return nil, errors.New("unimplemented")
}
