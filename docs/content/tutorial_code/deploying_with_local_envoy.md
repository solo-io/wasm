---
title: "Deploying Filters to Local Envoy"
weight: 3
description: Deploy a wasm filter using Gloo as the control plane.
---

In this tutorial we'll deploy an existing WebAssembly (WASM) module from [the WebAssembly Hub](https://webassemblyhub.io) directly to Envoy running locally in `docker`.

## Prepare Envoy to run Locally

Let's create a configuration for running Envoy locally. Paste the following into an `envoy.yaml` file:

```bash
admin:
  access_log_path: /dev/null
  address:
    socket_address:
      address: 0.0.0.0
      port_value: 19000
static_resources:
  listeners:
  - name: listener_0
    address:
      socket_address: { address: 0.0.0.0, port_value: 8080 }
    filter_chains:
    - filters:
      - name: envoy.http_connection_manager
        config:
          codec_type: AUTO
          stat_prefix: ingress_http
          route_config:
            name: test
            virtual_hosts:
            - name: jsonplaceholder
              domains: ["*"]
              routes:
              - match: { prefix: "/" }
                route:
                  cluster: static-cluster
                  auto_host_rewrite: true
          http_filters:
          - name: envoy.router
  clusters:
  - name: static-cluster
    connect_timeout: 0.25s
    type: LOGICAL_DNS
    lb_policy: ROUND_ROBIN
    dns_lookup_family: V4_ONLY
    tls_context:
      sni: jsonplaceholder.typicode.com
    hosts: [{ socket_address: { address: jsonplaceholder.typicode.com, port_value: 443, ipv4_compat: true } }]
```

Now let's run Envoy using docker:

```bash

docker run --rm --name e2e_envoy -p 8080:8080 -p 8443:8443 -p 19001:19001 \
    --entrypoint=envoy \
    quay.io/solo-io/gloo-envoy-wasm-wrapper:1.2.12 \
    --disable-hot-restart --log-level debug --config-yaml "$(cat envoy.yaml)"

```

In another terminal, test that you can query the proxied service:

```bash
curl localhost:8080/posts/1 -v
```

```
< HTTP/1.1 200 OK
< date: Mon, 23 Dec 2019 14:51:18 GMT
< content-type: application/json; charset=utf-8
< content-length: 292
< set-cookie: __cfduid=d4c04bb71af3e2c5340fe137a162315d21577112678; expires=Wed, 22-Jan-20 14:51:18 GMT; path=/; domain=.typicode.com; HttpOnly; SameSite=Lax
< x-powered-by: Express
< vary: Origin, Accept-Encoding
< access-control-allow-credentials: true
< cache-control: max-age=14400
< pragma: no-cache
< expires: -1
< x-content-type-options: nosniff
< etag: W/"124-yiKdLzqO5gfBrJFrcdJ8Yq0LGnU"
< via: 1.1 vegur
< cf-cache-status: HIT
< age: 7189
< accept-ranges: bytes
< expect-ct: max-age=604800, report-uri="https://report-uri.cloudflare.com/cdn-cgi/beacon/expect-ct"
< server: envoy
< cf-ray: 549b2721ea21c5c4-EWR
< x-envoy-upstream-service-time: 78
<
{
  "userId": 1,
  "id": 1,
  "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
  "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
* Connection #0 to host localhost left intact
}
```

This is the standard response if all of our setup steps worked correctly. Let's now add the `hello` filter to our Envoy.

## Pull the filter

Next we'll pull the filter with `wasme`. For this example we'll use the `webassemblyhub.io/ilackarms/hello:v0.1` filter,
which appends a `hello: World!` header to response requests. 

To pull the filter:

```shell
wasme pull webassemblyhub.io/ilackarms/hello:v0.1 -o hello.wasm
```

```
webassemblyhub.io/ilackarms/hello:v0.1 [{MediaType:application/vnd.io.solo.wasm.config.v1+json Digest:sha256:c6c060cce61aedc5b04c92f62b0b3238897958e448e4c2498f8302dd3af03b55 Size:39 URLs:[] Annotations:map[] Platform:<nil>} {MediaType:application/vnd.io.solo.wasm.code.v1+wasm Digest:sha256:d23cdeb8e7096cc5b2b2f7959a9be9412f840bd0a5ff56f364005efd5fc41c66 Size:1042994 URLs:[] Annotations:map[org.opencontainers.image.title:code.wasm] Platform:<nil>}] {MediaType:application/vnd.oci.image.manifest.v1+json Digest:sha256:d64b8187e3f71922dbbb332d4f8136519e0768ae9ecf150b15150b5a02eb4d63 Size:409 URLs:[] Annotations:map[] Platform:<nil>} <nil>
INFO[0000] Pulled filter image webassemblyhub.io/ilackarms/hello:v0.1
```

