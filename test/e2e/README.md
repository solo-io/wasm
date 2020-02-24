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
