# Bailer
Bail out your leaky cluster.

Bailer is an event driven service that can run a kubernetes job in response to a Prometheus AlertManager payload. 

## Arch
At a high level, Bailer consists of the following components

### Prometheus AlertManager
Bailer does not attempt to handle any of the actual alerting logic (retries, grouping, inhibition, silencing etc.) when bailing you out, and assumes that you've configured AlertManager to your specifications. 

[More informtation on AlertManager can be found here.](https://prometheus.io/docs/alerting/alertmanager/)

[More inforomation on configuring AlertManger to send payloads to Bailer can be found in the **Configuration** section below](https://github.com/nalbury/bailer#alertmanager)

### Bailer
Bailer's primary app is a simple API that accepts POST requests on a `/alert` route. It uses [gin](https://github.com/gin-gonic/gin to handle the actual http server, [cobra](https://github.com/spf13/cobra) for creating the command line interface, and [viper](https://github.com/spf13/viper) for reading it's configuration file(described in more detail below). 

A bailer is a command + container that should be run when an alert matches certain labels. Once the app receives a payload that matches a configred bailer, it creates a task in the Faktory queue and returns `200`. 

[More information on configuring Bailer can be found in the **Configuration** section below](https://github.com/nalbury/bailer#bailer-2)

### Faktory
Bailer uses Faktory to manage it's async task processing. It uses Redis to provide a queue for background jobs in Bailer, allowing the api to receive additional alert payloads, while the worker handles bailing your kubernetes cluster out. Each Faktory job has set of arguments and are placed into queues for workers to fetch and execute. 

[More info on Faktory itself can be found here.](https://github.com/contribsys/faktory/wiki)

### Bailer Worker
The Bailer Worker deployment runs bailer with the `faktory` subcommand, which will retrieve tasks from the Faktory queue, instead of pushing them into it. Once a job is retrieved, it's associated config is used to create a kubernetes job with a backOffLimit of 0, meaning it will only run once. This ensures that  there is only one bailer job run for each alert that comes in.

## Configuration

### AlertManager
To Configure and start using bailer you first must have a working Prometheus Server and AlertManager instance. I highly recommend using the [prometheus operator helm chart](https://github.com/helm/charts/tree/master/stable/prometheus-operator) to manage this, as it will bring up prometheus, alertmanager and grafana all in one chart. Once you have Prometheus/AlertManager Running, you simply need to add a receiver for bailers API endpoint, and a route for any alerts you want to send.

For Example, to route the alert `PodCrashLooping` to bailer:

```
route:
  routes:
  - match:
      alertname: PodCrashLooping
    receiver: 'bailer'
receivers:
- name: 'bailer'
  webhook_configs:
  - url: http://bailer.bailer.svc/alert/
```

[Prometheus Alerting Rules Docs](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules)
[AlertManager Configuration Docs](https://prometheus.io/docs/alerting/configuration/#configuration-file)

### Bailer

Each Alert from alert manager includes the prometheus labels from the underlying query that triggered it. We use a config.yaml file to tell Bailer which alerts to take action on, and how to take action on them. For example the following config will run a kubernetes job that echos "grab a bucket and start bailing!!!" in a busybox container whenever an alert of `PodCrashLooping` fires with a `pod` label matching `^nginx-.\*`, and a namespace label matching `default`. 

```
bailers:
- alert: PodCrashLooping
  command:
  - /bin/bash
  - -c
  - echo "grab a bucket and start bailing!!!"
  container:
    image: busybox
    tag: latest
  labels:
    namespace: default
    pod: ^nginx-.*
  serviceAccountName: bailer-admin
  ttlSecondsAfterFinished: 900
```

In addtion to using the Prometheus labels for filtering the alerts Bailer will act on, it also makes them available as environment variables in the resulting kubernetes job, prepended with the prefix `ALERT_`. Using this we can do something a little more useful with Bailer than echoing a string. 

The config below will actually delete the pod which fired the `PodCrashLooping` alert, potentially resolving an issue without user intervention. (**NOTE** The config below assumes there is a preexisting ServiceAccount called `bailer-admin` with permissions to delete pods in the `default` namespace)

```
bailers:
- alert: PodCrashLooping
  command:
  - /bin/bash
  - -c
  - kubectl delete pod -n $ALERT_NAMESPACE $ALERT_POD
  container:
    image: bitnami/kubectl
    tag: latest
  labels:
    namespace: default
    pod: ^nginx-.*
  serviceAccountName: bailer-admin
  ttlSecondsAfterFinished: 900
```

The config file should be a list of bailers
```
bailers:
- ...
```

With the follwing options:

| Parameter | Description |
| --------- | ----------- |
| `alert` | The name of the AlertManager alert to run this Bailer for |
| `command` | The Bailer container command that will be run to bail you out |
| `container` |  A map of the container image and tag to use |
| `labels` | A map of prometheus labels that must be matched on an alert payload before taking action. Each key must be an exact match, each value will be compiled to regular expression before matching, and each key value pair in the config must be present in the alert payload. |
| `serviceAccountName` | The Kubernetes service account that will be used by this bailer. Allows the bailer container to authenticate against the cluster when bailing it out |
| `ttlSecondsAfterFinished` | A TTL (in seconds) for each kubernetes job |

All config options are required for each bailer. 

## Installation
The recommended install method is using the projects helm chart:

```
helm repo add bailer https://bailer-charts.storage.googleapis.com
helm repo update
helm upgrade --install bailer bailer/bailer --namespace bailer
```

The Bailer config described above can be managed using the charts values.yaml file. All documentation for the helm chart can be found [here](./charts/bailer/README.md)





