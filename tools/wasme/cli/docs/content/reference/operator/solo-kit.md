
---
title: "core.solo.iogithub.com/solo-io/solo-kit/api/v1/solo-kit.proto"
---

## Package : `core.solo.io`



<a name="top"></a>

<a name="API Reference for github.com/solo-io/solo-kit/api/v1/solo-kit.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/solo-io/solo-kit/api/v1/solo-kit.proto


## Table of Contents
  - [Resource](#core.solo.io.Resource)


  - [File-level Extensions](#github.com/solo-io/solo-kit/api/v1/solo-kit.proto-extensions)





<a name="core.solo.io.Resource"></a>

### Resource



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| short_name | [string](#string) |  | becomes the kubernetes short name for the generated crd |
| plural_name | [string](#string) |  | becomes the kubernetes plural name for the generated crd |
| cluster_scoped | [bool](#bool) |  | the resource lives at the cluster level, namespace is ignored by the server |
| skip_docs_gen | [bool](#bool) |  | indicates whether documentation generation has to be skipped for the given resource, defaults to false |
| skip_hashing_annotations | [bool](#bool) |  | indicates whether annotations should be excluded from the resource&#39;s generated hash function.
if set to true, changes in annotations will not cause a new snapshot to be emitted |





 

 


<a name="github.com/solo-io/solo-kit/api/v1/solo-kit.proto-extensions"></a>

### File-level Extensions
| Extension | Type | Base | Number | Description |
| --------- | ---- | ---- | ------ | ----------- |
| resource | Resource | .google.protobuf.MessageOptions | 10000 | options for a message that&#39;s intended to become a solo-kit resource |

 

 

