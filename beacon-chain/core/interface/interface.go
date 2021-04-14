package interfaces

import (
	types "github.com/prysmaticlabs/eth2-types"
	v1 "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

type SignedBeaconBlock interface {
	GetSignature() []byte
	GetBlock() BeaconBlock
}

type BeaconBlock interface {
	GetSlot() types.Slot
	GetProposerIndex() types.ValidatorIndex
	GetParentRoot() []byte
	GetStateRoot() []byte
	GetBody() BeaconBlockBody
	HashTreeRoot() ([32]byte, error)
}

type BeaconBlockBody interface {
	GetRandaoReveal() []byte
	GetEth1Data() *v1.Eth1Data
	GetProposerSlashings() []*v1.ProposerSlashing
	GetAttesterSlashings() []*v1.AttesterSlashing
	GetAttestations() []*v1.Attestation
	GetDeposits() []*v1.Deposit
	GetVoluntaryExits() []*v1.SignedVoluntaryExit
	HashTreeRoot() ([32]byte, error)
}

type SignedBeaconBlockWrapper struct {
	Sbb interface{}
}

func (ssbw *SignedBeaconBlockWrapper) GetBlock() BeaconBlock {
	switch ssbw.Sbb.(type) {
	case *v1.SignedBeaconBlock:
		sbb := ssbw.Sbb.(*v1.SignedBeaconBlock)
		return sbb.GetBlock()
	default:
		panic("block is not a known type")
	}

}

func (ssbw *SignedBeaconBlockWrapper) GetBlockBody() BeaconBlockBody {
	switch ssbw.Sbb.(type) {
	case *v1.SignedBeaconBlock:
		sbb := ssbw.Sbb.(*v1.SignedBeaconBlock)
		return BeaconBlockBody(sbb.GetBlock().GetBody())
	default:
		panic("block body is not a known type")
	}
}
