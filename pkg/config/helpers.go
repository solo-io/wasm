package config

import (
	"io"
	"io/ioutil"

	"github.com/solo-io/go-utils/protoutils"
)

func FromBytes(b []byte) (*Config, error) {
	var cfg Config
	return &cfg, protoutils.UnmarshalBytes(b, &cfg)
}

func FromReader(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return FromBytes(b)
}
