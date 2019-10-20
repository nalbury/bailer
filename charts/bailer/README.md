bailer
======
A Helm chart for Bailer

Current chart version is `0.1.0`





## Chart Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| bailer.affinity | object | `{}` | Pod affinity spec for the bailer backend deployment |
| bailer.image.pullPolicy | string | `"Always"` | Pull policy to use for bailer container image |
| bailer.image.repository | string | `"quay.io/nalbury/bailer"` | Container image to use for bailer |
| bailer.image.tag | string | `nil` | Container image tag to use for bailer, defaults to the appVersion set in the chart's Chart.yaml (usually matches the chart version) |
| bailer.ingress.annotations | object | `{}` | A map of ingress annotations for bailer |
| bailer.ingress.enabled | bool | `false` | Enable/Disable an ingress definition for bailer |
| bailer.ingress.hosts | list | `[]` | A list of igress host definitions for bailer |
| bailer.ingress.tls | list | `[]` | A list of ingress tls secrets |
| bailer.nodeSelector | object | `{}` | A map of node selector labels for the bailer backend deployment |
| bailer.replicaCount | int | `1` | Replica's for the primary api backend deployment |
| bailer.resources | object | `{}` | A map of container resource limits and requests |
| bailer.service.annotations | object | `{}` | A map of service annotations for bailer |
| bailer.service.ports | list | `[80]` | List of TCP ports for the bailer service, the target port is hardcoded in the service definition to ensure it's connection to the backend api deployment |
| bailer.service.type | string | `"ClusterIP"` | service type for bailer |
| bailer.tolerations | list | `[]` | A list of node tolerations for the bailer backend deployment |
| bailerWorker.affinity | object | `{}` | Pod affinity spec for the bailer worker deployment |
| bailerWorker.fullnameOverride | string | `""` | Override the default naming for the bailer worker deployment |
| bailerWorker.nodeSelector | object | `{}` | A map of node selector labels for the bailer worker deployment |
| bailerWorker.replicaCount | int | `1` | Replica's for the bailer worker deployment |
| bailerWorker.resources | object | `{}` | A map of container resource limits and requests for the bailer worker contatiners |
| bailerWorker.tolerations | list | `[]` | A list of node tolerations for the bailer worker deployment |
| config | object | `{"bailers":[]}` | Bailer config yaml, see https://github.com/nalbury/bailer for more info. |
| faktory.affinity | object | `{}` | Pod affinity spec for the faktory deployment |
| faktory.dash.enabled | bool | `false` | Enable/Disable the faktory dashboard |
| faktory.dash.ingress.annotations | object | `{}` | A map of ingress annotations for the faktory dashboard |
| faktory.dash.ingress.enabled | bool | `false` | Enable/Disable an ingress definition for the faktory dashboard |
| faktory.dash.ingress.hosts | string | `nil` | A list of igress host definitions for the faktory dashboard |
| faktory.dash.ingress.tls | list | `[]` | A list of ingress tls secrets |
| faktory.dash.service.annotations | object | `{}` | A map of service annotations for the faktory fashbaord |
| faktory.dash.service.ports | list | `[80]` | List of TCP ports for the faktory dashboard service, the target port is hardcoded in the service definition to ensure it's connection to the backend api deployment |
| faktory.dash.service.type | string | `"ClusterIP"` | Service type for the faktory dashbaord |
| faktory.enabled | bool | `true` | Enables disables an in cluster faktory deployment |
| faktory.externalUrl | string | `""` | If faktory.enabled is false, use this to set an external Faktory URL for bailer, ignored if faktory is enabled |
| faktory.fullnameOverride | string | `""` | Override the default naming for the faktory deployment |
| faktory.image.pullPolicy | string | `"Always"` | Pull policy to use for faktory container image |
| faktory.image.repository | string | `"contribsys/faktory"` | Container image to use for faktory |
| faktory.image.tag | string | `"latest"` | Container image tag to use for faktory |
| faktory.nodeSelector | object | `{}` | A map of node selector labels for the faktory deployment |
| faktory.resources | object | `{}` | A map of container resource limits and requests |
| faktory.service.annotations | object | `{}` | A map of service annotations for faktory |
| faktory.service.port | int | `7419` | TCP port for the faktory service, the target port is hardcoded in the service definition to ensure it's connection to the faktory deployment |
| faktory.service.type | string | `"ClusterIP"` | service type for faktory |
| faktory.tolerations | list | `[]` | A list of node tolerations for the faktory deployment |
| fullnameOverride | string | `""` | Override the the name prepended to the installed kubernetes resources |
| nameOverride | string | `""` | Override the installed chart's name |
| rbac.enabled | bool | `true` | Enable/Disable creation of a cluster role, service account and rolebinding for bailer |
| rbac.serviceAccount | string | `""` | Service account for bailer to use, ignored if rbac is enabled |
