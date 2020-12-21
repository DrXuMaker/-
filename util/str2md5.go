package util

import (
	"crypto/md5"
	"fmt"
)

func Str2m5(str string) string {

	str2b := []byte(str)
	ret  :=md5.Sum(str2b)
	retf  :=fmt.Sprintf("%x",ret)
	return retf
}
