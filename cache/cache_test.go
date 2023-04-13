package cache

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestCacheLogic(t *testing.T) {
	manager := InitLocalCacheObject(20, 3)
	path := "/Users/rockey-lyy/ad-tencent/ytmp3api/musicsource/"
	fileLen := 100
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
	manager.Clear()

	fmt.Println("exec end")
	time.Sleep(100 * time.Second)
}
