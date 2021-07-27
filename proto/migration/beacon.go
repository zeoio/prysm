package migration

import (
	v1Alpha1 "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	v2 "github.com/prysmaticlabs/prysm/proto/prysm/v2"
)

// V1Alpha1ToV2ActiveSetChanges --
func V1Alpha1ToV2ActiveSetChanges(src *v1Alpha1.ActiveSetChanges) *v2.ActiveSetChanges {
	if src == nil {
		return &v2.ActiveSetChanges{}
	}
	return &v2.ActiveSetChanges{

		Epoch:               src.Epoch,
		ActivatedPublicKeys: src.ActivatedPublicKeys,
		ActivatedIndices:    src.ActivatedIndices,
		ExitedPublicKeys:    src.ExitedPublicKeys,
		ExitedIndices:       src.ExitedIndices,
		SlashedPublicKeys:   src.SlashedPublicKeys,
		SlashedIndices:      src.SlashedIndices,
		EjectedPublicKeys:   src.EjectedPublicKeys,
		EjectedIndices:      src.EjectedIndices,
	}
}

// V2ToV1Alpha1ActiveSetChanges --
func V2ToV1Alpha1ActiveSetChanges(src *v2.ActiveSetChanges) *v1Alpha1.ActiveSetChanges {
	if src == nil {
		return &v1Alpha1.ActiveSetChanges{}
	}
	return &v1Alpha1.ActiveSetChanges{

		Epoch:               src.Epoch,
		ActivatedPublicKeys: src.ActivatedPublicKeys,
		ActivatedIndices:    src.ActivatedIndices,
		ExitedPublicKeys:    src.ExitedPublicKeys,
		ExitedIndices:       src.ExitedIndices,
		SlashedPublicKeys:   src.SlashedPublicKeys,
		SlashedIndices:      src.SlashedIndices,
		EjectedPublicKeys:   src.EjectedPublicKeys,
		EjectedIndices:      src.EjectedIndices,
	}
}

// V1Alpha1ToV2AttestationData --
func V1Alpha1ToV2AttestationData(src *v1Alpha1.AttestationData) *v2.AttestationData {
	if src == nil {
		return &v2.AttestationData{}
	}
	return &v2.AttestationData{

		Slot:            src.Slot,
		CommitteeIndex:  src.CommitteeIndex,
		BeaconBlockRoot: src.BeaconBlockRoot,
		Source:          src.Source,
		Target:          src.Target,
	}
}

// V2ToV1Alpha1AttestationData --
func V2ToV1Alpha1AttestationData(src *v2.AttestationData) *v1Alpha1.AttestationData {
	if src == nil {
		return &v1Alpha1.AttestationData{}
	}
	return &v1Alpha1.AttestationData{

		Slot:            src.Slot,
		CommitteeIndex:  src.CommitteeIndex,
		BeaconBlockRoot: src.BeaconBlockRoot,
		Source:          src.Source,
		Target:          src.Target,
	}
}

// V1Alpha1ToV2AttestationPoolRequest --
func V1Alpha1ToV2AttestationPoolRequest(src *v1Alpha1.AttestationPoolRequest) *v2.AttestationPoolRequest {
	if src == nil {
		return &v2.AttestationPoolRequest{}
	}
	return &v2.AttestationPoolRequest{

		PageSize:  src.PageSize,
		PageToken: src.PageToken,
	}
}

// V2ToV1Alpha1AttestationPoolRequest --
func V2ToV1Alpha1AttestationPoolRequest(src *v2.AttestationPoolRequest) *v1Alpha1.AttestationPoolRequest {
	if src == nil {
		return &v1Alpha1.AttestationPoolRequest{}
	}
	return &v1Alpha1.AttestationPoolRequest{

		PageSize:  src.PageSize,
		PageToken: src.PageToken,
	}
}

// V1Alpha1ToV2AttesterSlashing --
func V1Alpha1ToV2AttesterSlashing(src *v1Alpha1.AttesterSlashing) *v2.AttesterSlashing {
	if src == nil {
		return &v2.AttesterSlashing{}
	}
	return &v2.AttesterSlashing{

		Attestation_1: V1Alpha1ToV2IndexedAttestation(src.Attestation_1),
		Attestation_2: V1Alpha1ToV2IndexedAttestation(src.Attestation_2),
	}
}

// V2ToV1Alpha1AttesterSlashing --
func V2ToV1Alpha1AttesterSlashing(src *v2.AttesterSlashing) *v1Alpha1.AttesterSlashing {
	if src == nil {
		return &v1Alpha1.AttesterSlashing{}
	}
	return &v1Alpha1.AttesterSlashing{

		Attestation_1: V2ToV1Alpha1IndexedAttestation(src.Attestation_1),
		Attestation_2: V2ToV1Alpha1IndexedAttestation(src.Attestation_2),
	}
}

