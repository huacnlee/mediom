package app

import (
	"fmt"
	. "github.com/huacnlee/mediom/app/models"
	"github.com/huacnlee/timeago"
	"github.com/revel/revel"
	"html/template"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

var (
	shareSites     = []string{"twitter", "weibo", "facebook", "google_plus", "email"}
	shareSiteIcons = map[string]string{
		"twitter":     "twitter",
		"weibo":       "weibo",
		"facebook":    "facebook-square",
		"google_plus": "google-plus-square",
		"email":       "envelope",
	}
)

func init() {
	revel.TemplateFuncs["plus"] = func(a, b int) int {
		return a + b
	}

	revel.TemplateFuncs["join"] = func(args []string, split string) string {
		return strings.Join(args, split)
	}

	revel.TemplateFuncs["is_owner"] = func(u User, obj interface{}) bool {
		if u.IsAdmin() {
			return true
		}

		switch obj.(type) {
		case User:
			u1 := obj.(User)
			return u1.Id == u.Id
		case Topic:
			t := obj.(Topic)
			return u.Id == t.UserId
		case Reply:
			r := obj.(Reply)
			return u.Id == r.UserId
		}

		return false

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
			out = fmt.Sprintf("<a href='/%v' class='uname'>%v</a>", template.HTMLEscapeString(u.Login), template.HTMLEscapeString(u.Login))
		default:
			login := fmt.Sprintf("%v", obj)
			out = fmt.Sprintf(`<a href="/%v" class="uname">%v</a>`, template.HTMLEscapeString(login), template.HTMLEscapeString(login))

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

			out = fmt.Sprintf("<a href=\"/%v\" class=\"uname\"><img src=\"%v\" class=\"media-object avatar-%v\" /></a>", template.HTMLEscapeString(u.Login), u.GavatarURL(size), size)
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
			out = fmt.Sprintf("<a href='/topics/node/%v' class='node-name'>%v</a>", n.Id, template.HTMLEscapeString(n.Name))
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
			html += fmt.Sprintf(`<li class="previous"><a href="%s%spage=%d"><i class="fa fa-arrow-left" aria-hidden="true"></i> 上一页</a></li>`, pageInfo.Path, linkFlag, pageInfo.Page-1)
		} else {
			html += fmt.Sprintf(`<li class="previous disabled"><a href="#""><i class="fa fa-arrow-left" aria-hidden="true"></i> 上一页</a></li>`)
		}

		html += fmt.Sprintf(`<li class="info"><samp>%d</samp> / <samp>%d</samp></li>`, pageInfo.Page, pageInfo.TotalPages)

		if pageInfo.Page < pageInfo.TotalPages {
			html += fmt.Sprintf(`<li class="next"><a href="%s%spage=%d">下一页 <i class="fa fa-arrow-right" aria-hidden="true"></i></a></li>`, pageInfo.Path, linkFlag, pageInfo.Page+1)
		} else {
			html += fmt.Sprintf(`<li class="next disabled"><a href="#">下一页 <i class="fa fa-arrow-right" aria-hidden="true"></i></a></li>`)
		}
		html += "</ul>"

		return template.HTML(html)
	}

	revel.TemplateFuncs["watch_tag"] = func(t Topic, u User) interface{} {
		out := ""
		if t.NewRecord() {
			return out
		}
		out = fmt.Sprintf(`<a href="#" data-id="%v" class="watch" title="关注此话题，当有新回帖的时候会收到通知"><i class="fa fa-eye"></i> 关注</a>`, t.Id)

		if u.NewRecord() {
			return template.HTML(out)
		}

		if u.IsWatched(t) {
			out = fmt.Sprintf(`<a href="#" data-id="%v" class="watch followed" title="点击取消关注"><i class="fa fa-eye"></i> 已关注</a>`, t.Id)
		}

		return template.HTML(out)
	}

	revel.TemplateFuncs["star_tag"] = func(t Topic, u User) interface{} {
		out := ""
		if t.NewRecord() {
			return out
		}
		label := fmt.Sprintf("%v 人收藏", t.StarsCount)
		out = fmt.Sprintf(`<a href="#" data-id="%v" data-count="%v" class="star"><i class="fa fa-star-o"></i> %v</a>`, t.Id, t.StarsCount, label)

		if u.NewRecord() {
			return template.HTML(out)
		}

		if u.IsStared(t) {
			out = fmt.Sprintf(`<a href="#" data-id="%v" data-count="%v" class="star followed"><i class="fa fa-star"></i> %v</a>`, t.Id, t.StarsCount, label)
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

	revel.TemplateFuncs["active_class"] = func(a string, b string) string {
		if strings.EqualFold(a, b) {
			return " active "
		}

		return ""
	}

	revel.TemplateFuncs["node_list"] = func() interface{} {
		roots := FindAllNodeRoots()
		outs := []string{}
		subs := []string{}
		outs = append(outs, `<div class="row node-list">`)
		for _, root := range roots {
			subs = []string{
				`<div class="node media clearfix">`,
				fmt.Sprintf(`<label class="media-left">%v</label>`, root.Name),
				`<div class="nodes media-body">`,
			}
			for _, node := range root.Children {
				subs = append(subs, fmt.Sprintf(`<span class="name"><a href="/topics/node/%v">%v</a></span>`, node.Id, node.Name))
			}
			subs = append(subs, "</div></div>")

			outs = append(outs, strings.Join(subs, ""))
		}
		outs = append(outs, "</div>")
		return template.HTML(strings.Join(outs, ""))
	}

	revel.TemplateFuncs["select_tag"] = func(objs interface{}, nameKey, valueKey, formName string, defaultValue interface{}) interface{} {
		objsVal := reflect.ValueOf(objs)
		if objsVal.Kind() != reflect.Slice {
			fmt.Println("Give a bad params, objs need to be a Slice")
			return ""
		}

		outs := []string{}

		subs := []string{}
		var nameField reflect.Value
		var valueField reflect.Value

		defaultName := "请选择"

		for i := 0; i < objsVal.Len(); i++ {
			val := objsVal.Index(i)
			nameField = val.FieldByName(nameKey)
			valueField = val.FieldByName(valueKey)
			subs = append(subs, fmt.Sprintf(`
               <li data-id="%v"><a href="#">%v</a></li>
            `, valueField.Int(), nameField.String()))

			// check current name
			if strings.EqualFold(fmt.Sprintf("%v", valueField.Int()), fmt.Sprintf("%v", defaultValue)) {
				defaultName = nameField.String()
			}
		}

		outs = append(outs, `<div class="input-group-btn md-dropdown">`)
		outs = append(outs, fmt.Sprintf(`
        <button class="btn btn-default dropdown-toggle" type="button" data-toggle="dropdown" aria-expanded="false">
            <span data-bind="label">%v</span> <span class="caret"></span>
        </button>
        <input type="hidden" data-bind="value" value="%v" name="%v" />`,
			defaultName, defaultValue, formName))

		outs = append(outs, `<ul class="dropdown-menu" role="menu">`)
		outs = append(outs, strings.Join(subs, ""))
		outs = append(outs, `</ul>`)
		outs = append(outs, `</div>`)

		return template.HTML(strings.Join(outs, ""))
	}

	revel.TemplateFuncs["total"] = func(key string) interface{} {
		switch key {
		case "users":
			return UsersCountCached()
		case "topics":
			return TopicsCountCached()
		case "replies":
			return RepliesCountCached()
		}

		return nil
	}

	revel.TemplateFuncs["setting"] = func(key string) interface{} {
		return template.HTML(GetSetting(key))
	}

	revel.TemplateFuncs["random_tip"] = func() interface{} {
		tipText := GetSetting("tips")
		tips := strings.Split(tipText, "\n")
		return template.HTML(tips[rand.Intn(len(tips))])
	}

	revel.TemplateFuncs["share_button"] = func(title, url string) interface{} {
		results := []string{}
		results = append(results, fmt.Sprintf(`<div class="social-share-button" data-via="gochina" data-title="%v" data-url="%v">`, title, url))
		for _, siteName := range shareSites {
			link := fmt.Sprintf(`<a rel="nofollow" data-site="%v" href="#"><i class="fa fa-%v"></i></a>`, siteName, shareSiteIcons[siteName])
			results = append(results, link)
		}
		results = append(results, "</div>")
		results = append(results, `<a href="#" class="share-button"><i class="fa fa-share-square-o"></i> 转发</a>`)
		return template.HTML(strings.Join(results, ""))
	}
}
