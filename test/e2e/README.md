# Set up CI Cluster

The e2e tests depend on Istio and Gloo being installed to the target cluster. To set up the cluster:

* Install Gloo to cluster with WASM enabled:
```bash
 k create ns gloo-system; helm install --namespace gloo-system --set global.wasm.enabled=true gloo gloo/gloo 
```

* Install petstore to the cluster:
```bash
kubectl apply -oyaml -f https://raw.githubusercontent.com/solo-io/gloo/master/example/petstore/petstore.yaml --dry-run
```

* Create a Gloo route to the petstore:
```bash
cat <<EOF | kubectl apply -f -                    
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: default        
  namespace: gloo-system  
spec:                                                     
  virtualHost:                                         
    domains:                
    - '*'                            
    routes:
    - matchers:
      - prefix: /         
      routeAction:
        single:
          upstream:
            name: default-petstore-8080
            namespace: gloo-system
EOF
```

* Install Istio to the cluster:

```bash
 curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.5.0-beta.2 sh -
istio-1.5.0-beta.2/bin/istioctl manifest apply --set profile=demo
```

# Run locally

Have you context pointed to a cluster with istio installed.

Update and export these envrionment variaables:
```
WASM_IMAGE_REGISTRY=either webassemblyhub.io or yuval-test.solo.io
WASME_LOGIN_USERNAME=user name for above registry; we usually use the staging deployment for e2e tests
WASME_LOGIN_PASSWORD=the password for user above
TAGGED_VERSION=starts with a v. example: "vilackarms". This will be the tag for the operator docker image. if not running tests in parallel this can be a constant
IMAGE_REGISTRY=image registry to push operator image to; should be accessible to current cluster. example: gcr.io/<PROJECT_ID>
FILTER_IMAGE_TAG=the filter tag to try push,pull and deploy. example: yuval-test.solo.io/ilackarms/test-image:v0.0.1
IMAGE_PUSH=1
```

Example:
```
export WASM_IMAGE_REGISTRY=yuval-test.solo.io
export WASME_LOGIN_USERNAME=yuval
export WASME_LOGIN_PASSWORD=yuval
export TAGGED_VERSION=vyuval
export IMAGE_REGISTRY=gcr.io/myproject
export FILTER_IMAGE_TAG=yuval-test.solo.io/ilackarms/test-image:v0.0.1
export IMAGE_PUSH=1
```

Then run:
```
make run-tests
```

Or to run a specific test
```
make run-tests TEST_PKG=test/e2e/operator
```