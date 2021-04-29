package gateway

import (
	"bytes"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	ethpbv1 "github.com/prysmaticlabs/ethereumapis/eth/v1"
)

func TestJSONPbHex_MarshalUnmarshalPb(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    []byte
		wantErr bool
	}{
		{
			name: "attestation",
			input: &ethpbv1.Attestation{
				Data: &ethpbv1.AttestationData{
					Source: &ethpbv1.Checkpoint{
						Root: []byte("chicken"),
					},
					Target: &ethpbv1.Checkpoint{
						Root: []byte("pork"),
					},
					BeaconBlockRoot: []byte("beef"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JSONPbHex{}
			got, err := j.Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var i *ethpbv1.Attestation
			if err := j.Unmarshal(got, &i); err != nil {
				t.Error(err)
			}
			if !proto.Equal(tt.input.(*ethpbv1.Attestation), i) {
				t.Errorf("Marshal() got = %s, want %s", got, i)
			}
		})
	}
}

func TestJSONPbHex_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    []byte
		wantErr bool
	}{
		{
			name: "req",
			input: &ethpbv1.StateRequest{
				StateId: []byte("head"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JSONPbHex{}
			got, err := j.Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var i *ethpbv1.StateRequest
			if err := j.Unmarshal(got, &i); err != nil {
				t.Error(err)
			}
			if !proto.Equal(tt.input.(*ethpbv1.StateRequest), i) {
				t.Errorf("Marshal() got = %s, want %s", got, i)
			}
		})
	}
}

func TestJSONPbHex_Encoder(t *testing.T) {
	j := &JSONPbHex{}
	var buf bytes.Buffer
	enc := j.NewEncoder(&buf)
	var s = &ethpbv1.Attestation{Signature: []byte("beef")}
	if err := enc.Encode(&s); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "0x62656566") {
		t.Fatal(buf.String())
	}
}

func TestJSONPbHex_Decoder(t *testing.T) {
	j := &JSONPbHex{}
	var buf = strings.NewReader("{\"signature\":\"0x62656566\"}")
	dec := j.NewDecoder(buf)
	var s = &ethpbv1.Attestation{}
	if err := dec.Decode(&s); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(s, &ethpbv1.Attestation{Signature: []byte("beef")}) {
		t.Fatal(s)
	}
}
