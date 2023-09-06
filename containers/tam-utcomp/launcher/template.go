package main

import (
	"os"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"env": func(key string, defaults ...string) string {
		def := ""
		if len(defaults) > 0 {
			def = defaults[0]
		}

		value, found := os.LookupEnv(key)
		if !found {
			return def
		}

		return value
	},
}

func Evaluate(expr string) (string, error) {
	tpl := template.New("expr")
	tpl.Funcs(funcMap)

	tpl, err := tpl.Parse(expr)
	if err != nil {
		return "", err
	}

	var buf strings.Builder

	err = tpl.Execute(&buf, nil)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
