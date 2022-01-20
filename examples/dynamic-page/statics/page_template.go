package statics

const PageTemplate = `
{{- define "page" }}

<!DOCTYPE html>
<html>
    {{- template "header" . }}
<body>

<p>&nbsp;&nbsp;🚀 <a href="https://github.com/bygui86/go-csv-view/examples/dynamic-page"><b>Dynamic page example</b></a> <em>is a real-time Golang runtime viewer example</em></p>

<style> .box { justify-content:center; display:flex; flex-wrap:wrap } </style>
<div class="box"> {{- range .Charts }} {{ template "base" . }} {{- end }} </div>

</body>
</html>

{{ end }}
`
