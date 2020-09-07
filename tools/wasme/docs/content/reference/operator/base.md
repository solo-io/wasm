
---
title: "envoy.api.v2.coreenvoy/api/v2/core/base.proto"
---

## Package : `envoy.api.v2.core`



<a name="top"></a>

<a name="API Reference for envoy/api/v2/core/base.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## envoy/api/v2/core/base.proto


## Table of Contents
  - [AsyncDataSource](#envoy.api.v2.core.AsyncDataSource)
  - [ControlPlane](#envoy.api.v2.core.ControlPlane)
  - [DataSource](#envoy.api.v2.core.DataSource)
  - [HeaderMap](#envoy.api.v2.core.HeaderMap)
  - [HeaderValue](#envoy.api.v2.core.HeaderValue)
  - [HeaderValueOption](#envoy.api.v2.core.HeaderValueOption)
  - [Locality](#envoy.api.v2.core.Locality)
  - [Metadata](#envoy.api.v2.core.Metadata)
  - [Metadata.FilterMetadataEntry](#envoy.api.v2.core.Metadata.FilterMetadataEntry)
  - [Node](#envoy.api.v2.core.Node)
  - [RemoteDataSource](#envoy.api.v2.core.RemoteDataSource)
  - [RuntimeFeatureFlag](#envoy.api.v2.core.RuntimeFeatureFlag)
  - [RuntimeFractionalPercent](#envoy.api.v2.core.RuntimeFractionalPercent)
  - [RuntimeUInt32](#envoy.api.v2.core.RuntimeUInt32)
  - [SocketOption](#envoy.api.v2.core.SocketOption)
  - [TransportSocket](#envoy.api.v2.core.TransportSocket)

  - [RequestMethod](#envoy.api.v2.core.RequestMethod)
  - [RoutingPriority](#envoy.api.v2.core.RoutingPriority)
  - [SocketOption.SocketState](#envoy.api.v2.core.SocketOption.SocketState)
  - [TrafficDirection](#envoy.api.v2.core.TrafficDirection)






<a name="envoy.api.v2.core.AsyncDataSource"></a>

### AsyncDataSource
Async data source which support async data fetch.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| local | [DataSource](#envoy.api.v2.core.DataSource) |  | Local async data source. |
| remote | [RemoteDataSource](#envoy.api.v2.core.RemoteDataSource) |  | Remote async data source. |






<a name="envoy.api.v2.core.ControlPlane"></a>

### ControlPlane
Identifies a specific ControlPlane instance that Envoy is connected to.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| identifier | [string](#string) |  | An opaque control plane identifier that uniquely identifies an instance
of control plane. This can be used to identify which control plane instance,
the Envoy is connected to. |






<a name="envoy.api.v2.core.DataSource"></a>

### DataSource
Data source consisting of either a file or an inline value.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| filename | [string](#string) |  | Local filesystem data source. |
| inline_bytes | [bytes](#bytes) |  | Bytes inlined in the configuration. |
| inline_string | [string](#string) |  | String inlined in the configuration. |






<a name="envoy.api.v2.core.HeaderMap"></a>

### HeaderMap
Wrapper for a set of headers.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| headers | [][HeaderValue](#envoy.api.v2.core.HeaderValue) | repeated |  |






<a name="envoy.api.v2.core.HeaderValue"></a>

### HeaderValue
Header name/value pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  | Header name. |
| value | [string](#string) |  | Header value.

The same `format specifier (config_access_log_format)` as used for
`HTTP access logging (config_access_log)` applies here, however
unknown header values are replaced with the empty string instead of `-`. |






<a name="envoy.api.v2.core.HeaderValueOption"></a>

### HeaderValueOption
Header name/value pair plus option to control append behavior.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| header | [HeaderValue](#envoy.api.v2.core.HeaderValue) |  | Header name/value pair that this option applies to. |
| append | [google.protobuf.BoolValue](#google.protobuf.BoolValue) |  | Should the value be appended? If true (default), the value is appended to
existing values. |






<a name="envoy.api.v2.core.Locality"></a>

### Locality
Identifies location of where either Envoy runs or where upstream hosts run.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| region | [string](#string) |  | Region this `zone (envoy_api_field_core.Locality.zone)` belongs to. |
| zone | [string](#string) |  | Defines the local service zone where Envoy is running. Though optional, it
should be set if discovery service routing is used and the discovery
service exposes `zone data (envoy_api_field_endpoint.LocalityLbEndpoints.locality)`,
either in this message or via :option:`--service-zone`. The meaning of zone
is context dependent, e.g. `Availability Zone (AZ)
&lt;https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html&gt;`_
on AWS, `Zone &lt;https://cloud.google.com/compute/docs/regions-zones/&gt;`_ on
GCP, etc. |
| sub_zone | [string](#string) |  | When used for locality of upstream hosts, this field further splits zone
into smaller chunks of sub-zones so they can be load balanced
independently. |






<a name="envoy.api.v2.core.Metadata"></a>

### Metadata
Metadata provides additional inputs to filters based on matched listeners,
filter chains, routes and endpoints. It is structured as a map, usually from
filter name (in reverse DNS format) to metadata specific to the filter. Metadata
key-values for a filter are merged as connection and request handling occurs,
with later values for the same key overriding earlier values.

An example use of metadata is providing additional values to
http_connection_manager in the envoy.http_connection_manager.access_log
namespace.

Another example use of metadata is to per service config info in cluster metadata, which may get
consumed by multiple filters.

For load balancing, Metadata provides a means to subset cluster endpoints.
Endpoints have a Metadata object associated and routes contain a Metadata
object to match against. There are some well defined metadata used today for
this purpose:

* ``{&#34;envoy.lb&#34;: {&#34;canary&#34;: &lt;bool&gt; }}`` This indicates the canary status of an
  endpoint and is also used during header processing
  (x-envoy-upstream-canary) and for stats purposes.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| filter_metadata | [][Metadata.FilterMetadataEntry](#envoy.api.v2.core.Metadata.FilterMetadataEntry) | repeated | Key is the reverse DNS filter name, e.g. com.acme.widget. The envoy.*
namespace is reserved for Envoy&#39;s built-in filters. |






<a name="envoy.api.v2.core.Metadata.FilterMetadataEntry"></a>

### Metadata.FilterMetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [google.protobuf.Struct](#google.protobuf.Struct) |  |  |






<a name="envoy.api.v2.core.Node"></a>

### Node
Identifies a specific Envoy instance. The node identifier is presented to the
management server, which may use this identifier to distinguish per Envoy
configuration for serving.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | An opaque node identifier for the Envoy node. This also provides the local
service node name. It should be set if any of the following features are
used: `statsd (arch_overview_statistics)`, `CDS
(config_cluster_manager_cds)`, and `HTTP tracing
(arch_overview_tracing)`, either in this message or via
:option:`--service-node`. |
| cluster | [string](#string) |  | Defines the local service cluster name where Envoy is running. Though
optional, it should be set if any of the following features are used:
`statsd (arch_overview_statistics)`, `health check cluster
verification (envoy_api_field_core.HealthCheck.HttpHealthCheck.service_name)`,
`runtime override directory (envoy_api_msg_config.bootstrap.v2.Runtime)`,
`user agent addition
(envoy_api_field_config.filter.network.http_connection_manager.v2.HttpConnectionManager.add_user_agent)`,
`HTTP global rate limiting (config_http_filters_rate_limit)`,
`CDS (config_cluster_manager_cds)`, and `HTTP tracing
(arch_overview_tracing)`, either in this message or via
:option:`--service-cluster`. |
| metadata | [google.protobuf.Struct](#google.protobuf.Struct) |  | Opaque metadata extending the node identifier. Envoy will pass this
directly to the management server. |
| locality | [Locality](#envoy.api.v2.core.Locality) |  | Locality specifying where the Envoy instance is running. |
| build_version | [string](#string) |  | This is motivated by informing a management server during canary which
version of Envoy is being tested in a heterogeneous fleet. This will be set
by Envoy in management server RPCs. |






<a name="envoy.api.v2.core.RemoteDataSource"></a>

### RemoteDataSource
The message specifies how to fetch data from remote and how to verify it.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| http_uri | [HttpUri](#envoy.api.v2.core.HttpUri) |  | The HTTP URI to fetch the remote data. |
| sha256 | [string](#string) |  | SHA256 string for verifying data. |






<a name="envoy.api.v2.core.RuntimeFeatureFlag"></a>

### RuntimeFeatureFlag
Runtime derived bool with a default when not specified.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_value | [google.protobuf.BoolValue](#google.protobuf.BoolValue) |  | Default value if runtime value is not available. |
| runtime_key | [string](#string) |  | Runtime key to get value for comparison. This value is used if defined. The boolean value must
be represented via its
`canonical JSON encoding &lt;https://developers.google.com/protocol-buffers/docs/proto3#json&gt;`_. |






<a name="envoy.api.v2.core.RuntimeFractionalPercent"></a>

### RuntimeFractionalPercent
Runtime derived FractionalPercent with defaults for when the numerator or denominator is not
specified via a runtime key.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_value | [envoy.type.FractionalPercent](#envoy.type.FractionalPercent) |  | Default value if the runtime value&#39;s for the numerator/denominator keys are not available. |
| runtime_key | [string](#string) |  | Runtime key for a YAML representation of a FractionalPercent. |






<a name="envoy.api.v2.core.RuntimeUInt32"></a>

### RuntimeUInt32
Runtime derived uint32 with a default when not specified.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| default_value | [uint32](#uint32) |  | Default value if runtime value is not available. |
| runtime_key | [string](#string) |  | Runtime key to get value for comparison. This value is used if defined. |






<a name="envoy.api.v2.core.SocketOption"></a>

### SocketOption
Generic socket option message. This would be used to set socket options that
might not exist in upstream kernels or precompiled Envoy binaries.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| description | [string](#string) |  | An optional name to give this socket option for debugging, etc.
Uniqueness is not required and no special meaning is assumed. |
| level | [int64](#int64) |  | Corresponding to the level value passed to setsockopt, such as IPPROTO_TCP |
| name | [int64](#int64) |  | The numeric name as passed to setsockopt |
| int_value | [int64](#int64) |  | Because many sockopts take an int value. |
| buf_value | [bytes](#bytes) |  | Otherwise it&#39;s a byte buffer. |
| state | [SocketOption.SocketState](#envoy.api.v2.core.SocketOption.SocketState) |  | The state in which the option will be applied. When used in BindConfig
STATE_PREBIND is currently the only valid value. |






<a name="envoy.api.v2.core.TransportSocket"></a>

### TransportSocket
Configuration for transport socket in `listeners (config_listeners)` and
`clusters (envoy_api_msg_Cluster)`. If the configuration is
empty, a default transport socket implementation and configuration will be
chosen based on the platform and existence of tls_context.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The name of the transport socket to instantiate. The name must match a supported transport
socket implementation. |
| config | [google.protobuf.Struct](#google.protobuf.Struct) |  |  |
| typed_config | [google.protobuf.Any](#google.protobuf.Any) |  |  |





 


<a name="envoy.api.v2.core.RequestMethod"></a>

### RequestMethod
HTTP request method.

| Name | Number | Description |
| ---- | ------ | ----------- |
| METHOD_UNSPECIFIED | 0 |  |
| GET | 1 |  |
| HEAD | 2 |  |
| POST | 3 |  |
| PUT | 4 |  |
| DELETE | 5 |  |
| CONNECT | 6 |  |
| OPTIONS | 7 |  |
| TRACE | 8 |  |
| PATCH | 9 |  |



<a name="envoy.api.v2.core.RoutingPriority"></a>

### RoutingPriority
Envoy supports `upstream priority routing
(arch_overview_http_routing_priority)` both at the route and the virtual
cluster level. The current priority implementation uses different connection
pool and circuit breaking settings for each priority level. This means that
even for HTTP/2 requests, two physical connections will be used to an
upstream host. In the future Envoy will likely support true HTTP/2 priority
over a single upstream connection.

| Name | Number | Description |
| ---- | ------ | ----------- |
| DEFAULT | 0 |  |
| HIGH | 1 |  |



<a name="envoy.api.v2.core.SocketOption.SocketState"></a>

### SocketOption.SocketState


| Name | Number | Description |
| ---- | ------ | ----------- |
| STATE_PREBIND | 0 | Socket options are applied after socket creation but before binding the socket to a port |
| STATE_BOUND | 1 | Socket options are applied after binding the socket to a port but before calling listen() |
| STATE_LISTENING | 2 | Socket options are applied after calling listen() |



<a name="envoy.api.v2.core.TrafficDirection"></a>

### TrafficDirection
Identifies the direction of the traffic relative to the local Envoy.

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNSPECIFIED | 0 | Default option is unspecified. |
| INBOUND | 1 | The transport is used for incoming traffic. |
| OUTBOUND | 2 | The transport is used for outgoing traffic. |


 

 

 

