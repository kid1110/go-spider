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

//εΊεε
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
	if err := encoder.Encode(value); err != nil { //ηΌη 
		return nil, err
	}
	return b.Bytes(), nil
}

func (c *RedisStore) RedisSet(key string, value interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()
	return c.invoke(conn.Do, key, value, expires)
}

//ζ·»ε ζ°ζ?
func (c *RedisStore) RedisAdd(key string, value interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()
	b, err := redis.Bool(conn.Do("exists", key))
	if err != nil {
		return err
	}
	if b { //ε¦ζkey ε·²η»ε­ε¨
		return nil
	} else {
		return c.invoke(conn.Do, key, value, expires)
	}

}

//ε ι€ζ°ζ?
func (c *RedisStore) RedisDelete(key string) error {
	conn := c.pool.Get()
	defer conn.Close()
	b, err := redis.Bool(conn.Do("exists", key))
	if err != nil {
		return err
	}
	if !b { //keyεΌδΈε­ε¨
		return nil
	}
	_, err2 := conn.Do("del", key) //ε ι€keyεΌ
	return err2
}
func (c *RedisStore) RedisReplace(key string, value interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	defer conn.Close()
	b, err := redis.Bool(conn.Do("exists", key))
	if err != nil {
		return err
	}
	if !b { //keyεΌδΈε­ε¨
		return nil
	}
	err = c.invoke(conn.Do, key, value, expires)
	if value == nil { //η©ΊεΌδΈθ½δΏε­
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
	b, err := serialize(value) //εΊεεζδ½οΌεΊεεε―δ»₯δΏε­ε―Ήθ±‘
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

//εεΊεε
func deserialize(byt []byte, ptr interface{}) (err error) {
	if bytes, ok := ptr.(*[]byte); ok {
		*bytes = byt

		return nil
	}

	if v := reflect.ValueOf(ptr); v.Kind() == reflect.Ptr { // ιθΏεε°εΎε°ptrη±»εοΌε€ζ­ptrζ―ζιη±»ε
		switch p := v.Elem(); p.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: //η¬¦ε·ζ΄ε
			var i int64
			i, err = strconv.ParseInt(string(byt), 10, 64)
			if err != nil {
				return err
			} else {
				p.SetInt(i)
			}
			return nil

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: //ζ η¬¦ε·ζ΄ε
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
	if err = decoder.Decode(ptr); err != nil { //θ§£η 
		return err
	}
	return nil
}
