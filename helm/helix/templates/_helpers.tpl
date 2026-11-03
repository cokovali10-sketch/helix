{{ define "helix.labels" }}
app.kubernetes.io/name: {{ include "helix.fullname" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{ end }}

{{ define "helix.selectorLabels" }}
app.kubernetes.io/name: {{ include "helix.fullname" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{ end }}

{{ define "helix.fullname" }}
{{ if .Values.fullnameOverride }}-{{ .Values.fullnameOverride }}{{ else }}-{{ .Release.Name }}-{{ .Chart.Name }}{{ end }}
{{ end }}
