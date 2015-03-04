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
	revel.TemplateFuncs["plus"] = func(a, b int) int {
		return a + b
	}

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
		bytes := []byte(template.HTMLEscapeString(text))
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
			out = fmt.Sprintf("<a href='/u/%v' class='uname'>%v</a>", template.HTMLEscapeString(u.Login), template.HTMLEscapeString(u.Login))
		default:
			login := fmt.Sprintf("%v", obj)
			out = fmt.Sprintf(`<a href="/u/%v" class="uname">%v</a>`, template.HTMLEscapeString(login), template.HTMLEscapeString(login))

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

			out = fmt.Sprintf("<a href=\"/u/%v\" class=\"uname\"><img src=\"%v\" class=\"avatar-%v\" /></a>", template.HTMLEscapeString(u.Login), u.GavatarURL(size), size)
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
			out = fmt.Sprintf("<a href='/topics/n%v' class='node-name'>%v</a>", n.Id, template.HTMLEscapeString(n.Name))
		}

		return template.HTML(out)
	}

	revel.TemplateFuncs["paginate"] = func(pageInfo Pagination) interface{} {
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

	revel.TemplateFuncs["watch_tag"] = func(t Topic, u User) interface{} {
		out := ""
		if t.NewRecord() {
			return out
		}
		out = fmt.Sprintf(`<a href="/topics/%v/watch" data-method="post" title="关注此话题，当有新回帖的时候会收到通知"><i class="fa fa-eye"></i> 关注</a>`, t.Id)

		if u.NewRecord() {
			return template.HTML(out)
		}

		if u.IsWatched(t) {
			out = fmt.Sprintf(`<a href="/topics/%v/unwatch" data-method="post" class="followed" title="点击取消关注"><i class="fa fa-eye"></i> 已关注</a>`, t.Id)
		}

		return template.HTML(out)
	}

	revel.TemplateFuncs["star_tag"] = func(t Topic, u User) interface{} {
		out := ""
		if t.NewRecord() {
			return out
		}
		label := fmt.Sprintf("%v 人收藏", t.StarsCount)
		out = fmt.Sprintf(`<a href="/topics/%v/star" data-method="post"><i class="fa fa-star-o"></i> %v</a>`, t.Id, label)

		if u.NewRecord() {
			return template.HTML(out)
		}

		if u.IsStared(t) {
			out = fmt.Sprintf(`<a href="/topics/%v/unstar" data-method="post" class="followed"><i class="fa fa-star"></i> %v</a>`, t.Id, label)
		}

		return template.HTML(out)
	}

	revel.TemplateFuncs["awesome_icon_tag"] = func(t Topic) interface{} {
		out := ""
		if !t.IsAwesome() {
			return out
		}

		out = `<i class="fa fa-diamond awesome" title="精华帖标记"></i>`
		return template.HTML(out)
	}
}
