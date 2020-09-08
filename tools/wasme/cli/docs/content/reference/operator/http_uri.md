
---
title: "envoy.api.v2.coreenvoy/api/v2/core/http_uri.proto"
---

## Package : `envoy.api.v2.core`



<a name="top"></a>

<a name="API Reference for envoy/api/v2/core/http_uri.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## envoy/api/v2/core/http_uri.proto


## Table of Contents
  - [HttpUri](#envoy.api.v2.core.HttpUri)







<a name="envoy.api.v2.core.HttpUri"></a>

### HttpUri
Envoy external URI descriptor


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uri | [string](#string) |  | The HTTP server URI. It should be a full FQDN with protocol, host and path.

Example:

.. code-block:: yaml

   uri: https://www.googleapis.com/oauth2/v1/certs |
| cluster | [string](#string) |  | A cluster is created in the Envoy &#34;cluster_manager&#34; config
section. This field specifies the cluster name.

Example:

.. code-block:: yaml

   cluster: jwks_cluster |
| timeout | [google.protobuf.Duration](#google.protobuf.Duration) |  | Sets the maximum duration in milliseconds that a response can take to arrive upon request. |





 

 

 

 

