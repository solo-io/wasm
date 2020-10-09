package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

// Other examples can be found at https://github.com/tetratelabs/proxy-wasm-go-sdk/tree/v0.0.8/examples

func main() {
	proxywasm.SetNewHttpContext(newHttpContext)
}

type httpHeaders struct {
	// we must embed the default context so that you need not to reimplement all the methods by yourself
	proxywasm.DefaultHttpContext
	contextID uint32
}

func newHttpContext(rootContextID, contextID uint32) proxywasm.HttpContext {
	return &httpHeaders{contextID: contextID}
}

// override
func (ctx *httpHeaders) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	hs, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
	}

	for _, h := range hs {
		proxywasm.LogInfof("request header: %s: %s", h[0], h[1])
	}
	return types.ActionContinue
}

// override
func (ctx *httpHeaders) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	if err := proxywasm.SetHttpResponseHeader("hello", "world"); err != nil {
		proxywasm.LogCriticalf("failed to set response header: %v", err)
	}
	return types.ActionContinue
}

// override
func (ctx *httpHeaders) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}
