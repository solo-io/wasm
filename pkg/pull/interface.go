package pull

import "io"

type FilterConfig interface {
	RootId() string
	Schema() []byte
}

type Filter interface {
	Code() io.ReadCloser
	Configs() []FilterConfig
}
