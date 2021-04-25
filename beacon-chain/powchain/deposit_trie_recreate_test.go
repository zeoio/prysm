package powchain

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	common2 "github.com/protolambda/zrnt/eth2/beacon/common"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/fileutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
	"github.com/prysmaticlabs/prysm/shared/trieutil"
)

type depositsMetadata struct {
	DepositIndex            uint64
	DepositDataRoot         []byte
	DepositTreeRoot         []byte
	DepositTreeContentsRoot []byte
	BlockHash               []byte
}

type depositFileContents struct {
	DepositData  common2.DepositData  `json:"data"`
	DepositIndex common2.DepositIndex `json:"deposit_index"`
	BlockHash    common.Hash          `json:"block_hash"`
	BlockNum     uint64               `json:"block_num"`
}

type depositContainer struct {
	DepositIndex    uint64
	DepositData     *ethpb.Deposit_Data
	DepositDataRoot [32]byte
}

func TestRecreateDepositTrie(t *testing.T) {
	base := "" // Change me
	depositsMetadataCSVPath := filepath.Join(base, "deposit_tree_roots.csv")
	depositsMetadataCSV, err := os.Open(depositsMetadataCSVPath)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, depositsMetadataCSV.Close())
	}()
	trieMetadataByDepositDataRoot := readCSVDepositsMetadata(t, depositsMetadataCSV)

	trie, err := trieutil.NewTrie(params.BeaconConfig().DepositContractTreeDepth)
	require.NoError(t, err)

	depositsPath := filepath.Join(base, "deposits")
	files, err := ioutil.ReadDir(depositsPath)
	require.NoError(t, err)
	require.Equal(t, true, len(files) > 0)
	filesInDir := make([]string, 0)
	for i := 0; i < len(files); i++ {
		if files[i].IsDir() {
			continue
		}
		filesInDir = append(filesInDir, files[i].Name())
	}

	deposits := make([]*depositContainer, 0)
	for i := 0; i < len(filesInDir); i++ {
		enc, err := fileutil.ReadFileAsBytes(filepath.Join(depositsPath, filesInDir[i]))
		require.NoError(t, err)
		var deposit *depositFileContents
		err = json.Unmarshal(enc, &deposit)
		require.NoError(t, err)
		pubKeyStr := deposit.DepositData.Pubkey.String()
		pubKeyBytes := hexDecodeOrDie(t, pubKeyStr)
		sigStr := deposit.DepositData.Signature.String()
		sigBytes := hexDecodeOrDie(t, sigStr)
		withdrawalStr := deposit.DepositData.WithdrawalCredentials.String()
		withdrawalBytes := hexDecodeOrDie(t, withdrawalStr)
		depositData := &ethpb.Deposit_Data{
			Amount:                uint64(deposit.DepositData.Amount),
			PublicKey:             pubKeyBytes,
			Signature:             sigBytes,
			WithdrawalCredentials: withdrawalBytes,
		}
		depositRoot, err := depositData.HashTreeRoot()
		require.NoError(t, err)

		key := fmt.Sprintf("%x:%d", depositRoot, deposit.DepositIndex)
		// Find it in the map.
		_, ok := trieMetadataByDepositDataRoot[key]
		if !ok {
			// Not a valid deposit.
			continue
		}

		deposits = append(deposits, &depositContainer{
			DepositDataRoot: depositRoot,
			DepositIndex:    uint64(deposit.DepositIndex),
			DepositData:     depositData,
		})
	}

	sort.Slice(deposits, func(i, j int) bool {
		return deposits[i].DepositIndex < deposits[j].DepositIndex
	})

	for _, depositCntr := range deposits {
		// Find it in the map.
		key := fmt.Sprintf("%x:%d", depositCntr.DepositDataRoot, depositCntr.DepositIndex)
		metadata, ok := trieMetadataByDepositDataRoot[key]
		require.Equal(t, true, ok)

		t.Logf("Inserting deposit with index %d and deposit data root %#x", depositCntr.DepositIndex, depositCntr.DepositDataRoot)
		trie.Insert(depositCntr.DepositDataRoot[:], int(depositCntr.DepositIndex))
		receivedTrieRoot := trie.Root()
		wantedRoot := metadata.DepositTreeRoot
		if !bytes.Equal(wantedRoot, receivedTrieRoot[:]) {
			t.Fatalf("Wanted deposit trie root %#x for deposit index %d and deposit data root %#x, received %#x as the deposit trie root", wantedRoot, depositCntr.DepositIndex, depositCntr.DepositDataRoot, receivedTrieRoot)
		}
	}
}

func readCSVDepositsMetadata(t *testing.T, rs io.ReadSeeker) map[string]*depositsMetadata {
	// Skip first row of headers.
	row1, err := bufio.NewReader(rs).ReadSlice('\n')
	require.NoError(t, err)
	_, err = rs.Seek(int64(len(row1)), io.SeekStart)
	require.NoError(t, err)

	// Read remaining rows
	r := csv.NewReader(rs)
	rows, err := r.ReadAll()
	require.NoError(t, err)
	trieMetadataByDepositDataRoot := make(map[string]*depositsMetadata)
	for i := 0; i < len(rows); i++ {
		// deposit_index,deposit_data_root,deposit_tree_root,deposit_tree_contents_root,block_hash
		data := rows[i]
		depositIdx, err := strconv.Atoi(data[0])
		require.NoError(t, err)

		key := fmt.Sprintf("%x:%d", bytesutil.ToBytes32(hexDecodeOrDie(t, data[1])), depositIdx)
		trieMetadataByDepositDataRoot[key] = &depositsMetadata{
			DepositIndex:            uint64(depositIdx),
			DepositDataRoot:         hexDecodeOrDie(t, data[1]),
			DepositTreeRoot:         hexDecodeOrDie(t, data[2]),
			DepositTreeContentsRoot: hexDecodeOrDie(t, data[3]),
			BlockHash:               hexDecodeOrDie(t, data[4]),
		}
	}
	return trieMetadataByDepositDataRoot
}

func hexDecodeOrDie(t *testing.T, hexStr string) []byte {
	res, err := hex.DecodeString(strings.TrimPrefix(hexStr, "0x"))
	require.NoError(t, err)
	return res
}
