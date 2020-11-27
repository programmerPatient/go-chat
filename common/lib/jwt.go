package lib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey ,_= ecdsa.GenerateKey(elliptic.P256(),rand.Reader)
//var jwtKey = []byte("-----BEGIN PUBLIC KEY-----MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDdlatRjRjogo3WojgGHFHYLugdUWAY9iR3fy4arWNA1KoS8kVw33cJibXr8bvwUAUparCwlvdbH6dvEOfou0/gCFQsHUfQrSDv+MuSUMAe8jzKE4qW+jK+xQU9a03GUnKHkkle+Q0pX/g6jXZ7r1/xAK5Do2kQ+X5xK9cipRgEKwIDAQAB-----END PUBLIC KEY-----")
//const (
//	jwtKeys =  jwtKey
//)

type Jwt struct {
	Data interface{}
	jwt.StandardClaims
}

func JwtNew(data interface{}) *Jwt{
	return &Jwt{
		Data:           data,
		StandardClaims: jwt.StandardClaims{},
	}
}

/**
数据加密
 */
func(j *Jwt) JwtEncode() (string,bool){
	fmt.Print("机密的的key\n");
	fmt.Printf("%v",jwtKey)
	expireTime := time.Now().Add(7*24*time.Hour)
	jwts := &Jwt{
		Data:          j.Data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),//过期时间
			IssuedAt:  time.Now().Unix(),//颁发时间
			Issuer:    "127.0.0.1",//签名颁发人
			Subject:   "用户信息存贮",//签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256,jwts)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Printf("生成令牌出粗")
		fmt.Println(err)
		return "",false
	}
	return tokenString,true
}

/**
数据解密
 */
func JwtDecode(data string) (interface{},bool){
	if data == "" {
		return "令牌不可为空",false
	}
	claims := &Jwt{}
	token , err := jwt.ParseWithClaims(data,claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey,nil
	})
	if token != nil {
		if claims ,ok := token.Claims.(*Jwt);ok&&token.Valid{
			return claims,true
		}
	}
	fmt.Printf("%v",err)
	return "令牌有误",false
}


