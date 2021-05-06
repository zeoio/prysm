package shard

import (
	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	state "github.com/prysmaticlabs/prysm/beacon-chain/state/stateV0"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/hashutil"
	"github.com/prysmaticlabs/prysm/shared/params"
)

// ShardProposerIndex returns the validator index for the proposer of a given shard block in a given slot.
// Randomly samples from the shard proposer committee that changes once per day.
//
// Spec code:
// def get_shard_proposer_index(beacon_state: BeaconState, slot: Slot, shard: Shard) -> ValidatorIndex:
//    """
//    Return the proposer's index of shard block at ``slot``.
//    """
//    epoch = compute_epoch_at_slot(slot)
//    committee = get_shard_committee(beacon_state, epoch, shard)
//    seed = hash(get_seed(beacon_state, epoch, DOMAIN_BEACON_PROPOSER) + uint_to_bytes(slot))
//
//    # Proposer must have sufficient balance to pay for worst case fee burn
//    EFFECTIVE_BALANCE_MAX_DOWNWARD_DEVIATION = (
//        (EFFECTIVE_BALANCE_INCREMENT - EFFECTIVE_BALANCE_INCREMENT)
//        * HYSTERESIS_DOWNWARD_MULTIPLIER // HYSTERESIS_QUOTIENT
//    )
//    min_effective_balance = (
//        beacon_state.shard_gasprice * MAX_SAMPLES_PER_BLOCK // TARGET_SAMPLES_PER_BLOCK
//        + EFFECTIVE_BALANCE_MAX_DOWNWARD_DEVIATION
//    )
//    return compute_proposer_index(beacon_state, committee, seed, min_effective_balance)
func ShardProposerIndex(beaconState *state.BeaconState, slot types.Slot, shard uint64) (types.ValidatorIndex, error) {
	shardCommittee, err := ShardCommittee(beaconState, helpers.SlotToEpoch(slot), shard)
	if err != nil {
		return 0, err
	}

	seed, err := helpers.Seed(beaconState, helpers.CurrentEpoch(beaconState), params.BeaconConfig().DomainBeaconProposer)
	if err != nil {
		return 0, err
	}
	seedWithSlot := append(seed[:], bytesutil.Bytes8(uint64(slot))...)
	seedWithSlotHash := hashutil.Hash(seedWithSlot)

	// Proposer must have sufficient balance to pay for worst case fee burn.

	inc := params.BeaconConfig().EffectiveBalanceIncrement
	downwardDeviation := inc - inc*params.BeaconConfig().HysteresisDownwardMultiplier/params.BeaconConfig().HysteresisQuotient
	minBalance := beaconState.ShardGasPrice()*params.BeaconConfig().MaxSamplesPerBlock/params.BeaconConfig().TargetSamplesPerBlock + downwardDeviation

	return helpers.ComputeProposerIndex(beaconState, shardCommittee, seedWithSlotHash, minBalance)
}