// V1Alpha1ToV2BeaconCommittees --
func V1Alpha1ToV2BeaconCommittees(src *v1Alpha1.BeaconCommittees) *v2.BeaconCommittees {
	if src == nil {
		return &v2.BeaconCommittees{}
	}
	return &v2.BeaconCommittees{

		Epoch:                src.Epoch,
		Committees:           src.Committees,
		ActiveValidatorCount: src.ActiveValidatorCount,
	}
}

// V2ToV1Alpha1BeaconCommittees --
func V2ToV1Alpha1BeaconCommittees(src *v2.BeaconCommittees) *v1Alpha1.BeaconCommittees {
	if src == nil {
		return &v1Alpha1.BeaconCommittees{}
	}
	return &v1Alpha1.BeaconCommittees{

		Epoch:                src.Epoch,
		Committees:           src.Committees,
		ActiveValidatorCount: src.ActiveValidatorCount,
	}
}

// V1Alpha1ToV2BeaconConfig --
func V1Alpha1ToV2BeaconConfig(src *v1Alpha1.BeaconConfig) *v2.BeaconConfig {
	if src == nil {
		return &v2.BeaconConfig{}
	}
	return &v2.BeaconConfig{

		Config: src.Config,
	}
}

// V2ToV1Alpha1BeaconConfig --
func V2ToV1Alpha1BeaconConfig(src *v2.BeaconConfig) *v1Alpha1.BeaconConfig {
	if src == nil {
		return &v1Alpha1.BeaconConfig{}
	}
	return &v1Alpha1.BeaconConfig{

		Config: src.Config,
	}
}

// V1Alpha1ToV2ChainHead --
func V1Alpha1ToV2ChainHead(src *v1Alpha1.ChainHead) *v2.ChainHead {
	if src == nil {
		return &v2.ChainHead{}
	}
	return &v2.ChainHead{

		HeadSlot:                   src.HeadSlot,
		HeadEpoch:                  src.HeadEpoch,
		HeadBlockRoot:              src.HeadBlockRoot,
		FinalizedSlot:              src.FinalizedSlot,
		FinalizedEpoch:             src.FinalizedEpoch,
		FinalizedBlockRoot:         src.FinalizedBlockRoot,
		JustifiedSlot:              src.JustifiedSlot,
		JustifiedEpoch:             src.JustifiedEpoch,
		JustifiedBlockRoot:         src.JustifiedBlockRoot,
		PreviousJustifiedSlot:      src.PreviousJustifiedSlot,
		PreviousJustifiedEpoch:     src.PreviousJustifiedEpoch,
		PreviousJustifiedBlockRoot: src.PreviousJustifiedBlockRoot,
	}
}

// V2ToV1Alpha1ChainHead --
func V2ToV1Alpha1ChainHead(src *v2.ChainHead) *v1Alpha1.ChainHead {
	if src == nil {
		return &v1Alpha1.ChainHead{}
	}
	return &v1Alpha1.ChainHead{

		HeadSlot:                   src.HeadSlot,
		HeadEpoch:                  src.HeadEpoch,
		HeadBlockRoot:              src.HeadBlockRoot,
		FinalizedSlot:              src.FinalizedSlot,
		FinalizedEpoch:             src.FinalizedEpoch,
		FinalizedBlockRoot:         src.FinalizedBlockRoot,
		JustifiedSlot:              src.JustifiedSlot,
		JustifiedEpoch:             src.JustifiedEpoch,
		JustifiedBlockRoot:         src.JustifiedBlockRoot,
		PreviousJustifiedSlot:      src.PreviousJustifiedSlot,
		PreviousJustifiedEpoch:     src.PreviousJustifiedEpoch,
		PreviousJustifiedBlockRoot: src.PreviousJustifiedBlockRoot,
	}
}

// V1Alpha1ToV2GetValidatorActiveSetChangesRequest --
func V1Alpha1ToV2GetValidatorActiveSetChangesRequest(src *v1Alpha1.GetValidatorActiveSetChangesRequest) *v2.GetValidatorActiveSetChangesRequest {
	if src == nil {
		return &v2.GetValidatorActiveSetChangesRequest{}
	}
	return &v2.GetValidatorActiveSetChangesRequest{

		QueryFilter: src.QueryFilter,
	}
}

// V2ToV1Alpha1GetValidatorActiveSetChangesRequest --
func V2ToV1Alpha1GetValidatorActiveSetChangesRequest(src *v2.GetValidatorActiveSetChangesRequest) *v1Alpha1.GetValidatorActiveSetChangesRequest {
	if src == nil {
		return &v1Alpha1.GetValidatorActiveSetChangesRequest{}
	}
	return &v1Alpha1.GetValidatorActiveSetChangesRequest{

		QueryFilter: src.QueryFilter,
	}
}

// V1Alpha1ToV2GetValidatorParticipationRequest --
func V1Alpha1ToV2GetValidatorParticipationRequest(src *v1Alpha1.GetValidatorParticipationRequest) *v2.GetValidatorParticipationRequest {
	if src == nil {
		return &v2.GetValidatorParticipationRequest{}
	}
	return &v2.GetValidatorParticipationRequest{

		QueryFilter: src.QueryFilter,
	}
}

