{{- if .ExactName }}
service "{{ .ExactName }}" {
  policy = "write"
}

{{- end }}
service_prefix "{{ .Name }}" {
  policy = "write"
}
