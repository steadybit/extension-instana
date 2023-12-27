{{/* vim: set filetype=mustache: */}}

{{/*
Expand the name of the chart.
*/}}
{{- define "instana.secret.name" -}}
{{- default "steadybit-extension-instana" .Values.instana.existingSecret -}}
{{- end -}}
