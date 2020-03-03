---
title: "Pushing your first WASM filter"
weight: 1
description: "Pull, tag, and push a WASM image."
---

`wasme` enables building and distributing WASM modules as OCI images. In this tutorial,
we detail the steps in pushing an image to a remote registry. This is used to 
distribute WASM modules to remote environments, such as Kubernetes clusters where they 
can be deployed to Envoy proxies.

In this tutorial, we will:

1. (Optionally) Pull a WASM image from `webassemblyhub.io`. 
1. Sign up a user on [`https://webassemblyhub.io`](https://webassemblyhub.io).
1. Tag our image with our new username. 
1. Push a image to `webassemblyhub.io`. 
1. Display our pushed image with `wasme list`

## Pulling a new WASM image

If you've already got a WASM module you'd like to push, you can skip this step. 

{{% notice note %}}
Note that WASM modules must be stored as OCI images in your local cache directory. To see cached images, run `wasme list`.
{{% /notice %}}

Let's pull a published image to our local cache, which we can then re-tag and push:

```bash
# pull the webassemblyhub.io/demo/assemblyscript-test:istio-1.5.0-alpha.0 image
wasme pull webassemblyhub.io/demo/assemblyscript-test:istio-1.5.0-alpha.0
```

```
INFO[0000] Pulling image webassemblyhub.io/demo/assemblyscript-test:istio-1.5.0-alpha.0
INFO[0000] Image: webassemblyhub.io/demo/assemblyscript-test:istio-1.5.0-alpha.0
INFO[0000] Digest: sha256:8b74e9b0bbc5ff674c49cde904669a775a939b4d8f7f72aba88c184d527dfc30
```

## Create a User on [`webassemblyhub.io`](https://webassemblyhub.io)

Pushing images with `wasme` requires a compatible OCI registry. In this tutorial, we'll use [`webassemblyhub.io`](https://webassemblyhub.io) as our remote registry. 

Let's open [`webassemblyhub.io`](https://webassemblyhub.io) in the browser to create an account. 

1. Click **Log In** in the top right:

    ![](../log-in-1.png)

1. Choose **Sign up now** under the login form:

    ![](../log-in-2.png)
    
1. Fill out the sign-up form and click **Sign Up**:

    ![](../log-in-3.png)
    
1. You should now be logged in as a new user:

    ![](../log-in-4.png)

## Tagging our image

Now that we have a user account which will allow us to push images to a remote registry, let's tag the image we pulled:

```bash
wasme tag webassemblyhub.io/demo/assemblyscript-test:istio-1.5.0-alpha.0  webassemblyhub.io/$YOUR_USERNAME/assemblyscript-test:istio-1.5.0-alpha.0 
```

```
INFO[0000] tagged image                                  digest="sha256:8b74e9b0bbc5ff674c49cde904669a775a939b4d8f7f72aba88c184d527dfc30" image="webassemblyhub.io/demo/assemblyscript-test:istio-1.5.0-alpha.0"
```

We should now see both images with `wasme list`:

```
NAME                                             TAG                 SIZE    SHA      UPDATED
webassemblyhub.io/demo/assemblyscript-test      istio-1.5.0-alpha.0 12.5 kB 8b74e9b0 03 Mar 20 09:37 EST
webassemblyhub.io/ilackarms/assemblyscript-test istio-1.5.0-alpha.0 12.5 kB 8b74e9b0 03 Mar 20 09:40 EST
```

## Log In from the `wasme` command line

In order to push images under our new username, we'll need to store our credentials where `wasme` can access them.

Let's do that now with `wasme login`:

```bash
 wasme login -u $YOUR_USERNAME -p $YOUR_PASSWORD
```

```
INFO[0000] Successfully logged in as ilackarms (Scott Weiss)
INFO[0000] stored credentials in /Users/ilackarms/.wasme/credentials.json
```

Great! We're logged in and ready to push our image.

## Push the image

Pushing the image is done with a single command:

```bash
wasme push webassemblyhub.io/$YOUR_USERNAME/assemblyscript-test:istio-1.5.0-alpha.0
```

```
INFO[0000] Pushing image webassemblyhub.io/ilackarms/assemblyscript-test:istio-1.5.0-alpha.0
INFO[0005] Pushed webassemblyhub.io/ilackarms/assemblyscript-test:istio-1.5.0-alpha.0
INFO[0005] Digest: sha256:297527f8740818bd18514a066f936927383a688f827a9573b108609d1a411beb
```
 
Awesome! Our image should be pushed and ready to deploy.
 
## View our published image 

Let's confirm the image has now appeared in our registry. Return to the [`webassemblyhub.io`](https://webassemblyhub.io/user) user home page.

We should be able to see the `assemblyscript-test` image listed under our user:

![](../log-in-5.png) 

We can also verify the image was pushed via the command-line:

```bash
wasme list --search $YOUR_USERNAME
```

```
NAME                                             TAG                 SIZE    SHA      UPDATED
webassemblyhub.io/ilackarms/assemblyscript-test istio-1.5.0-alpha.0 13.6 kB 297527f8 03 Mar 20 14:42 UTC
```

Now that you've practiced pushing an image, you can start deploying your own filters to Envoy with `wasme`!

## Deploying images

For instructions on deploying wasm filters, see [the deployment documentation](../deploy_tutorials)
