// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/solo-io/wasm/tools/wasme/cli/operator/api/wasme/v1/filter_deployment.proto

package v1

import (
	fmt "fmt"
	math "math"

	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// the state of the filter deployment
type WorkloadStatus_State int32

const (
	WorkloadStatus_Pending   WorkloadStatus_State = 0
	WorkloadStatus_Succeeded WorkloadStatus_State = 1
	WorkloadStatus_Failed    WorkloadStatus_State = 2
)

var WorkloadStatus_State_name = map[int32]string{
	0: "Pending",
	1: "Succeeded",
	2: "Failed",
}

var WorkloadStatus_State_value = map[string]int32{
	"Pending":   0,
	"Succeeded": 1,
	"Failed":    2,
}

func (x WorkloadStatus_State) String() string {
	return proto.EnumName(WorkloadStatus_State_name, int32(x))
}

func (WorkloadStatus_State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_24d13e575ab7b28c, []int{6, 0}
}

// A FilterDeployment tells the Wasme Operator
// to deploy a filter with the provided configuration
// to the target workloads.
// Currently FilterDeployments support Wasm filters on Istio
type FilterDeploymentSpec struct {
	// the spec of the filter to deploy
	Filter *FilterSpec `protobuf:"bytes,1,opt,name=filter,proto3" json:"filter,omitempty"`
	// Spec that selects one or more target workloads in the FilterDeployment namespace
	Deployment           *DeploymentSpec `protobuf:"bytes,2,opt,name=deployment,proto3" json:"deployment,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *FilterDeploymentSpec) Reset()         { *m = FilterDeploymentSpec{} }
func (m *FilterDeploymentSpec) String() string { return proto.CompactTextString(m) }
func (*FilterDeploymentSpec) ProtoMessage()    {}
func (*FilterDeploymentSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_24d13e575ab7b28c, []int{0}
}
func (m *FilterDeploymentSpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilterDeploymentSpec.Unmarshal(m, b)
}
func (m *FilterDeploymentSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilterDeploymentSpec.Marshal(b, m, deterministic)
}
func (m *FilterDeploymentSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilterDeploymentSpec.Merge(m, src)
}
func (m *FilterDeploymentSpec) XXX_Size() int {
	return xxx_messageInfo_FilterDeploymentSpec.Size(m)
}
func (m *FilterDeploymentSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_FilterDeploymentSpec.DiscardUnknown(m)
}

var xxx_messageInfo_FilterDeploymentSpec proto.InternalMessageInfo

func (m *FilterDeploymentSpec) GetFilter() *FilterSpec {
	if m != nil {
		return m.Filter
	}
	return nil
}

func (m *FilterDeploymentSpec) GetDeployment() *DeploymentSpec {
	if m != nil {
		return m.Deployment
	}
	return nil
}

// the filter to deploy
type FilterSpec struct {
	// unique identifier that will be used
	// to remove the filter as well as for logging.
	// if id is not set, it will be set automatically to be the name.namespace
	// of the FilterDeployment resource
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// name of image which houses the compiled wasm filter
	Image string `protobuf:"bytes,2,opt,name=image,proto3" json:"image,omitempty"`
	// Filter/service configuration used to configure or reconfigure a plugin
	// (proxy_on_configuration).
	// `google.protobuf.Struct` is serialized as JSON before
	// passing it to the plugin. `google.protobuf.BytesValue` and
	// `google.protobuf.StringValue` are passed directly without the wrapper.
	Config *types.Any `protobuf:"bytes,3,opt,name=config,proto3" json:"config,omitempty"`
	// the root id must match the root id
	// defined inside the filter.
	// if the user does not provide this field,
	// wasme will attempt to pull the image
	// and set it from the filter_conf
	// the first time it must pull the image and inspect it
	// second time it will cache it locally
	// if the user provides
	RootID string `protobuf:"bytes,4,opt,name=rootID,proto3" json:"rootID,omitempty"`
	// custom options if pulling from private / custom repositories
	ImagePullOptions     *ImagePullOptions `protobuf:"bytes,5,opt,name=imagePullOptions,proto3" json:"imagePullOptions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *FilterSpec) Reset()         { *m = FilterSpec{} }
func (m *FilterSpec) String() string { return proto.CompactTextString(m) }
func (*FilterSpec) ProtoMessage()    {}
func (*FilterSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_24d13e575ab7b28c, []int{1}
}
func (m *FilterSpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilterSpec.Unmarshal(m, b)
}
func (m *FilterSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilterSpec.Marshal(b, m, deterministic)
}
func (m *FilterSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilterSpec.Merge(m, src)
}
func (m *FilterSpec) XXX_Size() int {
	return xxx_messageInfo_FilterSpec.Size(m)
}
func (m *FilterSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_FilterSpec.DiscardUnknown(m)
}

var xxx_messageInfo_FilterSpec proto.InternalMessageInfo

func (m *FilterSpec) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *FilterSpec) GetImage() string {
	if m != nil {
		return m.Image
	}
	return ""
}

func (m *FilterSpec) GetConfig() *types.Any {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *FilterSpec) GetRootID() string {
	if m != nil {
		return m.RootID
	}
	return ""
}

func (m *FilterSpec) GetImagePullOptions() *ImagePullOptions {
	if m != nil {
		return m.ImagePullOptions
	}
	return nil
}

type ImagePullOptions struct {
	// if a username/password is required,
	// specify here the name of a secret:
	// with keys:
	// * username: <username>
	// * password: <password>
	//
	// the secret must live in the same namespace
	// as the FilterDeployment
	PullSecret string `protobuf:"bytes,1,opt,name=pullSecret,proto3" json:"pullSecret,omitempty"`
	// skip verifying the image server's TLS certificate
	InsecureSkipVerify bool `protobuf:"varint,2,opt,name=insecureSkipVerify,proto3" json:"insecureSkipVerify,omitempty"`
	// use HTTP instead of HTTPS
	PlainHttp            bool     `protobuf:"varint,3,opt,name=plainHttp,proto3" json:"plainHttp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImagePullOptions) Reset()         { *m = ImagePullOptions{} }
func (m *ImagePullOptions) String() string { return proto.CompactTextString(m) }
func (*ImagePullOptions) ProtoMessage()    {}
func (*ImagePullOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_24d13e575ab7b28c, []int{2}
}
func (m *ImagePullOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImagePullOptions.Unmarshal(m, b)
}
func (m *ImagePullOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImagePullOptions.Marshal(b, m, deterministic)
}
func (m *ImagePullOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImagePullOptions.Merge(m, src)
}
func (m *ImagePullOptions) XXX_Size() int {
	return xxx_messageInfo_ImagePullOptions.Size(m)
}
func (m *ImagePullOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_ImagePullOptions.DiscardUnknown(m)
}

var xxx_messageInfo_ImagePullOptions proto.InternalMessageInfo

func (m *ImagePullOptions) GetPullSecret() string {
	if m != nil {
		return m.PullSecret
	}
	return ""
}

func (m *ImagePullOptions) GetInsecureSkipVerify() bool {
	if m != nil {
		return m.InsecureSkipVerify
	}
	return false
}

func (m *ImagePullOptions) GetPlainHttp() bool {
	if m != nil {
		return m.PlainHttp
	}
	return false
}

// how to deploy the filter
type DeploymentSpec struct {
	// Types that are valid to be assigned to DeploymentType:
	//	*DeploymentSpec_Istio
	DeploymentType       isDeploymentSpec_DeploymentType `protobuf_oneof:"deploymentType"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *DeploymentSpec) Reset()         { *m = DeploymentSpec{} }
func (m *DeploymentSpec) String() string { return proto.CompactTextString(m) }
func (*DeploymentSpec) ProtoMessage()    {}
func (*DeploymentSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_24d13e575ab7b28c, []int{3}
}
func (m *DeploymentSpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeploymentSpec.Unmarshal(m, b)
}
func (m *DeploymentSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeploymentSpec.Marshal(b, m, deterministic)
}
func (m *DeploymentSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeploymentSpec.Merge(m, src)
}
func (m *DeploymentSpec) XXX_Size() int {
	return xxx_messageInfo_DeploymentSpec.Size(m)
}
func (m *DeploymentSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_DeploymentSpec.DiscardUnknown(m)
}

var xxx_messageInfo_DeploymentSpec proto.InternalMessageInfo

type isDeploymentSpec_DeploymentType interface {
	isDeploymentSpec_DeploymentType()
}

type DeploymentSpec_Istio struct {
	Istio *IstioDeploymentSpec `protobuf:"bytes,2,opt,name=istio,proto3,oneof" json:"istio,omitempty"`
}

func (*DeploymentSpec_Istio) isDeploymentSpec_DeploymentType() {}

func (m *DeploymentSpec) GetDeploymentType() isDeploymentSpec_DeploymentType {
	if m != nil {
		return m.DeploymentType
	}
	return nil
}

func (m *DeploymentSpec) GetIstio() *IstioDeploymentSpec {
	if x, ok := m.GetDeploymentType().(*DeploymentSpec_Istio); ok {
		return x.Istio
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*DeploymentSpec) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*DeploymentSpec_Istio)(nil),
	}
}

