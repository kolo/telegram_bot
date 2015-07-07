package main

import (
	"bytes"
	"strconv"
	"text/template"
)

var response = `
Bug ID: {{.ID}}
Summary: {{.Summary}}
Assignee: {{.Assignee}}
Status: {{.Status}}
`

var (
	cmdBug = &Command{
		Name: "bug",
		Run:  getBug,
	}

	responseTemplate = template.Must(template.New("response").Parse(response))
)

func getBug(cmd *Command, args string) (string, error) {
	id, err := strconv.Atoi(args)
	if err != nil {
		return "", err
	}

	bug, err := bzClient.GetBug(id)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	if err = responseTemplate.Execute(&b, bug); err != nil {
		return "", err
	}

	return b.String(), nil
}
