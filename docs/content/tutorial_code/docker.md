---
title: "Envoy WASM in docker"
weight: 2
description: Run/Test your WebAssembly filter locally with Envoy in docker
---

`wasme` makes creating Envoy filters easier than ever, now the question becomes how do we run and test them.
This tutorial will explain in detail how to run and test Envoy wasm filters locally, before deploying them.

For the purposes of this example we will be using a simple stats filter `webassemblyhub.io/yuval-k/metrics:v1`. However, 
the techniques used here can be applied to any and all wasm filters build this way.

#### Prerequisites
* [docker](https://www.docker.com/)
* [wasme](https://github.com/solo-io/wasme)

## Configuration

The first step in testing the filter locally is creating the static E config which will serve out new filter.
Documentation on E's API can be found [here](https://www.envoyproxy.io/docs/envoy/v1.12.0/api/api). 

Below is the full yaml necessary to test. We will unpack it step by step after.
```yaml
admin:
  access_log_path: /dev/null
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 19001
static_resources:
  listeners:
  - name: listener_0
    address:
      socket_address: { address: 0.0.0.0, port_value: 8080 }
    filter_chains:
    - filters:
      - name: envoy.http_connection_manager
        config:
          codec_type: auto
          stat_prefix: http
          route_config:
            name: test
            virtual_hosts:
            - name: test
              domains: ["*"]
              routes:
              - match: { prefix: "/" }
                route: { cluster: static-cluster }
          http_filters:
          - name: envoy.filters.http.wasm
            config: 
              config: 
                name: "test"
                root_id: "stats_root_id"
                vm_config:
                  vm_id: test
                  runtime: envoy.wasm.runtime.v8
                  allow_precompiled: false
                  code:
                    local:
                      filename: /etc/filter.wasm
                configuration: |
                  {}
          - name: envoy.router
  clusters:
  - name: static-cluster
    connect_timeout: 0.25s
    type: logical_dns
    lb_policy: round_robin
    hosts: [{ socket_address: { address: host.docker.internal, port_value: 10101 } }]
```

The first part of the config is not particularly interesting. The admin section simply tells Envoy which port to use for 
the admin interface.

Below that is where things start to get interesting. Specifically the `http_filters` section of the `envoy.http_connection_manager`.
This section features a new type of config, specifically the envoy-wasm config. This filter functions similarly to other Envoy filters 
with one distinct difference. It does not handle the request, but rather sends the data to the wasm module specified by the config.
```yaml
  - name: envoy.filters.http.wasm
    config: 
      config: 
        name: "test"
        root_id: "stats_root_id"
        vm_config:
          vm_id: test
          runtime: envoy.wasm.runtime.v8
          allow_precompiled: false
          code:
            local:
              filename: /etc/filter.wasm
        configuration: |
          {}
```
The full API for the above config can be found [here](https://github.com/envoyproxy/envoy-wasm/blob/master/api/envoy/config/filter/http/wasm/v2/wasm.proto).
There are a few things here which are worth highlighting.

1) `root_id: "stats_root_id"`: the root_id is a new concept to wasm filters, but very important. Similar to the "filter_name" in traditional
Envoy filters, this is how Envoy knows which wasm filter to use. If this id does not match any loaded wasm filter, it will cause Envoy to crash.
2) `runtime: envoy.wasm.runtime.v8`: the runtime is the type of wasm vm with which Envoy will run the wasm module. Currently the 2 options
are V8, and WAVM.
3) The code section in this example loads the wasm filter from a local code source, but this can also be configured to load from a remote source.
```yaml
  code:
    local:
      filename: /etc/filter.wasm
``` 

4) This last section will be most familiar to anyone familiar with configuring Envoy filters. The configuration is passed 
to the wasm filter as a json blob, exactly as defined below.
```yaml
configuration: |
  {}
```

Now that we have the config, we should save it to a local file called `config.yml`, we are going to need it soon.

## Filter

This step is much simpler than the previous, only a single command in this case.
```shell script
wasme pull webassemblyhub.io/yuval-k/metrics:v1
```

This command is relatively simple, but it's worth quickly explaining. It uses the WASM filter registry provided by wasme, and pulls
the given filter. After pulling it, it saves it to a local file called `filter.wasm`. Now that the filter is saved locally we can load
it into Envoy.

## Running Envoy

Now it is time to actually run Envoy. 
```shell script
docker run  --rm --name e2e_envoy -p 8080:8080 -p 8443:8443 -p 19001:19001 -v $(pwd)/filter.wasm:/etc/filter.wasm --entrypoint=envoy quay.io/solo-io/gloo-envoy-wasm-wrapper:v1.2.5 --disable-hot-restart --log-level debug --config-yaml "$(cat config.yml)"
```

Run the above command and then open up a new terminal to communicate with Envoy.
In the second terminal run: 
```shell script
curl localhost:8080

upstream connect error or disconnect/reset before headers. reset reason: connection failure
```
The response will not be happy, but that is intentional for now. The `static_cluster` we created earlier does not actually 
point anywhere, but still allows us to illustrate the usefulness of wasm for building filters.

Now navigate to http://localhost:19001/stats/prometheus in a browser, and our new Envoy wasm stats will be waiting :)
