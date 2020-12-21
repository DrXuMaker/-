package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

//第一种方法

////声明一个全局的redisDb变量
//var redisdb *redis.Client
//
////初始化连接
//func InitClient() (err error){
//
//	redisdb = redis.NewClient(&redis.Options{
//		Addr: "localhost:6379",
//		Password: "",
//		DB: 0,
//	})
//	//心跳
//	_,err = redisdb.Ping().Result()
//	if err != nil {
//		return err
//	}
//	fmt.Println("已连接...")
//	return nil
//}
//
//// set/get
//func redisExample(){
//	err := redisdb.Set("mykey",123,0).Err()
//	if err != nil {
//		return
//	}
//	res,_  := redisdb.Get("mykey").Result()
//	fmt.Println(res)
//}


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
		fmt.Println("hset err = ", err)
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
		fmt.Println("hget err = ", err)
		return "false"
	}
	return r
}
func main(){
	con := 0
	Save2Redis("123","1234",strconv.Itoa(con+1))
	res := ReadRedis("123","1234")
	fmt.Println(res)
}