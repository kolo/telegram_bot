package main

import (
	"bytes"
	"fmt"
	"text/template"
)

var templates = template.New("")
var welcomeMsg = `
{{.}} is up!
`
var bugMsg = `
Bug ID: {{.ID}}
Summary: {{.Summary}
Assignee: {{.Assignee}}
Status: {{.Status}}
`

func init() {
	template.Must(templates.New("welcome").Parse(welcomeMsg))
	template.Must(templates.New("bug").Parse(bugMsg))
}

func renderTemplate(name string, data interface{}) (string, error) {
	t := templates.Lookup(name)
	if t == nil {
		return "", fmt.Errorf("template %s is not defined", name)
	}

	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		return "", err
	}

	return b.String(), nil
}
