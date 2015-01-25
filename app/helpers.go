package app

import (
	"fmt"
	"github.com/revel/revel"
	"html/template"
	"strings"
)

func init() {
	revel.TemplateFuncs["error_messages"] = func(errors map[string]*revel.ValidationError) interface{} {
		out := ""
		var parts []string
		if len(errors) == 0 {
			return out
		}

		parts = append(parts, "<div class=\"alert alert-block alert-warning\" role=\"alert\"><ul>")
		for _, err := range errors {
			parts = append(parts, fmt.Sprintf("<li>%s %s</li>", err.Key, template.HTMLEscaper(err.Message)))
		}

		parts = append(parts, "</ul></div>")
		out = strings.Join(parts, "")

		return template.HTML(out)
	}
}
