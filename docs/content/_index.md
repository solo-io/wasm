---
title: "Introduction"
weight: 1
---

<h1 align="center">
    <img src="https://github.com/solo-io/autopilot/blob/master/docs/content/img/logo.png?raw=true" alt="Web Assembly Hub" width="260" height="242">
  <br>
  The Web Assembly Hub
</h1>

placeholder text

[**Installation**](https://docs.solo.io/web-assembly-hub/latest/installation/) &nbsp; |
&nbsp; [**Documentation**](https://docs.solo.io/web-assembly-hub/latest) &nbsp; |
&nbsp; [**Blog**](TODO LINK) &nbsp; |
&nbsp; [**Slack**](https://slack.solo.io) &nbsp; |
&nbsp; [**Twitter**](https://twitter.com/soloio_inc)

<center>
<img src="https://github.com/solo-io/autopilot/blob/master/docs/content/img/autopilot-workflow.png?raw=true" alt="Autopilot">
</center>

### How does it work?

Developers define an `autopilot.yaml` and `autopilot-operator.yaml` which specify the skeleton and configuration of an *Autopilot Operator*.

Autopilot makes use of these files to (re-)generate the project skeleton, build, deploy, and manage the lifecycle of the operator via the `ap` CLI.

Users place their API in a generated `spec.go` file, and business logic in generated `worker.go` files. Once these files have been modified, they will not be overwritten by `ap generate`.

### How is it different from SDKs like Operator Framework and Kubebuilder?

The [Operator Framework](https://github.com/operator-framework) and [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) are open-ended SDKs that take a far less opinionated approach to building Kubernetes software.

**Autopilot** provides a more opinionated control loop via a generated *scheduler* that implements the [Controller-Runtime Reconciler interface](https://github.com/kubernetes-sigs/controller-runtime/blob/master/pkg/reconcile/reconcile.go#L80), for which users write stateless Work functions for various states of their top-level CRD. State information is stored
 on the *status* of the CRD, promoting a stateless design for Autopilot operators.
 
**Autopilot** additionally provides primitives, generated code, and helper functions for interacting with a variery of service meshes. While Autopilot can be used to build operators that do not configure or monitor a mesh, much of *Autopilot*'s design has been oriented to facilitate easy integration with popular service meshes.

Finally, **Autopilot** favors simplicity over flexibility, though it is the intention of the project to support the vast majority of DevOps workflows built on top of Kubernetes+Service mesh.

### Getting Started

The [Getting Started Tutorial]({{< versioned_link_path fromRoot="/tutorial_code/getting_started_1">}}) provides the best entrypoint to begin understanding and using 
Autopilot.

### Next Steps
- Join us on our Slack channel: [https://slack.solo.io/](https://slack.solo.io/)
- Follow us on Twitter: [https://twitter.com/soloio_inc](https://twitter.com/soloio_inc)
- Check out the docs: [https://docs.solo.io/autopilot/latest](https://docs.solo.io/autopilot/latest)
- Contribute to the [Docs](https://github.com/solo-io/solo-docs)

### Thanks

**Autopilot** would not be possible without the valuable open-source work of projects in the community. 

Autopilot has leveraged inspiration and libraries from the following Kubernetes projects:

- [Flagger](https://flagger.app/) - a robust, feature-rich service mesh operator which deploys canaries. Flagger has helped pioneer the service mesh operator space.
- [Controller Runtime](https://github.com/kubernetes-sigs/controller-runtime) - Excellent libraries for building k8s controllers. Many of 
- [Operator Framework](https://github.com/operator-framework) - An SDK for building generalized k8s operators. The source of much inspiration for Autopilot.

### Roadmap
- Support for managing multiple (remote) clusters.
- GitOps integrations ootb
- Support opaque user config added in autopilot-operator.yaml
- validate method for project config
    - check operatorName is kube compliant
    - apiVerson, kind, phases are correct
    - customParameters
    - final phase with i/o
- Builder funcs for service mesh types (VirtualServices, Gateways, etc.)
- `ap undeploy` command (undeploy / delete all deployed resources)
    - includes label all resources for easy list/delete
- Expose garbage collection func to workers
    - rollback the phase when something ensure fails? (option in config)
- Support Operators with multiple top-level crds
- Language-agnostic gRPC Worker interface
- OpenAPI schema generation
- interactive cli
- automatic metrics for worker syncs
- automatic traces for worker syncs
- option to make workers persistent
