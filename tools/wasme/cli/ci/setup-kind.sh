#!/bin/bash

set -e

if [ "$1" == "cleanup" ]; then
  kind get clusters | grep wasme-$2 | while read -r r; do kind delete cluster --name "$r"; done
  exit 0
fi

make clean

# generate names: $1 allows to make several envs in parallel
cluster=wasme-$1

# set up each cluster
(cat <<EOF | kind create cluster --name $cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
kubeadmConfigPatches:
- |
  kind: InitConfiguration
  nodeRegistration:
    kubeletExtraArgs:
      authorization-mode: "AlwaysAllow"
      feature-gates: "EphemeralContainers=true"
- |
  kind: KubeletConfiguration
  featureGates:
    EphemeralContainers: true
- |
  kind: KubeProxyConfiguration
  featureGates:
    EphemeralContainers: true
- |
  kind: ClusterConfiguration
  metadata:
    name: config
  apiServer:
    extraArgs:
      "feature-gates": "EphemeralContainers=true"
  scheduler:
    extraArgs:
      "feature-gates": "EphemeralContainers=true"
  controllerManager:
    extraArgs:
      "feature-gates": "EphemeralContainers=true"
EOF
)

printf "\n\n---\n"
echo "Finished setting up cluster $cluster"

# build once to fail script if fails
make wasme-image -B
# make all the docker images again - everything is cached so this should be fast
# grab the image names out of the `make docker` output
make wasme-image -B | sed -nE 's|Successfully tagged (.*$)|\1|p' | while read f; do kind load docker-image --name $cluster $f; done

istioctl manifest apply --set profile=minimal
kubectl create ns gloo-system-test
helm install --namespace gloo-system-test --set global.wasm.enabled=true gloo https://storage.googleapis.com/solo-public-helm/charts/gloo-1.5.0-beta11.tgz

kubectl apply -f https://raw.githubusercontent.com/solo-io/gloo/master/example/petstore/petstore.yaml
cat <<EOF | kubectl apply -f -
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: default
  namespace: gloo-system-test
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
            namespace: gloo-system-test
EOF

kubectl -n istio-system rollout status deployment istiod
kubectl -n gloo-system-test rollout status deployment gloo
kubectl -n gloo-system-test rollout status deployment gateway
kubectl -n gloo-system-test rollout status deployment discovery
kubectl -n gloo-system-test rollout status deployment gateway-proxy
kubectl -n default rollout status deployment petstore

# not sure why, creating it during the test doesn't work for istio,
# so create it here.
kubectl create namespace bookinfo
kubectl label namespace bookinfo istio-injection=enabled
# setup local registry
docker run -d -p 5000:5000 --name registry registry:2

echo setup success
