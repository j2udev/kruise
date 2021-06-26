# Kruise(ctl)
[![Build Status](https://github.com/j2udevelopment/kruise/workflows/build/badge.svg?branch=master)](https://github.com/j2udevelopment/kruise/actions?query=workflow%3Abuild+branch%3Amaster)
[![GoReportCard](https://goreportcard.com/badge/github.com/j2udevelopment/kruise)](https://goreportcard.com/report/github.com/j2udevelopment/kruise)
[![Go Reference](https://pkg.go.dev/badge/github.com/j2udevelopment/kruise.svg)](https://pkg.go.dev/github.com/j2udevelopment/kruise)

Kruise is a CLI that aims to streamline the local development experience. Modern
software development can involve an overwhelming number of tools and sometimes
it's just difficult to keep up with it all. You can think of kruise as a CLI
wrapper that abstracts (but doesn't hide) the finer details of using many other
CLIs that commonly make their way into a software engineers tool kit. It also
comes with robust configuration management which allows you to to tailor
commands to perform exactly the way you need them too! For example, consider you
are consulting for two separate product teams. Team A deploys MongoDB and Kafka
Helm charts one way, while team B deploys them another. Each team can define
their own configuration to be applied to kruise.

Team A:

```yaml
deploy:
  mongodb:
    releaseName: mongodb
    chartPath: bitnami/mongodb
    namespace: mongodb
    version: 7.14.8
    values:
      - /foo/mongodb/values.yaml
  kafka:
    releaseName: kafka
    chartPath: bitnami/kafka
    namespace: kafka
    version: 11.8.9
    values:
      - /foo/kafka/values.yaml
```

Team B:

```yaml
deploy:
  mongodb:
    namespace: database
    version: 10.20.0
    values:
      - /bar/values.yaml
  kafka:
    namespace: event-broker
    version: 12.20.0
    values:
      - /baz/values.yaml
```

Now you, as the unfortunate soul trying to juggle 2 different projects, can
simply choose which team's config file to apply to kruise.

```zsh
kruise deploy mongodb kafka --config ~/.team-a-config.yaml
```

or

```zsh
kruise deploy mongodb kafka --config ~/.team-b-config.yaml
```

## Installation

Assuming your `$GOROOT`, `$GOBIN`, and `$GOPATH` are set up, just:

```zsh
make install
```
