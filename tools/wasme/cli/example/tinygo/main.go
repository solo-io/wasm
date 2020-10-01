package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetNewHttpContext(newContext)
}

type httpHeaders struct {
	proxywasm.DefaultContext
	contextID uint32
}

func newContext(contextID uint32) proxywasm.HttpContext {
	return &httpHeaders{contextID: contextID}
}

// override
func (ctx *httpHeaders) OnHttpRequestHeaders(int, bool) types.Action {
	hs, err := proxywasm.HostCallGetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
	}

	for _, h := range hs {
		proxywasm.LogInfof("request header: %s: %s", h[0], h[1])
	}
	return types.ActionContinue
}

// override
func (ctx *httpHeaders) OnHttpResponseHeaders(int, bool) types.Action {
	if err := proxywasm.HostCallSetHttpResponseHeader("hello", "world"); err != nil {
		proxywasm.LogCriticalf("failed to set response header: %v", err)
	}
	return types.ActionContinue
}

// override
func (ctx *httpHeaders) OnLog() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}
