# Kruise

[![Build Status](https://github.com/j2udev/kruise/workflows/build/badge.svg?branch=master)](https://github.com/j2udev/kruise/actions?query=workflow%3Abuild+branch%3Amaster)
[![GoReportCard](https://goreportcard.com/badge/github.com/j2udev/kruise)](https://goreportcard.com/report/github.com/j2udev/kruise)
[![Go Reference](https://pkg.go.dev/badge/github.com/j2udev/kruise.svg)](https://pkg.go.dev/github.com/j2udev/kruise)

Kruise is a [black box](https://en.wikipedia.org/wiki/Black_box) CLI that's
meant to simplify the deployment of... things.

> Things?

Yeah. Things. You know, Kubernetes things. Like Helm charts/repositories and
Kubernetes manifests. It even has extra utilities for creating k8s secrets with
masked prompts as opposed to having sensitive information under source control
or in your shell history.

> Ok, so how do I work this thing?

Without a manifest to drive the CLI, Kruise is just an empty black box. So the
first thing you'll need to do is define a `kruise.json/toml/yaml` file. The
manifests you'll find in the [examples folder](examples) should help with
crafting your own kruise configuration file.

By default, Kruise will check the following locations for config, in this order:

-   `$PWD/kruise.json/toml/yaml`

-   `$XDG_CONFIG_HOME/kruise/kruise.json/toml/yaml`

-   `$PWD/.kruise.json/toml/yaml`

-   `$XDG_CONFIG_HOME/.kruise.json/toml/yaml`

-   `$HOME/.kruise.json/toml/yaml`

`$PWD` is your current working directory. If the `XDG_CONFIG_HOME` environment
variable is not set The default location on Mac is
`~/Library/Application Support`, on Unix is `~/.config`, and on Windows is
`LocalAppData` which falls back to `%LOCALAPPDATA%`. Finally, `$HOME` will be
checked for the existence of a `.kruise.json/toml/yaml` file.

The path to your config can be overridden with the `KRUISE_CONFIG` environment
variable. You can set this in your bashrc, zshrc, etc for a persistent override,
or you can set it inline for quickly targetting different configuration files.
You can also set this to the URL of a raw kruise manifest you have under source
control.

```sg
export KRUISE_CONFIG=https://raw.githubusercontent.com/j2udev/kruise/master/examples/observability/kruise.yaml
kruise deploy -h
```

```sg
# absolute path to custom config
KRUISE_CONFIG=/path/to/foo.yaml kruise deploy -h
# relative path to custom config
KRUISE_CONFIG=bar.yaml kruise deploy -h
```

> Whoa, slow down. How do I even install it?

Check out the [Installation](#installation) section below!

> How about a TLDR?

Install Kruise and use the examples as inspiration for crafting your own k8s
deployment CLI.

```sg
cp -r examples/observability /somewhere/else/observability
cd /somewhere/else/observability
kruise deploy --help
```

or

```sh
export KRUISE_CONFIG=https://raw.githubusercontent.com/j2udev/kruise/master/examples/observability/kruise.yaml
kruise deploy -h
```

## Installation

Kruise is still very young so, for now, install from source. You'll need at
least Go 1.18.

[For your convenience](https://go.dev/dl/)

Once you've got at least Go 1.18 installed,
`go build -o /somewhere/on/your/PATH/kruise`:

```zsh
go build -o /usr/local/bin/kruise
```

## Abstract Deployments

The notion of a "deployment" can take many forms and can have varying degrees of
complexity. As much as possible, Kruise tries to strike a balance between
flexibility and functionality. The implementation details of deployments are
contained in the manifest, but the CLI will display an option that abstracts
away the deployment details.

For example, let's consider a team that has decided to leverage
[Istio](https://istio.io/) as their ingress and service mesh solution. The team
is composed of frontend devs, backend devs, and a devops engineer. The devops
engineer can determine what an abstract Istio "deployment" looks like:

```yaml
apiVersion: v1alpha1
kind: Config
deploy:
    deployments:
        - name: istio
          description:
              deploy: deploy Istio to your k8s cluster
              delete: delete Istio from your k8s cluster
          kubectl:
              manifests:
                  - namespace: istio-system
                    priority: 2
                    paths:
                        - manifests/istio-gateway.yaml
          helm:
              repositories:
                  - name: istio
                    url: https://istio-release.storage.googleapis.com/charts
                    init: true
              charts:
                  - priority: 1
                    chartName: base
                    releaseName: istio-base
                    namespace: istio-system
                    chartPath: istio/base
                    version: 1.14.1
                    values:
                        - values/istio-base-values.yaml
                    installArgs:
                        - --create-namespace
                  - priority: 1
                    chartName: istiod
                    releaseName: istiod
                    namespace: istio-system
                    chartPath: istio/istiod
                    version: 1.14.1
                    values:
                        - values/istiod-values.yaml
                    installArgs:
                        - --create-namespace
                  - priority: 2
                    chartName: gateway
                    releaseName: istio-ingressgateway
                    namespace: istio-system
                    chartPath: istio/gateway
                    version: 1.14.1
                    values:
                        - values/istio-gateway-values.yaml
                    setValues:
                        - service.externalIPs[0]=CHANGE_ME # Set this to your node InternalIP (minikube ip, or describe your node if you're not using minikube)
                    installArgs:
                        - --create-namespace
```

and the front and backend devs just need to know that they have an `istio`
option to deploy to their cluster without needing to know that it requires
adding a Helm repository, installing three Helm charts in a particular order,
and kubectl applying a Gateway custom resource.

```text
╰─❯ kruise deploy -h
INFO Using config file: /path/to/kruise.yaml
Usage:
  kruise deploy [flags] [options]

Aliases:
  deploy, dep

Options:
  istio   deploy Istio to your k8s cluster

Flags:
  -c, --concurrent        deploy the arguments concurrently (deploys in order based on the 'priority' of each deployment passed)
  -h, --help              help for deploy
  -i, --init              add Helm repositories and create Kubernetes secrets for the specified options
  -d, --dry-run           output the command being performed under the hood

Global Flags:
  -V, --verbosity string   specify the log level to be used (debug, info, warn, error) (default "warn")
```

## Transparent Abstractions

Abstraction is a double-edged sword. When you know what you're doing it's
synonymous with efficiency. When you don't, it can become a crutch.

As much as possible, Kruise adds `dry-run` capability to each command (this just
means the command being performed under the hood is printed to stdout).

In the previous example, a backend dev is asked by another team how they are
deploying Istio. Rather than having to tap the devops engineer on the shoulder,
they can simply run:

```txt
╰─❯ kruise deploy istio -id
INFO Using config file: /path/to/kruise.yaml
/usr/local/bin/helm repo add istio https://istio-release.storage.googleapis.com/charts --force-update
/usr/local/bin/helm repo update
/usr/local/bin/helm upgrade --install istio-base istio/base --namespace istio-system --version 1.14.1 -f values/istio-base-values.yaml --create-namespace
/usr/local/bin/helm upgrade --install istiod istio/istiod --namespace istio-system --version 1.14.1 -f values/istiod-values.yaml --create-namespace
/usr/local/bin/helm upgrade --install istio-ingressgateway istio/gateway --namespace istio-system --version 1.14.1 -f values/istio-gateway-values.yaml --set service.externalIPs[0]=CHANGE_ME --create-namespace
/usr/local/bin/kubectl create namespace istio-system
/usr/local/bin/kubectl apply --namespace istio-system -f manifests/istio-gateway.yaml
```

## Deployment Initialization

Preparing a new deployment can often require a few initialization steps. Whether
that's creating a secret, adding a Helm repository, etc. Often these steps don't
need to be executed in subsequent deployments so we wouldn't want to execute
them everytime. Instead, we can opt into them with the `--init` flag. If the
`init` property is added to an item that makes up a deployment in the
`kruise.yaml`, that item will only be deployed if the `--init` flag is passed on
the command line.

Some common items that make sense to be deployed only during initialization
might be Helm repositories and Kubectl secrets or namespaces. Kruise will prompt
you for any sensitive information associated with secrets and private Helm
repositories. This can help keep credentials out of source control and out of
your shell history! Due to the prompt, secrets and Helm repositories won't be
deployed concurrently if the `--concurrent` flag is used. Check out the
[secrets example](examples/secrets/kruise.yaml).

## Priority Deployments

Kruise can execute batches of deployments in parallel, at the cost of more
complexity. If the potential speed up of concurrent execution is not worth the
complexity, don't use the `--concurrent` flag and Kruise will deploy the options
passed syncronously, in the order they were given. If you have a set of simple
deployments in which deployment A does not depend on deployment B in any way,
you don't need to worry about priorities at all! Just execute kruise with the
`--concurrent` flag and your deployments will be executed in parallel; however,
let's say deployment A needs to apply a CustomResource that's defined in
deployment B. In this case, you would want to specify in the `kruise.yaml` that
the deployment responsible for creating the CustomResourceDefinition in B has a
higher priority than the application of that CustomResource in A. This can be
achieved with the `priority` field, which can be applied to any deployment item
(other than Kubectl secrets or Helm repositories) in the `kruise.yaml`.

Let's revisit the Istio example from earlier:

```yaml
apiVersion: v1alpha3
kind: Config
deploy:
    deployments:
        - name: istio
          description:
              deploy: deploy Istio to your k8s cluster
              delete: delete Istio from your k8s cluster
          kubectl:
              manifests:
                  - namespace: istio-system
                    priority: 2
                    paths:
                        - manifests/istio-gateway.yaml
          helm:
              repositories:
                  - name: istio
                    url: https://istio-release.storage.googleapis.com/charts
                    init: true
              charts:
                  - priority: 1
                    chartName: base
                    releaseName: istio-base
                    namespace: istio-system
                    chartPath: istio/base
                    version: 1.14.1
                    values:
                        - values/istio-base-values.yaml
                    installArgs:
                        - --create-namespace
                  - priority: 1
                    chartName: istiod
                    releaseName: istiod
                    namespace: istio-system
                    chartPath: istio/istiod
                    version: 1.14.1
                    values:
                        - values/istiod-values.yaml
                    installArgs:
                        - --create-namespace
                  - priority: 2
                    chartName: gateway
                    releaseName: istio-ingressgateway
                    namespace: istio-system
                    chartPath: istio/gateway
                    version: 1.14.1
                    values:
                        - values/istio-gateway-values.yaml
                    setValues:
                        - service.externalIPs[0]=CHANGE_ME # Set this to your node InternalIP (minikube ip, or describe your node if you're not using minikube)
                    installArgs:
                        - --create-namespace
```

in this example, you can see that the `istio-ingressgateway` is prioritized
_after_ `istio-base` and `istiod`. If the priorities were all weighted the same
(and you use the `--concurrent` flag), you'll find that the
`istio-ingressgateway` gets stuck in an `ImagePullBackOff`.

```txt
╰─❯ kubectl get pods -n istio-system
NAME                                    READY   STATUS             RESTARTS   AGE
istio-ingressgateway-8584fcc547-jg8rb   0/1     ImagePullBackOff   0          5m26s
istiod-766b78ccb7-2bjdx                 1/1     Running            0          5m26s
```

if we describe the pod, we can see in the events that k8s failed to pull the
`auto` image, which is a bit out of the ordinary.

```txt
Events:
  Type     Reason     Age                    From               Message
  ----     ------     ----                   ----               -------
  Normal   Scheduled  6m12s                  default-scheduler  Successfully assigned istio-system/istio-ingressgateway-8584fcc547-jg8rb to minikube
  Normal   Pulling    4m43s (x4 over 6m11s)  kubelet            Pulling image "auto"
  Warning  Failed     4m42s (x4 over 6m11s)  kubelet            Failed to pull image "auto": rpc error: code = Unknown desc = Error response from daemon: pull access denied for auto, repository does not exist or may require 'docker login': denied: requested access to the resource is denied
  Warning  Failed     4m42s (x4 over 6m11s)  kubelet            Error: ErrImagePull
  Warning  Failed     4m15s (x6 over 6m10s)  kubelet            Error: ImagePullBackOff
  Normal   BackOff    62s (x20 over 6m10s)   kubelet            Back-off pulling image "auto"
```

This seems like Istio magic at work... let's check
[ArtifactHub](https://artifacthub.io/packages/helm/istio-official/gateway?modal=template&template=deployment.yaml)

```yaml
containers:
    - name: istio-proxy
      # "auto" will be populated at runtime by the mutating webhook. See https://istio.io/latest/docs/setup/additional-setup/sidecar-injection/#customizing-injection
      image: auto
```

Ah so the image is set by the mutating webhook... but apparently deploying the
`istio-ingressgateway` in parallel with the Istio control plane (`istiod`)
results in the webhook failing to modify the `istio-proxy` image. If we bump the
priority of the `istio-ingressgateway` deployment, our problem should hopefully
go away.

```txt
╰─❯ kubectl get pods -n istio-system
NAME                                    READY   STATUS    RESTARTS   AGE
istio-ingressgateway-8584fcc547-rw882   1/1     Running   0          26s
istiod-766b78ccb7-cpnvd                 1/1     Running   0          34s
```

Much better, and we can see that mutating webhook did indeed change our image
when it was given the chance.

```txt
Events:
  Type     Reason     Age                 From               Message
  ----     ------     ----                ----               -------
  Normal   Scheduled  80s                 default-scheduler  Successfully assigned istio-system/istio-ingressgateway-8584fcc547-rw882 to minikube
  Normal   Pulling    80s                 kubelet            Pulling image "docker.io/istio/proxyv2:1.14.1"
  Normal   Pulled     79s                 kubelet            Successfully pulled image "docker.io/istio/proxyv2:1.14.1" in 523.8445ms
  Normal   Created    79s                 kubelet            Created container istio-proxy
  Normal   Started    79s                 kubelet            Started container istio-proxy
```

So this was a very long-winded and complicated explanation. Why?

To explain that just because you _can_ deploy things concurrently doesn't mean
you have to or even should. It introduces extra complexity and should be used
with caution. It will serve you better when deploying larger tech stacks in
which you are familiar with the interdepencies of the stack.

## Deployment Profiles

Kruise supports deployment profiles, which are essentially just bundles of other
generic deployments. The
[observability example](./examples/observability/kruise.yaml) defines multiple
example profiles as can be seen below (the deployments have been replaced with
`...` for readability):

```yaml
apiVersion: v1alpha1
kind: Config
deploy:
    profiles:
        - name: observability
          aliases:
              - telemetry
          description:
              deploy: "deploy an observability stack to the cluster"
              delete: "delete an observability stack from the cluster"
          items:
              - istio
              - jaeger
              - loki
              - prometheus-operator
        - name: logging
          description:
              deploy: "deploy a logging stack to the cluster"
              delete: "delete a logging stack from the cluster"
          items:
              - istio
              - loki
        - name: metrics
          description:
              deploy: "deploy a metrics stack to the cluster"
              delete: "delete a metrics stack from the cluster"
          items:
              - istio
              - prometheus-operator
    deployments:
        - name: istio ...
        - name: jaeger ...
        - name: loki ...
        - name: prometheus-operator ...
```

Each profile object has an `items` parameter. Each item in that list corresponds
to a name in the `deployments` list.

The help text for the above example is shown below:

```txt
╰─❯ kruise deploy -h
Using config file: /workspaces/kruise/examples/observability/kruise.yaml
Usage:
  kruise deploy [flags] [options]|[profiles]

Aliases:
  deploy, dep

Options:
  prometheus-operator, prom-op   deploy Prometheus Operator to your k8s cluster
  loki                           deploy Loki to your k8s cluster
  istio                          deploy Istio to your k8s cluster
  jaeger                         deploy Jaeger to your k8s cluster

Profiles:
  metrics,                   deploy a metrics stack to the cluster
    ↳ Options:               istio prometheus-operator
  logging,                   deploy a logging stack to the cluster
    ↳ Options:               istio loki
  observability, telemetry   deploy an observability stack to the cluster
    ↳ Options:               istio jaeger loki prometheus-operator

Flags:
  -c, --concurrent        deploy the arguments concurrently (deploys in order based on the 'priority' of each deployment passed)
  -h, --help              help for deploy
  -i, --init              add Helm repositories and create Kubernetes secrets for the specified options
  -d, --dry-run   output the command being performed under the hood

Global Flags:
  -V, --verbosity string   specify the log level to be used (trace, debug, info, warn, error) (default "error")
```

and the output of a dry-run for the `observability` profile is shown below:

```txt
╰─❯ kruise deploy -d observability
Using config file: /workspaces/kruise/examples/observability/kruise.yaml
/usr/local/bin/helm upgrade --install istio-base istio/base --namespace istio-system --version 1.14.1 -f values/istio-base-values.yaml --create-namespace
/usr/local/bin/helm upgrade --install istiod istio/istiod --namespace istio-system --version 1.14.1 -f values/istiod-values.yaml --create-namespace
/usr/local/bin/helm upgrade --install istio-ingressgateway istio/gateway --namespace istio-system --version 1.14.1 -f values/istio-gateway-values.yaml --set service.externalIPs[0]=CHANGE_ME --create-namespace
/usr/local/bin/kubectl create namespace istio-system
/usr/local/bin/kubectl apply --namespace istio-system -f manifests/istio-gateway.yaml
/usr/local/bin/helm upgrade --install jaeger jaegertracing/jaeger --namespace tracing --version 0.57.1 -f values/jaeger-values.yaml --create-namespace
/usr/local/bin/kubectl create namespace tracing
/usr/local/bin/kubectl apply --namespace tracing -f manifests/jaeger-virtual-service.yaml
/usr/local/bin/helm upgrade --install loki grafana/loki-stack --namespace logging --version 2.6.5 --create-namespace
/usr/local/bin/helm upgrade --install prometheus-operator prometheus-community/kube-prometheus-stack --namespace monitoring --version 36.0.2 -f values/prometheus-operator-values.yaml --create-namespace
/usr/local/bin/kubectl create namespace monitoring
/usr/local/bin/kubectl apply --namespace monitoring -f manifests/grafana-virtual-service.yaml
```

When a deployment or profile is passed without the `--concurrent` flag, order is
preserved. This means that individual deployments will be executed in the order
that they were given and profiles will execute deployments in the order that
they appear in the `items` list.

When the `--concurrent` flag is used (and the `--verbosity` is set at an
appropriate level), we can see the `observability` deployments executed
according to their priority.

```txt
╰─❯ kruise deploy observability -dci --verbosity info
Using config file: /workspaces/kruise/examples/observability/kruise.yaml
/usr/local/bin/helm repo add prometheus-community https://prometheus-community.github.io/helm-charts --force-update
/usr/local/bin/helm repo add istio https://istio-release.storage.googleapis.com/charts --force-update
/usr/local/bin/helm repo add jaegertracing https://jaegertracing.github.io/helm-charts --force-update
/usr/local/bin/helm repo add grafana https://grafana.github.io/helm-charts --force-update
/usr/local/bin/helm repo update
INFO Priority 1 waitgroup starting
/usr/local/bin/helm upgrade --install istiod istio/istiod --namespace istio-system --version 1.14.1 -f values/istiod-values.yaml --create-namespace
/usr/local/bin/helm upgrade --install istio-base istio/base --namespace istio-system --version 1.14.1 -f values/istio-base-values.yaml --create-namespace
INFO Priority 2 waitgroup starting
/usr/local/bin/kubectl create namespace monitoring
/usr/local/bin/kubectl apply --namespace monitoring -f manifests/grafana-virtual-service.yaml
/usr/local/bin/kubectl create namespace istio-system
/usr/local/bin/helm upgrade --install istio-ingressgateway istio/gateway --namespace istio-system --version 1.14.1 -f values/istio-gateway-values.yaml --set service.externalIPs[0]=CHANGE_ME --create-namespace
/usr/local/bin/kubectl apply --namespace istio-system -f manifests/istio-gateway.yaml
/usr/local/bin/helm upgrade --install prometheus-operator prometheus-community/kube-prometheus-stack --namespace monitoring --version 36.0.2 -f values/prometheus-operator-values.yaml --create-namespace
INFO Priority 3 waitgroup starting
/usr/local/bin/helm upgrade --install loki grafana/loki-stack --namespace logging --version 2.6.5 --create-namespace
/usr/local/bin/kubectl create namespace tracing
/usr/local/bin/kubectl apply --namespace tracing -f manifests/jaeger-virtual-service.yaml
/usr/local/bin/helm upgrade --install jaeger jaegertracing/jaeger --namespace tracing --version 0.57.1 -f values/jaeger-values.yaml --create-namespace
```
