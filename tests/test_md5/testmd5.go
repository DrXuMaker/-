package main

import (
	"crypto/md5"
	"fmt"
)

func str2m5(str string) string {

	str2b := []byte(str)
	ret  :=md5.Sum(str2b)
	retf  :=fmt.Sprintf("%x",ret)
	return retf
}


func main() {
	str := str2m5("hello world")
	fmt.Println(str)

}