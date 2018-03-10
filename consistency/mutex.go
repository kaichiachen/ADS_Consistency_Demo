package consistency

import (
	"fmt"
	"time"
	. "redsync"

        "github.com/garyburd/redigo/redis"
)

var addrs = []string{
        "127.0.0.1:6666",
        "127.0.0.1:6667",
        "127.0.0.1:6668",
        //"127.0.0.1:6669",
        //"127.0.0.1:6670",
}
var mutex *Mutex
var pools []Pool

func newMockPools() []Pool {
        pools = []Pool{}
        for _, addr := range addrs {
                pools = append(pools, &redis.Pool{
                                MaxIdle:     3,
                                IdleTimeout: 240 * time.Second,
                                Dial: func() (redis.Conn, error) {
                                        return redis.Dial("unix", addr) //TODO. maybe TCP or unix is better?
                                },
                                TestOnBorrow: func(c redis.Conn, t time.Time) error {
                                        _, err := c.Do("PING")
                                        return err
                                },
                        })
        }
        return pools
}

func MutexInit() {
	pools = newMockPools();
        rs := New(pools)
        mutex = rs.NewMutex("test-redsync")
        if mutex != nil {
                //just for test
                fmt.Println("get mutex done!")
        }
}
