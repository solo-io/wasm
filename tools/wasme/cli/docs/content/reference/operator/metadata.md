
---
title: "core.solo.iogithub.com/solo-io/solo-kit/api/v1/metadata.proto"
---

## Package : `core.solo.io`



<a name="top"></a>

<a name="API Reference for github.com/solo-io/solo-kit/api/v1/metadata.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/solo-io/solo-kit/api/v1/metadata.proto


## Table of Contents
  - [Metadata](#core.solo.io.Metadata)
  - [Metadata.AnnotationsEntry](#core.solo.io.Metadata.AnnotationsEntry)
  - [Metadata.LabelsEntry](#core.solo.io.Metadata.LabelsEntry)
  - [Metadata.OwnerReference](#core.solo.io.Metadata.OwnerReference)







<a name="core.solo.io.Metadata"></a>

### Metadata
Metadata contains general properties of resources for purposes of versioning, annotating, and namespacing.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the resource.

Names must be unique and follow the following syntax rules:

One or more lowercase rfc1035/rfc1123 labels separated by &#39;.&#39; with a maximum length of 253 characters. |
| namespace | [string](#string) |  | Namespace is used for the namespacing of resources. |
| cluster | [string](#string) |  | Cluster indicates the cluster this resource belongs to
Cluster is only applicable in certain contexts, e.g. Kubernetes
An empty string here refers to the local cluster |
| resource_version | [string](#string) |  | An opaque value that represents the internal version of this object that can
be used by clients to determine when objects have changed. |
| labels | [][Metadata.LabelsEntry](#core.solo.io.Metadata.LabelsEntry) | repeated | Map of string keys and values that can be used to organize and categorize
(scope and select) objects. Some resources contain `selectors` which
can be linked with other resources by their labels |
| annotations | [][Metadata.AnnotationsEntry](#core.solo.io.Metadata.AnnotationsEntry) | repeated | Annotations is an unstructured key value map stored with a resource that may be
set by external tools to store and retrieve arbitrary metadata. |
| generation | [int64](#int64) |  | A sequence number representing a specific generation of the desired state.
Currently only populated for resources backed by Kubernetes |
| owner_references | [][Metadata.OwnerReference](#core.solo.io.Metadata.OwnerReference) | repeated | List of objects depended by this object.
Currently only populated for resources backed by Kubernetes |






<a name="core.solo.io.Metadata.AnnotationsEntry"></a>

### Metadata.AnnotationsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="core.solo.io.Metadata.LabelsEntry"></a>

### Metadata.LabelsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="core.solo.io.Metadata.OwnerReference"></a>

### Metadata.OwnerReference
proto message representing kubernertes owner reference
https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.15/#ownerreference-v1-meta


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| api_version | [string](#string) |  |  |
| block_owner_deletion | [google.protobuf.BoolValue](#google.protobuf.BoolValue) |  |  |
| controller | [google.protobuf.BoolValue](#google.protobuf.BoolValue) |  |  |
| kind | [string](#string) |  |  |
| name | [string](#string) |  |  |
| uid | [string](#string) |  |  |





 

 

 

 