// how to deploy to Istio
type IstioDeploymentSpec struct {
	// the kind of workload to deploy the filter to
	// can either be Deployment or DaemonSet
	Kind string `protobuf:"bytes,1,opt,name=kind,proto3" json:"kind,omitempty"`
	// deploy the filter to workloads with these labels
	// the workload must live in the same namespace as the FilterDeployment
	// if empty, the filter will be deployed to all workloads in the namespace
	Labels map[string]string `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// the namespace where the Istio control plane is installed.
	// defaults to `istio-system`.
	IstioNamespace       string   `protobuf:"bytes,3,opt,name=istioNamespace,proto3" json:"istioNamespace,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IstioDeploymentSpec) Reset()         { *m = IstioDeploymentSpec{} }
func (m *IstioDeploymentSpec) String() string { return proto.CompactTextString(m) }
func (*IstioDeploymentSpec) ProtoMessage()    {}
func (*IstioDeploymentSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_24d13e575ab7b28c, []int{4}
}
func (m *IstioDeploymentSpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IstioDeploymentSpec.Unmarshal(m, b)
}
func (m *IstioDeploymentSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IstioDeploymentSpec.Marshal(b, m, deterministic)
}
func (m *IstioDeploymentSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IstioDeploymentSpec.Merge(m, src)
}
func (m *IstioDeploymentSpec) XXX_Size() int {
	return xxx_messageInfo_IstioDeploymentSpec.Size(m)
}
func (m *IstioDeploymentSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_IstioDeploymentSpec.DiscardUnknown(m)
}

var xxx_messageInfo_IstioDeploymentSpec proto.InternalMessageInfo

func (m *IstioDeploymentSpec) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func (m *IstioDeploymentSpec) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *IstioDeploymentSpec) GetIstioNamespace() string {
	if m != nil {
		return m.IstioNamespace
	}
	return ""
}

