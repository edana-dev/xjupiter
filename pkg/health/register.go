package health

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
)

var (
	RegisterMap = make(map[string]Checker)
	lock        sync.RWMutex
)

type Checker func() (bool, map[string]interface{})

func Register(name string, check Checker) error {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := RegisterMap[name]; ok {
		return errors.New(fmt.Sprintf("health check already exists: %s", name))
	}
	RegisterMap[name] = check
	return nil
}

const RedisInfoSection = "server"

func RegisterRedis(redisClient redis.Cmdable) {
	Register("redis", func() (bool, map[string]interface{}) {
		info, err := redisClient.Info(RedisInfoSection).Result()
		attrs := make(map[string]interface{})
		if err != nil {
			attrs["err"] = err.Error()
			return false, attrs
		}
		attrs[RedisInfoSection] = info
		return true, attrs
	})
}

const CheckSQL = "select 1"

func RegisterDB(db *sql.DB) {
	Register("redis", func() (bool, map[string]interface{}) {
		_, err := db.Query(CheckSQL)
		attrs := make(map[string]interface{})
		attrs["sql"] = CheckSQL
		if err != nil {
			attrs["err"] = err.Error()
		}
		return err == nil, attrs
	})
}