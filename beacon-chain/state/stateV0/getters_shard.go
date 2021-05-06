package stateV0

import pbp2p "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"

// CurrentEpochStartShard of the current beacon chain state.
func (b *BeaconState) CurrentEpochStartShard() uint64 {
	if !b.hasInnerState() {
		return 0
	}
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.state.CurrentEpochStartShard
}

// ShardGasPrice of the current beacon chain state.
func (b *BeaconState) ShardGasPrice() uint64 {
	if !b.hasInnerState() {
		return 0
	}
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.state.ShardGasPrice
}

// CurrentEpochShardPendingHeaders of the beacon chain.
func (b *BeaconState) CurrentEpochPendingShardHeaders() []*pbp2p.PendingShardHeader {
	if !b.hasInnerState() {
		return nil
	}
	if b.state.CurrentEpochPendingShardHeaders == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.currentEpochPendingShardHeaders()
}

// currentEpochPendingShardHeaders of the beacon chain.
// This assumes that a lock is already held on BeaconState.
func (b *BeaconState) currentEpochPendingShardHeaders() []*pbp2p.PendingShardHeader {
	if !b.hasInnerState() {
		return nil
	}

	return b.safeCopyPendingShardHeaderSlice(b.state.CurrentEpochPendingShardHeaders)
}

// PreviousEpochShardPendingHeaders of the beacon chain.
func (b *BeaconState) PreviousEpochPendingShardHeaders() []*pbp2p.PendingShardHeader {
	if !b.hasInnerState() {
		return nil
	}
	if b.state.PreviousEpochPendingShardHeaders == nil {
		return nil
	}

	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.previousEpochPendingShardHeaders()
}

// previousEpochPendingShardHeaders of the beacon chain.
// This assumes that a lock is already held on BeaconState.
func (b *BeaconState) previousEpochPendingShardHeaders() []*pbp2p.PendingShardHeader {
	if !b.hasInnerState() {
		return nil
	}

	return b.safeCopyPendingShardHeaderSlice(b.state.PreviousEpochPendingShardHeaders)
}

func (b *BeaconState) safeCopyPendingShardHeaderSlice(input []*pbp2p.PendingShardHeader) []*pbp2p.PendingShardHeader {
	if input == nil {
		return nil
	}

	res := make([]*pbp2p.PendingShardHeader, len(input))
	for i := 0; i < len(res); i++ {
		res[i] = CopyPendingShardHeader(input[i])
	}
	return res
}
