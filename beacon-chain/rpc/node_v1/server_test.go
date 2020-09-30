package node_v1

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
	ptypes "github.com/gogo/protobuf/types"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1"
	"github.com/prysmaticlabs/prysm/beacon-chain/p2p"
	mockP2p "github.com/prysmaticlabs/prysm/beacon-chain/p2p/testing"
	"github.com/prysmaticlabs/prysm/shared"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
	"github.com/prysmaticlabs/prysm/shared/version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//func TestNodeServer_GetSyncStatus(t *testing.T) {
//	mSync := &mockSync.Sync{IsSyncing: false}
//	ns := &Server{
//		SyncChecker: mSync,
//	}
//	res, err := ns.GetSyncStatus(context.Background(), &ptypes.Empty{})
//	require.NoError(t, err)
//	assert.Equal(t, false, res.Syncing)
//	ns.SyncChecker = &mockSync.Sync{IsSyncing: true}
//	res, err = ns.GetSyncStatus(context.Background(), &ptypes.Empty{})
//	require.NoError(t, err)
//	assert.Equal(t, true, res.Syncing)
//}

//func TestNodeServer_GetGenesis(t *testing.T) {
//	db, _ := dbutil.SetupDB(t)
//	ctx := context.Background()
//	addr := common.Address{1, 2, 3}
//	require.NoError(t, db.SaveDepositContractAddress(ctx, addr))
//	st := testutil.NewBeaconState()
//	genValRoot := bytesutil.ToBytes32([]byte("I am root"))
//	ns := &Server{
//		BeaconDB:           db,
//		GenesisTimeFetcher: &mock.ChainService{},
//		GenesisFetcher: &mock.ChainService{
//			State:          st,
//			ValidatorsRoot: genValRoot,
//		},
//	}
//	res, err := ns.GetGenesis(context.Background(), &ptypes.Empty{})
//	require.NoError(t, err)
//	assert.DeepEqual(t, addr.Bytes(), res.DepositContractAddress)
//	pUnix, err := ptypes.TimestampProto(time.Unix(0, 0))
//	require.NoError(t, err)
//	assert.Equal(t, true, res.GenesisTime.Equal(pUnix))
//	assert.DeepEqual(t, genValRoot[:], res.GenesisValidatorsRoot)
//
//	ns.GenesisTimeFetcher = &mock.ChainService{Genesis: time.Unix(10, 0)}
//	res, err = ns.GetGenesis(context.Background(), &ptypes.Empty{})
//	require.NoError(t, err)
//	pUnix, err = ptypes.TimestampProto(time.Unix(10, 0))
//	require.NoError(t, err)
//	assert.Equal(t, true, res.GenesisTime.Equal(pUnix))
//}

func TestNodeServer_GetVersion(t *testing.T) {
	v := version.GetVersion()
	ns := &Server{}
	res, err := ns.GetVersion(context.Background(), &ptypes.Empty{})
	require.NoError(t, err)
	assert.Equal(t, v, res.Data.Version)
}

func TestNodeServer_GetIdentity(t *testing.T) {
	server := grpc.NewServer()
	peersProvider := &mockP2p.MockPeersProvider{}
	mP2P := mockP2p.NewTestP2P(t)
	key, err := crypto.GenerateKey()
	db, err := enode.OpenDB("")
	require.NoError(t, err)
	lNode := enode.NewLocalNode(db, key)
	record := lNode.Node().Record()
	stringENR, err := p2p.SerializeENR(record)
	require.NoError(t, err)
	ns := &Server{
		PeerManager:  &mockP2p.MockPeerManager{BHost: mP2P.BHost, Enr: record, PID: mP2P.BHost.ID()},
		PeersFetcher: peersProvider,
	}
	ethpb.RegisterBeaconNodeServer(server, ns)
	reflection.Register(server)
	h, err := ns.GetIdentity(context.Background(), &ptypes.Empty{})
	require.NoError(t, err)
	assert.Equal(t, mP2P.PeerID().String(), h.Data.PeerId)
	assert.Equal(t, stringENR, h.Data.Enr)
}

func TestNodeServer_GetPeer(t *testing.T) {
	server := grpc.NewServer()
	peersProvider := &mockP2p.MockPeersProvider{}
	ns := &Server{
		PeersFetcher: peersProvider,
	}
	ethpb.RegisterBeaconNodeServer(server, ns)
	reflection.Register(server)
	firstPeer := peersProvider.Peers().All()[0]

	res, err := ns.GetPeer(context.Background(), &ethpb.PeerRequest{PeerId: firstPeer.String()})
	require.NoError(t, err)
	assert.Equal(t, firstPeer.String(), res.Data.PeerId, "Unexpected peer ID")
	assert.Equal(t, int(ethpb.PeerDirection_INBOUND), int(res.Data.Direction), "Expected 1st peer to be an inbound connection")
	assert.Equal(t, ethpb.ConnectionState_CONNECTED, res.Data.State, "Expected peer to be connected")
}

func TestNodeServer_ListPeers(t *testing.T) {
	server := grpc.NewServer()
	peersProvider := &mockP2p.MockPeersProvider{}
	ns := &Server{
		PeersFetcher: peersProvider,
	}
	ethpb.RegisterBeaconNodeServer(server, ns)
	reflection.Register(server)

	res, err := ns.ListPeers(context.Background(), &ptypes.Empty{})
	require.NoError(t, err)
	assert.Equal(t, 2, len(res.Data))
	assert.Equal(t, int(ethpb.PeerDirection_INBOUND), int(res.Data[0].Direction))
	assert.Equal(t, ethpb.PeerDirection_OUTBOUND, res.Data[1].Direction)
}

type mockService struct {
	status error
}

func (m *mockService) Start() {
}

func (m *mockService) Stop() error {
	return nil
}

func (m *mockService) Status() error {
	return m.status
}

func TestGetHealth_Healthy(t *testing.T) {
	registry := shared.NewServiceRegistry()
	m := &mockService{}
	require.NoError(t, registry.RegisterService(m))
	ns := &Server{
		svcRegistry: registry,
	}
	ctx := context.Background()
	_, err := ns.GetHealth(ctx, &ptypes.Empty{})
	require.NoError(t, err)
}

func TestGetHealth_NotHealthy(t *testing.T) {
	registry := shared.NewServiceRegistry()
	m := &mockService{}
	m.status = errors.New("something is wrong")
	require.NoError(t, registry.RegisterService(m))
	ns := &Server{
		svcRegistry: registry,
	}
	ctx := context.Background()
	_, err := ns.GetHealth(ctx, &ptypes.Empty{})
	require.ErrorContains(t, m.status.Error(), err)
}
