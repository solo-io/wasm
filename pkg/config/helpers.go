package config

import (
	"io"
	"io/ioutil"

	"github.com/solo-io/go-utils/protoutils"
)

func (cfg *Runtime) ToBytes() ([]byte, error) {
	return protoutils.MarshalBytes(cfg)
}

func FromBytes(b []byte) (*Runtime, error) {
	var cfg Runtime
	return &cfg, protoutils.UnmarshalBytes(b, &cfg)
}

func FromReader(r io.Reader) (*Runtime, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return FromBytes(b)
}
