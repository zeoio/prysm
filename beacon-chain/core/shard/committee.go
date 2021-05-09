package shard

import (
	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	state "github.com/prysmaticlabs/prysm/beacon-chain/state/stateV0"
	"github.com/prysmaticlabs/prysm/shared/params"
)

// CommitteeCountPerSlot returns the number of crosslink committees of a slot. The
// active validator count is provided as an argument rather than a imported implementation
// from the spec definition. Having the active validator count as an argument allows for
// cheaper computation, instead of retrieving head state, one can retrieve the validator
// count.
//
//
// Spec pseudocode definition:
//   def get_committee_count_per_slot(state: BeaconState, epoch: Epoch) -> uint64:
//    """
//    Return the number of committees in each slot for the given ``epoch``.
//    """
//    return max(uint64(1), min(
//        get_active_shard_count(state, epoch),
//        uint64(len(get_active_validator_indices(state, epoch))) // SLOTS_PER_EPOCH // TARGET_COMMITTEE_SIZE,
//    ))
func CommitteeCountPerSlot(activeValidatorCount uint64) uint64 {
	committeePerSlot := activeValidatorCount / uint64(params.BeaconConfig().SlotsPerEpoch) /
		params.BeaconConfig().TargetCommitteeSize

	// Committee per slot can't be greater than shard count.
	if committeePerSlot > ActiveShardCount() {
		committeePerSlot = ActiveShardCount()
	}

	// Committee per slot can't be less than 1.
	if committeePerSlot < 1 {
		return 1
	}

	return committeePerSlot
}

// ActiveShardCount returns the active shard count.
// Currently 64, may be changed in the future.
//
// Spec code:
// def get_active_shard_count(state: BeaconState, epoch: Epoch) -> uint64:
//    """
//    Return the number of active shards.
//    Note that this puts an upper bound on the number of committees per slot.
//    """
//    return INITIAL_ACTIVE_SHARDS
func ActiveShardCount() uint64 {
	return params.BeaconConfig().InitialActiveShards
}

// ShardCommittee returns the shard committee of a given epoch and shard.
// The proposer of a shard block is randomly sampled from the shard committee,
// which changes only once per ~1 day (with committees being computable 1 day ahead of time).
// def get_shard_committee(beacon_state: BeaconState, epoch: Epoch, shard: Shard) -> Sequence[ValidatorIndex]:
//    """
//    Return the shard committee of the given ``epoch`` of the given ``shard``.
//    """
//    source_epoch = compute_committee_source_epoch(epoch, SHARD_COMMITTEE_PERIOD)
//    active_validator_indices = get_active_validator_indices(beacon_state, source_epoch)
//    seed = get_seed(beacon_state, source_epoch, DOMAIN_SHARD_COMMITTEE)
//    return compute_committee(
//        indices=active_validator_indices,
//        seed=seed,
//        index=shard,
//        count=get_active_shard_count(beacon_state, epoch),
//    )
func ShardCommittee(beaconState *state.BeaconState, epoch types.Epoch, shard uint64) ([]types.ValidatorIndex, error) {
	e := helpers.CommitteeSourceEpoch(epoch, params.BeaconConfig().ShardCommitteePeriod)
	activeValidatorIndices, err := helpers.ActiveValidatorIndices(beaconState, e)
	if err != nil {
		return nil, err
	}
	seed, err := helpers.Seed(beaconState, e, params.BeaconConfig().DomainShardCommittee)
	if err != nil {
		return nil, err
	}
	return helpers.ComputeCommittee(activeValidatorIndices, seed, shard, ActiveShardCount())
}

// ShardFromCommitteeIndex converts the index of a committee into which shard that committee is responsible for
// at the given slot.
//
// Spec code:
// def compute_shard_from_committee_index(state: BeaconState, slot: Slot, index: CommitteeIndex) -> Shard:
//    active_shards = get_active_shard_count(state, compute_epoch_at_slot(slot))
//    return Shard((index + get_start_shard(state, slot)) % active_shards)
func ShardFromCommitteeIndex(beaconState *state.BeaconState, slot types.Slot, committeeID types.CommitteeIndex) (uint64, error) {
	activeShards := ActiveShardCount()
	startShard, err := StartShard(beaconState, slot)
	if err != nil {
		return 0, err
	}
	return (startShard + uint64(committeeID)) % activeShards, nil
}

// CommitteeIndexFromShard converts shard into committee index  that is responsible for
// at the given slot.
//
// Spec code:
// def compute_committee_index_from_shard(state: BeaconState, slot: Slot, shard: Shard) -> CommitteeIndex:
//    active_shards = get_active_shard_count(state, compute_epoch_at_slot(slot))
//    return CommitteeIndex((active_shards + shard - get_start_shard(state, slot)) % active_shards)
func CommitteeIndexFromShard(beaconState *state.BeaconState, slot types.Slot, shard uint64) (types.CommitteeIndex, error) {
	activeShards := ActiveShardCount()
	startShard, err := StartShard(beaconState, slot)
	if err != nil {
		return 0, err
	}
	return types.CommitteeIndex((shard + activeShards - startShard) % activeShards), nil
}
