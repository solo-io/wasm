package deploy

import (
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/printers"
	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"io"
)

// used to back Dry-Run calls to gloo CRDs
type dryRunGatewayClient struct {
	// base gateway into which we inject the filter, for demo purposes
	defaultGateway *gatewayv1.Gateway

	// where to write the yaml
	out io.Writer
}

func newDryRunGatewayClient(out io.Writer) *dryRunGatewayClient {
	return &dryRunGatewayClient{out: out, defaultGateway: defaults.DefaultGateway(gloodefaults.GlooSystem)}
}

func (c *dryRunGatewayClient) Write(resource *gatewayv1.Gateway, opts clients.WriteOpts) (*gatewayv1.Gateway, error) {
	return resource, printers.PrintKubeCrd(resource, gatewayv1.GatewayCrd)
}

func (c *dryRunGatewayClient) List(namespace string, opts clients.ListOpts) (gatewayv1.GatewayList, error) {
	return gatewayv1.GatewayList{c.defaultGateway}, nil
}

func (c *dryRunGatewayClient) BaseClient() clients.ResourceClient {
	panic("not implemented")
}

func (c *dryRunGatewayClient) Register() error {
	panic("not implemented")
}

func (c *dryRunGatewayClient) Read(namespace, name string, opts clients.ReadOpts) (*gatewayv1.Gateway, error) {
	panic("not implemented")
}

func (c *dryRunGatewayClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	panic("not implemented")
}

func (c *dryRunGatewayClient) Watch(namespace string, opts clients.WatchOpts) (<-chan gatewayv1.GatewayList, <-chan error, error) {
	panic("not implemented")
}

var _ gatewayv1.GatewayClient = &dryRunGatewayClient{}
