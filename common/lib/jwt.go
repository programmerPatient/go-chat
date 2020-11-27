package lib

import (
	"../../config"
	"encoding/json"
	"time"
)

//var jwtKey ,_= ecdsa.GenerateKey(elliptic.P256(),rand.Reader)
//var jwtKey = []byte("-----BEGIN PUBLIC KEY-----MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDdlatRjRjogo3WojgGHFHYLugdUWAY9iR3fy4arWNA1KoS8kVw33cJibXr8bvwUAUparCwlvdbH6dvEOfou0/gCFQsHUfQrSDv+MuSUMAe8jzKE4qW+jK+xQU9a03GUnKHkkle+Q0pX/g6jXZ7r1/xAK5Do2kQ+X5xK9cipRgEKwIDAQAB-----END PUBLIC KEY-----")
//const (
//	jwtKeys =  jwtKey
//)

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


