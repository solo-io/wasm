// This file contains generated Deepcopy methods for
// Protobuf types with oneofs

package v1

import (
	fmt "fmt"

	proto "github.com/gogo/protobuf/proto"

	math "math"

	_ "github.com/gogo/protobuf/gogoproto"

	_ "github.com/gogo/protobuf/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// DeepCopyInto supports using AttributeManifest within kubernetes types, where deepcopy-gen is used.
func (in *DeploymentSpec) DeepCopyInto(out *DeploymentSpec) {
	p := proto.Clone(in).(*DeploymentSpec)
	*out = *p
}
