package redis

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

type RedisDb struct {
	C redis.Conn
}

/**
释放该链接
 */
func ( r *RedisDb) Release(){
	r.C.Close()
}

func (r *RedisDb) Set(key string,value interface{}){
	values,_ := json.Marshal(value)
	_,err := r.C.Do("Set",key,values)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (r *RedisDb) Get(key string) (da interface{}){
	res,err := redis.String(r.C.Do("Get",key))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}
func (r *RedisDb) SetHash(key string,field string,value interface{}) {
	_ , err := r.C.Do("HSet",key,field,value)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (r *RedisDb) GetHash(key string,field string) string {
	value , err := redis.String(r.C.Do("HGet",key,field))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return value
}
func (r *RedisDb) HGetAll(key string) interface{} {
	value , err := redis.StringMap(r.C.Do("HGetAll",key))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return value
}

func (r *RedisDb) SetSAdd(key string,value interface{}) bool{
	_ , err := r.C.Do("SAdd",key,value)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