// V2ToV1Alpha1GetValidatorParticipationRequest --
func V2ToV1Alpha1GetValidatorParticipationRequest(src *v2.GetValidatorParticipationRequest) *v1Alpha1.GetValidatorParticipationRequest {
	if src == nil {
		return &v1Alpha1.GetValidatorParticipationRequest{}
	}
	return &v1Alpha1.GetValidatorParticipationRequest{

		QueryFilter: src.QueryFilter,
	}
}

// V1Alpha1ToV2GetValidatorRequest --
func V1Alpha1ToV2GetValidatorRequest(src *v1Alpha1.GetValidatorRequest) *v2.GetValidatorRequest {
	if src == nil {
		return &v2.GetValidatorRequest{}
	}
	return &v2.GetValidatorRequest{

		QueryFilter: src.QueryFilter,
	}
}

// V2ToV1Alpha1GetValidatorRequest --
func V2ToV1Alpha1GetValidatorRequest(src *v2.GetValidatorRequest) *v1Alpha1.GetValidatorRequest {
	if src == nil {
		return &v1Alpha1.GetValidatorRequest{}
	}
	return &v1Alpha1.GetValidatorRequest{

		QueryFilter: src.QueryFilter,
	}
}

// V1Alpha1ToV2IndexedAttestation --
func V1Alpha1ToV2IndexedAttestation(src *v1Alpha1.IndexedAttestation) *v2.IndexedAttestation {
	if src == nil {
		return &v2.IndexedAttestation{}
	}
	return &v2.IndexedAttestation{

		AttestingIndices: src.AttestingIndices,
		Data:             V1Alpha1ToV2AttestationData(src.Data),
		Signature:        src.Signature,
	}
}

// V2ToV1Alpha1IndexedAttestation --
func V2ToV1Alpha1IndexedAttestation(src *v2.IndexedAttestation) *v1Alpha1.IndexedAttestation {
	if src == nil {
		return &v1Alpha1.IndexedAttestation{}
	}
	return &v1Alpha1.IndexedAttestation{

		AttestingIndices: src.AttestingIndices,
		Data:             V2ToV1Alpha1AttestationData(src.Data),
		Signature:        src.Signature,
	}
}

// V1Alpha1ToV2IndividualVotesRequest --
func V1Alpha1ToV2IndividualVotesRequest(src *v1Alpha1.IndividualVotesRequest) *v2.IndividualVotesRequest {
	if src == nil {
		return &v2.IndividualVotesRequest{}
	}
	return &v2.IndividualVotesRequest{

		Epoch:      src.Epoch,
		PublicKeys: src.PublicKeys,
		Indices:    src.Indices,
	}
}

// V2ToV1Alpha1IndividualVotesRequest --
func V2ToV1Alpha1IndividualVotesRequest(src *v2.IndividualVotesRequest) *v1Alpha1.IndividualVotesRequest {
	if src == nil {
		return &v1Alpha1.IndividualVotesRequest{}
	}
	return &v1Alpha1.IndividualVotesRequest{

		Epoch:      src.Epoch,
		PublicKeys: src.PublicKeys,
		Indices:    src.Indices,
	}
}

// V1Alpha1ToV2IndividualVotesRespond --
func V1Alpha1ToV2IndividualVotesRespond(src *v1Alpha1.IndividualVotesRespond) *v2.IndividualVotesRespond {
	if src == nil {
		return &v2.IndividualVotesRespond{}
	}
	return &v2.IndividualVotesRespond{

		IndividualVotes: src.IndividualVotes,
	}
}

// V2ToV1Alpha1IndividualVotesRespond --
func V2ToV1Alpha1IndividualVotesRespond(src *v2.IndividualVotesRespond) *v1Alpha1.IndividualVotesRespond {
	if src == nil {
		return &v1Alpha1.IndividualVotesRespond{}
	}
	return &v1Alpha1.IndividualVotesRespond{

		IndividualVotes: src.IndividualVotes,
	}
}

// V1Alpha1ToV2ListAttestationsRequest --
func V1Alpha1ToV2ListAttestationsRequest(src *v1Alpha1.ListAttestationsRequest) *v2.ListAttestationsRequest {
	if src == nil {
		return &v2.ListAttestationsRequest{}
	}
	filter := v2.ListAttestationsRequest.GetQueryFilter(src.QueryFilter)
	return &v2.ListAttestationsRequest{
		QueryFilter: filter,
		PageSize:    src.PageSize,
		PageToken:   src.PageToken,
	}
}

// V2ToV1Alpha1ListAttestationsRequest --
func V2ToV1Alpha1ListAttestationsRequest(src *v2.ListAttestationsRequest) *v1Alpha1.ListAttestationsRequest {
	if src == nil {
		return &v1Alpha1.ListAttestationsRequest{}
	}
	filter := v1Alpha1.ListAttestationsRequest.GetQueryFilter(src.QueryFilter)
	return &v1Alpha1.ListAttestationsRequest{

		QueryFilter: filter,
		PageSize:    src.PageSize,
		PageToken:   src.PageToken,
	}
}

