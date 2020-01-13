// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: wasme/v1/filter_deployment.proto

package v1

import (
	bytes "bytes"
	fmt "fmt"
	math "math"

	github_com_gogo_protobuf_jsonpb "github.com/gogo/protobuf/jsonpb"
	proto "github.com/gogo/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// MarshalJSON is a custom marshaler for FilterDeploymentSpec
func (this *FilterDeploymentSpec) MarshalJSON() ([]byte, error) {
	str, err := FilterDeploymentMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for FilterDeploymentSpec
func (this *FilterDeploymentSpec) UnmarshalJSON(b []byte) error {
	return FilterDeploymentUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

// MarshalJSON is a custom marshaler for FilterSpec
func (this *FilterSpec) MarshalJSON() ([]byte, error) {
	str, err := FilterDeploymentMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for FilterSpec
func (this *FilterSpec) UnmarshalJSON(b []byte) error {
	return FilterDeploymentUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

// MarshalJSON is a custom marshaler for DeploymentSpec
func (this *DeploymentSpec) MarshalJSON() ([]byte, error) {
	str, err := FilterDeploymentMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for DeploymentSpec
func (this *DeploymentSpec) UnmarshalJSON(b []byte) error {
	return FilterDeploymentUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

// MarshalJSON is a custom marshaler for IstioDeploymentSpec
func (this *IstioDeploymentSpec) MarshalJSON() ([]byte, error) {
	str, err := FilterDeploymentMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for IstioDeploymentSpec
func (this *IstioDeploymentSpec) UnmarshalJSON(b []byte) error {
	return FilterDeploymentUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

// MarshalJSON is a custom marshaler for FilterDeploymentStatus
func (this *FilterDeploymentStatus) MarshalJSON() ([]byte, error) {
	str, err := FilterDeploymentMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for FilterDeploymentStatus
func (this *FilterDeploymentStatus) UnmarshalJSON(b []byte) error {
	return FilterDeploymentUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

// MarshalJSON is a custom marshaler for WorkloadStatus
func (this *WorkloadStatus) MarshalJSON() ([]byte, error) {
	str, err := FilterDeploymentMarshaler.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for WorkloadStatus
func (this *WorkloadStatus) UnmarshalJSON(b []byte) error {
	return FilterDeploymentUnmarshaler.Unmarshal(bytes.NewReader(b), this)
}

var (
	FilterDeploymentMarshaler   = &github_com_gogo_protobuf_jsonpb.Marshaler{}
	FilterDeploymentUnmarshaler = &github_com_gogo_protobuf_jsonpb.Unmarshaler{}
)
