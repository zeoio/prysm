package kv

import (
	"context"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/beacon-chain/state"
	statev1 "github.com/prysmaticlabs/prysm/beacon-chain/state/v1"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"io"
)

// SaveStateToHead sets the current head state.
func (s *Store) SaveStateToHead(ctx context.Context, bs state.BeaconState) error {
	/*
	// TODO: fork version of the initial state will not match the fork version from genesis
	// do we need to do similar validation to the initialize-from-genesis-state code?
	// bail out early if the fork version doesn't match built-in genesis fork version
	if !bytes.Equal(bs.Fork().CurrentVersion, params.BeaconConfig().GenesisForkVersion) {
		return fmt.Errorf("loaded state's fork version (%#x) does not match config genesis "+
			"fork version (%#x)", bs.Fork().CurrentVersion, params.BeaconConfig().GenesisForkVersion)
	}
	 */

	blockRoot, err := bs.LatestBlockHeader().HashTreeRoot()
	if err != nil {
		return errors.Wrap(err, "could not compute HashTreeRoot of LatestBlockHeader")
	}
	if err = s.SaveState(ctx, bs, blockRoot); err != nil {
		return errors.Wrap(err, "could not save state")
	}
	if err = s.SaveStateSummary(ctx, &ethpb.StateSummary{
		Slot: bs.Slot(),
		Root: blockRoot[:],
	}); err != nil {
		return err
	}
	if err = s.SaveHeadBlockRoot(ctx, blockRoot); err != nil {
		return errors.Wrap(err, "could not save head block root")
	}
	if err = s.SaveWeakSubjectivityInitialBlockRoot(ctx, blockRoot); err != nil {
		return err
	}
	// TODO:
	// save head block root -- before or after sync?
	// genesis block root -- how is this used? when we don't have a genesis, what breaks?

	return nil
}

// SaveWeakSubjectivityState loads an ssz serialized BeaconState from an io.Reader
// (ex: an open file) and sets the given state to the head of the chain.
func (s *Store) SaveWeakSubjectivityState(ctx context.Context, r io.Reader) error {
	bs, err := statev1.InitializeFromSSZReader(r)
	if err != nil {
		return err
	}

	return s.SaveStateToHead(ctx, bs)
}
