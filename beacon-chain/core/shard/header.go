package shard

import (
	"bytes"
	"context"
	"fmt"

	"github.com/pkg/errors"
	types "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	state "github.com/prysmaticlabs/prysm/beacon-chain/state/stateV0"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/params"
)

// ProcessShardHeader of the beacon chain.
//
// Spec code:
// def process_shard_header(state: BeaconState,
//                         signed_header: SignedShardBlobHeader) -> None:
//    header = signed_header.message
//    # Verify the header is not 0, and not from the future.
//    assert Slot(0) < header.slot <= state.slot
//    header_epoch = compute_epoch_at_slot(header.slot)
//    # Verify that the header is within the processing time window
//    assert header_epoch in [get_previous_epoch(state), get_current_epoch(state)]
//    # Verify that the shard is active
//    assert header.shard < get_active_shard_count(state, header_epoch)
//    # Verify that the block root matches,
//    # to ensure the header will only be included in this specific Beacon Chain sub-tree.
//    assert header.body_summary.beacon_block_root == get_block_root_at_slot(state, header.slot - 1)
//    # Verify proposer
//    assert header.proposer_index == get_shard_proposer_index(state, header.slot, header.shard)
//    # Verify signature
//    signing_root = compute_signing_root(header, get_domain(state, DOMAIN_SHARD_PROPOSER))
//    assert bls.Verify(state.validators[header.proposer_index].pubkey, signing_root, signed_header.signature)
//
//    # Verify the length by verifying the degree.
//    body_summary = header.body_summary
//    if body_summary.commitment.length == 0:
//        assert body_summary.degree_proof == G1_SETUP[0]
//    assert (
//        bls.Pairing(body_summary.degree_proof, G2_SETUP[0])
//        == bls.Pairing(body_summary.commitment.point, G2_SETUP[-body_summary.commitment.length])
//    )
//
//    # Get the correct pending header list
//    if header_epoch == get_current_epoch(state):
//        pending_headers = state.current_epoch_pending_shard_headers
//    else:
//        pending_headers = state.previous_epoch_pending_shard_headers
//
//    header_root = hash_tree_root(header)
//    # Check that this header is not yet in the pending list
//    assert header_root not in [pending_header.root for pending_header in pending_headers]
//
//    # Include it in the pending list
//    index = compute_committee_index_from_shard(state, header.slot, header.shard)
//    committee_length = len(get_beacon_committee(state, header.slot, index))
//    pending_headers.append(PendingShardHeader(
//        slot=header.slot,
//        shard=header.shard,
//        commitment=body_summary.commitment,
//        root=header_root,
//        votes=Bitlist[MAX_VALIDATORS_PER_COMMITTEE]([0] * committee_length),
//        confirmed=False,
//    ))
func ProcessShardHeader(
	ctx context.Context,
	beaconState *state.BeaconState,
	header *ethpb.SignedShardBlobHeader,
) (*state.BeaconState, error) {

	h := header.Message

	if h.Slot < 1 || types.Slot(h.Slot) > beaconState.Slot() {
		return nil, errors.New("incorrect header slot")
	}
	epoch := helpers.SlotToEpoch(types.Slot(h.Slot))
	currEpoch := helpers.CurrentEpoch(beaconState)
	prevEpoch := helpers.PrevEpoch(beaconState)
	if epoch != prevEpoch && epoch != currEpoch {
		return nil, fmt.Errorf(
			"expected target epoch (%d) to be the previous epoch (%d) or the current epoch (%d)",
			epoch,
			prevEpoch,
			currEpoch,
		)
	}
	c := ActiveShardCount()
	if h.Shard >= c {
		return nil, fmt.Errorf("shard %d >= shard count %d", h.Shard, c)
	}
	r, err := helpers.BlockRootAtSlot(beaconState, types.Slot(h.Slot)-1)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(r, h.BodySummary.BeaconBlockRoot) {
		return nil, errors.New("incorrect beacon block root")
	}
	i, err := ShardProposerIndex(beaconState, types.Slot(h.Slot), h.Shard)
	if err != nil {
		return nil, err
	}
	if i != types.ValidatorIndex(h.ProposerIndex) {
		return nil, errors.New("incorrect proposer index")
	}
	if err := helpers.ComputeDomainVerifySigningRoot(beaconState, types.ValidatorIndex(h.ProposerIndex), currEpoch, h, params.BeaconConfig().DomainShardProposer, header.Signature); err != nil {
		return nil, err
	}

	// TODO: Verify the length by verifying the degree.

	var pendingHeaders []*pb.PendingShardHeader
	var currentEpochHeader bool
	if epoch == currEpoch {
		pendingHeaders = beaconState.CurrentEpochPendingShardHeaders()
		currentEpochHeader = true
	} else {
		pendingHeaders = beaconState.PreviousEpochPendingShardHeaders()
		currentEpochHeader = false
	}

	// Check header is not yet in the pending list
	headerRoot, err := h.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	for _, pendingHeader := range pendingHeaders {
		if bytes.Equal(pendingHeader.Root, headerRoot[:]) {
			return nil, errors.New("incorrect header root")
		}
	}

	ci, err := CommitteeIndexFromShard(beaconState, types.Slot(h.Slot), h.Shard)
	if err != nil {
		return nil, err
	}
	indices, err := helpers.BeaconCommitteeFromState(beaconState, types.Slot(h.Slot), ci)
	if err != nil {
		return nil, err
	}

	ph := &pb.PendingShardHeader{
		Slot:           types.Slot(h.Slot),
		Shard:          h.Shard,
		DataCommitment: h.BodySummary.Commitment,
		Root:           headerRoot[:],
		Votes:          bitfield.NewBitlist(uint64(len(indices))),
		Confirmed:      false,
	}

	if currentEpochHeader {
		beaconState.AppendCurrentEpochPendingShardHeader(ph)
	} else {
		beaconState.AppendPreviousEpochPendingShardHeader(ph)
	}

	return beaconState, nil
}
