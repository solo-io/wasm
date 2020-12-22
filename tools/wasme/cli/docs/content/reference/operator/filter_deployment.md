
---
title: "wasme.iogithub.com/solo-io/wasm/tools/wasme/cli/operator/api/wasme/v1/filter_deployment.proto"
---

## Package : `wasme.io`



<a name="top"></a>

<a name="API Reference for github.com/solo-io/wasm/tools/wasme/cli/operator/api/wasme/v1/filter_deployment.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/solo-io/wasm/tools/wasme/cli/operator/api/wasme/v1/filter_deployment.proto


## Table of Contents
  - [DeploymentSpec](#wasme.io.DeploymentSpec)
  - [FilterDeploymentSpec](#wasme.io.FilterDeploymentSpec)
  - [FilterDeploymentStatus](#wasme.io.FilterDeploymentStatus)
  - [FilterDeploymentStatus.WorkloadsEntry](#wasme.io.FilterDeploymentStatus.WorkloadsEntry)
  - [FilterSpec](#wasme.io.FilterSpec)
  - [ImagePullOptions](#wasme.io.ImagePullOptions)
  - [IstioDeploymentSpec](#wasme.io.IstioDeploymentSpec)
  - [IstioDeploymentSpec.LabelsEntry](#wasme.io.IstioDeploymentSpec.LabelsEntry)
  - [WorkloadStatus](#wasme.io.WorkloadStatus)

  - [WorkloadStatus.State](#wasme.io.WorkloadStatus.State)






<a name="wasme.io.DeploymentSpec"></a>

### DeploymentSpec
how to deploy the filter


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| istio | [IstioDeploymentSpec](#wasme.io.IstioDeploymentSpec) |  | Deploy to Istio |






<a name="wasme.io.FilterDeploymentSpec"></a>

### FilterDeploymentSpec
A FilterDeployment tells the Wasme Operator
to deploy a filter with the provided configuration
to the target workloads.
Currently FilterDeployments support Wasm filters on Istio


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| filter | [FilterSpec](#wasme.io.FilterSpec) |  | the spec of the filter to deploy |
| deployment | [DeploymentSpec](#wasme.io.DeploymentSpec) |  | Spec that selects one or more target workloads in the FilterDeployment namespace |






<a name="wasme.io.FilterDeploymentStatus"></a>

### FilterDeploymentStatus
the current status of the deployment


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| observedGeneration | [int64](#int64) |  | the observed generation of the FilterDeployment |
| workloads | [][FilterDeploymentStatus.WorkloadsEntry](#wasme.io.FilterDeploymentStatus.WorkloadsEntry) | repeated | for each workload, was the deployment successful? |
| reason | [string](#string) |  | a human-readable string explaining the error, if any |






<a name="wasme.io.FilterDeploymentStatus.WorkloadsEntry"></a>

### FilterDeploymentStatus.WorkloadsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [WorkloadStatus](#wasme.io.WorkloadStatus) |  |  |






<a name="wasme.io.FilterSpec"></a>

### FilterSpec
the filter to deploy


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | unique identifier that will be used
to remove the filter as well as for logging.
if id is not set, it will be set automatically to be the name.namespace
of the FilterDeployment resource |
| image | [string](#string) |  | name of image which houses the compiled wasm filter |
| config | [google.protobuf.Any](#google.protobuf.Any) |  | Filter/service configuration used to configure or reconfigure a plugin
(proxy_on_configuration).
`google.protobuf.Struct` is serialized as JSON before
passing it to the plugin. `google.protobuf.BytesValue` and
`google.protobuf.StringValue` are passed directly without the wrapper. |
| rootID | [string](#string) |  | the root id must match the root id
defined inside the filter.
if the user does not provide this field,
wasme will attempt to pull the image
and set it from the filter_conf
the first time it must pull the image and inspect it
second time it will cache it locally
if the user provides |
| imagePullOptions | [ImagePullOptions](#wasme.io.ImagePullOptions) |  | custom options if pulling from private / custom repositories |
| patchContext | [string](#string) |  | a class of configurations based on the traffic flow direction
and workload type.
defaults to `inbound`. |






<a name="wasme.io.ImagePullOptions"></a>

### ImagePullOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pullSecret | [string](#string) |  | if a username/password is required,
specify here the name of a secret:
with keys:
* username: &lt;username&gt;
* password: &lt;password&gt;

the secret must live in the same namespace
as the FilterDeployment |
| insecureSkipVerify | [bool](#bool) |  | skip verifying the image server&#39;s TLS certificate |
| plainHttp | [bool](#bool) |  | use HTTP instead of HTTPS |






<a name="wasme.io.IstioDeploymentSpec"></a>

### IstioDeploymentSpec
how to deploy to Istio


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [string](#string) |  | the kind of workload to deploy the filter to
can either be Deployment, DaemonSet or Statefulset |
| labels | [][IstioDeploymentSpec.LabelsEntry](#wasme.io.IstioDeploymentSpec.LabelsEntry) | repeated | deploy the filter to workloads with these labels
the workload must live in the same namespace as the FilterDeployment
if empty, the filter will be deployed to all workloads in the namespace |
| istioNamespace | [string](#string) |  | the namespace where the Istio control plane is installed.
defaults to `istio-system`. |






<a name="wasme.io.IstioDeploymentSpec.LabelsEntry"></a>

### IstioDeploymentSpec.LabelsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="wasme.io.WorkloadStatus"></a>

### WorkloadStatus



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state | [WorkloadStatus.State](#wasme.io.WorkloadStatus.State) |  |  |
| reason | [string](#string) |  | a human-readable string explaining the error, if any |





 


<a name="wasme.io.WorkloadStatus.State"></a>

### WorkloadStatus.State
the state of the filter deployment

| Name | Number | Description |
| ---- | ------ | ----------- |
| Pending | 0 |  |
| Succeeded | 1 |  |
| Failed | 2 |  |


 

 

 

