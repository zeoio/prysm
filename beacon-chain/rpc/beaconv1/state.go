package beaconv1

import (
	"context"

	ptypes "github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (bs *Server) getState(ctx context.Context, stateId string) (*state.BeaconState, error) {
	switch stateId {
	case "head":
		headState, err := bs.HeadFetcher.HeadState(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "could not get head state")
		}
		return headState, nil
	case "genesis":
		genesisState, err := bs.StateGen.StateByRoot(ctx, params.BeaconConfig().ZeroHash)
		if err != nil {
			return nil, errors.Wrap(err, "could not get genesis checkpoint")
		}
		return genesisState, nil
	case "finalized":
		finalizedCheckpoint := bs.FinalizationFetcher.FinalizedCheckpt()
		finalizedState, err := bs.StateGen.StateByRoot(ctx, bytesutil.ToBytes32(finalizedCheckpoint.Root))
		if err != nil {
			return nil, errors.Wrap(err, "could not get finalized checkpoint")
		}
		return finalizedState, nil
	case "justified":
		justifiedCheckpoint := bs.FinalizationFetcher.CurrentJustifiedCheckpt()
		justifiedState, err := bs.StateGen.StateByRoot(ctx, bytesutil.ToBytes32(justifiedCheckpoint.Root))
		if err != nil {
			return nil, errors.Wrap(err, "could not get justified checkpoint")
		}
		return justifiedState, nil
	default:
		if len([]byte(stateId)) == 32 {
			requestedState, err := bs.StateGen.StateByRoot(ctx, bytesutil.ToBytes32([]byte(stateId)))
			if err != nil {
				return nil, errors.Wrap(err, "could not get state")
			}
			return requestedState, nil
		} else {
			requestedSlot := bytesutil.FromBytes8([]byte(stateId))
			requestedEpoch := helpers.SlotToEpoch(requestedSlot)
			currentEpoch := helpers.SlotToEpoch(bs.GenesisTimeFetcher.CurrentSlot())
			if requestedEpoch > currentEpoch {
				return nil, status.Errorf(
					codes.InvalidArgument,
					"Cannot retrieve information about an epoch in the future, current epoch %d, requesting %d",
					currentEpoch,
					requestedEpoch,
				)
			}

			requestedState, err := bs.StateGen.StateBySlot(ctx, requestedSlot)
			if err != nil {
				return nil, errors.Wrap(err, "could not get state")
			}
			return requestedState, nil
		}
	}
}
