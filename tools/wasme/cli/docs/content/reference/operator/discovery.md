
---
title: "envoy.api.v2github.com/solo-io/solo-kit/api/external/envoy/api/v2/discovery.proto"
---

## Package : `envoy.api.v2`



<a name="top"></a>

<a name="API Reference for github.com/solo-io/solo-kit/api/external/envoy/api/v2/discovery.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/solo-io/solo-kit/api/external/envoy/api/v2/discovery.proto


## Table of Contents
  - [DeltaDiscoveryRequest](#envoy.api.v2.DeltaDiscoveryRequest)
  - [DeltaDiscoveryRequest.InitialResourceVersionsEntry](#envoy.api.v2.DeltaDiscoveryRequest.InitialResourceVersionsEntry)
  - [DeltaDiscoveryResponse](#envoy.api.v2.DeltaDiscoveryResponse)
  - [DiscoveryRequest](#envoy.api.v2.DiscoveryRequest)
  - [DiscoveryResponse](#envoy.api.v2.DiscoveryResponse)
  - [Resource](#envoy.api.v2.Resource)







<a name="envoy.api.v2.DeltaDiscoveryRequest"></a>

### DeltaDiscoveryRequest
DeltaDiscoveryRequest and DeltaDiscoveryResponse are used in a new gRPC
endpoint for Delta xDS.

With Delta xDS, the DeltaDiscoveryResponses do not need to include a full
snapshot of the tracked resources. Instead, DeltaDiscoveryResponses are a
diff to the state of a xDS client.
In Delta XDS there are per-resource versions, which allow tracking state at
the resource granularity.
An xDS Delta session is always in the context of a gRPC bidirectional
stream. This allows the xDS server to keep track of the state of xDS clients
connected to it.

In Delta xDS the nonce field is required and used to pair
DeltaDiscoveryResponse to a DeltaDiscoveryRequest ACK or NACK.
Optionally, a response message level system_version_info is present for
debugging purposes only.

DeltaDiscoveryRequest plays two independent roles. Any DeltaDiscoveryRequest
can be either or both of: [1] informing the server of what resources the
client has gained/lost interest in (using resource_names_subscribe and
resource_names_unsubscribe), or [2] (N)ACKing an earlier resource update from
the server (using response_nonce, with presence of error_detail making it a NACK).
Additionally, the first message (for a given type_url) of a reconnected gRPC stream
has a third role: informing the server of the resources (and their versions)
that the client already possesses, using the initial_resource_versions field.

As with state-of-the-world, when multiple resource types are multiplexed (ADS),
all requests/acknowledgments/updates are logically walled off by type_url:
a Cluster ACK exists in a completely separate world from a prior Route NACK.
In particular, initial_resource_versions being sent at the &#34;start&#34; of every
gRPC stream actually entails a message for each type_url, each with its own
initial_resource_versions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| node | [core.Node](#envoy.api.v2.core.Node) |  | The node making the request. |
| type_url | [string](#string) |  | Type of the resource that is being requested, e.g.
&#34;type.googleapis.com/envoy.api.v2.ClusterLoadAssignment&#34;. |
| resource_names_subscribe | [][string](#string) | repeated | DeltaDiscoveryRequests allow the client to add or remove individual
resources to the set of tracked resources in the context of a stream.
All resource names in the resource_names_subscribe list are added to the
set of tracked resources and all resource names in the resource_names_unsubscribe
list are removed from the set of tracked resources.

*Unlike* state-of-the-world xDS, an empty resource_names_subscribe or
resource_names_unsubscribe list simply means that no resources are to be
added or removed to the resource list.
*Like* state-of-the-world xDS, the server must send updates for all tracked
resources, but can also send updates for resources the client has not subscribed to.

NOTE: the server must respond with all resources listed in resource_names_subscribe,
even if it believes the client has the most recent version of them. The reason:
the client may have dropped them, but then regained interest before it had a chance
to send the unsubscribe message. See DeltaSubscriptionStateTest.RemoveThenAdd.

These two fields can be set in any DeltaDiscoveryRequest, including ACKs
and initial_resource_versions.

A list of Resource names to add to the list of tracked resources. |
| resource_names_unsubscribe | [][string](#string) | repeated | A list of Resource names to remove from the list of tracked resources. |
| initial_resource_versions | [][DeltaDiscoveryRequest.InitialResourceVersionsEntry](#envoy.api.v2.DeltaDiscoveryRequest.InitialResourceVersionsEntry) | repeated | Informs the server of the versions of the resources the xDS client knows of, to enable the
client to continue the same logical xDS session even in the face of gRPC stream reconnection.
It will not be populated: [1] in the very first stream of a session, since the client will
not yet have any resources,  [2] in any message after the first in a stream (for a given
type_url), since the server will already be correctly tracking the client&#39;s state.
(In ADS, the first message *of each type_url* of a reconnected stream populates this map.)
The map&#39;s keys are names of xDS resources known to the xDS client.
The map&#39;s values are opaque resource versions. |
| response_nonce | [string](#string) |  | When the DeltaDiscoveryRequest is a ACK or NACK message in response
to a previous DeltaDiscoveryResponse, the response_nonce must be the
nonce in the DeltaDiscoveryResponse.
Otherwise (unlike in DiscoveryRequest) response_nonce must be omitted. |
| error_detail | [google.rpc.Status](#google.rpc.Status) |  | This is populated when the previous `DiscoveryResponse (envoy_api_msg_DiscoveryResponse)`
failed to update configuration. The *message* field in *error_details*
provides the Envoy internal exception related to the failure. |






<a name="envoy.api.v2.DeltaDiscoveryRequest.InitialResourceVersionsEntry"></a>

### DeltaDiscoveryRequest.InitialResourceVersionsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="envoy.api.v2.DeltaDiscoveryResponse"></a>

### DeltaDiscoveryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| system_version_info | [string](#string) |  | The version of the response data (used for debugging). |
| resources | [][Resource](#envoy.api.v2.Resource) | repeated | The response resources. These are typed resources, whose types must match
the type_url field. |
| type_url | [string](#string) |  | Type URL for resources. Identifies the xDS API when muxing over ADS.
Must be consistent with the type_url in the Any within &#39;resources&#39; if &#39;resources&#39; is non-empty. |
| removed_resources | [][string](#string) | repeated | Resources names of resources that have be deleted and to be removed from the xDS Client.
Removed resources for missing resources can be ignored. |
| nonce | [string](#string) |  | The nonce provides a way for DeltaDiscoveryRequests to uniquely
reference a DeltaDiscoveryResponse when (N)ACKing. The nonce is required. |






<a name="envoy.api.v2.DiscoveryRequest"></a>

### DiscoveryRequest
A DiscoveryRequest requests a set of versioned resources of the same type for
a given Envoy node on some API.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version_info | [string](#string) |  | The version_info provided in the request messages will be the version_info
received with the most recent successfully processed response or empty on
the first request. It is expected that no new request is sent after a
response is received until the Envoy instance is ready to ACK/NACK the new
configuration. ACK/NACK takes place by returning the new API config version
as applied or the previous API config version respectively. Each type_url
(see below) has an independent version associated with it. |
| node | [core.Node](#envoy.api.v2.core.Node) |  | The node making the request. |
| resource_names | [][string](#string) | repeated | List of resources to subscribe to, e.g. list of cluster names or a route
configuration name. If this is empty, all resources for the API are
returned. LDS/CDS may have empty resource_names, which will cause all
resources for the Envoy instance to be returned. The LDS and CDS responses
will then imply a number of resources that need to be fetched via EDS/RDS,
which will be explicitly enumerated in resource_names. |
| type_url | [string](#string) |  | Type of the resource that is being requested, e.g.
&#34;type.googleapis.com/envoy.api.v2.ClusterLoadAssignment&#34;. This is implicit
in requests made via singleton xDS APIs such as CDS, LDS, etc. but is
required for ADS. |
| response_nonce | [string](#string) |  | nonce corresponding to DiscoveryResponse being ACK/NACKed. See above
discussion on version_info and the DiscoveryResponse nonce comment. This
may be empty only if 1) this is a non-persistent-stream xDS such as HTTP,
or 2) the client has not yet accepted an update in this xDS stream (unlike
delta, where it is populated only for new explicit ACKs). |
| error_detail | [google.rpc.Status](#google.rpc.Status) |  | This is populated when the previous `DiscoveryResponse (envoy_api_msg_DiscoveryResponse)`
failed to update configuration. The *message* field in *error_details* provides the Envoy
internal exception related to the failure. It is only intended for consumption during manual
debugging, the string provided is not guaranteed to be stable across Envoy versions. |






<a name="envoy.api.v2.DiscoveryResponse"></a>

### DiscoveryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version_info | [string](#string) |  | The version of the response data. |
| resources | [][google.protobuf.Any](#google.protobuf.Any) | repeated | The response resources. These resources are typed and depend on the API being called. |
| canary | [bool](#bool) |  | [#not-implemented-hide:]
Canary is used to support two Envoy command line flags:

* --terminate-on-canary-transition-failure. When set, Envoy is able to
  terminate if it detects that configuration is stuck at canary. Consider
  this example sequence of updates:
  - Management server applies a canary config successfully.
  - Management server rolls back to a production config.
  - Envoy rejects the new production config.
  Since there is no sensible way to continue receiving configuration
  updates, Envoy will then terminate and apply production config from a
  clean slate.
* --dry-run-canary. When set, a canary response will never be applied, only
  validated via a dry run. |
| type_url | [string](#string) |  | Type URL for resources. Identifies the xDS API when muxing over ADS.
Must be consistent with the type_url in the &#39;resources&#39; repeated Any (if non-empty). |
| nonce | [string](#string) |  | For gRPC based subscriptions, the nonce provides a way to explicitly ack a
specific DiscoveryResponse in a following DiscoveryRequest. Additional
messages may have been sent by Envoy to the management server for the
previous version on the stream prior to this DiscoveryResponse, that were
unprocessed at response send time. The nonce allows the management server
to ignore any further DiscoveryRequests for the previous version until a
DiscoveryRequest bearing the nonce. The nonce is optional and is not
required for non-stream based xDS implementations. |
| control_plane | [core.ControlPlane](#envoy.api.v2.core.ControlPlane) |  | [#not-implemented-hide:]
The control plane instance that sent the response. |






<a name="envoy.api.v2.Resource"></a>

### Resource



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The resource&#39;s name, to distinguish it from others of the same type of resource. |
| aliases | [][string](#string) | repeated | [#not-implemented-hide:]
The aliases are a list of other names that this resource can go by. |
| version | [string](#string) |  | The resource level version. It allows xDS to track the state of individual
resources. |
| resource | [google.protobuf.Any](#google.protobuf.Any) |  | The resource being tracked. |





 

 

 

 

