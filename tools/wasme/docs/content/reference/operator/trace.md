
---
title: "opencensus.proto.tracegithub.com/solo-io/solo-kit/api/external/trace.proto"
---

## Package : `opencensus.proto.trace`



<a name="top"></a>

<a name="API Reference for github.com/solo-io/solo-kit/api/external/trace.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/solo-io/solo-kit/api/external/trace.proto


## Table of Contents
  - [AttributeValue](#opencensus.proto.trace.AttributeValue)
  - [Module](#opencensus.proto.trace.Module)
  - [Span](#opencensus.proto.trace.Span)
  - [Span.Attributes](#opencensus.proto.trace.Span.Attributes)
  - [Span.Attributes.AttributeMapEntry](#opencensus.proto.trace.Span.Attributes.AttributeMapEntry)
  - [Span.Link](#opencensus.proto.trace.Span.Link)
  - [Span.Links](#opencensus.proto.trace.Span.Links)
  - [Span.TimeEvent](#opencensus.proto.trace.Span.TimeEvent)
  - [Span.TimeEvent.Annotation](#opencensus.proto.trace.Span.TimeEvent.Annotation)
  - [Span.TimeEvent.MessageEvent](#opencensus.proto.trace.Span.TimeEvent.MessageEvent)
  - [Span.TimeEvents](#opencensus.proto.trace.Span.TimeEvents)
  - [StackTrace](#opencensus.proto.trace.StackTrace)
  - [StackTrace.StackFrame](#opencensus.proto.trace.StackTrace.StackFrame)
  - [StackTrace.StackFrames](#opencensus.proto.trace.StackTrace.StackFrames)
  - [Status](#opencensus.proto.trace.Status)
  - [TruncatableString](#opencensus.proto.trace.TruncatableString)

  - [Span.Link.Type](#opencensus.proto.trace.Span.Link.Type)
  - [Span.SpanKind](#opencensus.proto.trace.Span.SpanKind)
  - [Span.TimeEvent.MessageEvent.Type](#opencensus.proto.trace.Span.TimeEvent.MessageEvent.Type)






<a name="opencensus.proto.trace.AttributeValue"></a>

### AttributeValue
The value of an Attribute.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| string_value | [TruncatableString](#opencensus.proto.trace.TruncatableString) |  | A string up to 256 bytes long. |
| int_value | [int64](#int64) |  | A 64-bit signed integer. |
| bool_value | [bool](#bool) |  | A Boolean value represented by `true` or `false`. |






<a name="opencensus.proto.trace.Module"></a>

### Module
A description of a binary module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| module | [TruncatableString](#opencensus.proto.trace.TruncatableString) |  | TODO: document the meaning of this field.
For example: main binary, kernel modules, and dynamic libraries
such as libc.so, sharedlib.so. |
| build_id | [TruncatableString](#opencensus.proto.trace.TruncatableString) |  | A unique identifier for the module, usually a hash of its
contents. |






<a name="opencensus.proto.trace.Span"></a>

### Span
A span represents a single operation within a trace. Spans can be
nested to form a trace tree. Often, a trace contains a root span
that describes the end-to-end latency, and one or more subspans for
its sub-operations. A trace can also contain multiple root spans,
or none at all. Spans do not need to be contiguous - there may be
gaps or overlaps between spans in a trace.

The next id is 15.
TODO(bdrutu): Add an example.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trace_id | [bytes](#bytes) |  | A unique identifier for a trace. All spans from the same trace share
the same `trace_id`. The ID is a 16-byte array.

This field is required. |
| span_id | [bytes](#bytes) |  | A unique identifier for a span within a trace, assigned when the span
is created. The ID is an 8-byte array.

This field is required. |
| parent_span_id | [bytes](#bytes) |  | The `span_id` of this span&#39;s parent span. If this is a root span, then this
field must be empty. The ID is an 8-byte array. |
| name | [TruncatableString](#opencensus.proto.trace.TruncatableString) |  | A description of the span&#39;s operation.

For example, the name can be a qualified method name or a file name
and a line number where the operation is called. A best practice is to use
the same display name at the same call point in an application.
This makes it easier to correlate spans in different traces.

This field is required. |
| kind | [Span.SpanKind](#opencensus.proto.trace.Span.SpanKind) |  | Distinguishes between spans generated in a particular context. For example,
two spans with the same name may be distinguished using `CLIENT`
and `SERVER` to identify queueing latency associated with the span. |
| start_time | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | The start time of the span. On the client side, this is the time kept by
the local machine where the span execution starts. On the server side, this
is the time when the server&#39;s application handler starts running. |
| end_time | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | The end time of the span. On the client side, this is the time kept by
the local machine where the span execution ends. On the server side, this
is the time when the server application handler stops running. |
| attributes | [Span.Attributes](#opencensus.proto.trace.Span.Attributes) |  | A set of attributes on the span. |
| stack_trace | [StackTrace](#opencensus.proto.trace.StackTrace) |  | A stack trace captured at the start of the span. |
| time_events | [Span.TimeEvents](#opencensus.proto.trace.Span.TimeEvents) |  | The included time events. |
| links | [Span.Links](#opencensus.proto.trace.Span.Links) |  | The inclued links. |
| status | [Status](#opencensus.proto.trace.Status) |  | An optional final status for this span. |
| same_process_as_parent_span | [google.protobuf.BoolValue](#google.protobuf.BoolValue) |  | A highly recommended but not required flag that identifies when a trace
crosses a process boundary. True when the parent_span belongs to the
same process as the current span. |
| child_span_count | [google.protobuf.UInt32Value](#google.protobuf.UInt32Value) |  | An optional number of child spans that were generated while this span
was active. If set, allows an implementation to detect missing child spans. |






<a name="opencensus.proto.trace.Span.Attributes"></a>

### Span.Attributes
A set of attributes, each with a key and a value.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| attribute_map | [][Span.Attributes.AttributeMapEntry](#opencensus.proto.trace.Span.Attributes.AttributeMapEntry) | repeated | The set of attributes. The value can be a string, an integer, or the
Boolean values `true` and `false`. For example:

    &#34;/instance_id&#34;: &#34;my-instance&#34;
    &#34;/http/user_agent&#34;: &#34;&#34;
    &#34;/http/server_latency&#34;: 300
    &#34;abc.com/myattribute&#34;: true |
| dropped_attributes_count | [int32](#int32) |  | The number of attributes that were discarded. Attributes can be discarded
because their keys are too long or because there are too many attributes.
If this value is 0, then no attributes were dropped. |






<a name="opencensus.proto.trace.Span.Attributes.AttributeMapEntry"></a>

### Span.Attributes.AttributeMapEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [AttributeValue](#opencensus.proto.trace.AttributeValue) |  |  |






<a name="opencensus.proto.trace.Span.Link"></a>

### Span.Link
A pointer from the current span to another span in the same trace or in a
different trace. For example, this can be used in batching operations,
where a single batch handler processes multiple requests from different
traces or when the handler receives a request from a different project.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trace_id | [bytes](#bytes) |  | A unique identifier for a trace. All spans from the same trace share
the same `trace_id`. The ID is a 16-byte array. |
| span_id | [bytes](#bytes) |  | A unique identifier for a span within a trace, assigned when the span
is created. The ID is an 8-byte array. |
| type | [Span.Link.Type](#opencensus.proto.trace.Span.Link.Type) |  | The relationship of the current span relative to the linked span. |
| attributes | [Span.Attributes](#opencensus.proto.trace.Span.Attributes) |  | A set of attributes on the link. |






<a name="opencensus.proto.trace.Span.Links"></a>

### Span.Links
A collection of links, which are references from this span to a span
in the same or different trace.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| link | [][Span.Link](#opencensus.proto.trace.Span.Link) | repeated | A collection of links. |
| dropped_links_count | [int32](#int32) |  | The number of dropped links after the maximum size was enforced. If
this value is 0, then no links were dropped. |






<a name="opencensus.proto.trace.Span.TimeEvent"></a>

### Span.TimeEvent
A time-stamped annotation or message event in the Span.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| time | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | The time the event occurred. |
| annotation | [Span.TimeEvent.Annotation](#opencensus.proto.trace.Span.TimeEvent.Annotation) |  | A text annotation with a set of attributes. |
| message_event | [Span.TimeEvent.MessageEvent](#opencensus.proto.trace.Span.TimeEvent.MessageEvent) |  | An event describing a message sent/received between Spans. |






<a name="opencensus.proto.trace.Span.TimeEvent.Annotation"></a>

### Span.TimeEvent.Annotation
A text annotation with a set of attributes.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| description | [TruncatableString](#opencensus.proto.trace.TruncatableString) |  | A user-supplied message describing the event. |
| attributes | [Span.Attributes](#opencensus.proto.trace.Span.Attributes) |  | A set of attributes on the annotation. |






<a name="opencensus.proto.trace.Span.TimeEvent.MessageEvent"></a>

### Span.TimeEvent.MessageEvent
An event describing a message sent/received between Spans.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [Span.TimeEvent.MessageEvent.Type](#opencensus.proto.trace.Span.TimeEvent.MessageEvent.Type) |  | The type of MessageEvent. Indicates whether the message was sent or
received. |
| id | [uint64](#uint64) |  | An identifier for the MessageEvent&#39;s message that can be used to match
SENT and RECEIVED MessageEvents. For example, this field could
represent a sequence ID for a streaming RPC. It is recommended to be
unique within a Span. |
| uncompressed_size | [uint64](#uint64) |  | The number of uncompressed bytes sent or received. |
| compressed_size | [uint64](#uint64) |  | The number of compressed bytes sent or received. If zero, assumed to
be the same size as uncompressed. |






<a name="opencensus.proto.trace.Span.TimeEvents"></a>

### Span.TimeEvents
A collection of `TimeEvent`s. A `TimeEvent` is a time-stamped annotation
on the span, consisting of either user-supplied key-value pairs, or
details of a message sent/received between Spans.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| time_event | [][Span.TimeEvent](#opencensus.proto.trace.Span.TimeEvent) | repeated | A collection of `TimeEvent`s. |
| dropped_annotations_count | [int32](#int32) |  | The number of dropped annotations in all the included time events.
If the value is 0, then no annotations were dropped. |
| dropped_message_events_count | [int32](#int32) |  | The number of dropped message events in all the included time events.
If the value is 0, then no message events were dropped. |






<a name="opencensus.proto.trace.StackTrace"></a>

### StackTrace
The call stack which originated this span.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| stack_frames | [StackTrace.StackFrames](#opencensus.proto.trace.StackTrace.StackFrames) |  | Stack frames in this stack trace. |
| stack_trace_hash_id | [uint64](#uint64) |  | The hash ID is used to conserve network bandwidth for duplicate
stack traces within a single trace.

Often multiple spans will have identical stack traces.
The first occurrence of a stack trace should contain both
`stack_frames` and a value in `stack_trace_hash_id`.

Subsequent spans within the same request can refer
to that stack trace by setting only `stack_trace_hash_id`.

TODO: describe how to deal with the case where stack_trace_hash_id is
zero because it was not set. |






<a name="opencensus.proto.trace.StackTrace.StackFrame"></a>

### StackTrace.StackFrame
A single stack frame in a stack trace.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| function_name | [TruncatableString](#opencensus.proto.trace.TruncatableString) |  | The fully-qualified name that uniquely identifies the function or
method that is active in this frame. |
| original_function_name | [TruncatableString](#opencensus.proto.trace.TruncatableString) |  | An un-mangled function name, if `function_name` is
[mangled](http://www.avabodh.com/cxxin/namemangling.html). The name can
be fully qualified. |
| file_name | [TruncatableString](#opencensus.proto.trace.TruncatableString) |  | The name of the source file where the function call appears. |
| line_number | [int64](#int64) |  | The line number in `file_name` where the function call appears. |
| column_number | [int64](#int64) |  | The column number where the function call appears, if available.
This is important in JavaScript because of its anonymous functions. |
| load_module | [Module](#opencensus.proto.trace.Module) |  | The binary module from where the code was loaded. |
| source_version | [TruncatableString](#opencensus.proto.trace.TruncatableString) |  | The version of the deployed source code. |






<a name="opencensus.proto.trace.StackTrace.StackFrames"></a>

### StackTrace.StackFrames
A collection of stack frames, which can be truncated.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| frame | [][StackTrace.StackFrame](#opencensus.proto.trace.StackTrace.StackFrame) | repeated | Stack frames in this call stack. |
| dropped_frames_count | [int32](#int32) |  | The number of stack frames that were dropped because there
were too many stack frames.
If this value is 0, then no stack frames were dropped. |






<a name="opencensus.proto.trace.Status"></a>

### Status
The `Status` type defines a logical error model that is suitable for different
programming environments, including REST APIs and RPC APIs. This proto&#39;s fields
are a subset of those of
[google.rpc.Status](https://github.com/googleapis/googleapis/blob/master/google/rpc/status.proto),
which is used by [gRPC](https://github.com/grpc).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [int32](#int32) |  | The status code. |
| message | [string](#string) |  | A developer-facing error message, which should be in English. |






<a name="opencensus.proto.trace.TruncatableString"></a>

### TruncatableString
A string that might be shortened to a specified length.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [string](#string) |  | The shortened string. For example, if the original string was 500 bytes long and
the limit of the string was 128 bytes, then this value contains the first 128
bytes of the 500-byte string. Note that truncation always happens on a
character boundary, to ensure that a truncated string is still valid UTF-8.
Because it may contain multi-byte characters, the size of the truncated string
may be less than the truncation limit. |
| truncated_byte_count | [int32](#int32) |  | The number of bytes removed from the original string. If this
value is 0, then the string was not shortened. |





 


<a name="opencensus.proto.trace.Span.Link.Type"></a>

### Span.Link.Type
The relationship of the current span relative to the linked span: child,
parent, or unspecified.

| Name | Number | Description |
| ---- | ------ | ----------- |
| TYPE_UNSPECIFIED | 0 | The relationship of the two spans is unknown, or known but other
than parent-child. |
| CHILD_LINKED_SPAN | 1 | The linked span is a child of the current span. |
| PARENT_LINKED_SPAN | 2 | The linked span is a parent of the current span. |



<a name="opencensus.proto.trace.Span.SpanKind"></a>

### Span.SpanKind
Type of span. Can be used to specify additional relationships between spans
in addition to a parent/child relationship.

| Name | Number | Description |
| ---- | ------ | ----------- |
| SPAN_KIND_UNSPECIFIED | 0 | Unspecified. |
| SERVER | 1 | Indicates that the span covers server-side handling of an RPC or other
remote network request. |
| CLIENT | 2 | Indicates that the span covers the client-side wrapper around an RPC or
other remote request. |



<a name="opencensus.proto.trace.Span.TimeEvent.MessageEvent.Type"></a>

### Span.TimeEvent.MessageEvent.Type
Indicates whether the message was sent or received.

| Name | Number | Description |
| ---- | ------ | ----------- |
| TYPE_UNSPECIFIED | 0 | Unknown event type. |
| SENT | 1 | Indicates a sent message. |
| RECEIVED | 2 | Indicates a received message. |


 

 

 

