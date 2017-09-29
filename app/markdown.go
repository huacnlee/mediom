package app

import (
	"github.com/russross/blackfriday"
	"regexp"
)

var (
	blankRegexp, _        = regexp.Compile(`>\s+<`)
	mentionRegexp, _      = regexp.Compile(`@([\w\-\_]{3,20})`)
	mentionFloorRegexp, _ = regexp.Compile(`#([0-9]+)楼`)
)

func MarkdownGitHub(input []byte) []byte {
	htmlFlags := blackfriday.UseXHTML
	htmlFlags |= blackfriday.SkipHTML
	htmlFlags |= blackfriday.SkipLinks
	htmlFlags |= blackfriday.NofollowLinks
	htmlFlags |= blackfriday.CompletePage

	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: htmlFlags,
	})

	// set up the parser
	extensions := 0 |
		blackfriday.NoIntraEmphasis |
		blackfriday.Tables |
		blackfriday.FencedCode |
		blackfriday.Autolink |
		blackfriday.Strikethrough |
		blackfriday.SpaceHeadings |
		blackfriday.HardLineBreak |
		blackfriday.NoEmptyLineBeforeBlock

	html := blackfriday.Run(input, blackfriday.WithRenderer(renderer), blackfriday.WithExtensions(extensions))
	html = RemoveBlankChars(html)
	html = LinkMentionUser(html)
	html = LinkMentionFloor(html)
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
