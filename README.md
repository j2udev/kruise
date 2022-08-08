# Kruise(ctl)

[![Build Status](https://github.com/j2udevelopment/kruise/workflows/build/badge.svg?branch=master)](https://github.com/j2udevelopment/kruise/actions?query=workflow%3Abuild+branch%3Amaster)
[![GoReportCard](https://goreportcard.com/badge/github.com/j2udevelopment/kruise)](https://goreportcard.com/report/github.com/j2udevelopment/kruise)
[![Go Reference](https://pkg.go.dev/badge/github.com/j2udevelopment/kruise.svg)](https://pkg.go.dev/github.com/j2udevelopment/kruise)

Kruise is a [black box](https://en.wikipedia.org/wiki/Black_box) CLI that's
meant to simplify the deployment of... things.

> Things?

Yeah. Things. You know, Kubernetes things. Like Helm charts and k8s manifests.

> Aren't there already lots of tools like this?

Well, "lots" is subjective, but I suppose you're not wrong. Indeed there are
other tools that do something similar. Wonderful tools like
[Tilt](https://tilt.dev/) and [Skaffold](https://skaffold.dev/) exist that aim
to streamline the devloop when working with Kubernetes.

Kruise is different in that instead of running a command that executes what a
manifest specifies (like a
[skaffold.yaml](https://skaffold.dev/docs/references/yaml/)), it gives the user
the choice to execute individual/multiple things defined by the manifest.

If a `skaffold.yaml` specifies that it _should_ deploy x, y, and z, running
`skaffold deploy` will do just that. If you want to do some subset of that
deployment, you'd need to make a separate `skaffold.yaml` (or use their
[profiles](https://skaffold.dev/docs/environment/profiles/) feature which I've
personally used to _some_ success in the past, but have also encountered quite a
few
[issues](https://github.com/GoogleContainerTools/skaffold/issues?q=is%3Aissue+profiles)
over their many releases). If a `kruise.yaml` specifies that it _can_ deploy x,
y, and z, it will present those options to the user when they execute the deploy
command. This means you can `kruise deploy x z` or `kruise deploy y` or
`kruise deploy x z y` or any combination of those options.

> Ok, so how do I work this thing?

Without a manifest to drive the CLI, Kruise is just an empty black box. So the
first thing you'll need to do is define a `kruise.yaml` file. In the
[examples folder](examples) you'll find, you guessed it, examples! The manifests
you'll find in the examples folder should help with crafting your own
`kruise.yaml`.

By default, Kruise will check three locations for the existence of a
`kruise.yaml`. In order of priority, first, it will check your current working
directory. Second it will check your XDG_CONFIG_HOME, which can be set with an
environment variable. If it is not set, you can refer to
[the xdg package](https://github.com/adrg/xdg/blob/master/README.md) which
explains what the different xdg paths default to based on operating system.
Finally, `$HOME` will be checked for the existence of a `.kruise.yaml` file. The
path to your config can be overridden with the `KRUISE_CONFIG` environment
variable. You can set this in your bashrc, zshrc, etc for a persistent override,
or you can set it inline for quickly targetting different configuration files.

```bash
# absolute path to custom config
KRUISE_CONFIG=/path/to/foo.yaml kruise deploy -h
# relative path to custom config
KRUISE_CONFIG=bar.yaml kruise deploy -h
```

> Whoa, slow down. How do I even install it?

Check out the [Installation](#installation) section below!

> How about a TLDR?

Get Go 1.18+, copy an example from the [examples](examples) folder as a starting
point, and start tinkering with Kruise right away!

```bash
cp -r examples/observability /somewhere/else
cd /somewhere/else/observability
kruise deploy --help
```

## Installation

Kruise is still very young so, for now, install from source. Kruise uses
generics so you'll need at least Go 1.18.

[For your convenience](https://go.dev/dl/)

Once you've got at least Go 1.18 installed:

```zsh
go build -o /somewhere/on/your/PATH/kruise
```

## Abstract Deployments

The notion of a "deployment" can take many forms and can have varying degrees of
complexity. As much as possible, Kruise tries to strike a balance between
flexibility and functionality. The implementation details of deployments are
contained in the manifest, but the CLI will display an option that abstracts
away the deployment details.

For example, let's consider a team that has decided to leverage Istio as their
ingress and service mesh solution. The team is composed of frontend devs,
backend devs, and a devops engineer. The devops engineer can determine what an
abstract Istio "deployment" looks like:

```yaml
apiVersion: v1alpha1
kind: Config
deploy:
    deployments:
        istio:
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
adding a Helm repository, installing three Helm charts, and kubectl applying a
Gateway custom resource.

```txt
╰─❯ kruise deploy -h
Using config file: /path/to/kruise.yaml
Usage:
  kruise deploy [flags] [options]

Aliases:
  deploy, dep

Available Options:
  istio   deploy Istio to your k8s cluster

Flags:
  -c, --concurrent        deploy the arguments concurrently (deploys in order based on the 'priority' of each deployment passed)
  -h, --help              help for deploy
  -i, --init              add Helm repositories and create Kubernetes secrets for the specified options
  -d, --shallow-dry-run   output the command being performed under the hood

Global Flags:
  -V, --verbosity string   specify the log level to be used (trace, debug, info, warn, error) (default "error")
```

## Transparent Abstractions

Abstraction is a double-edged sword. When you know what you're doing it's
synonymous with efficiency. When you don't, it's theft of a new skillset. In
either case it can become a crutch.

As much as possible, Kruise adds `shallow-dry-run` capability to each command
(this just means the command being performed under the hood is printed to
stdout).

In the previous example, a backend dev is asked by another team how they are
deploying Istio. Rather than having tap the devops engineer on the shoulder,
they can simply run:

```txt
╰─❯ kruise deploy istio -id
Using config file: /path/to/kruise.yaml
/usr/local/bin/helm repo add istio https://istio-release.storage.googleapis.com/charts --force-update
/usr/local/bin/helm repo update
/usr/local/bin/helm upgrade --install istio-base istio/base --namespace istio-system --version 1.14.1 -f values/istio-base-values.yaml --create-namespace
/usr/local/bin/helm upgrade --install istiod istio/istiod --namespace istio-system --version 1.14.1 -f values/istiod-values.yaml --create-namespace
/usr/local/bin/helm upgrade --install istio-ingressgateway istio/gateway --namespace istio-system --version 1.14.1 -f values/istio-gateway-values.yaml --set service.externalIPs[0]=<redacted> --create-namespace
/usr/local/bin/kubectl create namespace istio-system
/usr/local/bin/kubectl apply --namespace istio-system -f manifests/istio-gateway.yaml
```
