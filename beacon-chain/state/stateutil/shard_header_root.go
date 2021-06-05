package stateutil

import (
	"bytes"
	"encoding/binary"

	"github.com/pkg/errors"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	ethpb "github.com/prysmaticlabs/prysm/proto/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/hashutil"
	"github.com/prysmaticlabs/prysm/shared/htrutils"
)

func pendingShardHeaderRoot(headers []*pb.PendingShardHeader) ([32]byte, error) {
	hasher := hashutil.CustomSHA256Hasher()
	roots := make([][]byte, len(headers))
	for i := 0; i < len(headers); i++ {
		pendingRoot, err := shardHeaderRoot(hasher, headers[i])
		if err != nil {
			return [32]byte{}, errors.Wrap(err, "could not compute shard header merkleization")
		}
		roots[i] = pendingRoot[:]
	}

	headersRoot, err := htrutils.BitwiseMerkleize(
		hasher,
		roots,
		uint64(len(roots)),
		131072, // TODO: Use params
	)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not compute shard header merkleization")
	}
	headerLenBuf := new(bytes.Buffer)
	if err := binary.Write(headerLenBuf, binary.LittleEndian, uint64(len(headersRoot))); err != nil {
		return [32]byte{}, errors.Wrap(err, "could not marshal epoch shard header length")
	}
	// We need to mix in the length of the slice.
	headerLenRoot := make([]byte, 32)
	copy(headerLenRoot, headerLenBuf.Bytes())
	res := htrutils.MixInLength(headersRoot, headerLenRoot)
	return res, nil
}

func shardHeaderRoot(hasher htrutils.HashFn, header *pb.PendingShardHeader) ([32]byte, error) {
	fieldRoots := make([][]byte, 6)

	if header != nil {
		slotBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(slotBuf, uint64(header.Slot))
		slotRoot := bytesutil.ToBytes32(slotBuf)
		fieldRoots[0] = slotRoot[:]

		shardBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(shardBuf, header.Shard)
		shardRoot := bytesutil.ToBytes32(shardBuf)
		fieldRoots[1] = shardRoot[:]

		dataRoot, err := dataCommitmentRoot(hasher, header.Commitment)
		if err != nil {
			return [32]byte{}, err
		}
		fieldRoots[2] = dataRoot[:]

		headerRoot := bytesutil.ToBytes32(header.Root)
		fieldRoots[3] = headerRoot[:]

		votesRoot, err := htrutils.BitlistRoot(hasher, header.Votes, 2048)
		if err != nil {
			return [32]byte{}, err
		}
		fieldRoots[4] = votesRoot[:]

		boolBytes := make([]byte, 32)
		if header.Confirmed {
			boolBytes[0] = 1
			fieldRoots[5] = boolBytes
		} else {
			boolBytes[0] = 0
			fieldRoots[5] = boolBytes
		}
	}

	return htrutils.BitwiseMerkleize(hasher, fieldRoots, uint64(len(fieldRoots)), uint64(len(fieldRoots)))
}

func dataCommitmentRoot(hasher htrutils.HashFn, commitment *ethpb.DataCommitment) ([32]byte, error) {
	fieldRoots := make([][]byte, 2)

	if commitment != nil {
		c := bytesutil.ToBytes48(commitment.Point)
		chunks, err := htrutils.Pack([][]byte{c[:]})
		if err != nil {
			return [32]byte{}, err
		}
		root, err := htrutils.BitwiseMerkleize(hasher, chunks, uint64(len(chunks)), uint64(len(chunks)))
		if err != nil {
			return [32]byte{}, err
		}
		fieldRoots[0] = root[:]

		lengthRoot := make([]byte, 32)
		binary.LittleEndian.PutUint64(lengthRoot[:8], commitment.Length)
		fieldRoots[1] = lengthRoot
	}

	return htrutils.BitwiseMerkleize(hasher, fieldRoots, uint64(len(fieldRoots)), uint64(len(fieldRoots)))
}
