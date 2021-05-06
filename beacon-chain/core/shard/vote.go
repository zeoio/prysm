package shard

import (
	"bytes"
	"context"

	"github.com/pkg/errors"
	types "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	state "github.com/prysmaticlabs/prysm/beacon-chain/state/stateV0"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/attestationutil"
)

// UpdatePendingVotes for the shard header.
//
// Spec code:
// def update_pending_votes(state: BeaconState, attestation: Attestation) -> None:
//    # Find and update the PendingShardHeader object, invalid block if pending header not in state
//    if compute_epoch_at_slot(attestation.data.slot) == get_current_epoch(state):
//        pending_headers = state.current_epoch_pending_shard_headers
//    else:
//        pending_headers = state.previous_epoch_pending_shard_headers
//    pending_header = None
//    for header in pending_headers:
//        if header.root == attestation.data.shard_header_root:
//            pending_header = header
//    assert pending_header is not None
//    assert pending_header.slot == attestation.data.slot
//    assert pending_header.shard == compute_shard_from_committee_index(
//        state,
//        attestation.data.slot,
//        attestation.data.index,
//    )
//    for i in range(len(pending_header.votes)):
//        pending_header.votes[i] = pending_header.votes[i] or attestation.aggregation_bits[i]
//
//    # Check if the PendingShardHeader is eligible for expedited confirmation
//    # Requirement 1: nothing else confirmed
//    all_candidates = [
//        c for c in pending_headers if
//        (c.slot, c.shard) == (pending_header.slot, pending_header.shard)
//    ]
//    if True in [c.confirmed for c in all_candidates]:
//        return
//
//    # Requirement 2: >= 2/3 of balance attesting
//    participants = get_attesting_indices(state, attestation.data, pending_header.votes)
//    participants_balance = get_total_balance(state, participants)
//    full_committee = get_beacon_committee(state, attestation.data.slot, attestation.data.index)
//    full_committee_balance = get_total_balance(state, set(full_committee))
//    if participants_balance * 3 >= full_committee_balance * 2:
//        pending_header.confirmed = True
func UpdatePendingVote(
	ctx context.Context,
	beaconState *state.BeaconState,
	att *ethpb.Attestation,
) (*state.BeaconState, error) {
	var pendingHeaders []*pb.PendingShardHeader
	var currentEpoch bool
	data := att.Data
	if helpers.SlotToEpoch(data.Slot) == helpers.CurrentEpoch(beaconState) {
		pendingHeaders = beaconState.CurrentEpochPendingShardHeaders()
		currentEpoch = true
	} else {
		pendingHeaders = beaconState.PreviousEpochPendingShardHeaders()
		currentEpoch = false
	}

	var pendingHeader *pb.PendingShardHeader
	for _, header := range pendingHeaders {
		if bytes.Equal(header.Root, data.ShardHeaderRoot) {
			pendingHeader = header
		}
	}

	if pendingHeader == nil {
		return nil, errors.New("unknown shard header root")
	}
	if pendingHeader.Slot != data.Slot {
		return nil, errors.New("incorrect shard header slot")
	}
	shard, err := ShardFromCommitteeIndex(beaconState, data.Slot, data.CommitteeIndex)
	if err != nil {
		return nil, err
	}
	if pendingHeader.Shard != shard {
		return nil, errors.New("incorrect shard")
	}
	pendingHeader.Votes.Or(att.AggregationBits)

	for _, header := range pendingHeaders {
		if header.Slot == data.Slot && header.Shard == shard {
			if header.Confirmed {
				return beaconState, nil
			}
		}
	}

	if !pendingHeader.Confirmed {
		committeeIndices, err := helpers.BeaconCommitteeFromState(beaconState, data.Slot, data.CommitteeIndex)
		if err != nil {
			return nil, err
		}
		indices, err := attestationutil.AttestingIndices(pendingHeader.Votes, committeeIndices)
		if err != nil {
			return nil, err
		}
		voitedIndices := make([]types.ValidatorIndex, len(indices))
		for i, index := range indices {
			voitedIndices[i] = types.ValidatorIndex(index)
		}
		committeeBalance := helpers.TotalBalance(beaconState, committeeIndices)
		if err != nil {
			return nil, err
		}

		votedBalance := helpers.TotalBalance(beaconState, voitedIndices)
		if err != nil {
			return nil, err
		}
		if votedBalance*3 > committeeBalance*2 {
			pendingHeader.Confirmed = true
		}
	}

	if currentEpoch {
		if err := beaconState.SetCurrentEpochPendingShardHeader(pendingHeaders); err != nil {
			return nil, err
		}
	} else {
		if err := beaconState.SetPreviousEpochPendingShardHeader(pendingHeaders); err != nil {
			return nil, err
		}
	}

	return beaconState, nil
}
