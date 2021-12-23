package ezRedis

import (
	"context"
	"fmt"
	"github.com/Anveena/ezTools/ezPasswordEncoder"
	"github.com/go-redis/redis/v8"
	"runtime"
)

type Info struct {
	Host              string
	Port              string
	PasswordBase64Str string
	UserName          string
	DatabaseIndex     int
}

func NewRedisClient(redisInfo *Info) (*redis.Client, error) {
	if redisInfo.PasswordBase64Str == "" {
		return nil, fmt.Errorf("密码没配置")
	}
	password, err := ezPasswordEncoder.GetPasswordFromEncodedStr(redisInfo.PasswordBase64Str)
	if err != nil {
		return nil, fmt.Errorf("密码配的不合适,需要一个神秘的字符串才能解析,错误:\n\t%s", err.Error())
	}
	rs := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisInfo.Host, redisInfo.Port),
		Username: redisInfo.UserName,
		Password: password,
		DB:       redisInfo.DatabaseIndex,
		//There is also a function runtime.GOMAXPROCS, which reports (or sets) the user-specified number of cores that a Go program can have running simultaneously. It defaults to the value of runtime.NumCPU but can be overridden by setting the similarly named shell environment variable or by calling the function with a positive number. Calling it with zero just queries the value. Therefore if we want to honor the user's resource request, we should write
		//
		//var numCPU = runtime.GOMAXPROCS(0)
		MinIdleConns: runtime.GOMAXPROCS(0) * 8,
		PoolSize:     runtime.GOMAXPROCS(0) * 10,
	})
	return rs, rs.Ping(context.Background()).Err()
}
