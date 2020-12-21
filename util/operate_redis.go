package util

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)



func FlushRedis(){
	//1. 链接到redis
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis.Dial err =", err)
		return
	}
	err1 :=conn.Flush()
	if err1!=nil{
		return
	}
}

func Save2Redis(key string ,key1 string ,value1 string){
	//通过go向redis写入数据和读取数据
	//1. 链接到redis
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis.Dial err =", err)
		return
	}
	defer conn.Close() //关闭...

	//2. 通过go向redis写入数据string [key - val]
	_,err = conn.Do("HSet",key,key1,value1)
	if err != nil {
		fmt.Println("Save2Redis err = ", err)
		return
	}
}

func ReadRedis(key string,key1 string) string {
	//通过go向redis写入数据和读取数据
	//1. 链接到redis
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis.Dial err =", err)
		return "false"
	}
	defer conn.Close() //关闭...

	//2. 通过go向redis读取数据
	r, err := redis.String(conn.Do("HGet", key, key1))
	if err != nil {
		fmt.Println("ReadRedis err = ", err)
		return "false"
	}
	return r
}

func RedisIsExist(key string,key1 string) int64{
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis.Dial err =", err)
		return 10
	}
	defer conn.Close() //关闭...

	//2. 通过go向redis读取数据
	r, err := redis.Int64(conn.Do("HExists", key, key1))
	if err != nil {
		fmt.Println("RedisIsExist err = ", err)
		return 10
	}
	return r
}

func IncRedisValue(key string,key1 string, num int){

	//1. 链接到redis
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis.Dial err =", err)
		return
	}
	defer conn.Close() //关闭...

	//2.HIncrBy
	_,err = conn.Do("HIncrBy",key,key1,num)
	if err != nil {
		fmt.Println("HIncrBy err = ", err)
		return
	}
}