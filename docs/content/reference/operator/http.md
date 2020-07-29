
---
title: "google.apigithub.com/solo-io/solo-kit/api/external/google/api/http.proto"
---

## Package : `google.api`



<a name="top"></a>

<a name="API Reference for github.com/solo-io/solo-kit/api/external/google/api/http.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/solo-io/solo-kit/api/external/google/api/http.proto


## Table of Contents
  - [CustomHttpPattern](#google.api.CustomHttpPattern)
  - [Http](#google.api.Http)
  - [HttpRule](#google.api.HttpRule)







<a name="google.api.CustomHttpPattern"></a>

### CustomHttpPattern
A custom pattern is used for defining custom HTTP verb.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [string](#string) |  | The name of this custom HTTP verb. |
| path | [string](#string) |  | The path matched by this custom verb. |






<a name="google.api.Http"></a>

### Http
Defines the HTTP configuration for an API service. It contains a list of
[HttpRule][google.api.HttpRule], each specifying the mapping of an RPC method
to one or more HTTP REST API methods.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rules | [][HttpRule](#google.api.HttpRule) | repeated | A list of HTTP configuration rules that apply to individual API methods.

**NOTE:** All service configuration rules follow &#34;last one wins&#34; order. |
| fully_decode_reserved_expansion | [bool](#bool) |  | When set to true, URL path parmeters will be fully URI-decoded except in
cases of single segment matches in reserved expansion, where &#34;%2F&#34; will be
left encoded.

The default behavior is to not decode RFC 6570 reserved characters in multi
segment matches. |






<a name="google.api.HttpRule"></a>

### HttpRule
`HttpRule` defines the mapping of an RPC method to one or more HTTP
REST API methods. The mapping specifies how different portions of the RPC
request message are mapped to URL path, URL query parameters, and
HTTP request body. The mapping is typically specified as an
`google.api.http` annotation on the RPC method,
see &#34;google/api/annotations.proto&#34; for details.

The mapping consists of a field specifying the path template and
method kind.  The path template can refer to fields in the request
message, as in the example below which describes a REST GET
operation on a resource collection of messages:


    service Messaging {
      rpc GetMessage(GetMessageRequest) returns (Message) {
        option (google.api.http).get = &#34;/v1/messages/{message_id}/{sub.subfield}&#34;;
      }
    }
    message GetMessageRequest {
      message SubMessage {
        string subfield = 1;
      }
      string message_id = 1; // mapped to the URL
      SubMessage sub = 2;    // `sub.subfield` is url-mapped
    }
    message Message {
      string text = 1; // content of the resource
    }

The same http annotation can alternatively be expressed inside the
`GRPC API Configuration` YAML file.

    http:
      rules:
        - selector: &lt;proto_package_name&gt;.Messaging.GetMessage
          get: /v1/messages/{message_id}/{sub.subfield}

This definition enables an automatic, bidrectional mapping of HTTP
JSON to RPC. Example:

HTTP | RPC
-----|-----
`GET /v1/messages/123456/foo`  | `GetMessage(message_id: &#34;123456&#34; sub: SubMessage(subfield: &#34;foo&#34;))`

In general, not only fields but also field paths can be referenced
from a path pattern. Fields mapped to the path pattern cannot be
repeated and must have a primitive (non-message) type.

Any fields in the request message which are not bound by the path
pattern automatically become (optional) HTTP query
parameters. Assume the following definition of the request message:


    service Messaging {
      rpc GetMessage(GetMessageRequest) returns (Message) {
        option (google.api.http).get = &#34;/v1/messages/{message_id}&#34;;
      }
    }
    message GetMessageRequest {
      message SubMessage {
        string subfield = 1;
      }
      string message_id = 1; // mapped to the URL
      int64 revision = 2;    // becomes a parameter
      SubMessage sub = 3;    // `sub.subfield` becomes a parameter
    }


This enables a HTTP JSON to RPC mapping as below:

HTTP | RPC
-----|-----
`GET /v1/messages/123456?revision=2&amp;sub.subfield=foo` | `GetMessage(message_id: &#34;123456&#34; revision: 2 sub: SubMessage(subfield: &#34;foo&#34;))`

Note that fields which are mapped to HTTP parameters must have a
primitive type or a repeated primitive type. Message types are not
allowed. In the case of a repeated type, the parameter can be
repeated in the URL, as in `...?param=A&amp;param=B`.

For HTTP method kinds which allow a request body, the `body` field
specifies the mapping. Consider a REST update method on the
message resource collection:


    service Messaging {
      rpc UpdateMessage(UpdateMessageRequest) returns (Message) {
        option (google.api.http) = {
          put: &#34;/v1/messages/{message_id}&#34;
          body: &#34;message&#34;
        };
      }
    }
    message UpdateMessageRequest {
      string message_id = 1; // mapped to the URL
      Message message = 2;   // mapped to the body
    }


The following HTTP JSON to RPC mapping is enabled, where the
representation of the JSON in the request body is determined by
protos JSON encoding:

HTTP | RPC
-----|-----
`PUT /v1/messages/123456 { &#34;text&#34;: &#34;Hi!&#34; }` | `UpdateMessage(message_id: &#34;123456&#34; message { text: &#34;Hi!&#34; })`

The special name `*` can be used in the body mapping to define that
every field not bound by the path template should be mapped to the
request body.  This enables the following alternative definition of
the update method:

    service Messaging {
      rpc UpdateMessage(Message) returns (Message) {
        option (google.api.http) = {
          put: &#34;/v1/messages/{message_id}&#34;
          body: &#34;*&#34;
        };
      }
    }
    message Message {
      string message_id = 1;
      string text = 2;
    }


The following HTTP JSON to RPC mapping is enabled:

HTTP | RPC
-----|-----
`PUT /v1/messages/123456 { &#34;text&#34;: &#34;Hi!&#34; }` | `UpdateMessage(message_id: &#34;123456&#34; text: &#34;Hi!&#34;)`

Note that when using `*` in the body mapping, it is not possible to
have HTTP parameters, as all fields not bound by the path end in
the body. This makes this option more rarely used in practice of
defining REST APIs. The common usage of `*` is in custom methods
which don&#39;t use the URL at all for transferring data.

It is possible to define multiple HTTP methods for one RPC by using
the `additional_bindings` option. Example:

    service Messaging {
      rpc GetMessage(GetMessageRequest) returns (Message) {
        option (google.api.http) = {
          get: &#34;/v1/messages/{message_id}&#34;
          additional_bindings {
            get: &#34;/v1/users/{user_id}/messages/{message_id}&#34;
          }
        };
      }
    }
    message GetMessageRequest {
      string message_id = 1;
      string user_id = 2;
    }


This enables the following two alternative HTTP JSON to RPC
mappings:

HTTP | RPC
-----|-----
`GET /v1/messages/123456` | `GetMessage(message_id: &#34;123456&#34;)`
`GET /v1/users/me/messages/123456` | `GetMessage(user_id: &#34;me&#34; message_id: &#34;123456&#34;)`

# Rules for HTTP mapping

The rules for mapping HTTP path, query parameters, and body fields
to the request message are as follows:

1. The `body` field specifies either `*` or a field path, or is
   omitted. If omitted, it indicates there is no HTTP request body.
2. Leaf fields (recursive expansion of nested messages in the
   request) can be classified into three types:
    (a) Matched in the URL template.
    (b) Covered by body (if body is `*`, everything except (a) fields;
        else everything under the body field)
    (c) All other fields.
3. URL query parameters found in the HTTP request are mapped to (c) fields.
4. Any body sent with an HTTP request can contain only (b) fields.

The syntax of the path template is as follows:

    Template = &#34;/&#34; Segments [ Verb ] ;
    Segments = Segment { &#34;/&#34; Segment } ;
    Segment  = &#34;*&#34; | &#34;**&#34; | LITERAL | Variable ;
    Variable = &#34;{&#34; FieldPath [ &#34;=&#34; Segments ] &#34;}&#34; ;
    FieldPath = IDENT { &#34;.&#34; IDENT } ;
    Verb     = &#34;:&#34; LITERAL ;

The syntax `*` matches a single path segment. The syntax `**` matches zero
or more path segments, which must be the last part of the path except the
`Verb`. The syntax `LITERAL` matches literal text in the path.

The syntax `Variable` matches part of the URL path as specified by its
template. A variable template must not contain other variables. If a variable
matches a single path segment, its template may be omitted, e.g. `{var}`
is equivalent to `{var=*}`.

If a variable contains exactly one path segment, such as `&#34;{var}&#34;` or
`&#34;{var=*}&#34;`, when such a variable is expanded into a URL path, all characters
except `[-_.~0-9a-zA-Z]` are percent-encoded. Such variables show up in the
Discovery Document as `{var}`.

If a variable contains one or more path segments, such as `&#34;{var=foo/*}&#34;`
or `&#34;{var=**}&#34;`, when such a variable is expanded into a URL path, all
characters except `[-_.~/0-9a-zA-Z]` are percent-encoded. Such variables
show up in the Discovery Document as `{&#43;var}`.

NOTE: While the single segment variable matches the semantics of
[RFC 6570](https://tools.ietf.org/html/rfc6570) Section 3.2.2
Simple String Expansion, the multi segment variable **does not** match
RFC 6570 Reserved Expansion. The reason is that the Reserved Expansion
does not expand special characters like `?` and `#`, which would lead
to invalid URLs.

NOTE: the field paths in variables and in the `body` must not refer to
repeated fields or map fields.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| selector | [string](#string) |  | Selects methods to which this rule applies.

Refer to [selector][google.api.DocumentationRule.selector] for syntax details. |
| get | [string](#string) |  | Used for listing and getting information about resources. |
| put | [string](#string) |  | Used for updating a resource. |
| post | [string](#string) |  | Used for creating a resource. |
| delete | [string](#string) |  | Used for deleting a resource. |
| patch | [string](#string) |  | Used for updating a resource. |
| custom | [CustomHttpPattern](#google.api.CustomHttpPattern) |  | The custom pattern is used for specifying an HTTP method that is not
included in the `pattern` field, such as HEAD, or &#34;*&#34; to leave the
HTTP method unspecified for this rule. The wild-card rule is useful
for services that provide content to Web (HTML) clients. |
| body | [string](#string) |  | The name of the request field whose value is mapped to the HTTP body, or
`*` for mapping all fields not captured by the path pattern to the HTTP
body. NOTE: the referred field must not be a repeated field and must be
present at the top-level of request message type. |
| additional_bindings | [][HttpRule](#google.api.HttpRule) | repeated | Additional HTTP bindings for the selector. Nested bindings must
not contain an `additional_bindings` field themselves (that is,
the nesting may only be one level deep). |





 

 

 

 

