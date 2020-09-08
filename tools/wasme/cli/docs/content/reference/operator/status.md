
---
title: "core.solo.iogithub.com/solo-io/solo-kit/api/v1/status.proto"
---

## Package : `core.solo.io`



<a name="top"></a>

<a name="API Reference for github.com/solo-io/solo-kit/api/v1/status.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/solo-io/solo-kit/api/v1/status.proto


## Table of Contents
  - [Status](#core.solo.io.Status)
  - [Status.SubresourceStatusesEntry](#core.solo.io.Status.SubresourceStatusesEntry)

  - [Status.State](#core.solo.io.Status.State)






<a name="core.solo.io.Status"></a>

### Status
Status indicates whether a resource has been (in)validated by a reporter in the system.
Statuses are meant to be read-only by users


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| state | [Status.State](#core.solo.io.Status.State) |  | State is the enum indicating the state of the resource |
| reason | [string](#string) |  | Reason is a description of the error for Rejected resources. If the resource is pending or accepted, this field will be empty |
| reported_by | [string](#string) |  | Reference to the reporter who wrote this status |
| subresource_statuses | [][Status.SubresourceStatusesEntry](#core.solo.io.Status.SubresourceStatusesEntry) | repeated | Reference to statuses (by resource-ref string: &#34;Kind.Namespace.Name&#34;) of subresources of the parent resource |
| details | [google.protobuf.Struct](#google.protobuf.Struct) |  | Opaque details about status results |






<a name="core.solo.io.Status.SubresourceStatusesEntry"></a>

### Status.SubresourceStatusesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Status](#core.solo.io.Status) |  |  |





 


<a name="core.solo.io.Status.State"></a>

### Status.State


| Name | Number | Description |
| ---- | ------ | ----------- |
| Pending | 0 | Pending status indicates the resource has not yet been validated |
| Accepted | 1 | Accepted indicates the resource has been validated |
| Rejected | 2 | Rejected indicates an invalid configuration by the user
Rejected resources may be propagated to the xDS server depending on their severity |
| Warning | 3 | Warning indicates a partially invalid configuration by the user
Resources with Warnings may be partially accepted by a controller, depending on the implementation |


 

 

 

