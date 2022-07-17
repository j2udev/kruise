# Kruise(ctl)

[![Build Status](https://github.com/j2udevelopment/kruise/workflows/build/badge.svg?branch=master)](https://github.com/j2udevelopment/kruise/actions?query=workflow%3Abuild+branch%3Amaster)
[![GoReportCard](https://goreportcard.com/badge/github.com/j2udevelopment/kruise)](https://goreportcard.com/report/github.com/j2udevelopment/kruise)
[![Go Reference](https://pkg.go.dev/badge/github.com/j2udevelopment/kruise.svg)](https://pkg.go.dev/github.com/j2udevelopment/kruise)

Kruise is still very much a work in progress. It is essentially a
[black box](https://en.wikipedia.org/wiki/Black_box) CLI. It has a set of core
commands, whose options are dynamically configurable through a
[kruise manifest file](examples/kruise.yaml). Because the CLI is driven by a
config file, users can change kruise's behavior at runtime.

Consider you are consulting for two separate product teams. Team A deploys
MongoDB and Kafka Helm charts, while team B deploys MongoDB (at a different
version) and Jaeger. Each team can define their own configuration to be applied
to kruise.

Team A:

```yaml
deploy:
    helm:
        - option:
              arguments: "kafka"
              description: "Deploys Kafka to your Kubernetes cluster"
          chart:
              repository:
                  name: bitnami
                  url: https://charts.bitnami.com/bitnami
                  private: false
              releaseName: kafka
              chartPath: bitnami/kafka
              namespace: kafka
              version: 11.8.9
              values:
                  - /path/to/kafka/values.yaml
        - option:
              arguments: "mongodb, mongo"
              description: "Deploys MongoDB to your Kubernetes cluster"
          chart:
              repository:
                  name: bitnami
                  url: https://charts.bitnami.com/bitnami
                  private: false
              releaseName: mongodb
              chartPath: bitnami/mongodb
              namespace: mongodb
              version: 7.14.8
              values:
                  - /path/to/mongodb/values.yaml
delete:
    helm:
        - option:
              arguments: "kafka"
              description: "Deletes Kafka from your Kubernetes cluster"
          chart:
              releaseName: kafka
              chartPath: bitnami/kafka
              namespace: kafka
        - option:
              arguments: "mongodb, mongo"
              description: "Deletes MongoDB from your Kubernetes cluster"
          chart:
              releaseName: mongodb
              chartPath: bitnami/mongodb
              namespace: mongodb
```

Team B:

```yaml
deploy:
    helm:
        - option:
              arguments: "jaeger"
              description: "Deploys Jaeger to your Kubernetes cluster"
          chart:
              repository:
                  name: jaegertracing
                  url: https://jaegertracing.github.io/helm-charts
                  private: false
              releaseName: jaeger
              chartPath: jaegertracing/jaeger
              namespace: observability
              version: 0.39.5
              values:
                  - /path/to/jaeger/values.yaml
        - option:
              arguments: "mongodb, mongo"
              description: "Deploys MongoDB to your Kubernetes cluster"
          chart:
              repository:
                  name: bitnami
                  url: https://charts.bitnami.com/bitnami
                  private: false
              releaseName: mongodb
              chartPath: bitnami/mongodb
              namespace: mongodb
              version: 10.3.1
              values:
                  - /path/to/mongodb/values.yaml
delete:
    helm:
        - option:
              arguments: "jaeger"
              description: "Deletes Jaeger from your Kubernetes cluster"
          chart:
              releaseName: jaeger
              chartPath: jaegertracing/jaeger
              namespace: observability
        - option:
              arguments: "mongodb, mongo"
              description: "Deletes MongoDB from your Kubernetes cluster"
          chart:
              releaseName: mongodb
              chartPath: bitnami/mongodb
              namespace: mongodb
```

Now you, as the unfortunate soul trying to juggle 2 different projects, can
simply choose which team's config file to apply to kruise.

```zsh
kruise deploy mongodb kafka --config ~/.team-a-config.yaml
```

or

```zsh
kruise deploy mongodb jaeger --config ~/.team-b-config.yaml
```

## Installation

For now, install from source. Assuming your `$GOROOT`, `$GOBIN`, and `$GOPATH`
are set up, just:

```zsh
make install
```
