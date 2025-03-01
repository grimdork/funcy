package funcy

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strings"
)

const errorTemplate = `<!DOCTYPE html>
<html>

<head>
	<title>Template error</title>
</head>

<body>
	<script>
		setTimeout(function() {
			window.location.href = "/";
		}, 3000);
	</script>

	<h1>Template error</h1>
	<p>There was an error loading the template.</p>
	<p>You will be redirected to the index page.</p>
</body>

</html>
`

// LoadTemplate from S3 storage.
func (cl *Client) LoadTemplate(fn string, data any) []byte {
	base := filepath.Base(fn)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	input := cl.GetFile(fn)
	tpl, err := template.New(name).Parse(string(input))
	if err != nil {
		return []byte(errorTemplate)
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		return []byte(errorTemplate)
	}

	return buf.Bytes()
}
