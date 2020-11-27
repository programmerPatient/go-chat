package lib

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

//遍历文件夹
func FindDir(dir string) []string{

	fileInfo ,err := ioutil.ReadDir(dir)
	fmt.Println(fileInfo)
	if err != nil {
		panic(err)
	}
	lens := len(fileInfo)
	files := make([]string,lens)
	fmt.Println(len(files))
	 i := 0
	//遍历这个文件夹
	for _,fi := range fileInfo {
		//判断是不是目录
		if !fi.IsDir() {
			files[i] = dir + "/" +fi.Name()
			i++
		}
	}
	return files
}

func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}



/**
数据加密
 */
func GetToken(data interface{}){

}
