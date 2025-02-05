package apimiddleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/gateway"
	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/prysm/shared/testutil/require"
)

func TestWrapAttestationArray(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		endpoint := gateway.Endpoint{
			PostRequest: &submitAttestationRequestJson{},
		}
		unwrappedAtts := []*attestationJson{{AggregationBits: "1010"}}
		unwrappedAttsJson, err := json.Marshal(unwrappedAtts)
		require.NoError(t, err)

		var body bytes.Buffer
		_, err = body.Write(unwrappedAttsJson)
		require.NoError(t, err)
		request := httptest.NewRequest("POST", "http://foo.example", &body)

		errJson := wrapAttestationsArray(endpoint, nil, request)
		require.Equal(t, true, errJson == nil)
		wrappedAtts := &submitAttestationRequestJson{}
		require.NoError(t, json.NewDecoder(request.Body).Decode(wrappedAtts))
		require.Equal(t, 1, len(wrappedAtts.Data), "wrong number of wrapped items")
		assert.Equal(t, "1010", wrappedAtts.Data[0].AggregationBits)
	})

	t.Run("invalid_body", func(t *testing.T) {
		endpoint := gateway.Endpoint{
			PostRequest: &submitAttestationRequestJson{},
		}
		var body bytes.Buffer
		_, err := body.Write([]byte("invalid"))
		require.NoError(t, err)
		request := httptest.NewRequest("POST", "http://foo.example", &body)

		errJson := wrapAttestationsArray(endpoint, nil, request)
		require.Equal(t, false, errJson == nil)
		assert.Equal(t, true, strings.Contains(errJson.Msg(), "could not decode body"))
		assert.Equal(t, http.StatusInternalServerError, errJson.StatusCode())
	})
}

func TestWrapValidatorIndicesArray(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		endpoint := gateway.Endpoint{
			PostRequest: &attesterDutiesRequestJson{},
		}
		unwrappedIndices := []string{"1", "2"}
		unwrappedIndicesJson, err := json.Marshal(unwrappedIndices)
		require.NoError(t, err)

		var body bytes.Buffer
		_, err = body.Write(unwrappedIndicesJson)
		require.NoError(t, err)
		request := httptest.NewRequest("POST", "http://foo.example", &body)

		errJson := wrapValidatorIndicesArray(endpoint, nil, request)
		require.Equal(t, true, errJson == nil)
		wrappedIndices := &attesterDutiesRequestJson{}
		require.NoError(t, json.NewDecoder(request.Body).Decode(wrappedIndices))
		require.Equal(t, 2, len(wrappedIndices.Index), "wrong number of wrapped items")
		assert.Equal(t, "1", wrappedIndices.Index[0])
		assert.Equal(t, "2", wrappedIndices.Index[1])
	})
}

func TestWrapSignedAggregateAndProofArray(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		endpoint := gateway.Endpoint{
			PostRequest: &submitAggregateAndProofsRequestJson{},
		}
		unwrappedAggs := []*signedAggregateAttestationAndProofJson{{Signature: "sig"}}
		unwrappedAggsJson, err := json.Marshal(unwrappedAggs)
		require.NoError(t, err)

		var body bytes.Buffer
		_, err = body.Write(unwrappedAggsJson)
		require.NoError(t, err)
		request := httptest.NewRequest("POST", "http://foo.example", &body)

		errJson := wrapSignedAggregateAndProofArray(endpoint, nil, request)
		require.Equal(t, true, errJson == nil)
		wrappedAggss := &submitAggregateAndProofsRequestJson{}
		require.NoError(t, json.NewDecoder(request.Body).Decode(wrappedAggss))
		require.Equal(t, 1, len(wrappedAggss.Data), "wrong number of wrapped items")
		assert.Equal(t, "sig", wrappedAggss.Data[0].Signature)
	})

	t.Run("invalid_body", func(t *testing.T) {
		endpoint := gateway.Endpoint{
			PostRequest: &submitAggregateAndProofsRequestJson{},
		}
		var body bytes.Buffer
		_, err := body.Write([]byte("invalid"))
		require.NoError(t, err)
		request := httptest.NewRequest("POST", "http://foo.example", &body)

		errJson := wrapSignedAggregateAndProofArray(endpoint, nil, request)
		require.Equal(t, false, errJson == nil)
		assert.Equal(t, true, strings.Contains(errJson.Msg(), "could not decode body"))
		assert.Equal(t, http.StatusInternalServerError, errJson.StatusCode())
	})
}

