
---
title: "google.protobufgithub.com/solo-io/solo-kit/api/external/google/protobuf/descriptor.proto"
---

## Package : `google.protobuf`



<a name="top"></a>

<a name="API Reference for github.com/solo-io/solo-kit/api/external/google/protobuf/descriptor.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/solo-io/solo-kit/api/external/google/protobuf/descriptor.proto


## Table of Contents
  - [DescriptorProto](#google.protobuf.DescriptorProto)
  - [DescriptorProto.ExtensionRange](#google.protobuf.DescriptorProto.ExtensionRange)
  - [DescriptorProto.ReservedRange](#google.protobuf.DescriptorProto.ReservedRange)
  - [EnumDescriptorProto](#google.protobuf.EnumDescriptorProto)
  - [EnumOptions](#google.protobuf.EnumOptions)
  - [EnumValueDescriptorProto](#google.protobuf.EnumValueDescriptorProto)
  - [EnumValueOptions](#google.protobuf.EnumValueOptions)
  - [FieldDescriptorProto](#google.protobuf.FieldDescriptorProto)
  - [FieldOptions](#google.protobuf.FieldOptions)
  - [FileDescriptorProto](#google.protobuf.FileDescriptorProto)
  - [FileDescriptorSet](#google.protobuf.FileDescriptorSet)
  - [FileOptions](#google.protobuf.FileOptions)
  - [GeneratedCodeInfo](#google.protobuf.GeneratedCodeInfo)
  - [GeneratedCodeInfo.Annotation](#google.protobuf.GeneratedCodeInfo.Annotation)
  - [MessageOptions](#google.protobuf.MessageOptions)
  - [MethodDescriptorProto](#google.protobuf.MethodDescriptorProto)
  - [MethodOptions](#google.protobuf.MethodOptions)
  - [OneofDescriptorProto](#google.protobuf.OneofDescriptorProto)
  - [OneofOptions](#google.protobuf.OneofOptions)
  - [ServiceDescriptorProto](#google.protobuf.ServiceDescriptorProto)
  - [ServiceOptions](#google.protobuf.ServiceOptions)
  - [SourceCodeInfo](#google.protobuf.SourceCodeInfo)
  - [SourceCodeInfo.Location](#google.protobuf.SourceCodeInfo.Location)
  - [UninterpretedOption](#google.protobuf.UninterpretedOption)
  - [UninterpretedOption.NamePart](#google.protobuf.UninterpretedOption.NamePart)

  - [FieldDescriptorProto.Label](#google.protobuf.FieldDescriptorProto.Label)
  - [FieldDescriptorProto.Type](#google.protobuf.FieldDescriptorProto.Type)
  - [FieldOptions.CType](#google.protobuf.FieldOptions.CType)
  - [FieldOptions.JSType](#google.protobuf.FieldOptions.JSType)
  - [FileOptions.OptimizeMode](#google.protobuf.FileOptions.OptimizeMode)
  - [MethodOptions.IdempotencyLevel](#google.protobuf.MethodOptions.IdempotencyLevel)






<a name="google.protobuf.DescriptorProto"></a>

### DescriptorProto
Describes a message type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) | optional |  |
| field | [][FieldDescriptorProto](#google.protobuf.FieldDescriptorProto) | repeated |  |
| extension | [][FieldDescriptorProto](#google.protobuf.FieldDescriptorProto) | repeated |  |
| nested_type | [][DescriptorProto](#google.protobuf.DescriptorProto) | repeated |  |
| enum_type | [][EnumDescriptorProto](#google.protobuf.EnumDescriptorProto) | repeated |  |
| extension_range | [][DescriptorProto.ExtensionRange](#google.protobuf.DescriptorProto.ExtensionRange) | repeated |  |
| oneof_decl | [][OneofDescriptorProto](#google.protobuf.OneofDescriptorProto) | repeated |  |
| options | [MessageOptions](#google.protobuf.MessageOptions) | optional |  |
| reserved_range | [][DescriptorProto.ReservedRange](#google.protobuf.DescriptorProto.ReservedRange) | repeated |  |
| reserved_name | [][string](#string) | repeated | Reserved field names, which may not be used by fields in the same message.
A given name may only be reserved once. |






<a name="google.protobuf.DescriptorProto.ExtensionRange"></a>

### DescriptorProto.ExtensionRange



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| start | [int32](#int32) | optional |  |
| end | [int32](#int32) | optional |  |






<a name="google.protobuf.DescriptorProto.ReservedRange"></a>

### DescriptorProto.ReservedRange
Range of reserved tag numbers. Reserved tag numbers may not be used by
fields or extension ranges in the same message. Reserved ranges may
not overlap.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| start | [int32](#int32) | optional | Inclusive. |
| end | [int32](#int32) | optional | Exclusive. |






<a name="google.protobuf.EnumDescriptorProto"></a>

### EnumDescriptorProto
Describes an enum type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) | optional |  |
| value | [][EnumValueDescriptorProto](#google.protobuf.EnumValueDescriptorProto) | repeated |  |
| options | [EnumOptions](#google.protobuf.EnumOptions) | optional |  |






<a name="google.protobuf.EnumOptions"></a>

### EnumOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| allow_alias | [bool](#bool) | optional | Set this option to true to allow mapping different tag names to the same
value. |
| deprecated | [bool](#bool) | optional | Is this enum deprecated?
Depending on the target platform, this can emit Deprecated annotations
for the enum, or it will be completely ignored; in the very least, this
is a formalization for deprecating enums. Default: false |
| uninterpreted_option | [][UninterpretedOption](#google.protobuf.UninterpretedOption) | repeated | The parser stores options it doesn&#39;t recognize here. See above. |






<a name="google.protobuf.EnumValueDescriptorProto"></a>

### EnumValueDescriptorProto
Describes a value within an enum.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) | optional |  |
| number | [int32](#int32) | optional |  |
| options | [EnumValueOptions](#google.protobuf.EnumValueOptions) | optional |  |






<a name="google.protobuf.EnumValueOptions"></a>

### EnumValueOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| deprecated | [bool](#bool) | optional | Is this enum value deprecated?
Depending on the target platform, this can emit Deprecated annotations
for the enum value, or it will be completely ignored; in the very least,
this is a formalization for deprecating enum values. Default: false |
| uninterpreted_option | [][UninterpretedOption](#google.protobuf.UninterpretedOption) | repeated | The parser stores options it doesn&#39;t recognize here. See above. |






<a name="google.protobuf.FieldDescriptorProto"></a>

### FieldDescriptorProto
Describes a field within a message.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) | optional |  |
| number | [int32](#int32) | optional |  |
| label | [FieldDescriptorProto.Label](#google.protobuf.FieldDescriptorProto.Label) | optional |  |
| type | [FieldDescriptorProto.Type](#google.protobuf.FieldDescriptorProto.Type) | optional | If type_name is set, this need not be set.  If both this and type_name
are set, this must be one of TYPE_ENUM, TYPE_MESSAGE or TYPE_GROUP. |
| type_name | [string](#string) | optional | For message and enum types, this is the name of the type.  If the name
starts with a &#39;.&#39;, it is fully-qualified.  Otherwise, C&#43;&#43;-like scoping
rules are used to find the type (i.e. first the nested types within this
message are searched, then within the parent, on up to the root
namespace). |
| extendee | [string](#string) | optional | For extensions, this is the name of the type being extended.  It is
resolved in the same manner as type_name. |
| default_value | [string](#string) | optional | For numeric types, contains the original text representation of the value.
For booleans, &#34;true&#34; or &#34;false&#34;.
For strings, contains the default text contents (not escaped in any way).
For bytes, contains the C escaped value.  All bytes &gt;= 128 are escaped.
TODO(kenton):  Base-64 encode? |
| oneof_index | [int32](#int32) | optional | If set, gives the index of a oneof in the containing type&#39;s oneof_decl
list.  This field is a member of that oneof. |
| json_name | [string](#string) | optional | JSON name of this field. The value is set by protocol compiler. If the
user has set a &#34;json_name&#34; option on this field, that option&#39;s value
will be used. Otherwise, it&#39;s deduced from the field&#39;s name by converting
it to camelCase. |
| options | [FieldOptions](#google.protobuf.FieldOptions) | optional |  |






<a name="google.protobuf.FieldOptions"></a>

### FieldOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ctype | [FieldOptions.CType](#google.protobuf.FieldOptions.CType) | optional | The ctype option instructs the C&#43;&#43; code generator to use a different
representation of the field than it normally would.  See the specific
options below.  This option is not yet implemented in the open source
release -- sorry, we&#39;ll try to include it in a future version! Default: STRING |
| packed | [bool](#bool) | optional | The packed option can be enabled for repeated primitive fields to enable
a more efficient representation on the wire. Rather than repeatedly
writing the tag and type for each element, the entire array is encoded as
a single length-delimited blob. In proto3, only explicit setting it to
false will avoid using packed encoding. |
| jstype | [FieldOptions.JSType](#google.protobuf.FieldOptions.JSType) | optional | The jstype option determines the JavaScript type used for values of the
field.  The option is permitted only for 64 bit integral and fixed types
(int64, uint64, sint64, fixed64, sfixed64).  By default these types are
represented as JavaScript strings.  This avoids loss of precision that can
happen when a large value is converted to a floating point JavaScript
numbers.  Specifying JS_NUMBER for the jstype causes the generated
JavaScript code to use the JavaScript &#34;number&#34; type instead of strings.
This option is an enum to permit additional types to be added,
e.g. goog.math.Integer. Default: JS_NORMAL |
| lazy | [bool](#bool) | optional | Should this field be parsed lazily?  Lazy applies only to message-type
fields.  It means that when the outer message is initially parsed, the
inner message&#39;s contents will not be parsed but instead stored in encoded
form.  The inner message will actually be parsed when it is first accessed.

This is only a hint.  Implementations are free to choose whether to use
eager or lazy parsing regardless of the value of this option.  However,
setting this option true suggests that the protocol author believes that
using lazy parsing on this field is worth the additional bookkeeping
overhead typically needed to implement it.

This option does not affect the public interface of any generated code;
all method signatures remain the same.  Furthermore, thread-safety of the
interface is not affected by this option; const methods remain safe to
call from multiple threads concurrently, while non-const methods continue
to require exclusive access.


Note that implementations may choose not to check required fields within
a lazy sub-message.  That is, calling IsInitialized() on the outer message
may return true even if the inner message has missing required fields.
This is necessary because otherwise the inner message would have to be
parsed in order to perform the check, defeating the purpose of lazy
parsing.  An implementation which chooses not to check required fields
must be consistent about it.  That is, for any particular sub-message, the
implementation must either *always* check its required fields, or *never*
check its required fields, regardless of whether or not the message has
been parsed. Default: false |
| deprecated | [bool](#bool) | optional | Is this field deprecated?
Depending on the target platform, this can emit Deprecated annotations
for accessors, or it will be completely ignored; in the very least, this
is a formalization for deprecating fields. Default: false |
| weak | [bool](#bool) | optional | For Google-internal migration only. Do not use. Default: false |
| uninterpreted_option | [][UninterpretedOption](#google.protobuf.UninterpretedOption) | repeated | The parser stores options it doesn&#39;t recognize here. See above. |






<a name="google.protobuf.FileDescriptorProto"></a>

### FileDescriptorProto
Describes a complete .proto file.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) | optional | file name, relative to root of source tree |
| package | [string](#string) | optional | e.g. &#34;foo&#34;, &#34;foo.bar&#34;, etc. |
| dependency | [][string](#string) | repeated | Names of files imported by this file. |
| public_dependency | [][int32](#int32) | repeated | Indexes of the public imported files in the dependency list above. |
| weak_dependency | [][int32](#int32) | repeated | Indexes of the weak imported files in the dependency list.
For Google-internal migration only. Do not use. |
| message_type | [][DescriptorProto](#google.protobuf.DescriptorProto) | repeated | All top-level definitions in this file. |
| enum_type | [][EnumDescriptorProto](#google.protobuf.EnumDescriptorProto) | repeated |  |
| service | [][ServiceDescriptorProto](#google.protobuf.ServiceDescriptorProto) | repeated |  |
| extension | [][FieldDescriptorProto](#google.protobuf.FieldDescriptorProto) | repeated |  |
| options | [FileOptions](#google.protobuf.FileOptions) | optional |  |
| source_code_info | [SourceCodeInfo](#google.protobuf.SourceCodeInfo) | optional | This field contains optional information about the original source code.
You may safely remove this entire field without harming runtime
functionality of the descriptors -- the information is needed only by
development tools. |
| syntax | [string](#string) | optional | The syntax of the proto file.
The supported values are &#34;proto2&#34; and &#34;proto3&#34;. |






<a name="google.protobuf.FileDescriptorSet"></a>

### FileDescriptorSet
The protocol compiler can output a FileDescriptorSet containing the .proto
files it parses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| file | [][FileDescriptorProto](#google.protobuf.FileDescriptorProto) | repeated |  |






<a name="google.protobuf.FileOptions"></a>

### FileOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| java_package | [string](#string) | optional | Sets the Java package where classes generated from this .proto will be
placed.  By default, the proto package is used, but this is often
inappropriate because proto packages do not normally start with backwards
domain names. |
| java_outer_classname | [string](#string) | optional | If set, all the classes from the .proto file are wrapped in a single
outer class with the given name.  This applies to both Proto1
(equivalent to the old &#34;--one_java_file&#34; option) and Proto2 (where
a .proto always translates to a single class, but you may want to
explicitly choose the class name). |
| java_multiple_files | [bool](#bool) | optional | If set true, then the Java code generator will generate a separate .java
file for each top-level message, enum, and service defined in the .proto
file.  Thus, these types will *not* be nested inside the outer class
named by java_outer_classname.  However, the outer class will still be
generated to contain the file&#39;s getDescriptor() method as well as any
top-level extensions defined in the file. Default: false |
| java_generate_equals_and_hash | [bool](#bool) | optional | This option does nothing. |
| java_string_check_utf8 | [bool](#bool) | optional | If set true, then the Java2 code generator will generate code that
throws an exception whenever an attempt is made to assign a non-UTF-8
byte sequence to a string field.
Message reflection will do the same.
However, an extension field still accepts non-UTF-8 byte sequences.
This option has no effect on when used with the lite runtime. Default: false |
| optimize_for | [FileOptions.OptimizeMode](#google.protobuf.FileOptions.OptimizeMode) | optional |  Default: SPEED |
| go_package | [string](#string) | optional | Sets the Go package where structs generated from this .proto will be
placed. If omitted, the Go package will be derived from the following:
  - The basename of the package import path, if provided.
  - Otherwise, the package statement in the .proto file, if present.
  - Otherwise, the basename of the .proto file, without extension. |
| cc_generic_services | [bool](#bool) | optional | Should generic services be generated in each language?  &#34;Generic&#34; services
are not specific to any particular RPC system.  They are generated by the
main code generators in each language (without additional plugins).
Generic services were the only kind of service generation supported by
early versions of google.protobuf.

Generic services are now considered deprecated in favor of using plugins
that generate code specific to your particular RPC system.  Therefore,
these default to false.  Old code which depends on generic services should
explicitly set them to true. Default: false |
| java_generic_services | [bool](#bool) | optional |  Default: false |
| py_generic_services | [bool](#bool) | optional |  Default: false |
| deprecated | [bool](#bool) | optional | Is this file deprecated?
Depending on the target platform, this can emit Deprecated annotations
for everything in the file, or it will be completely ignored; in the very
least, this is a formalization for deprecating files. Default: false |
| cc_enable_arenas | [bool](#bool) | optional | Enables the use of arenas for the proto messages in this file. This applies
only to generated classes for C&#43;&#43;. Default: false |
| objc_class_prefix | [string](#string) | optional | Sets the objective c class prefix which is prepended to all objective c
generated classes from this .proto. There is no default. |
| csharp_namespace | [string](#string) | optional | Namespace for generated classes; defaults to the package. |
| swift_prefix | [string](#string) | optional | By default Swift generators will take the proto package and CamelCase it
replacing &#39;.&#39; with underscore and use that to prefix the types/symbols
defined. When this options is provided, they will use this value instead
to prefix the types/symbols defined. |
| php_class_prefix | [string](#string) | optional | Sets the php class prefix which is prepended to all php generated classes
from this .proto. Default is empty. |
| uninterpreted_option | [][UninterpretedOption](#google.protobuf.UninterpretedOption) | repeated | The parser stores options it doesn&#39;t recognize here. See above. |






<a name="google.protobuf.GeneratedCodeInfo"></a>

### GeneratedCodeInfo
Describes the relationship between generated code and its original source
file. A GeneratedCodeInfo message is associated with only one generated
source file, but may contain references to different source .proto files.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| annotation | [][GeneratedCodeInfo.Annotation](#google.protobuf.GeneratedCodeInfo.Annotation) | repeated | An Annotation connects some span of text in generated code to an element
of its generating .proto file. |






<a name="google.protobuf.GeneratedCodeInfo.Annotation"></a>

### GeneratedCodeInfo.Annotation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | [][int32](#int32) | repeated | Identifies the element in the original source .proto file. This field
is formatted the same as SourceCodeInfo.Location.path. |
| source_file | [string](#string) | optional | Identifies the filesystem path to the original source .proto. |
| begin | [int32](#int32) | optional | Identifies the starting offset in bytes in the generated code
that relates to the identified object. |
| end | [int32](#int32) | optional | Identifies the ending offset in bytes in the generated code that
relates to the identified offset. The end offset should be one past
the last relevant byte (so the length of the text = end - begin). |






<a name="google.protobuf.MessageOptions"></a>

### MessageOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| message_set_wire_format | [bool](#bool) | optional | Set true to use the old proto1 MessageSet wire format for extensions.
This is provided for backwards-compatibility with the MessageSet wire
format.  You should not use this for any other reason:  It&#39;s less
efficient, has fewer features, and is more complicated.

The message must be defined exactly as follows:
  message Foo {
    option message_set_wire_format = true;
    extensions 4 to max;
  }
Note that the message cannot have any defined fields; MessageSets only
have extensions.

All extensions of your type must be singular messages; e.g. they cannot
be int32s, enums, or repeated messages.

Because this is an option, the above two restrictions are not enforced by
the protocol compiler. Default: false |
| no_standard_descriptor_accessor | [bool](#bool) | optional | Disables the generation of the standard &#34;descriptor()&#34; accessor, which can
conflict with a field of the same name.  This is meant to make migration
from proto1 easier; new code should avoid fields named &#34;descriptor&#34;. Default: false |
| deprecated | [bool](#bool) | optional | Is this message deprecated?
Depending on the target platform, this can emit Deprecated annotations
for the message, or it will be completely ignored; in the very least,
this is a formalization for deprecating messages. Default: false |
| map_entry | [bool](#bool) | optional | Whether the message is an automatically generated map entry type for the
maps field.

For maps fields:
    map&lt;KeyType, ValueType&gt; map_field = 1;
The parsed descriptor looks like:
    message MapFieldEntry {
        option map_entry = true;
        optional KeyType key = 1;
        optional ValueType value = 2;
    }
    repeated MapFieldEntry map_field = 1;

Implementations may choose not to generate the map_entry=true message, but
use a native map in the target language to hold the keys and values.
The reflection APIs in such implementions still need to work as
if the field is a repeated message field.

NOTE: Do not set the option in .proto files. Always use the maps syntax
instead. The option should only be implicitly set by the proto compiler
parser. |
| uninterpreted_option | [][UninterpretedOption](#google.protobuf.UninterpretedOption) | repeated | The parser stores options it doesn&#39;t recognize here. See above. |






<a name="google.protobuf.MethodDescriptorProto"></a>

### MethodDescriptorProto
Describes a method of a service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) | optional |  |
| input_type | [string](#string) | optional | Input and output type names.  These are resolved in the same way as
FieldDescriptorProto.type_name, but must refer to a message type. |
| output_type | [string](#string) | optional |  |
| options | [MethodOptions](#google.protobuf.MethodOptions) | optional |  |
| client_streaming | [bool](#bool) | optional | Identifies if client streams multiple client messages Default: false |
| server_streaming | [bool](#bool) | optional | Identifies if server streams multiple server messages Default: false |






<a name="google.protobuf.MethodOptions"></a>

### MethodOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| deprecated | [bool](#bool) | optional | Is this method deprecated?
Depending on the target platform, this can emit Deprecated annotations
for the method, or it will be completely ignored; in the very least,
this is a formalization for deprecating methods. Default: false |
| idempotency_level | [MethodOptions.IdempotencyLevel](#google.protobuf.MethodOptions.IdempotencyLevel) | optional |  Default: IDEMPOTENCY_UNKNOWN |
| uninterpreted_option | [][UninterpretedOption](#google.protobuf.UninterpretedOption) | repeated | The parser stores options it doesn&#39;t recognize here. See above. |






<a name="google.protobuf.OneofDescriptorProto"></a>

### OneofDescriptorProto
Describes a oneof.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) | optional |  |
| options | [OneofOptions](#google.protobuf.OneofOptions) | optional |  |






<a name="google.protobuf.OneofOptions"></a>

### OneofOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uninterpreted_option | [][UninterpretedOption](#google.protobuf.UninterpretedOption) | repeated | The parser stores options it doesn&#39;t recognize here. See above. |






<a name="google.protobuf.ServiceDescriptorProto"></a>

### ServiceDescriptorProto
Describes a service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) | optional |  |
| method | [][MethodDescriptorProto](#google.protobuf.MethodDescriptorProto) | repeated |  |
| options | [ServiceOptions](#google.protobuf.ServiceOptions) | optional |  |






<a name="google.protobuf.ServiceOptions"></a>

### ServiceOptions



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| deprecated | [bool](#bool) | optional | Is this service deprecated?
Depending on the target platform, this can emit Deprecated annotations
for the service, or it will be completely ignored; in the very least,
this is a formalization for deprecating services. Default: false |
| uninterpreted_option | [][UninterpretedOption](#google.protobuf.UninterpretedOption) | repeated | The parser stores options it doesn&#39;t recognize here. See above. |






<a name="google.protobuf.SourceCodeInfo"></a>

### SourceCodeInfo
Encapsulates information about the original source file from which a
FileDescriptorProto was generated.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| location | [][SourceCodeInfo.Location](#google.protobuf.SourceCodeInfo.Location) | repeated | A Location identifies a piece of source code in a .proto file which
corresponds to a particular definition.  This information is intended
to be useful to IDEs, code indexers, documentation generators, and similar
tools.

For example, say we have a file like:
  message Foo {
    optional string foo = 1;
  }
Let&#39;s look at just the field definition:
  optional string foo = 1;
  ^       ^^     ^^  ^  ^^^
  a       bc     de  f  ghi
We have the following locations:
  span   path               represents
  [a,i)  [ 4, 0, 2, 0 ]     The whole field definition.
  [a,b)  [ 4, 0, 2, 0, 4 ]  The label (optional).
  [c,d)  [ 4, 0, 2, 0, 5 ]  The type (string).
  [e,f)  [ 4, 0, 2, 0, 1 ]  The name (foo).
  [g,h)  [ 4, 0, 2, 0, 3 ]  The number (1).

Notes:
- A location may refer to a repeated field itself (i.e. not to any
  particular index within it).  This is used whenever a set of elements are
  logically enclosed in a single code segment.  For example, an entire
  extend block (possibly containing multiple extension definitions) will
  have an outer location whose path refers to the &#34;extensions&#34; repeated
  field without an index.
- Multiple locations may have the same path.  This happens when a single
  logical declaration is spread out across multiple places.  The most
  obvious example is the &#34;extend&#34; block again -- there may be multiple
  extend blocks in the same scope, each of which will have the same path.
- A location&#39;s span is not always a subset of its parent&#39;s span.  For
  example, the &#34;extendee&#34; of an extension declaration appears at the
  beginning of the &#34;extend&#34; block and is shared by all extensions within
  the block.
- Just because a location&#39;s span is a subset of some other location&#39;s span
  does not mean that it is a descendent.  For example, a &#34;group&#34; defines
  both a type and a field in a single declaration.  Thus, the locations
  corresponding to the type and field and their components will overlap.
- Code which tries to interpret locations should probably be designed to
  ignore those that it doesn&#39;t understand, as more types of locations could
  be recorded in the future. |






<a name="google.protobuf.SourceCodeInfo.Location"></a>

### SourceCodeInfo.Location



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | [][int32](#int32) | repeated | Identifies which part of the FileDescriptorProto was defined at this
location.

Each element is a field number or an index.  They form a path from
the root FileDescriptorProto to the place where the definition.  For
example, this path:
  [ 4, 3, 2, 7, 1 ]
refers to:
  file.message_type(3)  // 4, 3
      .field(7)         // 2, 7
      .name()           // 1
This is because FileDescriptorProto.message_type has field number 4:
  repeated DescriptorProto message_type = 4;
and DescriptorProto.field has field number 2:
  repeated FieldDescriptorProto field = 2;
and FieldDescriptorProto.name has field number 1:
  optional string name = 1;

Thus, the above path gives the location of a field name.  If we removed
the last element:
  [ 4, 3, 2, 7 ]
this path refers to the whole field declaration (from the beginning
of the label to the terminating semicolon). |
| span | [][int32](#int32) | repeated | Always has exactly three or four elements: start line, start column,
end line (optional, otherwise assumed same as start line), end column.
These are packed into a single field for efficiency.  Note that line
and column numbers are zero-based -- typically you will want to add
1 to each before displaying to a user. |
| leading_comments | [string](#string) | optional | If this SourceCodeInfo represents a complete declaration, these are any
comments appearing before and after the declaration which appear to be
attached to the declaration.

A series of line comments appearing on consecutive lines, with no other
tokens appearing on those lines, will be treated as a single comment.

leading_detached_comments will keep paragraphs of comments that appear
before (but not connected to) the current element. Each paragraph,
separated by empty lines, will be one comment element in the repeated
field.

Only the comment content is provided; comment markers (e.g. //) are
stripped out.  For block comments, leading whitespace and an asterisk
will be stripped from the beginning of each line other than the first.
Newlines are included in the output.

Examples:

  optional int32 foo = 1;  // Comment attached to foo.
  // Comment attached to bar.
  optional int32 bar = 2;

  optional string baz = 3;
  // Comment attached to baz.
  // Another line attached to baz.

  // Comment attached to qux.
  //
  // Another line attached to qux.
  optional double qux = 4;

  // Detached comment for corge. This is not leading or trailing comments
  // to qux or corge because there are blank lines separating it from
  // both.

  // Detached comment for corge paragraph 2.

  optional string corge = 5;
  /* Block comment attached
   * to corge.  Leading asterisks
   * will be removed. */
  /* Block comment attached to
   * grault. */
  optional int32 grault = 6;

  // ignored detached comments. |
| trailing_comments | [string](#string) | optional |  |
| leading_detached_comments | [][string](#string) | repeated |  |






<a name="google.protobuf.UninterpretedOption"></a>

### UninterpretedOption
A message representing a option the parser does not recognize. This only
appears in options protos created by the compiler::Parser class.
DescriptorPool resolves these when building Descriptor objects. Therefore,
options protos in descriptor objects (e.g. returned by Descriptor::options(),
or produced by Descriptor::CopyTo()) will never have UninterpretedOptions
in them.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [][UninterpretedOption.NamePart](#google.protobuf.UninterpretedOption.NamePart) | repeated |  |
| identifier_value | [string](#string) | optional | The value of the uninterpreted option, in whatever type the tokenizer
identified it as during parsing. Exactly one of these should be set. |
| positive_int_value | [uint64](#uint64) | optional |  |
| negative_int_value | [int64](#int64) | optional |  |
| double_value | [double](#double) | optional |  |
| string_value | [bytes](#bytes) | optional |  |
| aggregate_value | [string](#string) | optional |  |






<a name="google.protobuf.UninterpretedOption.NamePart"></a>

### UninterpretedOption.NamePart
The name of the uninterpreted option.  Each string represents a segment in
a dot-separated name.  is_extension is true iff a segment represents an
extension (denoted with parentheses in options specs in .proto files).
E.g.,{ [&#34;foo&#34;, false], [&#34;bar.baz&#34;, true], [&#34;qux&#34;, false] } represents
&#34;foo.(bar.baz).qux&#34;.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name_part | [string](#string) | required |  |
| is_extension | [bool](#bool) | required |  |





 


<a name="google.protobuf.FieldDescriptorProto.Label"></a>

### FieldDescriptorProto.Label


| Name | Number | Description |
| ---- | ------ | ----------- |
| LABEL_OPTIONAL | 1 | 0 is reserved for errors |
| LABEL_REQUIRED | 2 |  |
| LABEL_REPEATED | 3 |  |



<a name="google.protobuf.FieldDescriptorProto.Type"></a>

### FieldDescriptorProto.Type


| Name | Number | Description |
| ---- | ------ | ----------- |
| TYPE_DOUBLE | 1 | 0 is reserved for errors.
Order is weird for historical reasons. |
| TYPE_FLOAT | 2 |  |
| TYPE_INT64 | 3 | Not ZigZag encoded.  Negative numbers take 10 bytes.  Use TYPE_SINT64 if
negative values are likely. |
| TYPE_UINT64 | 4 |  |
| TYPE_INT32 | 5 | Not ZigZag encoded.  Negative numbers take 10 bytes.  Use TYPE_SINT32 if
negative values are likely. |
| TYPE_FIXED64 | 6 |  |
| TYPE_FIXED32 | 7 |  |
| TYPE_BOOL | 8 |  |
| TYPE_STRING | 9 |  |
| TYPE_GROUP | 10 | Tag-delimited aggregate.
Group type is deprecated and not supported in proto3. However, Proto3
implementations should still be able to parse the group wire format and
treat group fields as unknown fields. |
| TYPE_MESSAGE | 11 | Length-delimited aggregate. |
| TYPE_BYTES | 12 | New in version 2. |
| TYPE_UINT32 | 13 |  |
| TYPE_ENUM | 14 |  |
| TYPE_SFIXED32 | 15 |  |
| TYPE_SFIXED64 | 16 |  |
| TYPE_SINT32 | 17 | Uses ZigZag encoding. |
| TYPE_SINT64 | 18 | Uses ZigZag encoding. |



<a name="google.protobuf.FieldOptions.CType"></a>

### FieldOptions.CType


| Name | Number | Description |
| ---- | ------ | ----------- |
| STRING | 0 | Default mode. |
| CORD | 1 |  |
| STRING_PIECE | 2 |  |



<a name="google.protobuf.FieldOptions.JSType"></a>

### FieldOptions.JSType


| Name | Number | Description |
| ---- | ------ | ----------- |
| JS_NORMAL | 0 | Use the default type. |
| JS_STRING | 1 | Use JavaScript strings. |
| JS_NUMBER | 2 | Use JavaScript numbers. |



<a name="google.protobuf.FileOptions.OptimizeMode"></a>

### FileOptions.OptimizeMode
Generated classes can be optimized for speed or code size.

| Name | Number | Description |
| ---- | ------ | ----------- |
| SPEED | 1 | Generate complete code for parsing, serialization, |
| CODE_SIZE | 2 | etc.

Use ReflectionOps to implement these methods. |
| LITE_RUNTIME | 3 | Generate code using MessageLite and the lite runtime. |



<a name="google.protobuf.MethodOptions.IdempotencyLevel"></a>

### MethodOptions.IdempotencyLevel
Is this method side-effect-free (or safe in HTTP parlance), or idempotent,
or neither? HTTP based RPC implementation may choose GET verb for safe
methods, and PUT verb for idempotent methods instead of the default POST.

| Name | Number | Description |
| ---- | ------ | ----------- |
| IDEMPOTENCY_UNKNOWN | 0 |  |
| NO_SIDE_EFFECTS | 1 | implies idempotent |
| IDEMPOTENT | 2 | idempotent, but may have side effects |


 

 

 

