package jobs

import (
	"hanyny/app/library"

	"gopkg.in/robfig/cron.v2"
)

// CronRemoveCache ...
func CronRemoveCache() {
	c := cron.New()
	c.AddFunc("@daily", func() {
		library.CacheCleanAll()
	})
	c.Start()
}
