package models

import (
	"fmt"
	"github.com/revel/revel/cache"
	"time"
)

type Setting struct {
	Id  int32
	Key string `sql:"not null"`
	Val string `sql:"type: text; not null"`
}

func (s Setting) AfterSave() {
	s.RewriteCache()
}

func settingCacheKey(key string) string {
	return fmt.Sprintf("setting/%v/v1", key)
}

func (s Setting) RewriteCache() {
	cache.Set(settingCacheKey(s.Key), s.Val, 7*24*time.Hour)
}

func FindSettingByKey(key string) (s Setting) {
	s.Key = key
	DB.Where("`key` = ?", key).First(&s)
	return s
}

func GetSetting(key string) (out string) {
	out = ""
	if err := cache.Get(settingCacheKey(key), &out); err != nil {
		s := FindSettingByKey(key)
		if s.Id <= 0 {
			db.Save(&s)
		}

		out = s.Val
		s.RewriteCache()
	}

	return
}
