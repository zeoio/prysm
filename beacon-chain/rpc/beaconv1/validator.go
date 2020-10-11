package beaconv1

import (
	"context"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetValidator returns a validator specified by state and id or public key along with status and balance.
func (bs *Server) GetValidator(ctx context.Context, req *ethpb.StateValidatorRequest) (*ethpb.StateValidatorResponse, error) {
	if bs.GenesisTimeFetcher == nil {
		return nil, status.Errorf(codes.Internal, "Nil genesis time fetcher")
	}

	requestedState, err := bs.getState(ctx, req.StateId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get state")
	}

	validators := requestedState.Validators()
	valCount := len(validators)
	id := req.ValidatorId
	index, err := indexFromValidatorId(requestedState, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get index from validator id: %s", id)
	}
	if index >= uint64(valCount) {
		return nil, status.Errorf(codes.OutOfRange, "Validator index %d >= balance list %d", index, valCount)
	}
	val, err := requestedState.ValidatorAtIndex(index)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get validator at index %d: %v", index, err)
	}

	validator := &ethpb.Validator{
		PublicKey:                  val.PublicKey,
		WithdrawalCredentials:      val.WithdrawalCredentials,
		EffectiveBalance:           val.EffectiveBalance,
		Slashed:                    val.Slashed,
		ActivationEligibilityEpoch: val.ActivationEligibilityEpoch,
		ActivationEpoch:            val.ActivationEpoch,
		ExitEpoch:                  val.ExitEpoch,
		WithdrawableEpoch:          val.WithdrawableEpoch,
	}
	validatorContainer := &ethpb.ValidatorContainer{
		Index:     index,
		Balance:   0,
		Status:    "",
		Validator: validator,
	}

	return &ethpb.StateValidatorResponse{
		Data: validatorContainer,
	}, nil
}

// ListValidators returns filterable list of validators with their balance, status and index.
func (bs *Server) ListValidators(ctx context.Context, req *ethpb.StateValidatorsRequest) (*ethpb.StateValidatorsResponse, error) {
	if bs.GenesisTimeFetcher == nil {
		return nil, status.Errorf(codes.Internal, "Nil genesis time fetcher")
	}

	validatorsResp := make([]*ethpb.ValidatorContainer, 0)

	requestedState, err := bs.getState(ctx, req.StateId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get state")
	}

	validators := requestedState.Validators()
	valCount := len(validators)
	for _, id := range req.Id {
		index, err := indexFromValidatorId(requestedState, id)
		if err != nil {
			log.Errorf("Could not get index from validator id: %v", err)
			continue
		}
		if index >= uint64(valCount) {
			return nil, status.Errorf(codes.OutOfRange, "Validator index %d >= balance list %d", index, valCount)
		}
		val, err := requestedState.ValidatorAtIndex(index)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Could not get validator at index %d: %v", index, err)
		}

		validator := &ethpb.Validator{
			PublicKey:                  val.PublicKey,
			WithdrawalCredentials:      val.WithdrawalCredentials,
			EffectiveBalance:           val.EffectiveBalance,
			Slashed:                    val.Slashed,
			ActivationEligibilityEpoch: val.ActivationEligibilityEpoch,
			ActivationEpoch:            val.ActivationEpoch,
			ExitEpoch:                  val.ExitEpoch,
			WithdrawableEpoch:          val.WithdrawableEpoch,
		}
		validatorsResp = append(validatorsResp, &ethpb.ValidatorContainer{
			Index:     index,
			Balance:   0,
			Status:    "",
			Validator: validator,
		})
	}

	// Depending on the indices and public keys given, results might not be sorted.
	sort.Slice(validatorsResp, func(i, j int) bool {
		return validatorsResp[i].Index < validatorsResp[j].Index
	})

	// If there are no balances, we simply return a response specifying this.
	// Otherwise, attempting to paginate 0 balances below would result in an error.
	if valCount == 0 {
		return &ethpb.StateValidatorsResponse{
			Data: make([]*ethpb.ValidatorContainer, 0),
		}, nil
	}

	return &ethpb.StateValidatorsResponse{
		Data: validatorsResp,
	}, nil
}

// ListValidatorBalances returns a filterable list of validator balances.
func (bs *Server) ListValidatorBalances(ctx context.Context, req *ethpb.StateValidatorsRequest) (*ethpb.ValidatorBalancesResponse, error) {
	if bs.GenesisTimeFetcher == nil {
		return nil, status.Errorf(codes.Internal, "Nil genesis time fetcher")
	}

	balancesResp := make([]*ethpb.ValidatorBalance, 0)

	requestedState, err := bs.getState(ctx, req.StateId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get state")
	}

	balances := requestedState.Balances()
	balancesCount := len(balances)
	for _, id := range req.Id {
		index, err := indexFromValidatorId(requestedState, id)
		if err != nil {
			log.Errorf("Could not get index from validator id: %v", err)
			continue
		}
		if index >= uint64(len(balances)) {
			return nil, status.Errorf(codes.OutOfRange, "Validator index %d >= balance list %d", index, len(balances))
		}

		balancesResp = append(balancesResp, &ethpb.ValidatorBalance{
			Index:   index,
			Balance: balances[index],
		})
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
	currentSlot := bs.GenesisTimeFetcher.CurrentSlot()

	requestedEpoch := helpers.SlotToEpoch(req.Slot)
	currentEpoch := helpers.SlotToEpoch(currentSlot)
	if requestedEpoch > currentEpoch {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Cannot retrieve information for an future epoch, current epoch %d, requesting %d",
			currentEpoch,
			requestedEpoch,
		)
	}

	committees, _, err := bs.retrieveCommitteesForEpoch(ctx, requestedEpoch)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Could not retrieve committees for epoch %d: %v",
			requestedEpoch,
			err,
		)
	}

	return &ethpb.StateCommitteesResponse{
		Data: []*ethpb.Committee{
			{},
		},
	}, nil
}

