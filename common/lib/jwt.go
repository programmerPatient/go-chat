package lib

import (
	"../../config"
	"encoding/json"
	"time"
)


type Jwt struct {
	Data string
	StartTime int64//开始时间
	EndTime int64 //失效时间
}

func JwtNew(data string) *Jwt{
	nowTime := time.Now().Unix()
	return &Jwt{
		StartTime: nowTime,
		EndTime: nowTime+config.JwtDuration,
		Data:           data,
	}
}

/**
数据加密
 */
func JwtEncode(data string) (string){
	jwts := JwtNew(data)
	str ,_:= json.Marshal(jwts)
	token := EnCryption(string(str))
	return token
}

/**
数据解密
 */
func JwtDecode(data string) (string,bool){
	jwt := &Jwt{}
	nowStr := DeCryption(data)
	_ =json.Unmarshal([]byte(nowStr),jwt)
	if jwt.EndTime < time.Now().Unix() {
		return jwt.Data,false
	}
	return jwt.Data,true
}


