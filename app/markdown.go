package app

import (
	"github.com/microcosm-cc/bluemonday"
	. "github.com/russross/blackfriday"
	"regexp"
)

var (
	blankRegexp, _        = regexp.Compile(`>\s+<`)
	mentionRegexp, _      = regexp.Compile(`@([\w\-\_]{3,20})`)
	mentionFloorRegexp, _ = regexp.Compile(`#([0-9]+)楼`)
)

func MarkdownGitHub(input []byte) []byte {
	htmlFlags := HTML_USE_XHTML

	renderer := HtmlRenderer(htmlFlags, "", "")

	// set up the parser
	extensions := 0 |
		EXTENSION_NO_INTRA_EMPHASIS |
		EXTENSION_FENCED_CODE |
		EXTENSION_AUTOLINK |
		EXTENSION_STRIKETHROUGH |
		EXTENSION_SPACE_HEADERS |
		EXTENSION_HEADER_IDS |
		EXTENSION_HARD_LINE_BREAK |
		EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK

	unsafe := Markdown(input, renderer, extensions)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	html = LinkMentionUser(html)
	html = LinkMentionFloor(html)
	html = RemoveBlankChars(html)
	return html
}

func RemoveBlankChars(input []byte) []byte {
	return blankRegexp.ReplaceAll(input, []byte("><"))
}

func LinkMentionUser(input []byte) []byte {
	return mentionRegexp.ReplaceAll(input, []byte(`<a href="/$1" class="mention"><b>@</b>$1</a>`))
}

func LinkMentionFloor(input []byte) []byte {
	return mentionFloorRegexp.ReplaceAll(input, []byte(`<a href="#reply${1}" class="mention-floor" data-floor="$1">#${1}楼</a>`))
}
