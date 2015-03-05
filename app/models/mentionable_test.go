package models

import (
	"testing"
)

func TestSearchMentionLogins(t *testing.T) {
	body := `@huacnlee 你好啊 @monster @huacnlee`
	logins := searchMentionLogins(body)
	if logins[0] != "huacnlee" && logins[1] != "monster" {
		t.Error("not match result:", logins)
	}
}
