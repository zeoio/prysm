package powchain

import (
	"bufio"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	common2 "github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/prysmaticlabs/prysm/shared/fileutil"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
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

func TestRecreateDepositTrie(t *testing.T) {
	base := "" // Change me
	depositsMetadataCSVPath := filepath.Join(base, "deposit_tree_roots.csv")
	depositsMetadataCSV, err := os.Open(depositsMetadataCSVPath)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, depositsMetadataCSV.Close())
	}()
	items := readCSVDepositsMetadata(t, depositsMetadataCSV)
	_ = items
	//
	//depositsPath := filepath.Join(base, "deposits")
	//trie, err := trieutil.NewTrie(params.BeaconConfig().DepositContractTreeDepth)
	//require.NoError(t, err)

	// Try the first 10 deposts before attempting a large computation
	// as a simple sanity test for our setup.
	depositsPath := filepath.Join(base, "deposits")
	readDeposits(t, depositsPath, 10)
}

func readDeposits(t *testing.T, depositsPath string, numDeposits int) {
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
	filesInDir = filesInDir[:numDeposits]
	for i := 0; i < len(filesInDir); i++ {
		enc, err := fileutil.ReadFileAsBytes(filepath.Join(depositsPath, filesInDir[i]))
		require.NoError(t, err)
		var item *depositFileContents
		err = json.Unmarshal(enc, &item)
		require.NoError(t, err)
		fmt.Printf("%+v", item)
		fmt.Println("")
	}
}

func readCSVDepositsMetadata(t *testing.T, rs io.ReadSeeker) []*depositsMetadata {
	// Skip first row of headers.
	row1, err := bufio.NewReader(rs).ReadSlice('\n')
	require.NoError(t, err)
	_, err = rs.Seek(int64(len(row1)), io.SeekStart)
	require.NoError(t, err)

	// Read remaining rows
	r := csv.NewReader(rs)
	rows, err := r.ReadAll()
	require.NoError(t, err)
	results := make([]*depositsMetadata, len(rows))
	for i := 0; i < len(rows); i++ {
		// deposit_index,deposit_data_root,deposit_tree_root,deposit_tree_contents_root,block_hash
		data := rows[i]
		depositIdx, err := strconv.Atoi(data[0])
		require.NoError(t, err)

		results[i] = &depositsMetadata{
			DepositIndex:            uint64(depositIdx),
			DepositDataRoot:         hexDecodeOrDie(t, data[1]),
			DepositTreeRoot:         hexDecodeOrDie(t, data[2]),
			DepositTreeContentsRoot: hexDecodeOrDie(t, data[3]),
			BlockHash:               hexDecodeOrDie(t, data[4]),
		}
	}
	return results
}

func hexDecodeOrDie(t *testing.T, hexStr string) []byte {
	res, err := hex.DecodeString(strings.TrimPrefix(hexStr, "0x"))
	require.NoError(t, err)
	return res
}
