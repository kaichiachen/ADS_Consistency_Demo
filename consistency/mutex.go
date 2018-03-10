package consistency

import (
	"fmt"
	"time"
	"math/rand"
	. "redsync"

        "github.com/garyburd/redigo/redis"
)

var addrs = [5]string{
        "127.0.0.1:6666",
        "127.0.0.1:6667",
        "127.0.0.1:6668",
        "127.0.0.1:6669",
        "127.0.0.1:6670",
}
var mutex *Mutex

func newMockPools() []Pool {
        pools := []Pool{}
        for i := 0; i < 5; i++ {
                pools = append(pools, &redis.Pool{
                                MaxIdle:     3,
                                IdleTimeout: 240 * time.Second,
                                Dial: func() (redis.Conn, error) {
                                        return redis.Dial("udp", addrs[i]) //TODO. maybe TCP or unix is better?
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
	pools := newMockPools();
        rs := New(pools)
        mutex = rs.NewMutex(fmt.Sprintf("test-redsync%d", rand.Intn(100000)))
        if mutex != nil {
                //just for test
                fmt.Println("get mutex done!")
        }
}
