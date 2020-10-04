package beaconv1

import (
	"context"
	"errors"

	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/state"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ptypes "github.com/gogo/protobuf/types"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1"
)

// GetGenesis retrieves details of the chain's genesis which can be used to identify chain.
func (bs *Server) GetGenesis(ctx context.Context, _ *ptypes.Empty) (*ethpb.GenesisResponse, error) {
	return nil, errors.New("unimplemented")
}

// GetStateRoot calculates HashTreeRoot for state with given 'stateId'. If stateId is root, same value will be returned.
func (bs *Server) GetStateRoot(ctx context.Context, req *ethpb.StateRequest) (*ethpb.StateRootResponse, error) {
	return nil, errors.New("unimplemented")
}

// GetStateFork returns Fork object for state with given 'stateId'.
func (bs *Server) GetStateFork(ctx context.Context, req *ethpb.StateRequest) (*ethpb.StateForkResponse, error) {
	return nil, errors.New("unimplemented")
}

// GetFinalityCheckpoints returns finality checkpoints for state with given 'stateId'. In case finality is
// not yet achieved, checkpoint should return epoch 0 and ZERO_HASH as root.
func (bs *Server) GetFinalityCheckpoints(ctx context.Context, req *ethpb.StateRequest) (*ethpb.StateFinalityCheckpointResponse, error) {
	return nil, errors.New("unimplemented")
}

State identifier. Can be one of: "head" (canonical head in node's view), "genesis", "finalized", "justified", <slot>, <hex encoded stateRoot with 0x prefix>.

Example : head

func (bs *Server) getState(stateId []byte) *state.BeaconState {
	switch string(stateId) {
	case "head":
	case "genesis":
	case "finalized":
	case "justified":
	default:
		if len(bytes) == 32 {

		} else {

		}
	}
	startSlot, err := helpers.StartSlot(requestedEpoch)
	if err != nil {
		return nil, err
	}
	requestedState, err := bs.StateGen.StateBySlot(ctx, startSlot)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get state")
	}
}