// V1Alpha1ToV2ListAttestationsResponse --
func V1Alpha1ToV2ListAttestationsResponse(src *v1Alpha1.ListAttestationsResponse) *v2.ListAttestationsResponse {
	if src == nil {
		return &v2.ListAttestationsResponse{}
	}
	return &v2.ListAttestationsResponse{

		Attestations:  src.Attestations,
		NextPageToken: src.NextPageToken,
		TotalSize:     src.TotalSize,
	}
}

// V2ToV1Alpha1ListAttestationsResponse --
func V2ToV1Alpha1ListAttestationsResponse(src *v2.ListAttestationsResponse) *v1Alpha1.ListAttestationsResponse {
	if src == nil {
		return &v1Alpha1.ListAttestationsResponse{}
	}
	return &v1Alpha1.ListAttestationsResponse{

		Attestations:  src.Attestations,
		NextPageToken: src.NextPageToken,
		TotalSize:     src.TotalSize,
	}
}

// V1Alpha1ToV2ListBlocksRequest --
func V1Alpha1ToV2ListBlocksRequest(src *v1Alpha1.ListBlocksRequest) *v2.ListBlocksRequest {
	if src == nil {
		return &v2.ListBlocksRequest{}
	}
	return &v2.ListBlocksRequest{

		QueryFilter: src.QueryFilter,
		PageSize:    src.PageSize,
		PageToken:   src.PageToken,
	}
}

// V2ToV1Alpha1ListBlocksRequest --
func V2ToV1Alpha1ListBlocksRequest(src *v2.ListBlocksRequest) *v1Alpha1.ListBlocksRequest {
	if src == nil {
		return &v1Alpha1.ListBlocksRequest{}
	}
	return &v1Alpha1.ListBlocksRequest{

		QueryFilter: src.QueryFilter,
		PageSize:    src.PageSize,
		PageToken:   src.PageToken,
	}
}

// V1Alpha1ToV2ListCommitteesRequest --
func V1Alpha1ToV2ListCommitteesRequest(src *v1Alpha1.ListCommitteesRequest) *v2.ListCommitteesRequest {
	if src == nil {
		return &v2.ListCommitteesRequest{}
	}
	return &v2.ListCommitteesRequest{

		QueryFilter: src.QueryFilter,
	}
}

// V2ToV1Alpha1ListCommitteesRequest --
func V2ToV1Alpha1ListCommitteesRequest(src *v2.ListCommitteesRequest) *v1Alpha1.ListCommitteesRequest {
	if src == nil {
		return &v1Alpha1.ListCommitteesRequest{}
	}
	return &v1Alpha1.ListCommitteesRequest{

		QueryFilter: src.QueryFilter,
	}
}

// V1Alpha1ToV2ListIndexedAttestationsRequest --
func V1Alpha1ToV2ListIndexedAttestationsRequest(src *v1Alpha1.ListIndexedAttestationsRequest) *v2.ListIndexedAttestationsRequest {
	if src == nil {
		return &v2.ListIndexedAttestationsRequest{}
	}
	return &v2.ListIndexedAttestationsRequest{

		QueryFilter: src.QueryFilter,
		PageSize:    src.PageSize,
		PageToken:   src.PageToken,
	}
}

// V2ToV1Alpha1ListIndexedAttestationsRequest --
func V2ToV1Alpha1ListIndexedAttestationsRequest(src *v2.ListIndexedAttestationsRequest) *v1Alpha1.ListIndexedAttestationsRequest {
	if src == nil {
		return &v1Alpha1.ListIndexedAttestationsRequest{}
	}
	return &v1Alpha1.ListIndexedAttestationsRequest{

		QueryFilter: src.QueryFilter,
		PageSize:    src.PageSize,
		PageToken:   src.PageToken,
	}
}

// V1Alpha1ToV2ListIndexedAttestationsResponse --
func V1Alpha1ToV2ListIndexedAttestationsResponse(src *v1Alpha1.ListIndexedAttestationsResponse) *v2.ListIndexedAttestationsResponse {
	if src == nil {
		return &v2.ListIndexedAttestationsResponse{}
	}
	return &v2.ListIndexedAttestationsResponse{

		IndexedAttestations: src.IndexedAttestations,
		NextPageToken:       src.NextPageToken,
		TotalSize:           src.TotalSize,
	}
}

// V2ToV1Alpha1ListIndexedAttestationsResponse --
func V2ToV1Alpha1ListIndexedAttestationsResponse(src *v2.ListIndexedAttestationsResponse) *v1Alpha1.ListIndexedAttestationsResponse {
	if src == nil {
		return &v1Alpha1.ListIndexedAttestationsResponse{}
	}
	return &v1Alpha1.ListIndexedAttestationsResponse{

		IndexedAttestations: src.IndexedAttestations,
		NextPageToken:       src.NextPageToken,
		TotalSize:           src.TotalSize,
	}
}

