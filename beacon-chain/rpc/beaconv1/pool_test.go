package beaconv1

import (
	"context"
	"testing"

	"github.com/prysmaticlabs/prysm/beacon-chain/operations/attestations"

	ptypes "github.com/gogo/protobuf/types"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1"
	ethpb_alpha "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	mock "github.com/prysmaticlabs/prysm/beacon-chain/blockchain/testing"
	"github.com/prysmaticlabs/prysm/beacon-chain/operations/slashings"
	"github.com/prysmaticlabs/prysm/beacon-chain/operations/voluntaryexits"
	"github.com/prysmaticlabs/prysm/proto/migration"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestServer_ListPoolAttestations(t *testing.T) {
	ctx := context.Background()
	beaconState, privKeys := testutil.DeterministicGenesisState(t, 64)
	v1Attestations := make([]*ethpb.Attestation, 4)
	atts1, err := testutil.GenerateAttestations(beaconState, privKeys, 2, 0, true)
	require.NoError(t, err)
	atts2, err := testutil.GenerateAttestations(beaconState, privKeys, 2, 1, true)
	require.NoError(t, err)
	attestationsInPool := append(atts1, atts2...)
	for i := 0; i < len(attestationsInPool); i++ {
		v1Attestations[i] = migration.V1Alpha1AttestationToV1(attestationsInPool[0])
	}
	tests := []struct {
		name    string
		pending []*ethpb_alpha.Attestation
		want    []*ethpb.Attestation
	}{
		{
			name:    "Empty list",
			pending: []*ethpb_alpha.Attestation{},
			want:    []*ethpb.Attestation{},
		},
		{
			name:    "One",
			pending: attestationsInPool[0:1],
			want:    v1Attestations[0:1],
		},
		{
			name:    "All",
			pending: attestationsInPool,
			want:    v1Attestations,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := attestations.NewPool()
			for _, att := range tt.pending {
				require.NoError(t, pool.SaveUnaggregatedAttestation(att))
			}
			p := &Server{
				ChainInfoFetcher: &mock.ChainService{State: beaconState},
				AttestationsPool: pool,
			}
			atts, err := p.ListPoolAttestations(ctx, &ethpb.AttestationsPoolRequest{})
			require.NoError(t, err)
			require.Equal(t, len(tt.want), len(atts.Data))
			for i, _ := range tt.want {
				assert.DeepEqual(t, tt.want[i], atts.Data[i])
			}
		})
	}
}

func TestServer_ListPoolAttesterSlashings(t *testing.T) {
	ctx := context.Background()
	beaconState, privKeys := testutil.DeterministicGenesisState(t, 64)
	count := uint64(10)
	slashingsInPool := make([]*ethpb_alpha.AttesterSlashing, count)
	v1Slashings := make([]*ethpb.AttesterSlashing, count)
	for i := 0; i < len(slashingsInPool); i++ {
		sl, err := testutil.GenerateAttesterSlashingForValidator(beaconState, privKeys[i], uint64(i))
		require.NoError(t, err)
		slashingsInPool[i] = sl
		v1Slashings[i] = migration.V1Alpha1AttSlashingToV1(sl)
	}
	tests := []struct {
		name    string
		pending []*ethpb_alpha.AttesterSlashing
		want    []*ethpb.AttesterSlashing
	}{
		{
			name:    "Empty list",
			pending: []*ethpb_alpha.AttesterSlashing{},
			want:    []*ethpb.AttesterSlashing{},
		},
		{
			name:    "One",
			pending: slashingsInPool[0:1],
			want:    v1Slashings[0:1],
		},
		{
			name:    "All",
			pending: slashingsInPool,
			want:    v1Slashings,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := &slashings.Pool{}
			for _, slashing := range tt.pending {
				require.NoError(t, pool.InsertAttesterSlashing(ctx, beaconState, slashing))
			}
			p := &Server{
				ChainInfoFetcher: &mock.ChainService{State: beaconState},
				SlashingsPool:    pool,
			}
			attSlashings, err := p.ListPoolAttesterSlashings(ctx, &ptypes.Empty{})
			require.NoError(t, err)
			assert.DeepEqual(t, tt.want, attSlashings.Data)
		})
	}
}

func TestServer_ListPoolProposerSlashings(t *testing.T) {
	ctx := context.Background()
	beaconState, privKeys := testutil.DeterministicGenesisState(t, 64)
	count := uint64(10)
	slashingsInPool := make([]*ethpb_alpha.ProposerSlashing, count)
	v1Slashings := make([]*ethpb.ProposerSlashing, count)
	for i := 0; i < len(slashingsInPool); i++ {
		sl, err := testutil.GenerateProposerSlashingForValidator(beaconState, privKeys[i], uint64(i))
		require.NoError(t, err)
		slashingsInPool[i] = sl
		v1Slashings[i] = migration.V1Alpha1ProposerSlashingToV1(sl)
	}
	tests := []struct {
		name    string
		pending []*ethpb_alpha.ProposerSlashing
		want    []*ethpb.ProposerSlashing
	}{
		{
			name:    "Empty list",
			pending: []*ethpb_alpha.ProposerSlashing{},
			want:    []*ethpb.ProposerSlashing{},
		},
		{
			name:    "One",
			pending: slashingsInPool[0:1],
			want:    v1Slashings[0:1],
		},
		{
			name:    "All",
			pending: slashingsInPool,
			want:    v1Slashings,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := &slashings.Pool{}
			for _, slashing := range tt.pending {
				require.NoError(t, pool.InsertProposerSlashing(ctx, beaconState, slashing))
			}
			p := &Server{
				ChainInfoFetcher: &mock.ChainService{State: beaconState},
				SlashingsPool:    pool,
			}
			attSlashings, err := p.ListPoolProposerSlashings(ctx, &ptypes.Empty{})
			require.NoError(t, err)
			assert.DeepEqual(t, tt.want, attSlashings.Data)
		})
	}
}

func TestServer_ListPoolVoluntaryExits(t *testing.T) {
	ctx := context.Background()
	beaconState, privKeys := testutil.DeterministicGenesisState(t, 64)
	count := uint64(10)
	exitsInPool := make([]*ethpb_alpha.SignedVoluntaryExit, count)
	v1Exits := make([]*ethpb.SignedVoluntaryExit, count)
	for i := 0; i < len(exitsInPool); i++ {
		exit, err := testutil.GenerateVoluntaryExit(beaconState, privKeys[i], uint64(i))
		require.NoError(t, err)
		exitsInPool[i] = exit
		v1Exits[i] = migration.V1Alpha1ExitToV1(exit)
	}
	tests := []struct {
		name    string
		pending []*ethpb_alpha.SignedVoluntaryExit
		want    []*ethpb.SignedVoluntaryExit
	}{
		{
			name:    "Empty list",
			pending: []*ethpb_alpha.SignedVoluntaryExit{},
			want:    []*ethpb.SignedVoluntaryExit{},
		},
		{
			name:    "One",
			pending: exitsInPool[0:1],
			want:    v1Exits[0:1],
		},
		{
			name:    "All",
			pending: exitsInPool,
			want:    v1Exits,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := &voluntaryexits.Pool{}
			for _, slashing := range tt.pending {
				pool.InsertVoluntaryExit(ctx, beaconState, slashing)
			}
			p := &Server{
				ChainInfoFetcher:   &mock.ChainService{State: beaconState},
				VoluntaryExitsPool: pool,
			}
			exits, err := p.ListPoolVoluntaryExits(ctx, &ptypes.Empty{})
			require.NoError(t, err)
			assert.DeepEqual(t, tt.want, exits.Data)
		})
	}
}