// the current status of the deployment
type FilterDeploymentStatus struct {
	// the observed generation of the FilterDeployment
	ObservedGeneration int64 `protobuf:"varint,1,opt,name=observedGeneration,proto3" json:"observedGeneration,omitempty"`
	// for each workload, was the deployment successful?
	Workloads map[string]*WorkloadStatus `protobuf:"bytes,2,rep,name=workloads,proto3" json:"workloads,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// a human-readable string explaining the error, if any
	Reason               string   `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FilterDeploymentStatus) Reset()         { *m = FilterDeploymentStatus{} }
func (m *FilterDeploymentStatus) String() string { return proto.CompactTextString(m) }
func (*FilterDeploymentStatus) ProtoMessage()    {}
func (*FilterDeploymentStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_24d13e575ab7b28c, []int{5}
}
func (m *FilterDeploymentStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilterDeploymentStatus.Unmarshal(m, b)
}
func (m *FilterDeploymentStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilterDeploymentStatus.Marshal(b, m, deterministic)
}
func (m *FilterDeploymentStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilterDeploymentStatus.Merge(m, src)
}
func (m *FilterDeploymentStatus) XXX_Size() int {
	return xxx_messageInfo_FilterDeploymentStatus.Size(m)
}
func (m *FilterDeploymentStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_FilterDeploymentStatus.DiscardUnknown(m)
}