// V1Alpha1ToV2ListValidatorAssignmentsRequest --
func V1Alpha1ToV2ListValidatorAssignmentsRequest(src *v1Alpha1.ListValidatorAssignmentsRequest) *v2.ListValidatorAssignmentsRequest {
	if src == nil {
		return &v2.ListValidatorAssignmentsRequest{}
	}
	return &v2.ListValidatorAssignmentsRequest{

		QueryFilter: src.QueryFilter,
		PublicKeys:  src.PublicKeys,
		Indices:     src.Indices,
		PageSize:    src.PageSize,
		PageToken:   src.PageToken,
	}
}

// V2ToV1Alpha1ListValidatorAssignmentsRequest --
func V2ToV1Alpha1ListValidatorAssignmentsRequest(src *v2.ListValidatorAssignmentsRequest) *v1Alpha1.ListValidatorAssignmentsRequest {
	if src == nil {
		return &v1Alpha1.ListValidatorAssignmentsRequest{}
	}
	return &v1Alpha1.ListValidatorAssignmentsRequest{

		QueryFilter: src.QueryFilter,
		PublicKeys:  src.PublicKeys,
		Indices:     src.Indices,
		PageSize:    src.PageSize,
		PageToken:   src.PageToken,
	}
}

// V1Alpha1ToV2ListValidatorBalancesRequest --
func V1Alpha1ToV2ListValidatorBalancesRequest(src *v1Alpha1.ListValidatorBalancesRequest) *v2.ListValidatorBalancesRequest {
	if src == nil {
		return &v2.ListValidatorBalancesRequest{}
	}
	return &v2.ListValidatorBalancesRequest{

		QueryFilter: src.QueryFilter,
		PublicKeys:  src.PublicKeys,
		Indices:     src.Indices,
		PageSize:    src.PageSize,
		PageToken:   src.PageToken,
	}
}

// V2ToV1Alpha1ListValidatorBalancesRequest --
func V2ToV1Alpha1ListValidatorBalancesRequest(src *v2.ListValidatorBalancesRequest) *v1Alpha1.ListValidatorBalancesRequest {
	if src == nil {
		return &v1Alpha1.ListValidatorBalancesRequest{}
	}
	return &v1Alpha1.ListValidatorBalancesRequest{

		QueryFilter: src.QueryFilter,
		PublicKeys:  src.PublicKeys,
		Indices:     src.Indices,
		PageSize:    src.PageSize,
		PageToken:   src.PageToken,
	}
}

// V1Alpha1ToV2ListValidatorsRequest --
func V1Alpha1ToV2ListValidatorsRequest(src *v1Alpha1.ListValidatorsRequest) *v2.ListValidatorsRequest {
	if src == nil {
		return &v2.ListValidatorsRequest{}
	}
	return &v2.ListValidatorsRequest{

		QueryFilter: src.QueryFilter,
		Active:      src.Active,
		PageSize:    src.PageSize,
		PageToken:   src.PageToken,
		PublicKeys:  src.PublicKeys,
		Indices:     src.Indices,
	}
}

// V2ToV1Alpha1ListValidatorsRequest --
func V2ToV1Alpha1ListValidatorsRequest(src *v2.ListValidatorsRequest) *v1Alpha1.ListValidatorsRequest {
	if src == nil {
		return &v1Alpha1.ListValidatorsRequest{}
	}
	return &v1Alpha1.ListValidatorsRequest{

		QueryFilter: src.QueryFilter,
		Active:      src.Active,
		PageSize:    src.PageSize,
		PageToken:   src.PageToken,
		PublicKeys:  src.PublicKeys,
		Indices:     src.Indices,
	}
}

// V1Alpha1ToV2ProposerSlashing --
func V1Alpha1ToV2ProposerSlashing(src *v1Alpha1.ProposerSlashing) *v2.ProposerSlashing {
	if src == nil {
		return &v2.ProposerSlashing{}
	}
	return &v2.ProposerSlashing{

		Header_1: V1Alpha1ToV2SignedBeaconBlockHeader(src.Header_1),
		Header_2: V1Alpha1ToV2SignedBeaconBlockHeader(src.Header_2),
	}
}

// V2ToV1Alpha1ProposerSlashing --
func V2ToV1Alpha1ProposerSlashing(src *v2.ProposerSlashing) *v1Alpha1.ProposerSlashing {
	if src == nil {
		return &v1Alpha1.ProposerSlashing{}
	}
	return &v1Alpha1.ProposerSlashing{

		Header_1: V2ToV1Alpha1SignedBeaconBlockHeader(src.Header_1),
		Header_2: V2ToV1Alpha1SignedBeaconBlockHeader(src.Header_2),
	}
}

