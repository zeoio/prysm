package transition

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/snappy"
	types "github.com/prysmaticlabs/eth2-types"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/altair"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	coreState "github.com/prysmaticlabs/prysm/beacon-chain/core/state"
	"github.com/prysmaticlabs/prysm/beacon-chain/state/stateV0"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/testutil"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
	"github.com/prysmaticlabs/prysm/spectest/utils"
)

type ForkConfig struct {
	PostFork    string `json:"post_fork"`
	ForkEpoch   int    `json:"fork_epoch"`
	ForkBlock   int    `json:"fork_block"`
	BlocksCount int    `json:"blocks_count"`
}

// RunCoreTests is a helper function that runs Altair's transition core tests.
func RunCoreTests(t *testing.T, config string) {
	ctx := context.Background()
	require.NoError(t, utils.SetConfig(t, config))

	testFolders, testsFolderPath := utils.TestFolders(t, config, "altair", "transition/core/pyspec_tests")
	for _, folder := range testFolders {
		t.Run(folder.Name(), func(t *testing.T) {
			helpers.ClearCache()

			file, err := testutil.BazelFileBytes(testsFolderPath, folder.Name(), "meta.yaml")
			require.NoError(t, err)

			metaYaml := &ForkConfig{}
			require.NoError(t, utils.UnmarshalYaml(file, metaYaml), "Failed to Unmarshal")

			preforkBlocks := make([]*ethpb.SignedBeaconBlock, 0)
			postforkBlocks := make([]*ethpb.SignedBeaconBlockAltair, 0)
			for i := 0; i <= metaYaml.ForkBlock; i++ {
				fileName := fmt.Sprint("blocks_", i, ".ssz_snappy")
				blockFile, err := testutil.BazelFileBytes(testsFolderPath, folder.Name(), fileName)
				require.NoError(t, err)
				blockSSZ, err := snappy.Decode(nil /* dst */, blockFile)
				require.NoError(t, err, "Failed to decompress")
				block := &ethpb.SignedBeaconBlock{}
				require.NoError(t, block.UnmarshalSSZ(blockSSZ), "Failed to unmarshal")
				preforkBlocks = append(preforkBlocks, block)
			}
			t.Log(preforkBlocks[0].Block.StateRoot)
			for i := metaYaml.ForkBlock + 1; i < metaYaml.BlocksCount; i++ {
				fileName := fmt.Sprint("blocks_", i, ".ssz_snappy")
				blockFile, err := testutil.BazelFileBytes(testsFolderPath, folder.Name(), fileName)
				require.NoError(t, err)
				blockSSZ, err := snappy.Decode(nil /* dst */, blockFile)
				require.NoError(t, err, "Failed to decompress")
				block := &ethpb.SignedBeaconBlockAltair{}
				require.NoError(t, block.UnmarshalSSZ(blockSSZ), "Failed to unmarshal")
				postforkBlocks = append(postforkBlocks, block)
			}

			helpers.ClearCache()
			preBeaconStateFile, err := testutil.BazelFileBytes(testsFolderPath, folder.Name(), "pre.ssz_snappy")
			require.NoError(t, err)
			preBeaconStateSSZ, err := snappy.Decode(nil /* dst */, preBeaconStateFile)
			require.NoError(t, err, "Failed to decompress")
			beaconStateBase := &pb.BeaconState{}
			require.NoError(t, beaconStateBase.UnmarshalSSZ(preBeaconStateSSZ), "Failed to unmarshal")
			beaconState, err := stateV0.InitializeFromProto(beaconStateBase)
			require.NoError(t, err)
			var ok bool
			for _, b := range preforkBlocks {
				for beaconState.Slot() < b.Block.Slot {
					processedState, err := coreState.ProcessSlot(ctx, beaconState)
					require.NoError(t, err)
					if coreState.CanProcessEpoch(processedState) {
						processedState, err = coreState.ProcessEpochPrecompute(ctx, processedState)
						require.NoError(t, err)
					}
					require.NoError(t, processedState.SetSlot(processedState.Slot()+1))
					if helpers.IsEpochStart(processedState.Slot()) && helpers.SlotToEpoch(processedState.Slot()) == types.Epoch(metaYaml.ForkEpoch) {
						processedState, err = altair.UpgradeToAltair(processedState)
						require.NoError(t, err)
					}
					beaconState, ok = processedState.(*stateV0.BeaconState)
					require.Equal(t, true, ok)
				}

				set, beaconState, err := coreState.ProcessBlockNoVerifyAnySig(ctx, beaconState, b)
				require.NoError(t, err)
				postStateRoot, err := beaconState.HashTreeRoot(ctx)
				require.NoError(t, err)
				require.Equal(t, postStateRoot, bytesutil.ToBytes32(b.Block.StateRoot))

				valid, err := set.Verify()
				require.NoError(t, err)
				require.Equal(t, true, valid)
			}
			t.Error(beaconState.Slot())
		})
	}
}
