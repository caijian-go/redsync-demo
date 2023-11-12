package redislock

import (
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"sync"
	"testing"
	"time"
)

func TestRedSyncGroup(t *testing.T) {
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     "localhost:6379",
		Password: "123456",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)

	//根据需求设置锁名称
	mutexName := "goods-1"
	var wg sync.WaitGroup
	wg.Add(2)

	// wg.Done something
	for i := 0; i < 2; i++ {

		// 2 go runtime
		go func(i int) {
			defer wg.Done()

			// NewMutex
			mutex := rs.NewMutex(mutexName)

			// get lock
			fmt.Printf("%d 开始获取锁", i)
			if err := mutex.Lock(); err != nil {
				fmt.Printf("%d 获取锁异常", i)
				panic(err)
			}
			fmt.Printf("%d 获取锁成功", i)

			// 业务处理
			for j := 0; j < 5; j++ {
				time.Sleep(1 * time.Second)
				fmt.Printf("%d end of %d\n", i, time.Now().Unix())
			}
			//var goods Goods
			//db.Where(Goods{ProductId: 1}).First(&goods)
			//result := db.Model(&Goods{}).Where("product_id=?", 1).Updates(Goods{Inventory: goods.Inventory - 1})
			//if result.RowsAffected == 0 {
			//	fmt.Println("更新失败")
			//}

			// unlock
			fmt.Printf("%d 开始释放锁", i)
			if ok, err := mutex.Unlock(); !ok || err != nil {
				panic(fmt.Sprintf("%d unlock failed", i))
			}
			fmt.Printf("%d 锁释放成功", i)

		}(i)

	}
	wg.Wait()
}

// go test -v -bench=. ./pkg/redislock -run TestRedSync
func TestRedSync(t *testing.T) {
	// 创建一个redis客户端连接
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     "localhost:6379",
		Password: "123456",
	})

	// 创建redsync的客户端连接池
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// 创建redsync实例
	rs := redsync.New(pool)

	// 通过相同的key值名获取同一个互斥锁.
	mutexname := "my-global-mutex"
	//创建基于key的互斥锁
	mutex := rs.NewMutex(mutexname, redsync.WithTimeoutFactor(1))

	// 对key进行
	if err := mutex.Lock(); err != nil {
		panic(err)
	}

	// 获取锁后的业务逻辑处理.
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		t.Logf("end of %d", time.Now().Unix())
	}

	// 释放互斥锁
	if ok, err := mutex.Unlock(); !ok || err != nil {
		panic("unlock failed")
	}
}
