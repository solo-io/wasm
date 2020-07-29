
---
title: "envoy.typeenvoy/type/percent.proto"
---

## Package : `envoy.type`



<a name="top"></a>

<a name="API Reference for envoy/type/percent.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## envoy/type/percent.proto


## Table of Contents
  - [FractionalPercent](#envoy.type.FractionalPercent)
  - [Percent](#envoy.type.Percent)

  - [FractionalPercent.DenominatorType](#envoy.type.FractionalPercent.DenominatorType)






<a name="envoy.type.FractionalPercent"></a>

### FractionalPercent
A fractional percentage is used in cases in which for performance reasons performing floating
point to integer conversions during randomness calculations is undesirable. The message includes
both a numerator and denominator that together determine the final fractional value.

* **Example**: 1/100 = 1%.
* **Example**: 3/10000 = 0.03%.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| numerator | [uint32](#uint32) |  | Specifies the numerator. Defaults to 0. |
| denominator | [FractionalPercent.DenominatorType](#envoy.type.FractionalPercent.DenominatorType) |  | Specifies the denominator. If the denominator specified is less than the numerator, the final
fractional percentage is capped at 1 (100%). |






<a name="envoy.type.Percent"></a>

### Percent
Identifies a percentage, in the range [0.0, 100.0].


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [double](#double) |  |  |





 


<a name="envoy.type.FractionalPercent.DenominatorType"></a>

### FractionalPercent.DenominatorType
Fraction percentages support several fixed denominator values.

| Name | Number | Description |
| ---- | ------ | ----------- |
| HUNDRED | 0 | 100.

**Example**: 1/100 = 1%. |
| TEN_THOUSAND | 1 | 10,000.

**Example**: 1/10000 = 0.01%. |
| MILLION | 2 | 1,000,000.

**Example**: 1/1000000 = 0.0001%. |


 

 

 

