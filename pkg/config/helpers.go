package config

import "github.com/solo-io/go-utils/protoutils"

func FromBytes(b []byte) (*Config, error) {
	var cfg Config
	return &cfg, protoutils.UnmarshalBytes(b, &cfg)
}
