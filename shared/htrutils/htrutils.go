// Package htrutils defines HashTreeRoot utility functions.
package htrutils

import (
	"bytes"
	"encoding/binary"
	"fmt"

	fssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/hashutil"
	"github.com/prysmaticlabs/prysm/shared/params"
)

// Uint64Root computes the HashTreeRoot Merkleization of
// a simple uint64 value according to the eth2
// Simple Serialize specification.
func Uint64Root(val uint64) [32]byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, val)
	root := bytesutil.ToBytes32(buf)
	return root
}

// ForkRoot computes the HashTreeRoot Merkleization of
// a Fork struct value according to the eth2
// Simple Serialize specification.
func ForkRoot(fork *pb.Fork) ([32]byte, error) {
	fieldRoots := make([][]byte, 3)
	if fork != nil {
		prevRoot := bytesutil.ToBytes32(fork.PreviousVersion)
		fieldRoots[0] = prevRoot[:]
		currRoot := bytesutil.ToBytes32(fork.CurrentVersion)
		fieldRoots[1] = currRoot[:]
		forkEpochBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(forkEpochBuf, uint64(fork.Epoch))
		epochRoot := bytesutil.ToBytes32(forkEpochBuf)
		fieldRoots[2] = epochRoot[:]
	}
	return BitwiseMerkleize(hashutil.CustomSHA256Hasher(), fieldRoots, uint64(len(fieldRoots)), uint64(len(fieldRoots)))
}

// CheckpointRoot computes the HashTreeRoot Merkleization of
// a InitWithReset struct value according to the eth2
// Simple Serialize specification.
func CheckpointRoot(hasher HashFn, checkpoint *ethpb.Checkpoint) ([32]byte, error) {
	fieldRoots := make([][]byte, 2)
	if checkpoint != nil {
		epochBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(epochBuf, uint64(checkpoint.Epoch))
		epochRoot := bytesutil.ToBytes32(epochBuf)
		fieldRoots[0] = epochRoot[:]
		ckpRoot := bytesutil.ToBytes32(checkpoint.Root)
		fieldRoots[1] = ckpRoot[:]
	}
	return BitwiseMerkleize(hasher, fieldRoots, uint64(len(fieldRoots)), uint64(len(fieldRoots)))
}

// HistoricalRootsRoot computes the HashTreeRoot Merkleization of
// a list of [32]byte historical block roots according to the eth2
// Simple Serialize specification.
func HistoricalRootsRoot(historicalRoots [][]byte) ([32]byte, error) {
	result, err := BitwiseMerkleize(hashutil.CustomSHA256Hasher(), historicalRoots, uint64(len(historicalRoots)), params.BeaconConfig().HistoricalRootsLimit)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not compute historical roots merkleization")
	}
	historicalRootsBuf := new(bytes.Buffer)
	if err := binary.Write(historicalRootsBuf, binary.LittleEndian, uint64(len(historicalRoots))); err != nil {
		return [32]byte{}, errors.Wrap(err, "could not marshal historical roots length")
	}
	// We need to mix in the length of the slice.
	historicalRootsOutput := make([]byte, 32)
	copy(historicalRootsOutput, historicalRootsBuf.Bytes())
	mixedLen := MixInLength(result, historicalRootsOutput)
	return mixedLen, nil
}

// SlashingsRoot computes the HashTreeRoot Merkleization of
// a list of uint64 slashing values according to the eth2
// Simple Serialize specification.
func SlashingsRoot(slashings []uint64) ([32]byte, error) {
	slashingMarshaling := make([][]byte, params.BeaconConfig().EpochsPerSlashingsVector)
	for i := 0; i < len(slashings) && i < len(slashingMarshaling); i++ {
		slashBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(slashBuf, slashings[i])
		slashingMarshaling[i] = slashBuf
	}
	slashingChunks, err := Pack(slashingMarshaling)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not pack slashings into chunks")
	}
	return BitwiseMerkleize(hashutil.CustomSHA256Hasher(), slashingChunks, uint64(len(slashingChunks)), uint64(len(slashingChunks)))
}

func TransactionsRoot(txs [][]byte) ([32]byte, error) {
	hasher := hashutil.CustomSHA256Hasher()
	listMarshaling := make([][]byte, 0)
	for i := 0; i < len(txs); i++ {
		rt, err := TransactionRoot(txs[i])
		if err != nil {
			return [32]byte{}, err
		}
		listMarshaling = append(listMarshaling, rt[:])
	}

	bytesRoot, err := BitwiseMerkleize(hasher, listMarshaling, uint64(len(listMarshaling)), params.BeaconConfig().MaxExecutionTransactions)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not compute  merkleization")
	}
	bytesRootBuf := new(bytes.Buffer)
	if err := binary.Write(bytesRootBuf, binary.LittleEndian, uint64(len(txs))); err != nil {
		return [32]byte{}, errors.Wrap(err, "could not marshal length")
	}
	bytesRootBufRoot := make([]byte, 32)
	copy(bytesRootBufRoot, bytesRootBuf.Bytes())
	return MixInLength(bytesRoot, bytesRootBufRoot), nil
}

