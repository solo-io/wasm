
---
title: "google.protobufgithub.com/solo-io/solo-kit/api/external/google/protobuf/any.proto"
---

## Package : `google.protobuf`



<a name="top"></a>

<a name="API Reference for github.com/solo-io/solo-kit/api/external/google/protobuf/any.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/solo-io/solo-kit/api/external/google/protobuf/any.proto


## Table of Contents
  - [Any](#google.protobuf.Any)







<a name="google.protobuf.Any"></a>

### Any
`Any` contains an arbitrary serialized protocol buffer message along with a
URL that describes the type of the serialized message.

Protobuf library provides support to pack/unpack Any values in the form
of utility functions or additional generated methods of the Any type.

Example 1: Pack and unpack a message in C&#43;&#43;.

    Foo foo = ...;
    Any any;
    any.PackFrom(foo);
    ...
    if (any.UnpackTo(&amp;foo)) {
      ...
    }

Example 2: Pack and unpack a message in Java.

    Foo foo = ...;
    Any any = Any.pack(foo);
    ...
    if (any.is(Foo.class)) {
      foo = any.unpack(Foo.class);
    }

 Example 3: Pack and unpack a message in Python.

    foo = Foo(...)
    any = Any()
    any.Pack(foo)
    ...
    if any.Is(Foo.DESCRIPTOR):
      any.Unpack(foo)
      ...

The pack methods provided by protobuf library will by default use
&#39;type.googleapis.com/full.type.name&#39; as the type URL and the unpack
methods only use the fully qualified type name after the last &#39;/&#39;
in the type URL, for example &#34;foo.bar.com/x/y.z&#34; will yield type
name &#34;y.z&#34;.


JSON
====
The JSON representation of an `Any` value uses the regular
representation of the deserialized, embedded message, with an
additional field `@type` which contains the type URL. Example:

    package google.profile;
    message Person {
      string first_name = 1;
      string last_name = 2;
    }

    {
      &#34;@type&#34;: &#34;type.googleapis.com/google.profile.Person&#34;,
      &#34;firstName&#34;: &lt;string&gt;,
      &#34;lastName&#34;: &lt;string&gt;
    }

If the embedded message type is well-known and has a custom JSON
representation, that representation will be embedded adding a field
`value` which holds the custom JSON in addition to the `@type`
field. Example (for message [google.protobuf.Duration][]):

    {
      &#34;@type&#34;: &#34;type.googleapis.com/google.protobuf.Duration&#34;,
      &#34;value&#34;: &#34;1.212s&#34;
    }


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_url | [string](#string) |  | A URL/resource name whose content describes the type of the
serialized protocol buffer message.

For URLs which use the scheme `http`, `https`, or no scheme, the
following restrictions and interpretations apply:

* If no scheme is provided, `https` is assumed.
* The last segment of the URL&#39;s path must represent the fully
  qualified name of the type (as in `path/google.protobuf.Duration`).
  The name should be in a canonical form (e.g., leading &#34;.&#34; is
  not accepted).
* An HTTP GET on the URL must yield a [google.protobuf.Type][]
  value in binary format, or produce an error.
* Applications are allowed to cache lookup results based on the
  URL, or have them precompiled into a binary to avoid any
  lookup. Therefore, binary compatibility needs to be preserved
  on changes to types. (Use versioned type names to manage
  breaking changes.)

Schemes other than `http`, `https` (or the empty scheme) might be
used with implementation specific semantics. |
| value | [bytes](#bytes) |  | Must be a valid serialized protocol buffer of the above specified type. |





 

 

 

 

