package util

import (
	"bytes"
	gogojsonpb "github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/types"
	"github.com/rotisserie/eris"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	structpb "github.com/golang/protobuf/ptypes/struct"
)

var (
	jsonpbMarshaler = &jsonpb.Marshaler{OrigName: false}
)

func MarshalStruct(m proto.Message) (*structpb.Struct, error) {
	data, err := MarshalBytes(m)
	if err != nil {
		return nil, err
	}
	var pb structpb.Struct
	err = jsonpb.UnmarshalString(string(data), &pb)
	return &pb, err
}

func MarshalGogoStruct(msg proto.Message) (*types.Struct, error) {
	if msg == nil {
		return nil, eris.New("nil message")
	}

	buf := &bytes.Buffer{}
	if err := (&jsonpb.Marshaler{OrigName: true}).Marshal(buf, msg); err != nil {
		return nil, err
	}

	pbs := &types.Struct{}
	if err := gogojsonpb.Unmarshal(buf, pbs); err != nil {
		return nil, err
	}

	return pbs, nil
}

func UnmarshalStruct(s *structpb.Struct, m proto.Message) error {
	data, err := MarshalBytes(s)
	if err != nil {
		return err
	}
	return jsonpb.UnmarshalString(string(data), m)
}

func MarshalBytes(pb proto.Message) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := jsonpbMarshaler.Marshal(buf, pb)
	return buf.Bytes(), err
}

func UnmarshalBytes(b []byte, pb proto.Message) error {
	return jsonpb.Unmarshal(bytes.NewBuffer(b), pb)
}

func StructPbToGogo(structuredData *structpb.Struct) (*types.Struct, error) {
	if structuredData == nil {
		return nil, eris.New("cannot unmarshal nil struct")
	}
	byt, err := proto.Marshal(structuredData)
	if err != nil {
		return nil, err
	}
	var st types.Struct
	if err := proto.Unmarshal(byt, &st); err != nil {
		return nil, err
	}
	return &st, nil
}
