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
adding a Helm repository, installing three Helm charts in a particular order,
and kubectl applying a Gateway custom resource.

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

## Deployment Initialization

Preparing a new deployment can often require a few initialization steps. Whether
that's creating a secret, adding a Helm repository, etc. Often these steps don't
need to be executed in subsequent deployments and can be opted into with the
`--init` flag. Applying this flag will determine if any of the passed arguments
need to have k8s secrets created or Helm repositories added. Kruise will prompt
you for credentials if needed (if k8s secret(s) or `private` Helm repositories
are associated with deployment used).

The following manifest will result in the user being prompted for credentials to
authenticate against the specified container registry. It will also prompt the
user to enter credentials for the private Helm repository.

```yaml
apiVersion: v1alpha1
kind: Config
deploy:
    deployments:
        custom-deployment1:
            aliases:
                - cd1
            description:
                deploy: custom cd1 deployment description
                delete: custom cd1 deployment description
            kubectl:
                secrets:
                    - type: docker-registry
                      name: custom-image-pull-secret
                      namespace: custom
                      registry: private-registry.com
            helm:
                repositories:
                    - name: custom-helm-repo
                      url: https://custom.helm.repo
                      private: true # Setting private to true will prompt the user for credentials when the --init flag is used
```

...

```txt
╰─❯ kruise deploy cd1 -i
Using config file: /path/to/kruise.yaml
Please enter your username for the private-registry.com container registry: foo
Please enter your password for the private-registry.com container registry: ***
✔ Please confirm your password: █
```

...

```txt
╰─❯ kruise deploy cd1 -i
Using config file: /Users/wardjl/DEC/temp/custom/kruise.yaml
Please enter your username for the private-registry.com container registry: foo
Please enter your password for the private-registry.com container registry: ***
Please confirm your password: ***
secret/custom-image-pull-secret created
Please enter your username for the custom-helm-repo Helm repository: foo
Please enter your password for the custom-helm-repo Helm repository: ***
✔ Please confirm your password: █
```

...

```txt
╰─❯ kruise deploy cd1 -i
Using config file: /Users/wardjl/DEC/temp/custom/kruise.yaml
Please enter your username for the private-registry.com container registry: foo
Please enter your password for the private-registry.com container registry: ***
Please confirm your password: ***
secret/custom-image-pull-secret created
Please enter your username for the custom-helm-repo Helm repository: foo
Please enter your password for the custom-helm-repo Helm repository: ***
Please confirm your password: ***
Hang tight while we grab the latest from your chart repositories...
...Successfully got an update from the "custom-helm-repo" chart repository
Update Complete. ⎈Happy Helming!⎈
```

This can help keep credentials out of source control and out your shell history!

## Priority Deployments

Kruise can execute batches of deployments in parallel. If you have a set of
simple deployments in which deployment A does not depend on deployment B in any
way, you don't need to worry about priorities at all! Just execute kruise with
the `--concurrent` flag and your deployments will be executed in parallel;
however, let's say deployment A needs to apply a CustomResource that's defined
in deployment B. In this case, you would want to specify in the `kruise.yaml`
that the deployment responsible for creating the CustomResourceDefinition in B
has a higher priority than the application of that CustomResource in A. This can
be achieved with the `priority` field, which can be applied to each type of
installer in Kruise.

Let's revisit the Istio example from earlier:

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
