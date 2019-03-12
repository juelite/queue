package service

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

type RedisTool struct{}

var client redis.Conn

func init() {
	client = getRedisConn()
}

/**
 * 获取redis链接句柄
 */
func getRedisConn() redis.Conn {
	var client redis.Conn
	host := beego.AppConfig.String("redis_host")
	pass := beego.AppConfig.String("redis_password")

	client , err :=  redis.Dial("tcp", host)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	if pass != "" {
		_ , err = client.Do("AUTH", pass)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
	}
	return client
}

func (r *RedisTool) Set(key, val string) error {
	_ , err := client.Do("SET" , key , val)
	return err
}

func (r *RedisTool) SetEx(key, val string, exp int64) error {
	_ , err := client.Do("SETEX" , key , exp , val)
	return err
}

func (r *RedisTool) Get(key string) (map[string]string , error) {
	rsp := make(map[string]string)
	res , err := client.Do("GET" , key)
	rsp[key] = string(res.([]byte))
	return rsp , err
}

func (r *RedisTool) LPush(key, val string) error {
	_, err := client.Do("lpush", key, val)
	return err
}

func (r *RedisTool) LPop(key string) (interface{}, error) {
	//rsp := make(map[string]string)
	res, err := client.Do("lpop", key)
	return res, err
}

func (r *RedisTool) RPush(key, val string) error {
	_, err := client.Do("rpush", key, val)
	return err
}

func (r *RedisTool) RPop(key string) (interface{}, error) {
	res, err := client.Do("rpop", key)
	return res, err
}

func (r *RedisTool) Del(key string) (error) {
	_ , err := client.Do("DEL" , key)
	return err
}