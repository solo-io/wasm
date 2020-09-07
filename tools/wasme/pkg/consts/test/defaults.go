package test

import "github.com/solo-io/wasm/tools/wasme/pkg/consts"

var (
	IstioAssemblyScriptImage = consts.HubDomain + "/ilackarms/assemblyscript-test:" + Istio15Tag
	Istio15Tag               = "istio-1.5"
)
