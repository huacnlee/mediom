package app

import (
	"mediom/app"
	"strings"
	"testing"
)

func TestRemoveBlankChars(t *testing.T) {
	source := `foo bar
<b>foo</b>
<a>dar<A>`

	expect := `foo bar
<b>foo</b><a>dar<A>`

	out := string(app.RemoveBlankChars([]byte(source)))

	if !strings.Contains(out, expect) {
		t.Errorf("== expect\n%v \n== but\n%v", expect, out)
	}
}

func TestMarkdownGitHub(t *testing.T) {
	source := `# foo
**bar**

__dar__
`

	expect := `<h1>foo</h1><p><strong>bar</strong></p><p><strong>dar</strong></p>`
	out := string(app.MarkdownGitHub([]byte(source)))
	if !strings.Contains(out, expect) {
		t.Fatalf("\n== expect \n%v\n== but \n%v", expect, out)
	}
}
