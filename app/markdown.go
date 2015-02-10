package app

import (
	"github.com/microcosm-cc/bluemonday"
	. "github.com/russross/blackfriday"
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
	return html
}
