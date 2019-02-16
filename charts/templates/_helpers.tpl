{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "name" -}}
{{- default .Values.service.name .Values.nameOverride | trunc 48 -}}
{{- end -}}

{{/*
Expand the service name.
*/}}
{{- define "servicename" -}}
{{- printf "%s" .Values.service.name | trunc 48 -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 48 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "fullname" -}}
{{- $name := default .Values.service.name .Values.nameOverride -}}
{{- printf "%s-%s" $name .Release.Namespace | trunc 48 -}}
{{- end -}}