// V1Alpha1ToV2SignedBeaconBlockHeader --
func V1Alpha1ToV2SignedBeaconBlockHeader(src *v1Alpha1.SignedBeaconBlockHeader) *v2.SignedBeaconBlockHeader {
	if src == nil {
		return &v2.SignedBeaconBlockHeader{}
	}
	return &v2.SignedBeaconBlockHeader{

		Header:    src.Header,
		Signature: src.Signature,
	}
}

// V2ToV1Alpha1SignedBeaconBlockHeader --
func V2ToV1Alpha1SignedBeaconBlockHeader(src *v2.SignedBeaconBlockHeader) *v1Alpha1.SignedBeaconBlockHeader {
	if src == nil {
		return &v1Alpha1.SignedBeaconBlockHeader{}
	}
	return &v1Alpha1.SignedBeaconBlockHeader{

		Header:    src.Header,
		Signature: src.Signature,
	}
}

// V1Alpha1ToV2SubmitSlashingResponse --
func V1Alpha1ToV2SubmitSlashingResponse(src *v1Alpha1.SubmitSlashingResponse) *v2.SubmitSlashingResponse {
	if src == nil {
		return &v2.SubmitSlashingResponse{}
	}
	return &v2.SubmitSlashingResponse{

		SlashedIndices: src.SlashedIndices,
	}
}

// V2ToV1Alpha1SubmitSlashingResponse --
func V2ToV1Alpha1SubmitSlashingResponse(src *v2.SubmitSlashingResponse) *v1Alpha1.SubmitSlashingResponse {
	if src == nil {
		return &v1Alpha1.SubmitSlashingResponse{}
	}
	return &v1Alpha1.SubmitSlashingResponse{

		SlashedIndices: src.SlashedIndices,
	}
}

// V1Alpha1ToV2Validator --
func V1Alpha1ToV2Validator(src *v1Alpha1.Validator) *v2.Validator {
	if src == nil {
		return &v2.Validator{}
	}
	return &v2.Validator{

		PublicKey:                  src.PublicKey,
		WithdrawalCredentials:      src.WithdrawalCredentials,
		EffectiveBalance:           src.EffectiveBalance,
		Slashed:                    src.Slashed,
		ActivationEligibilityEpoch: src.ActivationEligibilityEpoch,
		ActivationEpoch:            src.ActivationEpoch,
		ExitEpoch:                  src.ExitEpoch,
		WithdrawableEpoch:          src.WithdrawableEpoch,
	}
}

// V2ToV1Alpha1Validator --
func V2ToV1Alpha1Validator(src *v2.Validator) *v1Alpha1.Validator {
	if src == nil {
		return &v1Alpha1.Validator{}
	}
	return &v1Alpha1.Validator{

		PublicKey:                  src.PublicKey,
		WithdrawalCredentials:      src.WithdrawalCredentials,
		EffectiveBalance:           src.EffectiveBalance,
		Slashed:                    src.Slashed,
		ActivationEligibilityEpoch: src.ActivationEligibilityEpoch,
		ActivationEpoch:            src.ActivationEpoch,
		ExitEpoch:                  src.ExitEpoch,
		WithdrawableEpoch:          src.WithdrawableEpoch,
	}
}

// V1Alpha1ToV2ValidatorAssignments --
func V1Alpha1ToV2ValidatorAssignments(src *v1Alpha1.ValidatorAssignments) *v2.ValidatorAssignments {
	if src == nil {
		return &v2.ValidatorAssignments{}
	}
	return &v2.ValidatorAssignments{

		Epoch:         src.Epoch,
		Assignments:   src.Assignments,
		NextPageToken: src.NextPageToken,
		TotalSize:     src.TotalSize,
	}
}

// V2ToV1Alpha1ValidatorAssignments --
func V2ToV1Alpha1ValidatorAssignments(src *v2.ValidatorAssignments) *v1Alpha1.ValidatorAssignments {
	if src == nil {
		return &v1Alpha1.ValidatorAssignments{}
	}
	return &v1Alpha1.ValidatorAssignments{

		Epoch:         src.Epoch,
		Assignments:   src.Assignments,
		NextPageToken: src.NextPageToken,
		TotalSize:     src.TotalSize,
	}
}

// V1Alpha1ToV2ValidatorBalances --
func V1Alpha1ToV2ValidatorBalances(src *v1Alpha1.ValidatorBalances) *v2.ValidatorBalances {
	if src == nil {
		return &v2.ValidatorBalances{}
	}
	return &v2.ValidatorBalances{

		Epoch:         src.Epoch,
		Balances:      src.Balances,
		NextPageToken: src.NextPageToken,
		TotalSize:     src.TotalSize,
	}
}

// V2ToV1Alpha1ValidatorBalances --
func V2ToV1Alpha1ValidatorBalances(src *v2.ValidatorBalances) *v1Alpha1.ValidatorBalances {
	if src == nil {
		return &v1Alpha1.ValidatorBalances{}
	}
	return &v1Alpha1.ValidatorBalances{

		Epoch:         src.Epoch,
		Balances:      src.Balances,
		NextPageToken: src.NextPageToken,
		TotalSize:     src.TotalSize,
	}
}

