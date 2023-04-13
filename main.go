package main

import (
	"ytmp3api/api"
	"ytmp3api/cache"
)

func main() {
	cache.InitLocalCacheObject(50, 5)
	api.RunGinServer()
}
