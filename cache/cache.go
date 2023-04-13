package cache

import (
	"fmt"
	"github.com/golang/groupcache/lru"
	"github.com/jinzhu/copier"
	"sync"
	"time"
	convert "ytmp3api/covert"
)

var manager *LocalCacheManager
var once sync.Once

type LocalCacheManager struct {
	cache     *lru.Cache
	maxEntry  int
	trash     []string
	interval  int64
	timestamp int64
	lock      sync.Mutex
}

func InitLocalCacheObject(MaxEntry, interval int) *LocalCacheManager {
	once.Do(func() {
		cache := lru.New(MaxEntry)

		OnEvicted := func(key lru.Key, value interface{}) {
			fileName, ok := key.(string)
			if !ok {
				fmt.Println("OnEvicted key not string:", key)
				return
			}

			manager.trash = append(manager.trash, fileName)

			now := time.Now().Unix()
			if len(manager.trash) > manager.maxEntry/2 || now-manager.timestamp > manager.interval {
				delFileArray := make([]string, 0, manager.maxEntry)

				//加锁，copy数组 解锁
				manager.lock.Lock()
				copier.Copy(&delFileArray, &manager.trash)
				manager.trash = manager.trash[0:0]
				manager.timestamp = now
				manager.lock.Unlock()
				fmt.Println("start delFileArray operation, ", len(delFileArray))

				go func() {
					//del array file
					convert.DelTrashFile(delFileArray)
				}()
			}

		}

		cache.OnEvicted = OnEvicted
		lock := sync.Mutex{}
		manager = &LocalCacheManager{
			cache:     cache,
			maxEntry:  MaxEntry,
			trash:     make([]string, 0, 2*MaxEntry),
			timestamp: time.Now().Unix(),
			interval:  int64(interval),
			lock:      lock,
		}
	})

	return manager
}
func GetCacheManger() *LocalCacheManager {
	return manager
}

func (m *LocalCacheManager) DropDelFile() {
	////定时，定速将需要删除的文件写给job
	//timer := time.NewTimer(1 * time.Second)
	//for {
	//	select {
	//	case <-conn:
	//		if timer.Stop() {
	//			fmt.Println("timer.Stop()")
	//		}
	//	case <-timer.C: // timer 通道超时
	//		fmt.Println("timer Channel timeout!")
	//		timer.Reset(time.Duration(m.interval) * time.Second)
	//	}
	//
	//}
	//
	//timer.Stop()
}

func (m *LocalCacheManager) Add(key, value interface{}) {
	m.cache.Add(key, value)
}

func (m *LocalCacheManager) Get(key interface{}) (value interface{}, ok bool) {
	return m.cache.Get(key)
}

func (m *LocalCacheManager) Clear() {
	m.cache.Clear()
}
