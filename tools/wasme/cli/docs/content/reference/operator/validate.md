
---
title: "validategithub.com/envoyproxy/protoc-gen-validate/validate/validate.proto"
---

## Package : `validate`



<a name="top"></a>

<a name="API Reference for github.com/envoyproxy/protoc-gen-validate/validate/validate.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/envoyproxy/protoc-gen-validate/validate/validate.proto


## Table of Contents
  - [AnyRules](#validate.AnyRules)
  - [BoolRules](#validate.BoolRules)
  - [BytesRules](#validate.BytesRules)
  - [DoubleRules](#validate.DoubleRules)
  - [DurationRules](#validate.DurationRules)
  - [EnumRules](#validate.EnumRules)
  - [FieldRules](#validate.FieldRules)
  - [Fixed32Rules](#validate.Fixed32Rules)
  - [Fixed64Rules](#validate.Fixed64Rules)
  - [FloatRules](#validate.FloatRules)
  - [Int32Rules](#validate.Int32Rules)
  - [Int64Rules](#validate.Int64Rules)
  - [MapRules](#validate.MapRules)
  - [MessageRules](#validate.MessageRules)
  - [RepeatedRules](#validate.RepeatedRules)
  - [SFixed32Rules](#validate.SFixed32Rules)
  - [SFixed64Rules](#validate.SFixed64Rules)
  - [SInt32Rules](#validate.SInt32Rules)
  - [SInt64Rules](#validate.SInt64Rules)
  - [StringRules](#validate.StringRules)
  - [TimestampRules](#validate.TimestampRules)
  - [UInt32Rules](#validate.UInt32Rules)
  - [UInt64Rules](#validate.UInt64Rules)

  - [KnownRegex](#validate.KnownRegex)

  - [File-level Extensions](#github.com/envoyproxy/protoc-gen-validate/validate/validate.proto-extensions)
  - [File-level Extensions](#github.com/envoyproxy/protoc-gen-validate/validate/validate.proto-extensions)
  - [File-level Extensions](#github.com/envoyproxy/protoc-gen-validate/validate/validate.proto-extensions)





<a name="validate.AnyRules"></a>

### AnyRules
AnyRules describe constraints applied exclusively to the
`google.protobuf.Any` well-known type


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| required | [bool](#bool) | optional | Required specifies that this field must be set |
| in | [][string](#string) | repeated | In specifies that this field&#39;s `type_url` must be equal to one of the
specified values. |
| not_in | [][string](#string) | repeated | NotIn specifies that this field&#39;s `type_url` must not be equal to any of
the specified values. |






<a name="validate.BoolRules"></a>

### BoolRules
BoolRules describes the constraints applied to `bool` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [bool](#bool) | optional | Const specifies that this field must be exactly the specified value |






<a name="validate.BytesRules"></a>

### BytesRules
BytesRules describe the constraints applied to `bytes` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [bytes](#bytes) | optional | Const specifies that this field must be exactly the specified value |
| len | [uint64](#uint64) | optional | Len specifies that this field must be the specified number of bytes |
| min_len | [uint64](#uint64) | optional | MinLen specifies that this field must be the specified number of bytes
at a minimum |
| max_len | [uint64](#uint64) | optional | MaxLen specifies that this field must be the specified number of bytes
at a maximum |
| pattern | [string](#string) | optional | Pattern specifes that this field must match against the specified
regular expression (RE2 syntax). The included expression should elide
any delimiters. |
| prefix | [bytes](#bytes) | optional | Prefix specifies that this field must have the specified bytes at the
beginning of the string. |
| suffix | [bytes](#bytes) | optional | Suffix specifies that this field must have the specified bytes at the
end of the string. |
| contains | [bytes](#bytes) | optional | Contains specifies that this field must have the specified bytes
anywhere in the string. |
| in | [][bytes](#bytes) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][bytes](#bytes) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |
| ip | [bool](#bool) | optional | Ip specifies that the field must be a valid IP (v4 or v6) address in
byte format |
| ipv4 | [bool](#bool) | optional | Ipv4 specifies that the field must be a valid IPv4 address in byte
format |
| ipv6 | [bool](#bool) | optional | Ipv6 specifies that the field must be a valid IPv6 address in byte
format |






<a name="validate.DoubleRules"></a>

### DoubleRules
DoubleRules describes the constraints applied to `double` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [double](#double) | optional | Const specifies that this field must be exactly the specified value |
| lt | [double](#double) | optional | Lt specifies that this field must be less than the specified value,
exclusive |
| lte | [double](#double) | optional | Lte specifies that this field must be less than or equal to the
specified value, inclusive |
| gt | [double](#double) | optional | Gt specifies that this field must be greater than the specified value,
exclusive. If the value of Gt is larger than a specified Lt or Lte, the
range is reversed. |
| gte | [double](#double) | optional | Gte specifies that this field must be greater than or equal to the
specified value, inclusive. If the value of Gte is larger than a
specified Lt or Lte, the range is reversed. |
| in | [][double](#double) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][double](#double) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |






<a name="validate.DurationRules"></a>

### DurationRules
DurationRules describe the constraints applied exclusively to the
`google.protobuf.Duration` well-known type


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| required | [bool](#bool) | optional | Required specifies that this field must be set |
| const | [google.protobuf.Duration](#google.protobuf.Duration) | optional | Const specifies that this field must be exactly the specified value |
| lt | [google.protobuf.Duration](#google.protobuf.Duration) | optional | Lt specifies that this field must be less than the specified value,
exclusive |
| lte | [google.protobuf.Duration](#google.protobuf.Duration) | optional | Lt specifies that this field must be less than the specified value,
inclusive |
| gt | [google.protobuf.Duration](#google.protobuf.Duration) | optional | Gt specifies that this field must be greater than the specified value,
exclusive |
| gte | [google.protobuf.Duration](#google.protobuf.Duration) | optional | Gte specifies that this field must be greater than the specified value,
inclusive |
| in | [][google.protobuf.Duration](#google.protobuf.Duration) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][google.protobuf.Duration](#google.protobuf.Duration) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |






<a name="validate.EnumRules"></a>

### EnumRules
EnumRules describe the constraints applied to enum values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [int32](#int32) | optional | Const specifies that this field must be exactly the specified value |
| defined_only | [bool](#bool) | optional | DefinedOnly specifies that this field must be only one of the defined
values for this enum, failing on any undefined value. |
| in | [][int32](#int32) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][int32](#int32) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |






<a name="validate.FieldRules"></a>

### FieldRules
FieldRules encapsulates the rules for each type of field. Depending on the
field, the correct set should be used to ensure proper validations.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message | [MessageRules](#validate.MessageRules) | optional |  |
| float | [FloatRules](#validate.FloatRules) | optional | Scalar Field Types |
| double | [DoubleRules](#validate.DoubleRules) | optional |  |
| int32 | [Int32Rules](#validate.Int32Rules) | optional |  |
| int64 | [Int64Rules](#validate.Int64Rules) | optional |  |
| uint32 | [UInt32Rules](#validate.UInt32Rules) | optional |  |
| uint64 | [UInt64Rules](#validate.UInt64Rules) | optional |  |
| sint32 | [SInt32Rules](#validate.SInt32Rules) | optional |  |
| sint64 | [SInt64Rules](#validate.SInt64Rules) | optional |  |
| fixed32 | [Fixed32Rules](#validate.Fixed32Rules) | optional |  |
| fixed64 | [Fixed64Rules](#validate.Fixed64Rules) | optional |  |
| sfixed32 | [SFixed32Rules](#validate.SFixed32Rules) | optional |  |
| sfixed64 | [SFixed64Rules](#validate.SFixed64Rules) | optional |  |
| bool | [BoolRules](#validate.BoolRules) | optional |  |
| string | [StringRules](#validate.StringRules) | optional |  |
| bytes | [BytesRules](#validate.BytesRules) | optional |  |
| enum | [EnumRules](#validate.EnumRules) | optional | Complex Field Types |
| repeated | [RepeatedRules](#validate.RepeatedRules) | optional |  |
| map | [MapRules](#validate.MapRules) | optional |  |
| any | [AnyRules](#validate.AnyRules) | optional | Well-Known Field Types |
| duration | [DurationRules](#validate.DurationRules) | optional |  |
| timestamp | [TimestampRules](#validate.TimestampRules) | optional |  |






<a name="validate.Fixed32Rules"></a>

### Fixed32Rules
Fixed32Rules describes the constraints applied to `fixed32` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [fixed32](#fixed32) | optional | Const specifies that this field must be exactly the specified value |
| lt | [fixed32](#fixed32) | optional | Lt specifies that this field must be less than the specified value,
exclusive |
| lte | [fixed32](#fixed32) | optional | Lte specifies that this field must be less than or equal to the
specified value, inclusive |
| gt | [fixed32](#fixed32) | optional | Gt specifies that this field must be greater than the specified value,
exclusive. If the value of Gt is larger than a specified Lt or Lte, the
range is reversed. |
| gte | [fixed32](#fixed32) | optional | Gte specifies that this field must be greater than or equal to the
specified value, inclusive. If the value of Gte is larger than a
specified Lt or Lte, the range is reversed. |
| in | [][fixed32](#fixed32) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][fixed32](#fixed32) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |






<a name="validate.Fixed64Rules"></a>

### Fixed64Rules
Fixed64Rules describes the constraints applied to `fixed64` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [fixed64](#fixed64) | optional | Const specifies that this field must be exactly the specified value |
| lt | [fixed64](#fixed64) | optional | Lt specifies that this field must be less than the specified value,
exclusive |
| lte | [fixed64](#fixed64) | optional | Lte specifies that this field must be less than or equal to the
specified value, inclusive |
| gt | [fixed64](#fixed64) | optional | Gt specifies that this field must be greater than the specified value,
exclusive. If the value of Gt is larger than a specified Lt or Lte, the
range is reversed. |
| gte | [fixed64](#fixed64) | optional | Gte specifies that this field must be greater than or equal to the
specified value, inclusive. If the value of Gte is larger than a
specified Lt or Lte, the range is reversed. |
| in | [][fixed64](#fixed64) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][fixed64](#fixed64) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |






<a name="validate.FloatRules"></a>

### FloatRules
FloatRules describes the constraints applied to `float` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [float](#float) | optional | Const specifies that this field must be exactly the specified value |
| lt | [float](#float) | optional | Lt specifies that this field must be less than the specified value,
exclusive |
| lte | [float](#float) | optional | Lte specifies that this field must be less than or equal to the
specified value, inclusive |
| gt | [float](#float) | optional | Gt specifies that this field must be greater than the specified value,
exclusive. If the value of Gt is larger than a specified Lt or Lte, the
range is reversed. |
| gte | [float](#float) | optional | Gte specifies that this field must be greater than or equal to the
specified value, inclusive. If the value of Gte is larger than a
specified Lt or Lte, the range is reversed. |
| in | [][float](#float) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][float](#float) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |






<a name="validate.Int32Rules"></a>

### Int32Rules
Int32Rules describes the constraints applied to `int32` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [int32](#int32) | optional | Const specifies that this field must be exactly the specified value |
| lt | [int32](#int32) | optional | Lt specifies that this field must be less than the specified value,
exclusive |
| lte | [int32](#int32) | optional | Lte specifies that this field must be less than or equal to the
specified value, inclusive |
| gt | [int32](#int32) | optional | Gt specifies that this field must be greater than the specified value,
exclusive. If the value of Gt is larger than a specified Lt or Lte, the
range is reversed. |
| gte | [int32](#int32) | optional | Gte specifies that this field must be greater than or equal to the
specified value, inclusive. If the value of Gte is larger than a
specified Lt or Lte, the range is reversed. |
| in | [][int32](#int32) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][int32](#int32) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |






<a name="validate.Int64Rules"></a>

### Int64Rules
Int64Rules describes the constraints applied to `int64` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [int64](#int64) | optional | Const specifies that this field must be exactly the specified value |
| lt | [int64](#int64) | optional | Lt specifies that this field must be less than the specified value,
exclusive |
| lte | [int64](#int64) | optional | Lte specifies that this field must be less than or equal to the
specified value, inclusive |
| gt | [int64](#int64) | optional | Gt specifies that this field must be greater than the specified value,
exclusive. If the value of Gt is larger than a specified Lt or Lte, the
range is reversed. |
| gte | [int64](#int64) | optional | Gte specifies that this field must be greater than or equal to the
specified value, inclusive. If the value of Gte is larger than a
specified Lt or Lte, the range is reversed. |
| in | [][int64](#int64) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][int64](#int64) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |






<a name="validate.MapRules"></a>

### MapRules
MapRules describe the constraints applied to `map` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| min_pairs | [uint64](#uint64) | optional | MinPairs specifies that this field must have the specified number of
KVs at a minimum |
| max_pairs | [uint64](#uint64) | optional | MaxPairs specifies that this field must have the specified number of
KVs at a maximum |
| no_sparse | [bool](#bool) | optional | NoSparse specifies values in this field cannot be unset. This only
applies to map&#39;s with message value types. |
| keys | [FieldRules](#validate.FieldRules) | optional | Keys specifies the constraints to be applied to each key in the field. |
| values | [FieldRules](#validate.FieldRules) | optional | Values specifies the constraints to be applied to the value of each key
in the field. Message values will still have their validations evaluated
unless skip is specified here. |






<a name="validate.MessageRules"></a>

### MessageRules
MessageRules describe the constraints applied to embedded message values.
For message-type fields, validation is performed recursively.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| skip | [bool](#bool) | optional | Skip specifies that the validation rules of this field should not be
evaluated |
| required | [bool](#bool) | optional | Required specifies that this field must be set |






<a name="validate.RepeatedRules"></a>

### RepeatedRules
RepeatedRules describe the constraints applied to `repeated` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| min_items | [uint64](#uint64) | optional | MinItems specifies that this field must have the specified number of
items at a minimum |
| max_items | [uint64](#uint64) | optional | MaxItems specifies that this field must have the specified number of
items at a maximum |
| unique | [bool](#bool) | optional | Unique specifies that all elements in this field must be unique. This
contraint is only applicable to scalar and enum types (messages are not
supported). |
| items | [FieldRules](#validate.FieldRules) | optional | Items specifies the contraints to be applied to each item in the field.
Repeated message fields will still execute validation against each item
unless skip is specified here. |






<a name="validate.SFixed32Rules"></a>

### SFixed32Rules
SFixed32Rules describes the constraints applied to `sfixed32` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [sfixed32](#sfixed32) | optional | Const specifies that this field must be exactly the specified value |
| lt | [sfixed32](#sfixed32) | optional | Lt specifies that this field must be less than the specified value,
exclusive |
| lte | [sfixed32](#sfixed32) | optional | Lte specifies that this field must be less than or equal to the
specified value, inclusive |
| gt | [sfixed32](#sfixed32) | optional | Gt specifies that this field must be greater than the specified value,
exclusive. If the value of Gt is larger than a specified Lt or Lte, the
range is reversed. |
| gte | [sfixed32](#sfixed32) | optional | Gte specifies that this field must be greater than or equal to the
specified value, inclusive. If the value of Gte is larger than a
specified Lt or Lte, the range is reversed. |
| in | [][sfixed32](#sfixed32) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][sfixed32](#sfixed32) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |






<a name="validate.SFixed64Rules"></a>

### SFixed64Rules
SFixed64Rules describes the constraints applied to `sfixed64` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [sfixed64](#sfixed64) | optional | Const specifies that this field must be exactly the specified value |
| lt | [sfixed64](#sfixed64) | optional | Lt specifies that this field must be less than the specified value,
exclusive |
| lte | [sfixed64](#sfixed64) | optional | Lte specifies that this field must be less than or equal to the
specified value, inclusive |
| gt | [sfixed64](#sfixed64) | optional | Gt specifies that this field must be greater than the specified value,
exclusive. If the value of Gt is larger than a specified Lt or Lte, the
range is reversed. |
| gte | [sfixed64](#sfixed64) | optional | Gte specifies that this field must be greater than or equal to the
specified value, inclusive. If the value of Gte is larger than a
specified Lt or Lte, the range is reversed. |
| in | [][sfixed64](#sfixed64) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][sfixed64](#sfixed64) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |






<a name="validate.SInt32Rules"></a>

### SInt32Rules
SInt32Rules describes the constraints applied to `sint32` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [sint32](#sint32) | optional | Const specifies that this field must be exactly the specified value |
| lt | [sint32](#sint32) | optional | Lt specifies that this field must be less than the specified value,
exclusive |
| lte | [sint32](#sint32) | optional | Lte specifies that this field must be less than or equal to the
specified value, inclusive |
| gt | [sint32](#sint32) | optional | Gt specifies that this field must be greater than the specified value,
exclusive. If the value of Gt is larger than a specified Lt or Lte, the
range is reversed. |
| gte | [sint32](#sint32) | optional | Gte specifies that this field must be greater than or equal to the
specified value, inclusive. If the value of Gte is larger than a
specified Lt or Lte, the range is reversed. |
| in | [][sint32](#sint32) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][sint32](#sint32) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |






<a name="validate.SInt64Rules"></a>

### SInt64Rules
SInt64Rules describes the constraints applied to `sint64` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [sint64](#sint64) | optional | Const specifies that this field must be exactly the specified value |
| lt | [sint64](#sint64) | optional | Lt specifies that this field must be less than the specified value,
exclusive |
| lte | [sint64](#sint64) | optional | Lte specifies that this field must be less than or equal to the
specified value, inclusive |
| gt | [sint64](#sint64) | optional | Gt specifies that this field must be greater than the specified value,
exclusive. If the value of Gt is larger than a specified Lt or Lte, the
range is reversed. |
| gte | [sint64](#sint64) | optional | Gte specifies that this field must be greater than or equal to the
specified value, inclusive. If the value of Gte is larger than a
specified Lt or Lte, the range is reversed. |
| in | [][sint64](#sint64) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][sint64](#sint64) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |






<a name="validate.StringRules"></a>

### StringRules
StringRules describe the constraints applied to `string` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [string](#string) | optional | Const specifies that this field must be exactly the specified value |
| len | [uint64](#uint64) | optional | Len specifies that this field must be the specified number of
characters (Unicode code points). Note that the number of
characters may differ from the number of bytes in the string. |
| min_len | [uint64](#uint64) | optional | MinLen specifies that this field must be the specified number of
characters (Unicode code points) at a minimum. Note that the number of
characters may differ from the number of bytes in the string. |
| max_len | [uint64](#uint64) | optional | MaxLen specifies that this field must be the specified number of
characters (Unicode code points) at a maximum. Note that the number of
characters may differ from the number of bytes in the string. |
| len_bytes | [uint64](#uint64) | optional | LenBytes specifies that this field must be the specified number of bytes
at a minimum |
| min_bytes | [uint64](#uint64) | optional | MinBytes specifies that this field must be the specified number of bytes
at a minimum |
| max_bytes | [uint64](#uint64) | optional | MaxBytes specifies that this field must be the specified number of bytes
at a maximum |
| pattern | [string](#string) | optional | Pattern specifes that this field must match against the specified
regular expression (RE2 syntax). The included expression should elide
any delimiters. |
| prefix | [string](#string) | optional | Prefix specifies that this field must have the specified substring at
the beginning of the string. |
| suffix | [string](#string) | optional | Suffix specifies that this field must have the specified substring at
the end of the string. |
| contains | [string](#string) | optional | Contains specifies that this field must have the specified substring
anywhere in the string. |
| not_contains | [string](#string) | optional | NotContains specifies that this field cannot have the specified substring
anywhere in the string. |
| in | [][string](#string) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][string](#string) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |
| email | [bool](#bool) | optional | Email specifies that the field must be a valid email address as
defined by RFC 5322 |
| hostname | [bool](#bool) | optional | Hostname specifies that the field must be a valid hostname as
defined by RFC 1034. This constraint does not support
internationalized domain names (IDNs). |
| ip | [bool](#bool) | optional | Ip specifies that the field must be a valid IP (v4 or v6) address.
Valid IPv6 addresses should not include surrounding square brackets. |
| ipv4 | [bool](#bool) | optional | Ipv4 specifies that the field must be a valid IPv4 address. |
| ipv6 | [bool](#bool) | optional | Ipv6 specifies that the field must be a valid IPv6 address. Valid
IPv6 addresses should not include surrounding square brackets. |
| uri | [bool](#bool) | optional | Uri specifies that the field must be a valid, absolute URI as defined
by RFC 3986 |
| uri_ref | [bool](#bool) | optional | UriRef specifies that the field must be a valid URI as defined by RFC
3986 and may be relative or absolute. |
| address | [bool](#bool) | optional | Address specifies that the field must be either a valid hostname as
defined by RFC 1034 (which does not support internationalized domain
names or IDNs), or it can be a valid IP (v4 or v6). |
| uuid | [bool](#bool) | optional | Uuid specifies that the field must be a valid UUID as defined by
RFC 4122 |
| well_known_regex | [KnownRegex](#validate.KnownRegex) | optional | WellKnownRegex specifies a common well known pattern defined as a regex. |
| strict | [bool](#bool) | optional | This applies to regexes HTTP_HEADER_NAME and HTTP_HEADER_VALUE to enable
strict header validation.
By default, this is true, and HTTP header validations are RFC-compliant.
Setting to false will enable a looser validations that only disallows
\r\n\0 characters, which can be used to bypass header matching rules. Default: true |






<a name="validate.TimestampRules"></a>

### TimestampRules
TimestampRules describe the constraints applied exclusively to the
`google.protobuf.Timestamp` well-known type


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| required | [bool](#bool) | optional | Required specifies that this field must be set |
| const | [google.protobuf.Timestamp](#google.protobuf.Timestamp) | optional | Const specifies that this field must be exactly the specified value |
| lt | [google.protobuf.Timestamp](#google.protobuf.Timestamp) | optional | Lt specifies that this field must be less than the specified value,
exclusive |
| lte | [google.protobuf.Timestamp](#google.protobuf.Timestamp) | optional | Lte specifies that this field must be less than the specified value,
inclusive |
| gt | [google.protobuf.Timestamp](#google.protobuf.Timestamp) | optional | Gt specifies that this field must be greater than the specified value,
exclusive |
| gte | [google.protobuf.Timestamp](#google.protobuf.Timestamp) | optional | Gte specifies that this field must be greater than the specified value,
inclusive |
| lt_now | [bool](#bool) | optional | LtNow specifies that this must be less than the current time. LtNow
can only be used with the Within rule. |
| gt_now | [bool](#bool) | optional | GtNow specifies that this must be greater than the current time. GtNow
can only be used with the Within rule. |
| within | [google.protobuf.Duration](#google.protobuf.Duration) | optional | Within specifies that this field must be within this duration of the
current time. This constraint can be used alone or with the LtNow and
GtNow rules. |






<a name="validate.UInt32Rules"></a>

### UInt32Rules
UInt32Rules describes the constraints applied to `uint32` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [uint32](#uint32) | optional | Const specifies that this field must be exactly the specified value |
| lt | [uint32](#uint32) | optional | Lt specifies that this field must be less than the specified value,
exclusive |
| lte | [uint32](#uint32) | optional | Lte specifies that this field must be less than or equal to the
specified value, inclusive |
| gt | [uint32](#uint32) | optional | Gt specifies that this field must be greater than the specified value,
exclusive. If the value of Gt is larger than a specified Lt or Lte, the
range is reversed. |
| gte | [uint32](#uint32) | optional | Gte specifies that this field must be greater than or equal to the
specified value, inclusive. If the value of Gte is larger than a
specified Lt or Lte, the range is reversed. |
| in | [][uint32](#uint32) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][uint32](#uint32) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |






<a name="validate.UInt64Rules"></a>

### UInt64Rules
UInt64Rules describes the constraints applied to `uint64` values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| const | [uint64](#uint64) | optional | Const specifies that this field must be exactly the specified value |
| lt | [uint64](#uint64) | optional | Lt specifies that this field must be less than the specified value,
exclusive |
| lte | [uint64](#uint64) | optional | Lte specifies that this field must be less than or equal to the
specified value, inclusive |
| gt | [uint64](#uint64) | optional | Gt specifies that this field must be greater than the specified value,
exclusive. If the value of Gt is larger than a specified Lt or Lte, the
range is reversed. |
| gte | [uint64](#uint64) | optional | Gte specifies that this field must be greater than or equal to the
specified value, inclusive. If the value of Gte is larger than a
specified Lt or Lte, the range is reversed. |
| in | [][uint64](#uint64) | repeated | In specifies that this field must be equal to one of the specified
values |
| not_in | [][uint64](#uint64) | repeated | NotIn specifies that this field cannot be equal to one of the specified
values |





 


<a name="validate.KnownRegex"></a>

### KnownRegex
WellKnownRegex contain some well-known patterns.

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| HTTP_HEADER_NAME | 1 | HTTP header name as defined by RFC 7230. |
| HTTP_HEADER_VALUE | 2 | HTTP header value as defined by RFC 7230. |


 


<a name="github.com/envoyproxy/protoc-gen-validate/validate/validate.proto-extensions"></a>

### File-level Extensions
| Extension | Type | Base | Number | Description |
| --------- | ---- | ---- | ------ | ----------- |
| rules | FieldRules | .google.protobuf.FieldOptions | 1071 | Rules specify the validations to be performed on this field. By default,
no validation is performed against a field. |
| disabled | bool | .google.protobuf.MessageOptions | 1071 | Disabled nullifies any validation rules for this message, including any
message fields associated with it that do support validation. |
| required | bool | .google.protobuf.OneofOptions | 1071 | Required ensures that exactly one the field options in a oneof is set;
validation fails if no fields in the oneof are set. |

 

 

