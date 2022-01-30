package service

import (
	"fmt"
	"sync"

	"github.com/garyburd/redigo/redis"
	"github.com/solenovex/web-tutor/util"
)

type Bloom struct {
	Conn      *redis.Conn
	Key       string
	HashFuncs []util.F //保存hash函数
}

var (
	MyBloom   *Bloom
	InitBloom sync.Once
)

func BloomInit() {
	InitBloom.Do(func() {
		conn := MyRedis.RedisConn()
		NewBloom(conn)
	})
}
func NewBloom(con *redis.Conn) {
	MyBloom = &Bloom{Conn: con, Key: "bloom", HashFuncs: util.NewFunc()}
}
func (b *Bloom) Add(str string) error {
	var err error
	fmt.Println(str)
	for _, f := range b.HashFuncs {
		offset := f(str)
		_, err := (*b.Conn).Do("setbit", b.Key, offset, 1)
		if err != nil {
			return err
		}
	}
	return err
}

func (b *Bloom) Exist(str string) bool {
	var a int64 = 1
	for _, f := range b.HashFuncs {
		offset := f(str)
		bitValue, _ := (*b.Conn).Do("getbit", b.Key, offset)
		if bitValue != a {
			return false
		}
	}
	return true
}
