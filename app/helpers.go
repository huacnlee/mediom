package app

import (
	"fmt"
	"github.com/revel/revel"
	"html/template"
	"time"
	//"reflect"
	"github.com/huacnlee/timeago"
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
		switch obj.(type) {
		case User:
			u := obj.(User)
			if u.NewRecord() {
				return out
			}
			out = fmt.Sprintf("<a href='/u/%v' class='uname'>%v</a>", u.Login, u.Login)
		default:
			out = fmt.Sprintf(`<a href="/u/%v" class="uname">%v</a>`, obj, obj)

		}

		return template.HTML(out)
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

	revel.TemplateFuncs["node_name_tag"] = func(obj interface{}) interface{} {
		out := ""
		switch obj.(type) {
		case Node:
			n := obj.(Node)
			if n.NewRecord() {
				return out
			}
			out = fmt.Sprintf("<a href='/topics/n%v' class='node-name'>%v</a>", n.Id, n.Name)
		}

		return template.HTML(out)
	}

	revel.TemplateFuncs["paginate"] = func(pageInfo Pagination) interface{} {
		fmt.Println("--------- pageInfo:", pageInfo)
		if pageInfo.TotalPages < 2 {
			return ""
		}

		linkFlag := "?"

		if strings.ContainsAny(pageInfo.Path, "?") {
			linkFlag = "&"
		}

		html := `<ul class="pager">`
		if pageInfo.Page > 1 {
			html += fmt.Sprintf(`<li class="previous"><a href="%s%spage=%d"><span aria-hidden="true">&larr;</span> 上一页</a></li>`, pageInfo.Path, linkFlag, pageInfo.Page-1)
		} else {
			html += fmt.Sprintf(`<li class="previous disabled"><a href="%s%spage=1"><span aria-hidden="true">&larr;</span> 上一页</a></li>`, pageInfo.Path, linkFlag)
		}

		html += fmt.Sprintf(`<li class="number">%d/%d</li>`, pageInfo.Page, pageInfo.TotalPages)

		if pageInfo.Page < pageInfo.TotalPages {
			html += fmt.Sprintf(`<li class="next"><a href="%s%spage=%d">下一页 <span aria-hidden="true">&rarr;</span></a></li>`, pageInfo.Path, linkFlag, pageInfo.Page+1)
		} else {
			html += fmt.Sprintf(`<li class="next disabled"><a href="%s%spage=%s">下一页 <span aria-hidden="true">&rarr;</span></a></li>`, pageInfo.Path, linkFlag, pageInfo.TotalPages)
		}
		html += "</ul>"

		return template.HTML(html)

	}
}
