package util

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

type Captcha struct {
	Id string
	BS64 string
	Code int
}


var Store = base64Captcha.DefaultMemStore

func GetCaptcha() (id string, bs64 string, err error) {
	rgbaColor := color.RGBA{0,0,0,0}
	fonts := []string{"wqy-microhei.ttc"}
	driver := base64Captcha.NewDriverMath(50,140,0,0,&rgbaColor,fonts)


	captcha := base64Captcha.NewCaptcha(driver,Store)
	id, bs64, err = captcha.Generate()
	return id, bs64, err
}

func VerifyCaptcha(id string,ret string) bool {
	f :=Store.Verify(id,ret,true)
	return f
}