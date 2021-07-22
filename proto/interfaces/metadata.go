package interfaces

import (
	"github.com/prysmaticlabs/go-bitfield"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
)

type Metadata interface {
	SequenceNumber() uint64
	AttnetsBitfield() bitfield.Bitvector64
	InnerObject() interface{}
	IsNil() bool
	Copy() Metadata
	MetadataObjV0() *pb.MetaDataV0
	MetadataObjV1() *pb.MetaDataV1
}