var xxx_messageInfo_FilterDeploymentStatus proto.InternalMessageInfo

func (m *FilterDeploymentStatus) GetObservedGeneration() int64 {
	if m != nil {
		return m.ObservedGeneration
	}
	return 0
}

func (m *FilterDeploymentStatus) GetWorkloads() map[string]*WorkloadStatus {
	if m != nil {
		return m.Workloads
	}
	return nil
}

func (m *FilterDeploymentStatus) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

type WorkloadStatus struct {
	State WorkloadStatus_State `protobuf:"varint,1,opt,name=state,proto3,enum=wasme.io.WorkloadStatus_State" json:"state,omitempty"`
	// a human-readable string explaining the error, if any
	Reason               string   `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WorkloadStatus) Reset()         { *m = WorkloadStatus{} }
func (m *WorkloadStatus) String() string { return proto.CompactTextString(m) }
func (*WorkloadStatus) ProtoMessage()    {}
func (*WorkloadStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_24d13e575ab7b28c, []int{6}
}
func (m *WorkloadStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WorkloadStatus.Unmarshal(m, b)
}
func (m *WorkloadStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WorkloadStatus.Marshal(b, m, deterministic)
}
func (m *WorkloadStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WorkloadStatus.Merge(m, src)
}
func (m *WorkloadStatus) XXX_Size() int {
	return xxx_messageInfo_WorkloadStatus.Size(m)
}
func (m *WorkloadStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_WorkloadStatus.DiscardUnknown(m)
}

var xxx_messageInfo_WorkloadStatus proto.InternalMessageInfo

func (m *WorkloadStatus) GetState() WorkloadStatus_State {
	if m != nil {
		return m.State
	}
	return WorkloadStatus_Pending
}

func (m *WorkloadStatus) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

func init() {
	proto.RegisterEnum("wasme.io.WorkloadStatus_State", WorkloadStatus_State_name, WorkloadStatus_State_value)
	proto.RegisterType((*FilterDeploymentSpec)(nil), "wasme.io.FilterDeploymentSpec")
	proto.RegisterType((*FilterSpec)(nil), "wasme.io.FilterSpec")
	proto.RegisterType((*ImagePullOptions)(nil), "wasme.io.ImagePullOptions")
	proto.RegisterType((*DeploymentSpec)(nil), "wasme.io.DeploymentSpec")
	proto.RegisterType((*IstioDeploymentSpec)(nil), "wasme.io.IstioDeploymentSpec")
	proto.RegisterMapType((map[string]string)(nil), "wasme.io.IstioDeploymentSpec.LabelsEntry")
	proto.RegisterType((*FilterDeploymentStatus)(nil), "wasme.io.FilterDeploymentStatus")
	proto.RegisterMapType((map[string]*WorkloadStatus)(nil), "wasme.io.FilterDeploymentStatus.WorkloadsEntry")
	proto.RegisterType((*WorkloadStatus)(nil), "wasme.io.WorkloadStatus")
}

func init() {
	proto.RegisterFile("github.com/solo-io/wasm/tools/wasme/cli/operator/api/wasme/v1/filter_deployment.proto", fileDescriptor_24d13e575ab7b28c)
}

var fileDescriptor_24d13e575ab7b28c = []byte{
	// 641 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xad, 0x9d, 0x26, 0x34, 0x13, 0x11, 0x59, 0x4b, 0x55, 0x99, 0x08, 0xaa, 0xca, 0x07, 0x54,
	0x24, 0xb0, 0xd5, 0x02, 0x52, 0xe1, 0xd6, 0xaa, 0x84, 0x56, 0xe2, 0xa3, 0x72, 0xa0, 0x08, 0x2e,
	0x68, 0x63, 0x4f, 0xcc, 0x2a, 0x9b, 0x5d, 0xcb, 0x5e, 0xa7, 0xf2, 0x05, 0x71, 0xe3, 0xc8, 0x5f,
	0xe2, 0x1f, 0xf0, 0x97, 0x90, 0xd7, 0x4e, 0xed, 0xa4, 0x01, 0x71, 0xf2, 0xee, 0xf8, 0xbd, 0x37,
	0xf3, 0x76, 0x66, 0x17, 0x3e, 0x44, 0x4c, 0x7d, 0xcd, 0xc6, 0x6e, 0x20, 0x67, 0x5e, 0x2a, 0xb9,
	0x7c, 0xcc, 0xa4, 0x77, 0x45, 0xd3, 0x99, 0xa7, 0xa4, 0xe4, 0xa9, 0x5e, 0xa2, 0x17, 0x70, 0xe6,
	0xc9, 0x18, 0x13, 0xaa, 0x64, 0xe2, 0xd1, 0x98, 0x55, 0xe1, 0xf9, 0x81, 0x37, 0x61, 0x5c, 0x61,
	0xf2, 0x25, 0xc4, 0x98, 0xcb, 0x7c, 0x86, 0x42, 0xb9, 0x71, 0x22, 0x95, 0x24, 0x5b, 0x1a, 0xe1,
	0x32, 0x39, 0xb8, 0x1b, 0x49, 0x19, 0x71, 0xf4, 0x74, 0x7c, 0x9c, 0x4d, 0x3c, 0x2a, 0xf2, 0x12,
	0xe4, 0x7c, 0x83, 0xed, 0xa1, 0xe6, 0x9f, 0x5e, 0xd3, 0x47, 0x31, 0x06, 0xe4, 0x11, 0x74, 0x4a,
	0x5d, 0xdb, 0xd8, 0x33, 0xf6, 0x7b, 0x87, 0xdb, 0xee, 0x42, 0xcd, 0x2d, 0xf1, 0x05, 0xca, 0xaf,
	0x30, 0xe4, 0x08, 0xa0, 0x4e, 0x6f, 0x9b, 0x9a, 0x61, 0xd7, 0x8c, 0x65, 0x6d, 0xbf, 0x81, 0x75,
	0x7e, 0x19, 0x00, 0xb5, 0x20, 0xe9, 0x83, 0xc9, 0x42, 0x9d, 0xb2, 0xeb, 0x9b, 0x2c, 0x24, 0xdb,
	0xd0, 0x66, 0x33, 0x1a, 0xa1, 0xd6, 0xec, 0xfa, 0xe5, 0xa6, 0x28, 0x2e, 0x90, 0x62, 0xc2, 0x22,
	0xbb, 0x55, 0x15, 0x57, 0x1a, 0x74, 0x17, 0x06, 0xdd, 0x63, 0x91, 0xfb, 0x15, 0x86, 0xec, 0x40,
	0x27, 0x91, 0x52, 0x9d, 0x9f, 0xda, 0x9b, 0x5a, 0xa4, 0xda, 0x91, 0x21, 0x58, 0x5a, 0xee, 0x22,
	0xe3, 0xfc, 0x5d, 0xac, 0x98, 0x14, 0xa9, 0xdd, 0xd6, 0x7a, 0x83, 0xba, 0xf4, 0xf3, 0x15, 0x84,
	0x7f, 0x83, 0xe3, 0x7c, 0x37, 0xc0, 0x5a, 0x85, 0x91, 0x5d, 0x80, 0x38, 0xe3, 0x7c, 0x84, 0x41,
	0x82, 0xaa, 0x32, 0xd4, 0x88, 0x10, 0x17, 0x08, 0x13, 0x29, 0x06, 0x59, 0x82, 0xa3, 0x29, 0x8b,
	0x2f, 0x31, 0x61, 0x93, 0x5c, 0xbb, 0xdc, 0xf2, 0xd7, 0xfc, 0x21, 0xf7, 0xa0, 0x1b, 0x73, 0xca,
	0xc4, 0x99, 0x52, 0xb1, 0x76, 0xbd, 0xe5, 0xd7, 0x01, 0xe7, 0x13, 0xf4, 0x57, 0xfa, 0xf7, 0x0c,
	0xda, 0x2c, 0x55, 0x4c, 0x56, 0xcd, 0xb8, 0xdf, 0x70, 0x54, 0x84, 0x97, 0xd1, 0x67, 0x1b, 0x7e,
	0x89, 0x3e, 0xb1, 0xa0, 0x5f, 0x37, 0xe7, 0x7d, 0x1e, 0xa3, 0xf3, 0xdb, 0x80, 0x3b, 0x6b, 0x28,
	0x84, 0xc0, 0xe6, 0x94, 0x89, 0x45, 0xaf, 0xf4, 0x9a, 0x1c, 0x43, 0x87, 0xd3, 0x31, 0xf2, 0xd4,
	0x36, 0xf7, 0x5a, 0xfb, 0xbd, 0xc3, 0x87, 0xff, 0xcc, 0xea, 0xbe, 0xd6, 0xd8, 0x97, 0x42, 0x25,
	0xb9, 0x5f, 0x11, 0xc9, 0x03, 0xe8, 0xeb, 0x4a, 0xde, 0xd2, 0x19, 0xa6, 0x31, 0x0d, 0x50, 0x9b,
	0xed, 0xfa, 0x2b, 0xd1, 0xc1, 0x73, 0xe8, 0x35, 0xe8, 0xc4, 0x82, 0xd6, 0x14, 0xf3, 0xaa, 0x98,
	0x62, 0x59, 0x4c, 0xce, 0x9c, 0xf2, 0xec, 0x7a, 0x72, 0xf4, 0xe6, 0x85, 0x79, 0x64, 0x38, 0x3f,
	0x4c, 0xd8, 0xb9, 0x31, 0xf3, 0x8a, 0xaa, 0x2c, 0x2d, 0xba, 0x22, 0xc7, 0x29, 0x26, 0x73, 0x0c,
	0x5f, 0xa1, 0x28, 0x2e, 0x1b, 0x93, 0x42, 0xab, 0xb6, 0xfc, 0x35, 0x7f, 0xc8, 0x1b, 0xe8, 0x5e,
	0xc9, 0x64, 0xca, 0x25, 0x0d, 0x17, 0x9e, 0xbd, 0xd5, 0x8b, 0xb2, 0x9a, 0xc4, 0xfd, 0xb8, 0x60,
	0x94, 0xce, 0x6b, 0x05, 0x3d, 0xa9, 0x48, 0x53, 0x29, 0x2a, 0xd3, 0xd5, 0x6e, 0x70, 0x09, 0xfd,
	0x65, 0xd2, 0x1a, 0xbf, 0x6e, 0xd3, 0xef, 0xd2, 0xed, 0x5b, 0x50, 0xcb, 0xf4, 0xcd, 0x93, 0xf8,
	0x69, 0xd4, 0xc2, 0xd5, 0x09, 0x3c, 0x85, 0x76, 0xaa, 0xa8, 0x42, 0x2d, 0xdd, 0x3f, 0xdc, 0xfd,
	0x9b, 0x8c, 0x5b, 0x7c, 0xd0, 0x2f, 0xc1, 0x8d, 0xc2, 0xcd, 0x66, 0xe1, 0x8e, 0x07, 0x6d, 0x8d,
	0x23, 0x3d, 0xb8, 0x75, 0x81, 0x22, 0x64, 0x22, 0xb2, 0x36, 0xc8, 0x6d, 0xe8, 0x8e, 0xb2, 0x20,
	0x40, 0x0c, 0x31, 0xb4, 0x0c, 0x02, 0xd0, 0x19, 0x52, 0xc6, 0x31, 0xb4, 0xcc, 0x93, 0xe1, 0xe7,
	0xd3, 0xff, 0x7d, 0x0c, 0xe3, 0x69, 0xb4, 0xe6, 0x41, 0x74, 0x99, 0xf4, 0xe6, 0x07, 0xe3, 0x8e,
	0x7e, 0x09, 0x9e, 0xfc, 0x09, 0x00, 0x00, 0xff, 0xff, 0x82, 0x66, 0x05, 0xe8, 0x5b, 0x05, 0x00,
	0x00,
}
