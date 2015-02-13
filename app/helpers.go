package app

import (
	"fmt"
	"github.com/revel/revel"
	"html/template"
	"time"
	//"reflect"
	"github.com/huacnlee/timeago"
	"github.com/shaoshing/train"
	. "mediom/app/models"
	"strings"
)

func init() {
	revel.TemplateFuncs["error_messages"] = func(args ...interface{}) interface{} {
		out := ""

		if len(args) == 0 {
			return out
		}

		switch args[0].(type) {
		case string:
			return out
		case revel.Validation:
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
		default:
			return out
		}

		return template.HTML(out)
	}

	revel.TemplateFuncs["timeago"] = func(t time.Time) string {
		return timeago.Chinese.Format(t)
	}

	revel.TemplateFuncs["markdown"] = func(text string) interface{} {
		bytes := []byte(text)
		outBytes := MarkdownGitHub(bytes)
		htmlText := string(outBytes[:])
		return template.HTML(htmlText)
	}

	revel.TemplateFuncs["user_name_tag"] = func(obj interface{}) interface{} {
		out := "未知用户"
		if obj != nil {
			u := obj.(User)
			if u.NewRecord() {
				return out
			}
			out = fmt.Sprintf("<a href='/u/%v' class='uname'>%v</a>", u.Login, u.Login)
			return template.HTML(out)
		}

		return out
	}

	revel.TemplateFuncs["user_avatar_tag"] = func(obj interface{}, size string) interface{} {
		out := ""
		if obj != nil {
			u := obj.(User)
			if u.NewRecord() {
				return out
			}

			out = fmt.Sprintf("<a href=\"/u/%v\" class=\"uname\"><img src=\"%v\" class=\"avatar-%v\" /></a>", u.Login, u.GavatarURL(size), size)
		}

		return template.HTML(out)
	}

	revel.TemplateFuncs["javascript_tag"] = train.JavascriptTag
	revel.TemplateFuncs["stylesheet_tag"] = train.StylesheetTag
}
