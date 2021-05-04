package gateway

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// JSONPbHex is a Marshaler which marshals/unmarshals into/from JSON
// with the "google.golang.org/protobuf/encoding/protojson" marshaler.
// It supports the full functionality of protobuf unlike JSONBuiltin.
//
// The NewDecoder method returns a DecoderWrapper, so the underlying
// *json.Decoder methods can be used.
type JSONPbHex struct {
	protojson.MarshalOptions
	protojson.UnmarshalOptions
}

// ContentType always returns "application/json".
func (*JSONPbHex) ContentType(_ interface{}) string {
	return "application/json"
}

// Marshal marshals "v" into JSON.
func (j *JSONPbHex) Marshal(v interface{}) ([]byte, error) {
	//fmt.Println("Being used marshal")
	//fmt.Println(v)
	jsonPb := &gwruntime.JSONPb{
		MarshalOptions:   j.MarshalOptions,
		UnmarshalOptions: j.UnmarshalOptions,
	}
	marshaledBytes, err := jsonPb.Marshal(v)
	if err != nil {
		return nil, err
	}

	//fmt.Println("marshaled")
	//fmt.Println(string(marshaledBytes))
	r, err := j.convertJsonBytesBase64(marshaledBytes)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (j *JSONPbHex) convertJsonBytesBase64(marshaledBytes []byte) ([]byte, error) {
	var newV interface{}
	if err := json.Unmarshal(marshaledBytes, &newV); err != nil {
		return nil, err
	}
	convertBase64(newV)
	//fmt.Println("converted")
	r, err := json.Marshal(newV)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (j *JSONPbHex) marshalTo(w io.Writer, v interface{}) error {
	jsonPb := &gwruntime.JSONPb{
		MarshalOptions:   j.MarshalOptions,
		UnmarshalOptions: j.UnmarshalOptions,
	}
	marshaledBytes, err := jsonPb.Marshal(v)
	if err != nil {
		return err
	}
	r, err := j.convertJsonBytesBase64(marshaledBytes)
	if err != nil {
		return err
	}
	_, err = w.Write(r)
	return err
}

// Unmarshal unmarshals JSON "data" into "v"
func (j *JSONPbHex) Unmarshal(data []byte, v interface{}) error {
	fmt.Println("Being used unmarshal")

	jsonPb := &gwruntime.JSONPb{
		MarshalOptions:   j.MarshalOptions,
		UnmarshalOptions: j.UnmarshalOptions,
	}
	var newV interface{}
	if err := json.Unmarshal(data, &newV); err != nil {
		return err
	}
	convertHex(newV)
	r, err := json.MarshalIndent(newV, "", j.Indent)
	if err != nil {
		return err
	}
	return jsonPb.Unmarshal(r, v)
}

// NewDecoder returns a Decoder which reads JSON stream from "r".
func (j *JSONPbHex) NewDecoder(r io.Reader) gwruntime.Decoder {
	d := json.NewDecoder(r)
	return DecoderWrapper{
		Decoder:          d,
		UnmarshalOptions: j.UnmarshalOptions,
	}
}

// DecoderWrapper is a wrapper around a *json.Decoder that adds
// support for protos to the Decode method.
type DecoderWrapper struct {
	*json.Decoder
	protojson.UnmarshalOptions
}

// Decode wraps the embedded decoder's Decode method to support
// protos using a jsonpb.Unmarshaler.
func (d DecoderWrapper) Decode(v interface{}) error {
	fmt.Println("Being used decode")

	return decodeJSONPb(d.Decoder, d.UnmarshalOptions, v)
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (j *JSONPbHex) NewEncoder(w io.Writer) gwruntime.Encoder {
	return gwruntime.EncoderFunc(func(v interface{}) error {
		fmt.Println("Being used encode")

		if err := j.marshalTo(w, v); err != nil {
			return err
		}
		// mimic json.Encoder by adding a newline (makes output
		// easier to read when it contains multiple encoded items)
		_, err := w.Write(j.Delimiter())
		return err
	})
}

func unmarshalJSONPb(data []byte, unmarshaler protojson.UnmarshalOptions, v interface{}) error {
	d := json.NewDecoder(bytes.NewReader(data))
	return decodeJSONPb(d, unmarshaler, v)
}

func decodeJSONPb(d *json.Decoder, unmarshaler protojson.UnmarshalOptions, v interface{}) error {
	p, ok := v.(proto.Message)
	if !ok {
		return decodeNonProtoField(d, unmarshaler, v)
	}

	// Decode into bytes for marshalling
	var b json.RawMessage
	err := d.Decode(&b)
	if err != nil {
		return err
	}

	fmt.Println("penis")
	fmt.Println(&b)
	var newV interface{}
	if err := json.Unmarshal(b, &newV); err != nil {
		return err
	}
	convertHex(newV)
	r, err := json.Marshal(newV)
	if err != nil {
		return err
	}

	return unmarshaler.Unmarshal(r, p)
}

func decodeNonProtoField(d *json.Decoder, unmarshaler protojson.UnmarshalOptions, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("%T is not a pointer", v)
	}
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		if rv.Type().ConvertibleTo(typeProtoMessage) {
			// Decode into bytes for marshalling
			var b json.RawMessage
			err := d.Decode(&b)
			if err != nil {
				return err
			}
			var newV interface{}
			if err := json.Unmarshal(b, &newV); err != nil {
				return err
			}
			convertHex(newV)
			r, err := json.Marshal(newV)
			if err != nil {
				return err
			}

			return unmarshaler.Unmarshal(r, rv.Interface().(proto.Message))
		}
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Map {
		if rv.IsNil() {
			rv.Set(reflect.MakeMap(rv.Type()))
		}
		conv, ok := convFromType[rv.Type().Key().Kind()]
		if !ok {
			return fmt.Errorf("unsupported type of map field key: %v", rv.Type().Key())
		}

		m := make(map[string]*json.RawMessage)
		if err := d.Decode(&m); err != nil {
			return err
		}
		for k, v := range m {
			result := conv.Call([]reflect.Value{reflect.ValueOf(k)})
			if err := result[1].Interface(); err != nil {
				return err.(error)
			}
			bk := result[0]
			bv := reflect.New(rv.Type().Elem())
			if err := unmarshalJSONPb([]byte(*v), unmarshaler, bv.Interface()); err != nil {
				return err
			}
			rv.SetMapIndex(bk, bv.Elem())
		}
		return nil
	}
	if rv.Kind() == reflect.Slice {
		var sl []json.RawMessage
		if err := d.Decode(&sl); err != nil {
			return err
		}
		if sl != nil {
			rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
		}
		for _, item := range sl {
			bv := reflect.New(rv.Type().Elem())
			if err := unmarshalJSONPb([]byte(item), unmarshaler, bv.Interface()); err != nil {
				return err
			}
			rv.Set(reflect.Append(rv, bv.Elem()))
		}
		return nil
	}
	if _, ok := rv.Interface().(protoEnum); ok {
		var repr interface{}
		if err := d.Decode(&repr); err != nil {
			return err
		}
		switch v := repr.(type) {
		case string:
			// TODO(yugui) Should use proto.StructProperties?
			return fmt.Errorf("unmarshaling of symbolic enum %q not supported: %T", repr, rv.Interface())
		case float64:
			rv.Set(reflect.ValueOf(int32(v)).Convert(rv.Type()))
			return nil
		default:
			return fmt.Errorf("cannot assign %#v into Go type %T", repr, rv.Interface())
		}
	}
	return d.Decode(v)
}

type protoEnum interface {
	fmt.Stringer
	EnumDescriptor() ([]byte, []int)
}

var typeProtoMessage = reflect.TypeOf((*proto.Message)(nil)).Elem()

// Delimiter for newline encoded JSON streams.
func (j *JSONPbHex) Delimiter() []byte {
	return []byte("\n")
}

var (
	convFromType = map[reflect.Kind]reflect.Value{
		reflect.String:  reflect.ValueOf(gwruntime.String),
		reflect.Bool:    reflect.ValueOf(gwruntime.Bool),
		reflect.Float64: reflect.ValueOf(gwruntime.Float64),
		reflect.Float32: reflect.ValueOf(gwruntime.Float32),
		reflect.Int64:   reflect.ValueOf(gwruntime.Int64),
		reflect.Int32:   reflect.ValueOf(gwruntime.Int32),
		reflect.Uint64:  reflect.ValueOf(gwruntime.Uint64),
		reflect.Uint32:  reflect.ValueOf(gwruntime.Uint32),
		reflect.Slice:   reflect.ValueOf(gwruntime.Bytes),
	}
)

func convertBase64(data interface{}) {
	fmt.Printf("%v\n", data)
	switch d := data.(type) {
	case map[string]interface{}:
		for k, v := range d {
			switch tv := v.(type) {
			case string:
				fmt.Println(tv)
				fmt.Printf("%s\n", tv)
				decoded, err := base64.StdEncoding.DecodeString(tv)
				if err == nil {
					fmt.Printf("%#x\n", decoded)
					fmt.Printf("%s\n", decoded)
					d[k] = "0x" + hex.EncodeToString(decoded)
				}
			case map[string]interface{}:
				convertBase64(tv)
			case []interface{}:
				convertBase64(tv)
				//case nil:
				//	delete(d, k)
			}
		}
	case []interface{}:
		if len(d) > 0 {
			switch d[0].(type) {
			case string:
				for i, s := range d {
					fmt.Println(s)
					stringed, ok := s.(string)
					if !ok {
						continue
					}
					decoded, err := base64.StdEncoding.DecodeString(stringed)
					if err == nil {
						d[i] = "0x" + hex.EncodeToString(decoded)
					}
				}
			case map[string]interface{}:
				for _, t := range d {
					convertBase64(t)
				}
			case []interface{}:
				for _, t := range d {
					convertBase64(t)
				}
			}
		}
	}
}

func convertHex(data interface{}) {
	switch d := data.(type) {
	case map[string]interface{}:
		for k, v := range d {
			switch tv := v.(type) {
			case string:
				fmt.Println(tv)
				if strings.Contains(tv, "0x") {
					noPrefix := tv[2:]
					decoded, err := hex.DecodeString(noPrefix)
					if err == nil {
						d[k] = base64.StdEncoding.EncodeToString(decoded)
					}
				}
			case map[string]interface{}:
				convertHex(tv)
			case []interface{}:
				convertHex(tv)
			}
		}
	case []interface{}:
		if len(d) > 0 {
			switch d[0].(type) {
			case string:
				for i, s := range d {
					fmt.Println(s)
					stringed, ok := s.(string)
					if !ok {
						continue
					}
					if strings.Contains(stringed, "0x") {
						noPrefix := stringed[3 : len(stringed)-1]
						decoded, err := hex.DecodeString(noPrefix)
						if err == nil {
							d[i] = base64.StdEncoding.EncodeToString(decoded)
						}
					}
				}
			case map[string]interface{}:
				for _, t := range d {
					convertHex(t)
				}
			case []interface{}:
				for _, t := range d {
					convertHex(t)
				}
			}
		}
	}
}
