package cache

import (
	"context"
	"fmt"
	"github.com/golang/groupcache/lru"
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

			manager.lock.Lock()
			manager.trash = append(manager.trash, fileName)
			manager.lock.Unlock()
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

func (m *LocalCacheManager) DropDelFile(wg *sync.WaitGroup, ctx context.Context) {
	//定时，定速将需要删除的文件写给job
	timer := time.NewTimer(2 * time.Second)
	for {
		select {
		case <-ctx.Done():
			//程序退出时，需要处理再次check遗留的数据
			fmt.Println("Done...........")
			m.Clear()
			fmt.Println("exit deal:", len(m.trash))
			convert.DelTrashFile(m.trash)
			wg.Done()
			return
		case <-timer.C:
			fmt.Println("timer Channel timeout!")

			now := time.Now().Unix()
			if len(manager.trash) >= manager.maxEntry/2 || now-manager.timestamp > manager.interval {

				//copy用的汇编，需要指定复制的长度才行，按照最小的长度复制。copier用的反射。性能上一个copy更好。
				//如果使用copy的话，需要将make放到lock锁中。否则并发高的情况下，len的长度不够，会丢失数据。
				//加锁，copy数组 解锁
				manager.lock.Lock()
				delFileArray := make([]string, len(manager.trash), len(manager.trash))
				//copier.Copy(&delFileArray, &manager.trash)
				copy(delFileArray, manager.trash)
				manager.trash = manager.trash[0:0]
				manager.lock.Unlock()

				manager.timestamp = now

				fmt.Println("start delFileArray operation, ", len(delFileArray))

				go func() {
					//del array file
					convert.DelTrashFile(delFileArray)
				}()
			}

			timer.Reset(time.Duration(m.interval) * time.Second)
		}

	}

	timer.Stop()
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

func (m *LocalCacheManager) RunDelJob(wg *sync.WaitGroup, ctx context.Context) {
	m.DropDelFile(wg, ctx)
}
