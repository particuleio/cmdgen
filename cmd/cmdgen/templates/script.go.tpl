#!/bin/bash

{{ range .Items }}
  {{- range (split .Description) -}}
    {{- printf "# %s\n" . }}
  {{- end -}}
  {{- println .Cmd }}
{{ end }} 

