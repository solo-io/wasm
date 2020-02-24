---
title: "Deploying Filters to Local Envoy"
weight: 3
description: Run Envoy locally to test a WASM filter attached.
---

In this tutorial we'll deploy an existing WebAssembly (WASM) module from [the WebAssembly Hub](https://webassemblyhub.io) directly to Envoy running locally in `docker`.

Deploying filters locally with Envoy is a great way to develop and test custom filters. `wasme deploy envoy` runs a single instance of Envoy in a Docker container on the local machine. Envoy is started with a static configuration, which which defaults to a single route to `jsonplaceholder.typicode.com` unless supplied by the user.

# Tutorial

In this example, we'll be working with a filter pulled from the registry at [`webassemblyhub.io`](https://webassemblyhub.io). 
 
 If you are working with a filter you've built yourself, you can skip to [Run the Filter](#Run the Filter).

## Pull the filter

First we'll pull the filter with `wasme`. For this example we'll use the `webassemblyhub.io/ilackarms/assemblyscript-test:istio-1.5.0-alpha.0` filter,
which appends a header to HTTP responses. 

{{% notice note %}}
The `webassemblyhub.io/ilackarms/assemblyscript-test:istio-1.5.0-alpha.0` filter image is compatible with the versions of Envoy packaged with Gloo `1.3.x` and Istio `1.5.x`. 
{{% /notice %}}


To pull the filter:

```shell
wasme pull webassemblyhub.io/ilackarms/assemblyscript-test:istio-1.5.0-alpha.0
```

```
INFO[0000] Pulling image webassemblyhub.io/ilackarms/assemblyscript-test:istio-1.5.0-alpha.0
INFO[0000] Image: webassemblyhub.io/ilackarms/assemblyscript-test:istio-1.5.0-alpha.0
INFO[0000] Digest: sha256:8b74e9b0bbc5ff674c49cde904669a775a939b4d8f7f72aba88c184d527dfc30
```

We should see image `webassemblyhub.io/ilackarms/assemblyscript-test:istio-1.5.0-alpha.0` has been downloaded to local cache:

```
wasme list
```

```
NAME                                            TAG                 SIZE    SHA      UPDATED
webassemblyhub.io/ilackarms/assemblyscript-test istio-1.5.0-alpha.0 12.5 kB 8b74e9b0 13 Feb 20 13:59 EST
```

## Run the Filter

Running the filter with a local instance of Envoy is as done with a single command:

```bash
wasme deploy envoy webassemblyhub.io/ilackarms/assemblyscript-test:istio-1.5.0-alpha.0
```

{{% notice note %}}
The `wasme deploy envoy` command runs the filter with an Envoy image built for Gloo `1.3.5`. You can override this with the `--envoy-image` flag. 
{{% /notice %}}

This will start Envoy in a docker container. We should see logs printed in the current terminal:

```
INFO[0000] mounting filter file at /Users/ilackarms/.wasme/store/7bda74acb544159ac98f58e85d573d12/filter.wasm
INFO[0000] running envoy-in-docker                       container_name=add_header envoy_image="quay.io/solo-io/gloo-envoy-wasm-wrapper:1.3.5" filter_image="webassemblyhub.io/ilackarms/assemblyscript-test:istio-1.5.0-alpha.0"
[2020-02-17 17:50:52.050][1][info][main] [external/envoy/source/server/server.cc:252] initializing epoch 0 (hot restart version=disabled)
[2020-02-17 17:50:52.050][1][info][main] [external/envoy/source/server/server.cc:254] statically linked extensions:
[2020-02-17 17:50:52.050][1][info][main] [external/envoy/source/server/server.cc:256]   access_loggers: envoy.file_access_log, envoy.http_grpc_access_log, envoy.tcp_grpc_access_log, envoy.wasm_access_log
[2020-02-17 17:50:52.050][1][info][main] [external/envoy/source/server/server.cc:256]   clusters: envoy.cluster.eds, envoy.cluster.logical_dns, envoy.cluster.original_dst, envoy.cluster.static, envoy.cluster.strict_dns, envoy.clusters.aggregate, envoy.clusters.dynamic_forward_proxy, envoy.clusters.redis
```

`Ctrl+C` can be used at any time to terminate the container.

We can see the running container in another terminal session:

```bash
docker ps
```

```
CONTAINER ID        IMAGE                                           COMMAND                  CREATED              STATUS              PORTS                                              NAMES
ca6ee2f57522        quay.io/solo-io/gloo-envoy-wasm-wrapper:1.3.5   "envoy --disable-hot…"   About a minute ago   Up About a minute   0.0.0.0:8080->8080/tcp, 0.0.0.0:19000->19000/tcp   add_header
```

Docker is port-forwarding port `8080` of the container to our local machine. 

Let's try hitting the container with a request:

```bash
curl localhost:8080/posts/1 -v
```

{{< highlight yaml "hl_lines=30" >}}
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET /posts/1 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< date: Mon, 17 Feb 2020 17:54:57 GMT
< content-type: application/json; charset=utf-8
< content-length: 292
< set-cookie: __cfduid=de818e5f1056f5e426c15afc6c73380d21581962097; expires=Wed, 18-Mar-20 17:54:57 GMT; path=/; domain=.typicode.com; HttpOnly; SameSite=Lax
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
< age: 2717
< accept-ranges: bytes
< expect-ct: max-age=604800, report-uri="https://report-uri.cloudflare.com/cdn-cgi/beacon/expect-ct"
< server: envoy
< cf-ray: 5669a1240ab7ff90-BOS
< x-envoy-upstream-service-time: 76
< hello: world!
<
{
  "userId": 1,
  "id": 1,
  "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
  "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
* Connection #0 to host localhost left intact
}
{{< /highlight >}}

If everything worked correctly, we should see the `hello: world!` header appended in the `curl` response.

# Summary

Using `wasme deploy envoy`, we can locally test filters against Envoy. See [the CLI documentation]({{< versioned_link_path fromRoot="/reference/cli/wasme_deploy_envoy">}}) for all the supported options for this command. 

For more information and support using `wasme` and the Web Assembly Hub, visit the Solo.io slack channel at
https://slack.solo.io.
