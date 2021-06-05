package shard

import (
	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	state "github.com/prysmaticlabs/prysm/beacon-chain/state/stateV0"
)

// StartShard returns the start shard of the beacon state.
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
