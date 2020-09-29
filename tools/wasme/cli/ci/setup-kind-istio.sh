#!/bin/bash

# Set up a kind cluster with Istio installed

set -e

if [ "$1" == "cleanup" ]; then
  kind get clusters | grep wasme-$2 | while read -r r; do kind delete cluster --name "$r"; done
  exit 0
fi

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

kubectl apply -f https://raw.githubusercontent.com/solo-io/gloo/master/example/petstore/petstore.yaml

# manifest apply depreciated after Istio 1.5, use install for later versions of istioctl
if [[ "$ISTIO_VERSION" == *"1.5"* ]]; then
  istioctl manifest apply --set profile=minimal
else
  istioctl install --set profile=minimal
fi

kubectl -n istio-system rollout status deployment istiod

kubectl -n default rollout status deployment petstore

# creating it during the test doesn't work for istio,
# so create it here.
kubectl create namespace bookinfo
kubectl label namespace bookinfo istio-injection=enabled

# setup local registry
docker run -d -p 5000:5000 --name registry registry:2

echo setup success
