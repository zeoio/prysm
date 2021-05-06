package shard

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/validators"
	iface "github.com/prysmaticlabs/prysm/beacon-chain/state/interface"
	"github.com/prysmaticlabs/prysm/shared/params"
)

// ProcessShardProposerSlashing is one of the operations performed
// on each processed beacon block to slash shard proposers based on
// slashing conditions if any slashable events occurred.
//
// Spec pseudocode definition:
// def process_shard_proposer_slashing(state: BeaconState, proposer_slashing: ShardProposerSlashing) -> None:
//    reference_1 = proposer_slashing.signed_reference_1.message
//    reference_2 = proposer_slashing.signed_reference_2.message
//
//    # Verify header slots match
//    assert reference_1.slot == reference_2.slot
//    # Verify header shards match
//    assert reference_1.shard == reference_2.shard
//    # Verify header proposer indices match
//    assert reference_1.proposer_index == reference_2.proposer_index
//    # Verify the headers are different (i.e. different body)
//    assert reference_1 != reference_2
//    # Verify the proposer is slashable
//    proposer = state.validators[reference_1.proposer_index]
//    assert is_slashable_validator(proposer, get_current_epoch(state))
//    # Verify signatures
//    for signed_header in (proposer_slashing.signed_reference_1, proposer_slashing.signed_reference_2):
//        domain = get_domain(state, DOMAIN_SHARD_PROPOSER, compute_epoch_at_slot(signed_header.message.slot))
//        signing_root = compute_signing_root(signed_header.message, domain)
//        assert bls.Verify(proposer.pubkey, signing_root, signed_header.signature)
//
//    slash_validator(state, reference_1.proposer_index)
func ProcessShardProposerSlashing(
	_ context.Context,
	beaconState iface.BeaconState,
	slashings *ethpb.ShardProposerSlashing,
) (iface.BeaconState, error) {
	slashing1 := slashings.SignedReference_1.Message
	slashing2 := slashings.SignedReference_2.Message

	if slashing1.Slot != slashing2.Slot {
		return nil, errors.New("mismatch slots")
	}
	if slashing1.Shard != slashing2.Shard {
		return nil, errors.New("mismatch shards")
	}
	if slashing1.ProposerIndex != slashing2.ProposerIndex {
		return nil, errors.New("mismatch proposer indices")
	}
	if proto.Equal(slashing1, slashing2) {
		return nil, errors.New("expected slashing headers to differ")
	}

	proposer, err := beaconState.ValidatorAtIndexReadOnly(slashing1.ProposerIndex)
	if err != nil {
		return nil, err
	}
	if !helpers.IsSlashableValidatorUsingTrie(proposer, helpers.CurrentEpoch(beaconState)) {
		return nil, fmt.Errorf("validator with key %#x is not slashable", proposer.PublicKey())
	}

	for _, ref := range []*ethpb.SignedShardBlobReference{slashings.SignedReference_1, slashings.SignedReference_2} {
		if err := helpers.ComputeDomainVerifySigningRoot(beaconState, ref.Message.ProposerIndex, helpers.SlotToEpoch(ref.Message.Slot),
			ref.Message, params.BeaconConfig().DomainShardProposer, ref.Signature); err != nil {
			return nil, errors.Wrap(err, "could not verify header signature")
		}
	}

	beaconState, err = validators.SlashValidator(beaconState, slashing1.ProposerIndex)
	if err != nil {
		return nil, err
	}

	return beaconState, nil
}
