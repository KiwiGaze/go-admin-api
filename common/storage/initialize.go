package storage

import (
	"log"

	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/captcha"
)

func Setup() {
	setupCache()
	setupCaptcha()
	setupQueue()
}

func setupCache() {
	cacheAdapter, err := config.CacheConfig.Setup()
	if err != nil {
		log.Fatalf("cache setup error, %s\n", err.Error())
	}
	sdk.Runtime.SetCacheAdapter(cacheAdapter)
}

func setupCaptcha() {
	captcha.SetStore(captcha.NewCacheStore(sdk.Runtime.GetCacheAdapter(), 600))
}

func setupQueue() {
	if config.QueueConfig.Empty() {
		return
	}
	if q := sdk.Runtime.GetQueueAdapter(); q != nil {
		q.Shutdown()
	}
	queueAdapter, err := config.QueueConfig.Setup()
	if err != nil {
		log.Fatalf("queue setup error, %s\n", err.Error())
	}
	sdk.Runtime.SetQueueAdapter(queueAdapter)
	go queueAdapter.Run()
}
