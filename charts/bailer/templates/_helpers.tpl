{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "bailer.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "bailer.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "bailer.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "bailer.labels" -}}
chart: {{ include "bailer.chart" . }}
release: {{ .Release.Name }}
{{- if .Chart.AppVersion }}
version: {{ .Chart.AppVersion | quote }}
{{- end }}
managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Create worker name.
*/}}
{{- define "bailer.bailerWorker.fullname" -}}
{{- if .Values.bailerWorker.fullnameOverride }}
{{ .Values.bailerWorker.fullnameOverride }}
{{- else }}
{{- printf "%s-worker" ( include "bailer.fullname" . ) -}} 
{{- end }}
{{- end -}}

{{/*
Create faktory name.
*/}}
{{- define "bailer.faktory.fullname" -}}
{{- if .Values.faktory.fullnameOverride }}
{{ .Values.faktory.fullnameOverride }}
{{- else }}
{{- printf "%s-faktory" ( include "bailer.fullname" . ) -}} 
{{- end }}
{{- end -}}

{{/*
Create bailer service account name
*/}}
{{- define "bailer.serviceAccount" -}}
{{- if and ( .Values.rbac.serviceAccount ) ( eq .Values.rbac.enabled false ) }}
{{ .Values.rbac.serviceAccount }}
{{- else }}
{{- include "bailer.fullname" . }}
{{- end }}
{{- end -}}

{{/*
Create faktory url
*/}}
{{- define "bailer.faktory.url" -}}
{{- if eq .Values.faktory.enabled false }}
{{ .Values.faktory.externalUrl }}
{{- else }}
{{- printf "tcp://%s.%s.svc:%0.0f" ( include "bailer.faktory.fullname" . ) .Release.Namespace .Values.faktory.service.port -}}
{{- end }}
{{- end -}}
