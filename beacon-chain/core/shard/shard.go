package shard

import (
	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	state "github.com/prysmaticlabs/prysm/beacon-chain/state/stateV0"
)

// def get_start_shard(state: BeaconState, slot: Slot) -> Shard:
//    """
//    Return the start shard at ``slot``.
//    """
//    current_epoch_start_slot = compute_start_slot_at_epoch(get_current_epoch(state))
//    shard = state.current_epoch_start_shard
//    if slot > current_epoch_start_slot:
//        # Current epoch or the next epoch lookahead
//        for _slot in range(current_epoch_start_slot, slot):
//            committee_count = get_committee_count_per_slot(state, compute_epoch_at_slot(Slot(_slot)))
//            active_shard_count = get_active_shard_count(state, compute_epoch_at_slot(Slot(_slot)))
//            shard = (shard + committee_count) % active_shard_count
//    elif slot < current_epoch_start_slot:
//        # Previous epoch
//        for _slot in list(range(slot, current_epoch_start_slot))[::-1]:
//            committee_count = get_committee_count_per_slot(state, compute_epoch_at_slot(Slot(_slot)))
//            active_shard_count = get_active_shard_count(state, compute_epoch_at_slot(Slot(_slot)))
//            # Ensure positive
//            shard = (shard + active_shard_count - committee_count) % active_shard_count
//    return Shard(shard)
func StartShard(beaconState *state.BeaconState, slot types.Slot) (uint64, error) {
	currentEpoch := helpers.CurrentEpoch(beaconState)
	currentEpochStartSlot, err := helpers.StartSlot(currentEpoch)
	if err != nil {
		return 0, err
	}
	shard := beaconState.CurrentEpochStartShard()

	if slot > currentEpochStartSlot {
		for i := currentEpochStartSlot; i < slot; i++ {
			activeValidatorCount, err := helpers.ActiveValidatorCount(beaconState, helpers.SlotToEpoch(i))
			if err != nil {
				return 0, err
			}
			committeeCount := CommitteeCountPerSlot(activeValidatorCount)
			activeShardCount := ActiveShardCount()
			shard = (shard + committeeCount) % activeShardCount
		}
	} else if slot < currentEpochStartSlot {
		for i := currentEpochStartSlot; i > slot; i-- {
			activeValidatorCount, err := helpers.ActiveValidatorCount(beaconState, helpers.SlotToEpoch(i))
			if err != nil {
				return 0, err
			}
			committeeCount := CommitteeCountPerSlot(activeValidatorCount)
			activeShardCount := ActiveShardCount()
			shard = (shard + activeValidatorCount - committeeCount) % activeShardCount
		}
	}
	return shard, nil
}
