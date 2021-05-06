package stateV0

import (
	"github.com/prysmaticlabs/prysm/beacon-chain/state/stateutil"
	pbp2p "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
)

// SetEpochStartShard for the beacon state.
func (b *BeaconState) SetEpochStartShard(val uint64) error {
	if !b.hasInnerState() {
		return ErrNilInnerState
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	b.state.CurrentEpochStartShard = val
	b.markFieldAsDirty(currentEpochStartShard)
	return nil
}

// SetShardGasPrice for the beacon state.
func (b *BeaconState) SetShardGasPrice(val uint64) error {
	if !b.hasInnerState() {
		return ErrNilInnerState
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	b.state.ShardGasPrice = val
	b.markFieldAsDirty(shardGasPrice)
	return nil
}

// SetPreviousEpochPendingShardHeader for the beacon state. Updates the entire
// list to a new value by overwriting the previous one.
func (b *BeaconState) SetPreviousEpochPendingShardHeader(val []*pbp2p.PendingShardHeader) error {
	if !b.hasInnerState() {
		return ErrNilInnerState
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	b.sharedFieldReferences[previousEpochPendingShardHeader].MinusRef()
	b.sharedFieldReferences[previousEpochPendingShardHeader] = stateutil.NewRef(1)

	b.state.PreviousEpochPendingShardHeaders = val
	b.markFieldAsDirty(previousEpochPendingShardHeader)
	b.rebuildTrie[previousEpochPendingShardHeader] = true
	return nil
}

// SetCurrentEpochPendingShardHeader for the beacon state. Updates the entire
// list to a new value by overwriting the previous one.
func (b *BeaconState) SetCurrentEpochPendingShardHeader(val []*pbp2p.PendingShardHeader) error {
	if !b.hasInnerState() {
		return ErrNilInnerState
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	b.sharedFieldReferences[currentEpochPendingShardHeader].MinusRef()
	b.sharedFieldReferences[currentEpochPendingShardHeader] = stateutil.NewRef(1)

	b.state.CurrentEpochPendingShardHeaders = val
	b.markFieldAsDirty(currentEpochPendingShardHeader)
	b.rebuildTrie[currentEpochPendingShardHeader] = true
	return nil
}

// AppendCurrentEpochPendingShardHeader for the beacon state. Appends the new value
// to the the end of list.
func (b *BeaconState) AppendCurrentEpochPendingShardHeader(val *pbp2p.PendingShardHeader) error {
	if !b.hasInnerState() {
		return ErrNilInnerState
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	atts := b.state.CurrentEpochPendingShardHeaders
	if b.sharedFieldReferences[currentEpochPendingShardHeader].Refs() > 1 {
		// Copy elements in underlying array by reference.
		atts = make([]*pbp2p.PendingShardHeader, len(b.state.CurrentEpochPendingShardHeaders))
		copy(atts, b.state.CurrentEpochPendingShardHeaders)
		b.sharedFieldReferences[currentEpochPendingShardHeader].MinusRef()
		b.sharedFieldReferences[currentEpochPendingShardHeader] = stateutil.NewRef(1)
	}

	b.state.CurrentEpochPendingShardHeaders = append(atts, val)
	b.markFieldAsDirty(currentEpochPendingShardHeader)
	b.dirtyIndices[currentEpochPendingShardHeader] = append(b.dirtyIndices[currentEpochPendingShardHeader], uint64(len(b.state.CurrentEpochPendingShardHeaders)-1))
	return nil
}

// AppendPreviousEpochPendingShardHeader for the beacon state. Appends the new value
// to the the end of list.
func (b *BeaconState) AppendPreviousEpochPendingShardHeader(val *pbp2p.PendingShardHeader) error {
	if !b.hasInnerState() {
		return ErrNilInnerState
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	atts := b.state.PreviousEpochPendingShardHeaders
	if b.sharedFieldReferences[previousEpochPendingShardHeader].Refs() > 1 {
		// Copy elements in underlying array by reference.
		atts = make([]*pbp2p.PendingShardHeader, len(b.state.PreviousEpochPendingShardHeaders))
		copy(atts, b.state.PreviousEpochPendingShardHeaders)
		b.sharedFieldReferences[previousEpochPendingShardHeader].MinusRef()
		b.sharedFieldReferences[previousEpochPendingShardHeader] = stateutil.NewRef(1)
	}

	b.state.PreviousEpochPendingShardHeaders = append(atts, val)
	b.markFieldAsDirty(previousEpochPendingShardHeader)
	b.dirtyIndices[previousEpochPendingShardHeader] = append(b.dirtyIndices[previousEpochPendingShardHeader], uint64(len(b.state.PreviousEpochPendingShardHeaders)-1))
	return nil
}
