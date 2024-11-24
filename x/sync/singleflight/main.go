package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

var (
	cache        = make(map[string]*User) // 模拟缓存
	mu           sync.RWMutex             // 保护缓存
	requestGroup singleflight.Group       // SingleFlight 实例
)

type User struct {
	Id    int64
	Name  string
	Email string
}

// GetUserFromDB 模拟从数据库获取数据
func GetUserFromDB(username string) *User {
	fmt.Printf("Querying DB for key: %s\n", username)

	time.Sleep(1 * time.Second) // 模拟耗时操作

	id, _ := strconv.Atoi(username[len(username)-3:])
	fakeUser := &User{
		Id:    int64(id),
		Name:  username,
		Email: username + "@jianghushinian.cn",
	}
	return fakeUser
}

// GetUser 获取数据，先从缓存读取，若没有命中，则从数据库查询
func GetUser(key string) *User {
	// 先尝试从缓存获取
	mu.RLock()
	val, ok := cache[key]
	mu.RUnlock()
	if ok {
		return val
	}

	fmt.Printf("User %s not in cache\n", key)

	// 缓存未命中，使用 SingleFlight 防止重复查询
	result, _, _ := requestGroup.Do(key, func() (interface{}, error) {
		// 模拟从数据库获取数据
		val := GetUserFromDB(key)

		// 存入缓存
		mu.Lock()
		cache[key] = val
		mu.Unlock()

		return val, nil
	})

	return result.(*User)
}

func main() {
	var wg sync.WaitGroup
	keys := []string{"user_123", "user_123", "user_456"}

	// 第一轮并发查询，缓存中还没有数据，使用 SingleFlight 减少 DB 查询
	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			fmt.Printf("Get user for key: %s -> %+v\n", k, GetUser(k))
		}(key)
	}

	time.Sleep(2 * time.Second)
	fmt.Println("===================================")

	// 第二轮并发查询，缓存中有数据，直接读取缓存，不会查询 DB
	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			fmt.Printf("Get user for key: %s -> %+v\n", k, GetUser(k))
		}(key)
	}

	wg.Wait()
}
