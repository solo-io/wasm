
---
title: "io.prometheus.clientgithub.com/solo-io/solo-kit/api/external/metrics.proto"
---

## Package : `io.prometheus.client`



<a name="top"></a>

<a name="API Reference for github.com/solo-io/solo-kit/api/external/metrics.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/solo-io/solo-kit/api/external/metrics.proto


## Table of Contents
  - [Bucket](#io.prometheus.client.Bucket)
  - [Counter](#io.prometheus.client.Counter)
  - [Gauge](#io.prometheus.client.Gauge)
  - [Histogram](#io.prometheus.client.Histogram)
  - [LabelPair](#io.prometheus.client.LabelPair)
  - [Metric](#io.prometheus.client.Metric)
  - [MetricFamily](#io.prometheus.client.MetricFamily)
  - [Quantile](#io.prometheus.client.Quantile)
  - [Summary](#io.prometheus.client.Summary)
  - [Untyped](#io.prometheus.client.Untyped)

  - [MetricType](#io.prometheus.client.MetricType)






<a name="io.prometheus.client.Bucket"></a>

### Bucket



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cumulative_count | [uint64](#uint64) | optional | Cumulative in increasing order. |
| upper_bound | [double](#double) | optional | Inclusive. |






<a name="io.prometheus.client.Counter"></a>

### Counter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [double](#double) | optional |  |






<a name="io.prometheus.client.Gauge"></a>

### Gauge



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [double](#double) | optional |  |






<a name="io.prometheus.client.Histogram"></a>

### Histogram



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sample_count | [uint64](#uint64) | optional |  |
| sample_sum | [double](#double) | optional |  |
| bucket | [][Bucket](#io.prometheus.client.Bucket) | repeated | Ordered in increasing order of upper_bound, &#43;Inf bucket is optional. |






<a name="io.prometheus.client.LabelPair"></a>

### LabelPair



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) | optional |  |
| value | [string](#string) | optional |  |






<a name="io.prometheus.client.Metric"></a>

### Metric



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| label | [][LabelPair](#io.prometheus.client.LabelPair) | repeated |  |
| gauge | [Gauge](#io.prometheus.client.Gauge) | optional |  |
| counter | [Counter](#io.prometheus.client.Counter) | optional |  |
| summary | [Summary](#io.prometheus.client.Summary) | optional |  |
| untyped | [Untyped](#io.prometheus.client.Untyped) | optional |  |
| histogram | [Histogram](#io.prometheus.client.Histogram) | optional |  |
| timestamp_ms | [int64](#int64) | optional |  |






<a name="io.prometheus.client.MetricFamily"></a>

### MetricFamily



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) | optional |  |
| help | [string](#string) | optional |  |
| type | [MetricType](#io.prometheus.client.MetricType) | optional |  |
| metric | [][Metric](#io.prometheus.client.Metric) | repeated |  |






<a name="io.prometheus.client.Quantile"></a>

### Quantile



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| quantile | [double](#double) | optional |  |
| value | [double](#double) | optional |  |






<a name="io.prometheus.client.Summary"></a>

### Summary



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sample_count | [uint64](#uint64) | optional |  |
| sample_sum | [double](#double) | optional |  |
| quantile | [][Quantile](#io.prometheus.client.Quantile) | repeated |  |






<a name="io.prometheus.client.Untyped"></a>

### Untyped



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [double](#double) | optional |  |





 


<a name="io.prometheus.client.MetricType"></a>

### MetricType


| Name | Number | Description |
| ---- | ------ | ----------- |
| COUNTER | 0 |  |
| GAUGE | 1 |  |
| SUMMARY | 2 |  |
| UNTYPED | 3 |  |
| HISTOGRAM | 4 |  |


 

 

 