func TestWrapBeaconCommitteeSubscriptionsArray(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		endpoint := gateway.Endpoint{
			PostRequest: &submitBeaconCommitteeSubscriptionsRequestJson{},
		}
		unwrappedSubs := []*beaconCommitteeSubscribeJson{{
			ValidatorIndex:   "1",
			CommitteeIndex:   "1",
			CommitteesAtSlot: "1",
			Slot:             "1",
			IsAggregator:     true,
		}}
		unwrappedSubsJson, err := json.Marshal(unwrappedSubs)
		require.NoError(t, err)

		var body bytes.Buffer
		_, err = body.Write(unwrappedSubsJson)
		require.NoError(t, err)
		request := httptest.NewRequest("POST", "http://foo.example", &body)

		errJson := wrapBeaconCommitteeSubscriptionsArray(endpoint, nil, request)
		require.Equal(t, true, errJson == nil)
		wrappedAggss := &submitBeaconCommitteeSubscriptionsRequestJson{}
		require.NoError(t, json.NewDecoder(request.Body).Decode(wrappedAggss))
		require.Equal(t, 1, len(wrappedAggss.Data), "wrong number of wrapped items")
		assert.Equal(t, "1", wrappedAggss.Data[0].ValidatorIndex)
		assert.Equal(t, "1", wrappedAggss.Data[0].CommitteeIndex)
		assert.Equal(t, "1", wrappedAggss.Data[0].CommitteesAtSlot)
		assert.Equal(t, "1", wrappedAggss.Data[0].Slot)
		assert.Equal(t, true, wrappedAggss.Data[0].IsAggregator)
	})

	t.Run("invalid_body", func(t *testing.T) {
		endpoint := gateway.Endpoint{
			PostRequest: &submitBeaconCommitteeSubscriptionsRequestJson{},
		}
		var body bytes.Buffer
		_, err := body.Write([]byte("invalid"))
		require.NoError(t, err)
		request := httptest.NewRequest("POST", "http://foo.example", &body)

		errJson := wrapBeaconCommitteeSubscriptionsArray(endpoint, nil, request)
		require.Equal(t, false, errJson == nil)
		assert.Equal(t, true, strings.Contains(errJson.Msg(), "could not decode body"))
		assert.Equal(t, http.StatusInternalServerError, errJson.StatusCode())
	})
}

func TestPrepareGraffiti(t *testing.T) {
	endpoint := gateway.Endpoint{
		PostRequest: &beaconBlockContainerJson{
			Message: &beaconBlockJson{
				Body: &beaconBlockBodyJson{},
			},
		},
	}

	t.Run("32_bytes", func(t *testing.T) {
		endpoint.PostRequest.(*beaconBlockContainerJson).Message.Body.Graffiti = string(bytesutil.PadTo([]byte("foo"), 32))

		prepareGraffiti(endpoint, nil, nil)
		assert.Equal(
			t,
			"0x666f6f0000000000000000000000000000000000000000000000000000000000",
			endpoint.PostRequest.(*beaconBlockContainerJson).Message.Body.Graffiti,
		)
	})

	t.Run("less_than_32_bytes", func(t *testing.T) {
		endpoint.PostRequest.(*beaconBlockContainerJson).Message.Body.Graffiti = "foo"

		prepareGraffiti(endpoint, nil, nil)
		assert.Equal(
			t,
			"0x666f6f0000000000000000000000000000000000000000000000000000000000",
			endpoint.PostRequest.(*beaconBlockContainerJson).Message.Body.Graffiti,
		)
	})

	t.Run("more_than_32_bytes", func(t *testing.T) {
		endpoint.PostRequest.(*beaconBlockContainerJson).Message.Body.Graffiti = string(bytesutil.PadTo([]byte("foo"), 33))

		prepareGraffiti(endpoint, nil, nil)
		assert.Equal(
			t,
			"0x666f6f0000000000000000000000000000000000000000000000000000000000",
			endpoint.PostRequest.(*beaconBlockContainerJson).Message.Body.Graffiti,
		)
	})
}