// V1Alpha1ToV2ValidatorParticipation --
func V1Alpha1ToV2ValidatorParticipation(src *v1Alpha1.ValidatorParticipation) *v2.ValidatorParticipation {
	if src == nil {
		return &v2.ValidatorParticipation{}
	}
	return &v2.ValidatorParticipation{

		GlobalParticipationRate:          src.GlobalParticipationRate,
		VotedEther:                       src.VotedEther,
		EligibleEther:                    src.EligibleEther,
		CurrentEpochActiveGwei:           src.CurrentEpochActiveGwei,
		CurrentEpochAttestingGwei:        src.CurrentEpochAttestingGwei,
		CurrentEpochTargetAttestingGwei:  src.CurrentEpochTargetAttestingGwei,
		PreviousEpochActiveGwei:          src.PreviousEpochActiveGwei,
		PreviousEpochAttestingGwei:       src.PreviousEpochAttestingGwei,
		PreviousEpochTargetAttestingGwei: src.PreviousEpochTargetAttestingGwei,
		PreviousEpochHeadAttestingGwei:   src.PreviousEpochHeadAttestingGwei,
	}
}

// V2ToV1Alpha1ValidatorParticipation --
func V2ToV1Alpha1ValidatorParticipation(src *v2.ValidatorParticipation) *v1Alpha1.ValidatorParticipation {
	if src == nil {
		return &v1Alpha1.ValidatorParticipation{}
	}
	return &v1Alpha1.ValidatorParticipation{

		GlobalParticipationRate:          src.GlobalParticipationRate,
		VotedEther:                       src.VotedEther,
		EligibleEther:                    src.EligibleEther,
		CurrentEpochActiveGwei:           src.CurrentEpochActiveGwei,
		CurrentEpochAttestingGwei:        src.CurrentEpochAttestingGwei,
		CurrentEpochTargetAttestingGwei:  src.CurrentEpochTargetAttestingGwei,
		PreviousEpochActiveGwei:          src.PreviousEpochActiveGwei,
		PreviousEpochAttestingGwei:       src.PreviousEpochAttestingGwei,
		PreviousEpochTargetAttestingGwei: src.PreviousEpochTargetAttestingGwei,
		PreviousEpochHeadAttestingGwei:   src.PreviousEpochHeadAttestingGwei,
	}
}

// V1Alpha1ToV2ValidatorParticipationResponse --
func V1Alpha1ToV2ValidatorParticipationResponse(src *v1Alpha1.ValidatorParticipationResponse) *v2.ValidatorParticipationResponse {
	if src == nil {
		return &v2.ValidatorParticipationResponse{}
	}
	return &v2.ValidatorParticipationResponse{

		Epoch:         src.Epoch,
		Finalized:     src.Finalized,
		Participation: V1Alpha1ToV2ValidatorParticipation(src.Participation),
	}
}

// V2ToV1Alpha1ValidatorParticipationResponse --
func V2ToV1Alpha1ValidatorParticipationResponse(src *v2.ValidatorParticipationResponse) *v1Alpha1.ValidatorParticipationResponse {
	if src == nil {
		return &v1Alpha1.ValidatorParticipationResponse{}
	}
	return &v1Alpha1.ValidatorParticipationResponse{

		Epoch:         src.Epoch,
		Finalized:     src.Finalized,
		Participation: V2ToV1Alpha1ValidatorParticipation(src.Participation),
	}
}

// V1Alpha1ToV2ValidatorPerformanceRequest --
func V1Alpha1ToV2ValidatorPerformanceRequest(src *v1Alpha1.ValidatorPerformanceRequest) *v2.ValidatorPerformanceRequest {
	if src == nil {
		return &v2.ValidatorPerformanceRequest{}
	}
	return &v2.ValidatorPerformanceRequest{

		PublicKeys: src.PublicKeys,
		Indices:    src.Indices,
	}
}

// V2ToV1Alpha1ValidatorPerformanceRequest --
func V2ToV1Alpha1ValidatorPerformanceRequest(src *v2.ValidatorPerformanceRequest) *v1Alpha1.ValidatorPerformanceRequest {
	if src == nil {
		return &v1Alpha1.ValidatorPerformanceRequest{}
	}
	return &v1Alpha1.ValidatorPerformanceRequest{

		PublicKeys: src.PublicKeys,
		Indices:    src.Indices,
	}
}

// V1Alpha1ToV2ValidatorPerformanceResponse --
func V1Alpha1ToV2ValidatorPerformanceResponse(src *v1Alpha1.ValidatorPerformanceResponse) *v2.ValidatorPerformanceResponse {
	if src == nil {
		return &v2.ValidatorPerformanceResponse{}
	}
	return &v2.ValidatorPerformanceResponse{

		CurrentEffectiveBalances:      src.CurrentEffectiveBalances,
		InclusionSlots:                src.InclusionSlots,
		InclusionDistances:            src.InclusionDistances,
		CorrectlyVotedSource:          src.CorrectlyVotedSource,
		CorrectlyVotedTarget:          src.CorrectlyVotedTarget,
		CorrectlyVotedHead:            src.CorrectlyVotedHead,
		BalancesBeforeEpochTransition: src.BalancesBeforeEpochTransition,
		BalancesAfterEpochTransition:  src.BalancesAfterEpochTransition,
		MissingValidators:             src.MissingValidators,
		AverageActiveValidatorBalance: src.AverageActiveValidatorBalance,
		PublicKeys:                    src.PublicKeys,
	}
}