We should see the filter `hello.wasm` has been downloaded to our current directory:

```
ls -l
```

```
-rw-r--r--   1 ilackarms  staff   1.0M Dec 20 16:09 hello.wasm
```

# Update the Config

Now we'll need to update Envoy's configuration. 

For the purposes of this tutorial we are configuring Envoy statically, which means we'll need to restart Envoy in order to pick up the new WASM filter. 

{{% notice note %}}
In production environments, it's better to use dynamic configuration when possible.
{{% /notice %}}

To add our wasm filter to the config, simply run:

```bash
deploy envoy webassemblyhub.io/ilackarms/hello:v0.1 \
  --id=myfilter \
  --in=envoy.yaml \
  --out=envoy.yaml \
  --filter=/etc/hello.wasm
```

# Run and Test

Now we'll run Envoy again, this time mounting `hello.wasm` to `/etc` in the container:

```bash

docker run --rm --name e2e_envoy -p 8080:8080 -p 8443:8443 -p 19001:19001 \
    -v $(pwd)/hello.wasm:/etc/hello.wasm --entrypoint=envoy \
    quay.io/solo-io/gloo-envoy-wasm-wrapper:1.2.12 \
    --disable-hot-restart --log-level debug --config-yaml "$(cat envoy.yaml)"

```

In another terminal, test with another `curl`:

```bash
curl localhost:8080/posts/1 -v
```

We should see the `hello: World!` header in the response:

{{< highlight bash "hl_lines=22" >}}
< HTTP/1.1 200 OK
< date: Mon, 23 Dec 2019 15:46:42 GMT
< content-type: application/json; charset=utf-8
< content-length: 292
< set-cookie: __cfduid=d5b623b10179315552ba69f38674bb4351577116002; expires=Wed, 22-Jan-20 15:46:42 GMT; path=/; domain=.typicode.com; HttpOnly; SameSite=Lax
< x-powered-by: Express
< vary: Origin, Accept-Encoding
< access-control-allow-credentials: true
< cache-control: max-age=14400
< pragma: no-cache
< expires: -1
< x-content-type-options: nosniff
< etag: W/"124-yiKdLzqO5gfBrJFrcdJ8Yq0LGnU"
< via: 1.1 vegur
< cf-cache-status: HIT
< age: 2477
< accept-ranges: bytes
< expect-ct: max-age=604800, report-uri="https://report-uri.cloudflare.com/cdn-cgi/beacon/expect-ct"
< server: envoy
< cf-ray: 549b78488804e76c-EWR
< x-envoy-upstream-service-time: 84
< hello: World!
< location: envoy-wasm
<
{
  "userId": 1,
  "id": 1,
  "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
  "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
* Connection #0 to host localhost left intact
}
{{< /highlight >}}

Great! If everything worked correctly, we should see the 
above response. If you encountered an issue anywhere along the way, please report it to the `wasme` authors at https://github.com/solo-io/wasme/issues/new.

## Cleaning up

We can remove our filter from the static config with the `wasme undeploy` command:

```bash
wasme undeploy envoy \
    --id=myfilter \
    --in envoy.yaml \
    --out envoy.yaml
```

Restart Envoy:

```bash
docker run --rm --name e2e_envoy -p 8080:8080 -p 8443:8443 -p 19001:19001 \
    --entrypoint=envoy \
    quay.io/solo-io/gloo-envoy-wasm-wrapper:1.2.12 \
    --disable-hot-restart --log-level debug --config-yaml "$(cat envoy.yaml)"
```

Then re-try the `curl`:

```shell
curl -v $URL/posts/1
```

```
*   Trying 34.73.225.160...
* TCP_NODELAY set
* Connected to 34.73.225.160 (34.73.225.160) port 80 (#0)
> GET /api/pets HTTP/1.1
> Host: 34.73.225.160
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< content-type: application/xml
< date: Fri, 20 Dec 2019 19:19:13 GMT
< content-length: 86
< x-envoy-upstream-service-time: 1
< server: envoy
<
[{"id":1,"name":"Dog","status":"available"},{"id":2,"name":"Cat","status":"pending"}]
* Connection #0 to host 34.73.225.160 left intact
```

Cool! We've just seen how easy it is to dynamically add and remove filters from Envoy using `wasme`.

For more information and support using `wasme` and the Web Assembly Hub, visit the Solo.io slack channel at
https://slack.solo.io.
