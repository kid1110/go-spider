package service

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

type CacheStore interface {
	RedisGet(key string, value interface{}) error
	RedisSet(key string, value interface{}, expires time.Duration) error
	RedisAdd(key string, value interface{}, expires time.Duration) error
	RedisDelete(key string) error
	RedisReplace(key string, value interface{}, expire time.Duration) error
}
type RedisStore struct {
	pool              *redis.Pool
	defaultExpiration time.Duration
}

var (
	MyRedis       *RedisStore
	initRedisOnce sync.Once
)

func RedisInit() {
	initRedisOnce.Do(func() {
		ConnectRedis(3 * time.Second)
	})
}
func ConnectRedis(defaultExpiration time.Duration) {
	pool := &redis.Pool{
		MaxActive: 512,
		MaxIdle:   10,
		Wait:      false,

		IdleTimeout: 3 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", GetConfiguration().Redis.Addr)
			if err != nil {
				return nil, err

			}
			_, err = c.Do("AUTH", GetConfiguration().Redis.Password)
			if err != nil {

				return nil, nil
			}
			return c, err
		},
	}

	MyRedis = &RedisStore{pool, defaultExpiration}
}

func (c *RedisStore) RedisGet(key string, value interface{}) error {
	conn := c.pool.Get()
	defer conn.Close()
	reply, err := conn.Do("get", key)
	if err != nil {
		return err
	}
	if reply == nil {
		fmt.Printf("nothing")
		return nil
	}
	bytes, err := redis.Bytes(reply, err)
	return deserialize(bytes, value)

}
func (c *RedisStore) RedisConn() (con *redis.Conn) {
	conn := c.pool.Get()
	return &conn
}

//序列化
func serialize(value interface{}) ([]byte, error) {
	if bytes, ok := value.([]byte); ok {
		return bytes, nil
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(v.Uint(), 10)), nil
	}

	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	if err := encoder.Encode(value); err != nil { //编码
		return nil, err
	}
	return b.Bytes(), nil
}

func (c *RedisStore) RedisSet(key string, value interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()
	return c.invoke(conn.Do, key, value, expires)
}

//添加数据
func (c *RedisStore) RedisAdd(key string, value interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()
	b, err := redis.Bool(conn.Do("exists", key))
	if err != nil {
		return err
	}
	if b { //如果key 已经存在
		return nil
	} else {
		return c.invoke(conn.Do, key, value, expires)
	}

}

//删除数据
func (c *RedisStore) RedisDelete(key string) error {
	conn := c.pool.Get()
	defer conn.Close()
	b, err := redis.Bool(conn.Do("exists", key))
	if err != nil {
		return err
	}
	if !b { //key值不存在
		return nil
	}
	_, err2 := conn.Do("del", key) //删除key值
	return err2
}
func (c *RedisStore) RedisReplace(key string, value interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()
	b, err := redis.Bool(conn.Do("exists", key))
	if err != nil {
		return err
	}
	if !b { //key值不存在
		return nil
	}
	err = c.invoke(conn.Do, key, value, expires)
	if value == nil { //空值不能保存
		return nil
	} else {
		return err
	}
}

func (c *RedisStore) RedisClear() error {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHALL")
	if err != nil {
		return err
	}
	return nil

}

func (c *RedisStore) invoke(f func(string, ...interface{}) (interface{}, error),
	key string, value interface{}, expires time.Duration) error {
	b, err := serialize(value) //序列化操作，序列化可以保存对象
	if err != nil {
		return err
	}
	if expires > 0 {
		_, err := f("setex", key, int32(expires/time.Second), b)
		return err
	} else {
		_, err := f("set", key, b)
		return err
	}
}

//反序列化
func deserialize(byt []byte, ptr interface{}) (err error) {
	if bytes, ok := ptr.(*[]byte); ok {
		*bytes = byt

		return nil
	}

	if v := reflect.ValueOf(ptr); v.Kind() == reflect.Ptr { // 通过反射得到ptr类型，判断ptr是指针类型
		switch p := v.Elem(); p.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: //符号整型
			var i int64
			i, err = strconv.ParseInt(string(byt), 10, 64)
			if err != nil {
				return err
			} else {
				p.SetInt(i)
			}
			return nil

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: //无符号整型
			var i uint64
			i, err = strconv.ParseUint(string(byt), 10, 64)
			if err != nil {
				return err
			} else {
				p.SetUint(i)
			}
			return nil
		}
	}

	b := bytes.NewBuffer(byt)
	decoder := gob.NewDecoder(b)
	if err = decoder.Decode(ptr); err != nil { //解码
		return err
	}
	return nil
}
