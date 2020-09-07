
---
title: "google.protobufgithub.com/solo-io/solo-kit/api/external/google/protobuf/struct.proto"
---

## Package : `google.protobuf`



<a name="top"></a>

<a name="API Reference for github.com/solo-io/solo-kit/api/external/google/protobuf/struct.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/solo-io/solo-kit/api/external/google/protobuf/struct.proto


## Table of Contents
  - [ListValue](#google.protobuf.ListValue)
  - [Struct](#google.protobuf.Struct)
  - [Struct.FieldsEntry](#google.protobuf.Struct.FieldsEntry)
  - [Value](#google.protobuf.Value)

  - [NullValue](#google.protobuf.NullValue)






<a name="google.protobuf.ListValue"></a>

### ListValue
`ListValue` is a wrapper around a repeated field of values.

The JSON representation for `ListValue` is JSON array.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| values | [][Value](#google.protobuf.Value) | repeated | Repeated field of dynamically typed values. |






<a name="google.protobuf.Struct"></a>

### Struct
`Struct` represents a structured data value, consisting of fields
which map to dynamically typed values. In some languages, `Struct`
might be supported by a native representation. For example, in
scripting languages like JS a struct is represented as an
object. The details of that representation are described together
with the proto support for the language.

The JSON representation for `Struct` is JSON object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| fields | [][Struct.FieldsEntry](#google.protobuf.Struct.FieldsEntry) | repeated | Unordered map of dynamically typed values. |






<a name="google.protobuf.Struct.FieldsEntry"></a>

### Struct.FieldsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Value](#google.protobuf.Value) |  |  |






<a name="google.protobuf.Value"></a>

### Value
`Value` represents a dynamically typed value which can be either
null, a number, a string, a boolean, a recursive struct value, or a
list of values. A producer of value is expected to set one of that
variants, absence of any variant indicates an error.

The JSON representation for `Value` is JSON value.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| null_value | [NullValue](#google.protobuf.NullValue) |  | Represents a null value. |
| number_value | [double](#double) |  | Represents a double value. |
| string_value | [string](#string) |  | Represents a string value. |
| bool_value | [bool](#bool) |  | Represents a boolean value. |
| struct_value | [Struct](#google.protobuf.Struct) |  | Represents a structured value. |
| list_value | [ListValue](#google.protobuf.ListValue) |  | Represents a repeated `Value`. |





 


<a name="google.protobuf.NullValue"></a>

### NullValue
`NullValue` is a singleton enumeration to represent the null value for the
`Value` type union.

 The JSON representation for `NullValue` is JSON `null`.

| Name | Number | Description |
| ---- | ------ | ----------- |
| NULL_VALUE | 0 | Null value. |


 

 

 

