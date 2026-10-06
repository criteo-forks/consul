{{- if .ExactName }}
key_prefix "services-data/{{ .ExactName }}/" {
  policy = "write"
}

{{- end }}
key_prefix "services-data/{{ .Name }}" {
  policy = "write"
}
