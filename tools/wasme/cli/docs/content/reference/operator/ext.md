
---
title: "extprotoextproto/ext.proto"
---

## Package : `extproto`



<a name="top"></a>

<a name="API Reference for extproto/ext.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## extproto/ext.proto


## Table of Contents


  - [File-level Extensions](#extproto/ext.proto-extensions)
  - [File-level Extensions](#extproto/ext.proto-extensions)
  - [File-level Extensions](#extproto/ext.proto-extensions)
  - [File-level Extensions](#extproto/ext.proto-extensions)
  - [File-level Extensions](#extproto/ext.proto-extensions)
  - [File-level Extensions](#extproto/ext.proto-extensions)




 

 


<a name="extproto/ext.proto-extensions"></a>

### File-level Extensions
| Extension | Type | Base | Number | Description |
| --------- | ---- | ---- | ------ | ----------- |
| skip_hashing | bool | .google.protobuf.FieldOptions | 10071 | Rules specify the validations to be performed on this field. By default,
no validation is performed against a field. |
| skip_merging | bool | .google.protobuf.FieldOptions | 10072 | This field will not be merged when a message&#39;s Merge() method is called. |
| equal_all | bool | .google.protobuf.FileOptions | 10072 |  |
| hash_all | bool | .google.protobuf.FileOptions | 10071 | Disabled nullifies any validation rules for this message, including any
message fields associated with it that do support validation. |
| merge_all | bool | .google.protobuf.FileOptions | 10073 | enabling merge_all for a file will generate a Merge(m) method for all Messages in the file.
Merge(m) will replace non-nil fields from a Proto passed as an override. |
| skip_merging_oneof | bool | .google.protobuf.OneofOptions | 10072 | The fields in this oneof will not be merged when a message&#39;s Merge() method is called. |

 

 