func (bs *Server) retrieveCommitteesForEpoch(
	ctx context.Context,
	epoch uint64,
) (map[uint64][]*ethpb.Committee, []uint64, error) {
	startSlot, err := helpers.StartSlot(epoch)
	if err != nil {
		return nil, nil, err
	}
	requestedState, err := bs.StateGen.StateBySlot(ctx, startSlot)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, "Could not get state")
	}
	seed, err := helpers.Seed(requestedState, epoch, params.BeaconConfig().DomainBeaconAttester)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, "Could not get seed")
	}
	activeIndices, err := helpers.ActiveValidatorIndices(requestedState, epoch)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, "Could not get active indices")
	}

	committeesListsBySlot, err := computeCommittees(startSlot, activeIndices, seed)
	if err != nil {
		return nil, nil, status.Errorf(
			codes.InvalidArgument,
			"Could not compute committees for epoch %d: %v",
			helpers.SlotToEpoch(startSlot),
			err,
		)
	}
	return committeesListsBySlot, activeIndices, nil
}

// retrieveCommitteesForRoot uses the provided state root to get the current epoch committees.
// Note: This function is always recommended over retrieveCommitteesForEpoch as states are
// retrieved from the DB for this function, rather than generated.
func (bs *Server) retrieveCommitteesForRoot(
	ctx context.Context,
	root []byte,
) (map[uint64][]*ethpb.Committee, []uint64, error) {
	requestedState, err := bs.StateGen.StateByRoot(ctx, bytesutil.ToBytes32(root))
	if err != nil {
		return nil, nil, status.Error(codes.Internal, fmt.Sprintf("Could not get state: %v", err))
	}
	epoch := helpers.CurrentEpoch(requestedState)
	seed, err := helpers.Seed(requestedState, epoch, params.BeaconConfig().DomainBeaconAttester)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, "Could not get seed")
	}
	activeIndices, err := helpers.ActiveValidatorIndices(requestedState, epoch)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, "Could not get active indices")
	}

	startSlot, err := helpers.StartSlot(epoch)
	if err != nil {
		return nil, nil, err
	}
	committeesListsBySlot, err := computeCommittees(startSlot, activeIndices, seed)
	if err != nil {
		return nil, nil, status.Errorf(
			codes.InvalidArgument,
			"Could not compute committees for epoch %d: %v",
			epoch,
			err,
		)
	}
	return committeesListsBySlot, activeIndices, nil
}

// Compute committees given a start slot, active validator indices, and
// the attester seeds value.
func computeCommittees(
	startSlot uint64,
	activeIndices []uint64,
	attesterSeed [32]byte,
) (map[uint64][]*ethpb.Committee, error) {
	committeesListsBySlot := make(map[uint64][]*ethpb.Committee, params.BeaconConfig().SlotsPerEpoch)
	for slot := startSlot; slot < startSlot+params.BeaconConfig().SlotsPerEpoch; slot++ {
		var countAtSlot = uint64(len(activeIndices)) / params.BeaconConfig().SlotsPerEpoch / params.BeaconConfig().TargetCommitteeSize
		if countAtSlot > params.BeaconConfig().MaxCommitteesPerSlot {
			countAtSlot = params.BeaconConfig().MaxCommitteesPerSlot
		}
		if countAtSlot == 0 {
			countAtSlot = 1
		}
		committeeItems := make([]*ethpb.Committee, countAtSlot)
		for committeeIndex := uint64(0); committeeIndex < countAtSlot; committeeIndex++ {
			committee, err := helpers.BeaconCommittee(activeIndices, attesterSeed, slot, committeeIndex)
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					"Could not compute committee for slot %d: %v",
					slot,
					err,
				)
			}
			committeeItems[committeeIndex] = &ethpb.Committee{
				Index:      committeeIndex,
				Slot:       slot,
				Validators: committee,
			}
		}
		committeesListsBySlot[slot] = committeeItems
	}
	return committeesListsBySlot, nil
}

func indexFromValidatorId(state *state.BeaconState, valId string) (uint64, error) {
	// Skip empty public key.
	if valId == "" {
		return 0, errors.New("empty input")
	}
	var index uint64
	var ok bool
	// Skip empty public key.
	if strings.HasPrefix(valId, "0x") {
		pubkeyBytes, err := hex.DecodeString(valId)
		if err != nil {
			return 0, errors.Wrap(err, "could not decode string")
		}
		index, ok = state.ValidatorIndexByPubkey(bytesutil.ToBytes48(pubkeyBytes))
		if !ok {
			return 0, fmt.Errorf("validator public key %#x not found in state", pubkeyBytes)
		}
	} else {
		idx, err := strconv.Atoi(valId)
		if err != nil {
			return 0, errors.Wrap(err, "could not convert to number")
		}
		index = uint64(idx)
	}
	return index, nil
}
