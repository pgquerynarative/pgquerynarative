{{/*
Common labels
*/}}
{{- define "pgquerynarrative.labels" -}}
app.kubernetes.io/name: pgquerynarrative
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
