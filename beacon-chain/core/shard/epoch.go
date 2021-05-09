package shard

import (
	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	state "github.com/prysmaticlabs/prysm/beacon-chain/state/stateV0"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/attestationutil"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/params"
)

// ProcessPendingHeaders for beacon chain.
func ProcessPendingHeaders(state *state.BeaconState) (*state.BeaconState, error) {
	if helpers.CurrentEpoch(state) == params.BeaconConfig().GenesisEpoch {
		return state, nil
	}

	prevEpoch := helpers.PrevEpoch(state)
	prevEpochStartSlot, err := helpers.StartSlot(prevEpoch)
	if err != nil {
		return nil, err
	}
	for slot := prevEpochStartSlot; slot < slot+params.BeaconConfig().SlotsPerEpoch; slot++ {
		for shard := uint64(0); shard < ActiveShardCount(); shard++ {
			var candidates []*pb.PendingShardHeader
			confirmed := false
			for _, header := range state.PreviousEpochPendingShardHeaders() {
				if header.Slot == slot && header.Shard == shard {
					candidates = append(candidates, header)
					if header.Confirmed == true {
						confirmed = true
					}
				}
			}
			if confirmed {
				continue
			}

			index, err := CommitteeIndexFromShard(state, slot, shard)
			if err != nil {
				return nil, err
			}
			committee, err := helpers.BeaconCommitteeFromState(state, slot, index)
			if err != nil {
				return nil, err
			}
			var bestIndex int
			var bestBalance uint64
			for i, candidate := range candidates {
				committee, err := attestationutil.AttestingIndices(candidate.Votes, committee)
				if err != nil {
					return nil, err
				}
				attested := make([]types.ValidatorIndex, len(committee))
				for c, i := range committee {
					attested[i] = types.ValidatorIndex(c)
				}
				attestedBalance := helpers.TotalBalance(state, attested)
				if attestedBalance > bestBalance {
					bestBalance = attestedBalance
					bestIndex = i
				}
			}

			if bestBalance == 0 {
				for i, candidate := range candidates {
					if bytesutil.ToBytes32(candidate.Root) == params.BeaconConfig().ZeroHash {
						bestIndex = i
						break
					}
				}
			}
			candidates[bestIndex].Confirmed = true
		}
	}

	// Update grand parent epoch confirmed commitments

	return state, nil
}

// ChargeConfirmedHeaderFees for beacon chain.
func ChargeConfirmedHeaderFees(state *state.BeaconState) (*state.BeaconState, error) {
	newGasPrice := state.ShardGasPrice()
	adjustmentQuotient := ActiveShardCount() * uint64(params.BeaconConfig().SlotsPerEpoch) * params.BeaconConfig().GaspriceAdjustmentCoefficient
	prevEpoch := helpers.PrevEpoch(state)
	prevEpochStartSlot, err := helpers.StartSlot(prevEpoch)
	if err != nil {
		return nil, err
	}
	for slot := prevEpochStartSlot; slot < slot+params.BeaconConfig().SlotsPerEpoch; slot++ {
		for shard := uint64(0); shard < ActiveShardCount(); shard++ {
			for _, header := range state.PreviousEpochPendingShardHeaders() {
				if header.Slot == slot && header.Shard == shard && header.Confirmed {
					proposer, err := ShardProposerIndex(state, slot, shard)
					if err != nil {
						return nil, err
					}
					fee := (state.ShardGasPrice() * header.Commitment.Length) / params.BeaconConfig().TargetSamplesPerBlock
					if err := helpers.DecreaseBalance(state, proposer, fee); err != nil {
						return nil, err
					}
					newGasPrice = UpdatedGasPrice(newGasPrice, header.Commitment.Length, adjustmentQuotient)
				}
			}
		}
	}
	state.SetShardGasPrice(newGasPrice)

	return state, nil
}
