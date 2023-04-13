package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestCacheLogic(t *testing.T) {
	wg := &sync.WaitGroup{}
	manager := InitLocalCacheObject(20, 3)
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		manager.DropDelFile(ctx)
		wg.Done()
	}()

	path := "/Users/rockey-lyy/ad-tencent/ytmp3api/musicsource/"
	fileLen := 123
	//create file
	for i := 0; i < fileLen; i++ {
		fileName := path + "filename_" + strconv.Itoa(i)
		os.Create(fileName)
	}

	fmt.Println("create file down,file len:", fileLen)
	time.Sleep(2 * time.Second)

	for i := 0; i < fileLen; i++ {
		fileName := path + "filename_" + strconv.Itoa(i)
		manager.Add(fileName, "1")
		//fmt.Println("add op:", fileName)
		//time.Sleep(100 * time.Millisecond)
	}

	time.Sleep(5 * time.Second)
	fmt.Println("manager.Clear():")
	manager.Clear()

	time.Sleep(5 * time.Second)
	fmt.Println("exec end")
	cancel()
	wg.Wait()

}
