package captcha

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
	log "github.com/wonderivan/logger"
	"image/color"
	config "ticket-service/conf"
	"ticket-service/pkg/utils/redis"
	"time"
)

var Driver base64Captcha.DriverString

var store = base64Captcha.DefaultMemStore

type Data struct {
	CaptchaId string `json:"captcha_id"` //验证码id
	Data      string `json:"data"`       //验证码数据base64类型
	Answer    string `json:"answer"`     //验证码数字
}

//func NewCaptcha() *base64Captcha.DriverDigit {
//	driverConf := config.Conf.DigitDriver
//	return &base64Captcha.DriverDigit{
//		Height:   driverConf.Height,
//		Width:    driverConf.Width,
//		Length:   driverConf.Length,   //验证码长度
//		MaxSkew:  driverConf.MaxSkew,  //倾斜
//		DotCount: driverConf.DotCount, //背景的点数，越大，字体越模糊
//	}
//}

func InitCaptchaDriver() {
	driverConf := config.Conf.DigitDriver
	Driver = base64Captcha.DriverString{
		Height:          driverConf.Height,
		Width:           driverConf.Width,
		Length:          driverConf.Length,     //验证码字符数量
		NoiseCount:      driverConf.NoiseCount, //干扰点数量
		ShowLineOptions: 1 | 4,                 //干扰线
		Source:          "123456789abcdefghijklmnopqrstuvwxyz",
		BgColor:         &color.RGBA{R: 3, G: 102, B: 214, A: 125},
		Fonts:           []string{"wqy-microhei.ttc"},
	}
	Driver.ConvertFonts()
	return
}

// Generate 验证码生成
func Generate() (Data, error) {
	var ret Data
	expire := config.Conf.DigitDriver.ExpireTime
	c := base64Captcha.NewCaptcha(&Driver, store)
	id, b64s, answer, err := c.Generate()
	if err != nil {
		return ret, err
	}
	//存储进redis中,验证码有效时常
	_, err = redis.RedisClient.SetNX(id, answer, time.Minute*time.Duration(expire)).Result()
	if err != nil {
		log.Error("生成验证码.写入redis失败 err:[%v]", err)
		return Data{}, err
	}
	ret.CaptchaId = id
	ret.Data = b64s
	ret.Answer = answer
	fmt.Println(ret.Answer, "验证码answer")
	return ret, nil
}

func Verify(data Data) bool {
	//内存存储验证码
	//return store.Verify(data.CaptchaId, data.Answer, true)
	answer, _ := redis.RedisClient.Get(data.CaptchaId).Result()
	return answer != "" && data.Answer == answer

}
