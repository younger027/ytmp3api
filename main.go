package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"ytmp3api/api"
	"ytmp3api/cache"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	go cache.InitLocalCacheObject(50, 3).RunDelJob(wg, ctx)

	go api.RunGinServer()

	select {
	case s := <-c:
		fmt.Println("Got signal:", s)
		cancel()
	}

	wg.Wait()

}
