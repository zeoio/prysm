package beaconv1

import (
	"context"
	"encoding/hex"
	"sort"

	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetValidator returns a validator specified by state and id or public key along with status and balance.
func (bs *Server) GetValidator(ctx context.Context, req *ethpb.StateValidatorRequest) (*ethpb.StateValidatorResponse, error) {
	return nil, errors.New("unimplemented")
}

// ListValidators returns filterable list of validators with their balance, status and index.
func (bs *Server) ListValidators(ctx context.Context, req *ethpb.StateValidatorsRequest) (*ethpb.StateValidatorsResponse, error) {
	return nil, errors.New("unimplemented")
}

// ListValidatorBalances returns a filterable list of validator balances.
func (bs *Server) ListValidatorBalances(ctx context.Context, req *ethpb.StateValidatorsRequest) (*ethpb.ValidatorBalancesResponse, error) {
	if bs.GenesisTimeFetcher == nil {
		return nil, status.Errorf(codes.Internal, "Nil genesis time fetcher")
	}

	balancesResp := make([]*ethpb.ValidatorBalance, 0)
	filtered := map[uint64]bool{} // Track filtered validators to prevent duplication in the response.

	requestedState, err := bs.getState(ctx, req.StateId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get state")
	}

	balances := requestedState.Balances()
	balancesCount := len(balances)
	for _, id := range req.Id {
		// Skip empty public key.
		if len([]byte(id)) == 0 {
			continue
		}
		var index uint64
		var ok bool
		// Skip empty public key.
		if len([]byte(id)) == 48 {
			pubkeyBytes, err := hex.DecodeString(id)
			if err != nil {
				return nil, errors.Wrap(err, "could not decode string")
			}
			index, ok = requestedState.ValidatorIndexByPubkey(bytesutil.ToBytes48(pubkeyBytes))
			if !ok {
				continue
			}
		} else {
			index = bytesutil.FromBytes8([]byte(id))
		}

		if index >= uint64(len(balances)) {
			return nil, status.Errorf(codes.OutOfRange, "Validator index %d >= balance list %d",
				index, len(balances))
		}

		if !filtered[index] {
			balancesResp = append(balancesResp, &ethpb.ValidatorBalance{
				Index:   index,
				Balance: balances[index],
			})
		}
		filtered[index] = true
	}

	// Depending on the indices and public keys given, results might not be sorted.
	sort.Slice(balancesResp, func(i, j int) bool {
		return balancesResp[i].Index < balancesResp[j].Index
	})

	// If there are no balances, we simply return a response specifying this.
	// Otherwise, attempting to paginate 0 balances below would result in an error.
	if balancesCount == 0 {
		return &ethpb.ValidatorBalancesResponse{
			Data: make([]*ethpb.ValidatorBalance, 0),
		}, nil
	}

	return &ethpb.ValidatorBalancesResponse{
		Data: balancesResp,
	}, nil
}

// ListCommittees retrieves the committees for the given state at the given epoch.
func (bs *Server) ListCommittees(ctx context.Context, req *ethpb.StateCommitteesRequest) (*ethpb.StateCommitteesResponse, error) {
	return nil, errors.New("unimplemented")
}
