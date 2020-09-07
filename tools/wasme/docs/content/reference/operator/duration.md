
---
title: "google.protobufgithub.com/solo-io/solo-kit/api/external/google/protobuf/duration.proto"
---

## Package : `google.protobuf`



<a name="top"></a>

<a name="API Reference for github.com/solo-io/solo-kit/api/external/google/protobuf/duration.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## github.com/solo-io/solo-kit/api/external/google/protobuf/duration.proto


## Table of Contents
  - [Duration](#google.protobuf.Duration)







<a name="google.protobuf.Duration"></a>

### Duration
A Duration represents a signed, fixed-length span of time represented
as a count of seconds and fractions of seconds at nanosecond
resolution. It is independent of any calendar and concepts like &#34;day&#34;
or &#34;month&#34;. It is related to Timestamp in that the difference between
two Timestamp values is a Duration and it can be added or subtracted
from a Timestamp. Range is approximately &#43;-10,000 years.

# Examples

Example 1: Compute Duration from two Timestamps in pseudo code.

    Timestamp start = ...;
    Timestamp end = ...;
    Duration duration = ...;

    duration.seconds = end.seconds - start.seconds;
    duration.nanos = end.nanos - start.nanos;

    if (duration.seconds &lt; 0 &amp;&amp; duration.nanos &gt; 0) {
      duration.seconds &#43;= 1;
      duration.nanos -= 1000000000;
    } else if (durations.seconds &gt; 0 &amp;&amp; duration.nanos &lt; 0) {
      duration.seconds -= 1;
      duration.nanos &#43;= 1000000000;
    }

Example 2: Compute Timestamp from Timestamp &#43; Duration in pseudo code.

    Timestamp start = ...;
    Duration duration = ...;
    Timestamp end = ...;

    end.seconds = start.seconds &#43; duration.seconds;
    end.nanos = start.nanos &#43; duration.nanos;

    if (end.nanos &lt; 0) {
      end.seconds -= 1;
      end.nanos &#43;= 1000000000;
    } else if (end.nanos &gt;= 1000000000) {
      end.seconds &#43;= 1;
      end.nanos -= 1000000000;
    }

Example 3: Compute Duration from datetime.timedelta in Python.

    td = datetime.timedelta(days=3, minutes=10)
    duration = Duration()
    duration.FromTimedelta(td)

# JSON Mapping

In JSON format, the Duration type is encoded as a string rather than an
object, where the string ends in the suffix &#34;s&#34; (indicating seconds) and
is preceded by the number of seconds, with nanoseconds expressed as
fractional seconds. For example, 3 seconds with 0 nanoseconds should be
encoded in JSON format as &#34;3s&#34;, while 3 seconds and 1 nanosecond should
be expressed in JSON format as &#34;3.000000001s&#34;, and 3 seconds and 1
microsecond should be expressed in JSON format as &#34;3.000001s&#34;.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| seconds | [int64](#int64) |  | Signed seconds of the span of time. Must be from -315,576,000,000
to &#43;315,576,000,000 inclusive. Note: these bounds are computed from:
60 sec/min * 60 min/hr * 24 hr/day * 365.25 days/year * 10000 years |
| nanos | [int32](#int32) |  | Signed fractions of a second at nanosecond resolution of the span
of time. Durations less than one second are represented with a 0
`seconds` field and a positive or negative `nanos` field. For durations
of one second or more, a non-zero value for the `nanos` field must be
of the same sign as the `seconds` field. Must be from -999,999,999
to &#43;999,999,999 inclusive. |





 

 

 

 

