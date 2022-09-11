package cache 

import (
    "context"
    "github.com/go-redis/redis/v8"
    "filetransfer-back/internal/config"
)

var rc *redis.Client

func isConnAvailable(ctx context.Context) bool {
    status := false
    if rc != nil{
        _, err := rc.Ping(ctx).Result()

        if err == nil {
            status = true
        }
    }
    return status
}

func RedisInstance(ctx context.Context) *redis.Client {
    if isConnAvailable(ctx) != true {
        rc = redis.NewClient(&redis.Options{
            Addr:     config.RedisIp() + ":" + config.RedisPort(),
            Password: config.RedisPass(), // no password set
            DB:       config.RedisDbi(),  // use default DB
        })
    }
    return rc
}


