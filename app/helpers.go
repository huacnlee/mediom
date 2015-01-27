package app

import (
	"fmt"
	"github.com/revel/revel"
	"html/template"
	"strings"
)

func init() {
	revel.TemplateFuncs["error_messages"] = func(args ...interface{}) interface{} {
		out := ""
		if len(args) == 0 {
			return out
		}

		v := args[0].(revel.Validation)
		var parts []string
		if !v.HasErrors() {
			return out
		}

		parts = append(parts, "<div class=\"alert alert-block alert-warning\" role=\"alert\"><ul>")
		for _, err := range v.ErrorMap() {
			parts = append(parts, fmt.Sprintf("<li>%s %s</li>", err.Key, template.HTMLEscaper(err.Message)))
		}

		parts = append(parts, "</ul></div>")
		out = strings.Join(parts, "")

		return template.HTML(out)
	}
}
