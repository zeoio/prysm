package kv

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	dbIface "github.com/prysmaticlabs/prysm/beacon-chain/db/iface"
	"github.com/prysmaticlabs/prysm/beacon-chain/state"
	statev1 "github.com/prysmaticlabs/prysm/beacon-chain/state/v1"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/wrapper"
	"github.com/prysmaticlabs/prysm/shared/params"
	"io"
)

// SaveGenesisData bootstraps the beaconDB with a given genesis state.
func (s *Store) SaveGenesisData(ctx context.Context, genesisState state.BeaconState) error {
	stateRoot, err := genesisState.HashTreeRoot(ctx)
	if err != nil {
		return err
	}
	genesisBlk := blocks.NewGenesisBlock(stateRoot[:])
	genesisBlkRoot, err := genesisBlk.Block.HashTreeRoot()
	if err != nil {
		return errors.Wrap(err, "could not get genesis block root")
	}
	if err := s.SaveBlock(ctx, wrapper.WrappedPhase0SignedBeaconBlock(genesisBlk)); err != nil {
		return errors.Wrap(err, "could not save genesis block")
	}
	if err := s.SaveState(ctx, genesisState, genesisBlkRoot); err != nil {
		return errors.Wrap(err, "could not save genesis state")
	}
	if err := s.SaveStateSummary(ctx, &ethpb.StateSummary{
		Slot: 0,
		Root: genesisBlkRoot[:],
	}); err != nil {
		return err
	}

	if err := s.SaveHeadBlockRoot(ctx, genesisBlkRoot); err != nil {
		return errors.Wrap(err, "could not save head block root")
	}
	if err := s.SaveGenesisBlockRoot(ctx, genesisBlkRoot); err != nil {
		return errors.Wrap(err, "could not save genesis block root")
	}
	return nil
}

// LoadGenesis loads a genesis state from a given file path, if no genesis exists already.
func (s *Store) LoadGenesis(ctx context.Context, r io.Reader) error {
	gs, err := statev1.InitializeFromSSZReader(r)
	if err != nil {
		return err
	}

	// bail out early if the fork version doesn't match built-in genesis fork version
	if !bytes.Equal(gs.Fork().CurrentVersion, params.BeaconConfig().GenesisForkVersion) {
		return fmt.Errorf("loaded genesis fork version (%#x) does not match config genesis "+
			"fork version (%#x)", gs.Fork().CurrentVersion, params.BeaconConfig().GenesisForkVersion)
	}

	existing, err := s.GenesisState(ctx)
	if err != nil {
		return err
	}
	eq, err := gs.Equal(ctx, existing)
	if err != nil {
		return err
	}
	if eq {
		return dbIface.ErrExistingGenesisState
	}

	return s.SaveGenesisData(ctx, gs)
}

// EnsureEmbeddedGenesis checks that a genesis block has been generated when an embedded genesis
// state is used. If a genesis block does not exist, but a genesis state does, then we should call
// SaveGenesisData on the existing genesis state.
func (s *Store) EnsureEmbeddedGenesis(ctx context.Context) error {
	gb, err := s.GenesisBlock(ctx)
	if err != nil {
		return err
	}
	if gb != nil && !gb.IsNil() {
		return nil
	}
	gs, err := s.GenesisState(ctx)
	if err != nil {
		return err
	}
	if gs != nil && !gs.IsNil() {
		return s.SaveGenesisData(ctx, gs)
	}
	return nil
}