// V2ToV1Alpha1ValidatorPerformanceResponse --
func V2ToV1Alpha1ValidatorPerformanceResponse(src *v2.ValidatorPerformanceResponse) *v1Alpha1.ValidatorPerformanceResponse {
	if src == nil {
		return &v1Alpha1.ValidatorPerformanceResponse{}
	}
	return &v1Alpha1.ValidatorPerformanceResponse{

		CurrentEffectiveBalances:      src.CurrentEffectiveBalances,
		InclusionSlots:                src.InclusionSlots,
		InclusionDistances:            src.InclusionDistances,
		CorrectlyVotedSource:          src.CorrectlyVotedSource,
		CorrectlyVotedTarget:          src.CorrectlyVotedTarget,
		CorrectlyVotedHead:            src.CorrectlyVotedHead,
		BalancesBeforeEpochTransition: src.BalancesBeforeEpochTransition,
		BalancesAfterEpochTransition:  src.BalancesAfterEpochTransition,
		MissingValidators:             src.MissingValidators,
		AverageActiveValidatorBalance: src.AverageActiveValidatorBalance,
		PublicKeys:                    src.PublicKeys,
	}
}

// V1Alpha1ToV2ValidatorQueue --
func V1Alpha1ToV2ValidatorQueue(src *v1Alpha1.ValidatorQueue) *v2.ValidatorQueue {
	if src == nil {
		return &v2.ValidatorQueue{}
	}
	return &v2.ValidatorQueue{

		ChurnLimit:                 src.ChurnLimit,
		ActivationPublicKeys:       src.ActivationPublicKeys,
		ExitPublicKeys:             src.ExitPublicKeys,
		ActivationValidatorIndices: src.ActivationValidatorIndices,
		ExitValidatorIndices:       src.ExitValidatorIndices,
	}
}

// V2ToV1Alpha1ValidatorQueue --
func V2ToV1Alpha1ValidatorQueue(src *v2.ValidatorQueue) *v1Alpha1.ValidatorQueue {
	if src == nil {
		return &v1Alpha1.ValidatorQueue{}
	}
	return &v1Alpha1.ValidatorQueue{

		ChurnLimit:                 src.ChurnLimit,
		ActivationPublicKeys:       src.ActivationPublicKeys,
		ExitPublicKeys:             src.ExitPublicKeys,
		ActivationValidatorIndices: src.ActivationValidatorIndices,
		ExitValidatorIndices:       src.ExitValidatorIndices,
	}
}

// V1Alpha1ToV2Validators --
func V1Alpha1ToV2Validators(src *v1Alpha1.Validators) *v2.Validators {
	if src == nil {
		return &v2.Validators{}
	}
	return &v2.Validators{

		Epoch:         src.Epoch,
		ValidatorList: src.ValidatorList,
		NextPageToken: src.NextPageToken,
		TotalSize:     src.TotalSize,
	}
}

// V2ToV1Alpha1Validators --
func V2ToV1Alpha1Validators(src *v2.Validators) *v1Alpha1.Validators {
	if src == nil {
		return &v1Alpha1.Validators{}
	}
	return &v1Alpha1.Validators{

		Epoch:         src.Epoch,
		ValidatorList: src.ValidatorList,
		NextPageToken: src.NextPageToken,
		TotalSize:     src.TotalSize,
	}
}

// V1Alpha1ToV2WeakSubjectivityCheckpoint --
func V1Alpha1ToV2WeakSubjectivityCheckpoint(src *v1Alpha1.WeakSubjectivityCheckpoint) *v2.WeakSubjectivityCheckpoint {
	if src == nil {
		return &v2.WeakSubjectivityCheckpoint{}
	}
	return &v2.WeakSubjectivityCheckpoint{

		BlockRoot: src.BlockRoot,
		StateRoot: src.StateRoot,
		Epoch:     src.Epoch,
	}
}

// V2ToV1Alpha1WeakSubjectivityCheckpoint --
func V2ToV1Alpha1WeakSubjectivityCheckpoint(src *v2.WeakSubjectivityCheckpoint) *v1Alpha1.WeakSubjectivityCheckpoint {
	if src == nil {
		return &v1Alpha1.WeakSubjectivityCheckpoint{}
	}
	return &v1Alpha1.WeakSubjectivityCheckpoint{

		BlockRoot: src.BlockRoot,
		StateRoot: src.StateRoot,
		Epoch:     src.Epoch,
	}
}
