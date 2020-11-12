package util

import (
	"bytes"

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