func SingleTransactionChunkerFastSSZ(tx []byte) ([32]byte, error) {
	// Chunk the transaction into 32 byte pieces and merkleize.
	hh := fssz.DefaultHasherPool.Get()
	index := hh.Index()
	hh.PutBytes(tx)
	maxLength := params.BeaconConfig().MaxBytesPerOpaqueTransaction
	hh.MerkleizeWithMixin(index, uint64(len(tx)), maxLength)
	return hh.HashRoot()
}

func SingleTransactionChunker(tx []byte) ([32]byte, error) {
	hasher := hashutil.CustomSHA256Hasher()
	chunkedRoots, err := packChunks(tx)
	if err != nil {
		return [32]byte{}, err
	}
	fmt.Println(len(chunkedRoots))
	maxLength := (params.BeaconConfig().MaxBytesPerOpaqueTransaction + 31) / 32
	return BitwiseMerkleize(hasher, chunkedRoots, uint64(len(chunkedRoots)), maxLength)
}

func FastSSZTransactionsRoot(txs [][]byte) ([32]byte, error) {
	hh := fssz.DefaultHasherPool.Get()
	if err := FastSSZTransactionsRootInner(hh, txs); err != nil {
		fssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	fssz.DefaultHasherPool.Put(hh)
	return root, err
}

func FastSSZTransactionsRootInner(hh *fssz.Hasher, txs [][]byte) error {
	subIndx := hh.Index()
	num := uint64(len(txs))
	if num > 16384 {
		return fssz.ErrIncorrectListSize
	}
	for i := uint64(0); i < num; i++ {
		txSubIndx := hh.Index()
		innerHH := fssz.DefaultHasherPool.Get()
		chunks, err := packChunks(txs[i])
		if err != nil {
			fssz.DefaultHasherPool.Put(innerHH)
			return err
		}
		maxLength := uint64((1048576 + 31) / 32)
		startIdx := innerHH.Index()
		for _, chunk := range chunks {
			innerHH.Append(chunk)
		}
		numItems := uint64(len(chunks))
		innerHH.MerkleizeWithMixin(startIdx, numItems, maxLength)
		chunksRoot, err := innerHH.HashRoot()
		if err != nil {
			fmt.Println("got issue")
			fssz.DefaultHasherPool.Put(innerHH)
			return err
		}
		fssz.DefaultHasherPool.Put(innerHH)
		fmt.Printf("Got limit %d and num chunks %d, chunks root %#x\n", maxLength, numItems, chunksRoot)

		// Put in the chunks root and add length mixin.
		hh.Append(chunksRoot[:])
		numItems = uint64(len(txs[i]))
		hh.MerkleizeWithMixin(txSubIndx, numItems, maxLength)
	}
	hh.MerkleizeWithMixin(subIndx, num, 16384)
	return nil
}

func TransactionRoot(tx []byte) ([32]byte, error) {
	hasher := hashutil.CustomSHA256Hasher()
	chunkedRoots, err := packChunks(tx)
	if err != nil {
		return [32]byte{}, err
	}

	maxLength := (params.BeaconConfig().MaxBytesPerOpaqueTransaction + 31) / 32
	bytesRoot, err := BitwiseMerkleize(hasher, chunkedRoots, uint64(len(chunkedRoots)), maxLength)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not compute merkleization")
	}
	fmt.Printf("Chunk root %#x\n", bytesRoot)
	return bytesRoot, nil
}

// Pack a given byte array into chunks. It'll pad the last chunk with zero bytes if
// it does not have length bytes per chunk.
func packChunks(bytes []byte) ([][]byte, error) {
	numItems := len(bytes)
	var chunks [][]byte
	for i := 0; i < numItems; i += 32 {
		j := i + 32
		// We create our upper bound index of the chunk, if it is greater than numItems,
		// we set it as numItems itself.
		if j > numItems {
			j = numItems
		}
		// We create chunks from the list of items based on the
		// indices determined above.
		chunks = append(chunks, bytes[i:j])
	}

	if len(chunks) == 0 {
		return chunks, nil
	}

	// Right-pad the last chunk with zero bytes if it does not
	// have length bytes.
	lastChunk := chunks[len(chunks)-1]
	for len(lastChunk) < 32 {
		lastChunk = append(lastChunk, 0)
	}
	chunks[len(chunks)-1] = lastChunk
	return chunks, nil
}
