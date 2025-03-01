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
	tpl, err := template.New(name).Parse(string(cl.GetFile(fn)))
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

// LoadTemplates from S3 storage.
func (cl *Client) LoadTemplates(files []string, data any) []byte {
	var buf bytes.Buffer
	for _, fn := range files {
		base := filepath.Base(fn)
		name := strings.TrimSuffix(base, filepath.Ext(base))
		tpl, err := template.New(name).Parse(string(cl.GetFile(fn)))
		if err != nil {
			return []byte(errorTemplate)
		}

		err = tpl.Execute(&buf, data)
		if err != nil {
			return []byte(errorTemplate)
		}
	}

	return buf.Bytes()
}
