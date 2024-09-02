package gocaptcha

import (
	"context"
	"fmt"
	goredis "github.com/gif-gif/go.io/go-db/go-redis"
	"github.com/mojocn/base64Captcha"
	"time"
)

// 默认内存，分布式用redis等
type Config struct {
	Name        string          `json:"name,optional" yaml:"Name"`
	RedisConfig *goredis.Config `json:"redisConfig,optional" yaml:"RedisConfig"`
}

type configJsonBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

type CaptchaData struct {
	Data      string `json:"data"`
	CaptchaId string `json:"captchaId"`
	Answer    string `json:"answer"`
}

type GoCaptcha struct {
	store base64Captcha.Store //验证码信息自定义存储
}

type RedisStore struct {
	redis   *goredis.GoRedis
	Context context.Context
}

func (r *RedisStore) Set(id string, value string) error {
	r.redis.SetEx(id, value, 10*time.Minute)
	return nil
}

func (r *RedisStore) Get(id string, clear bool) string {
	rst := r.redis.Get(id).Val()
	if clear {
		r.redis.Del(id)
	}
	return rst
}

func (r *RedisStore) Verify(id, answer string, clear bool) bool {
	rst := r.Get(id, clear)
	return rst == answer
}

// new other store
func NewRedis(config goredis.Config) (*GoCaptcha, error) {
	err := goredis.Init(config)
	if err != nil {
		return nil, err
	}

	redis := goredis.GetClient(config.Name)
	if redis == nil {
		return nil, fmt.Errorf("failed to connect to redis %v", config)
	}

	return New(&RedisStore{
		redis:   redis,
		Context: context.Background(),
	}), nil
}

func New(store base64Captcha.Store) *GoCaptcha {
	return &GoCaptcha{
		store: store,
	}
}

func NewDefault() *GoCaptcha {
	return New(base64Captcha.DefaultMemStore)
}

// 返回不同类型的验证码
func (g *GoCaptcha) GetCaptcha(param configJsonBody) (*CaptchaData, error) {
	var driver base64Captcha.Driver

	//create base64 encoding captcha
	switch param.CaptchaType {
	case "audio":
		driver = param.DriverAudio
	case "string":
		driver = param.DriverString.ConvertFonts()
	case "math":
		driver = param.DriverMath.ConvertFonts()
	case "chinese":
		driver = param.DriverChinese.ConvertFonts()
	default:
		driver = param.DriverDigit
	}

	c := base64Captcha.NewCaptcha(driver, g.store)
	id, b64s, answer, err := c.Generate()
	if err != nil {
		return nil, err
	}

	data := CaptchaData{
		Data:      b64s,
		CaptchaId: id,
		Answer:    answer,
	}

	return &data, nil
}

// 验证
func (g *GoCaptcha) CaptchaVerify(id, code string) bool {
	if id == "" || code == "" {
		return false
	}
	return g.store.Verify(id, code, true)
}

func (g *GoCaptcha) DigitCaptcha(width, height, length int) (*CaptchaData, error) {
	if width == 0 {
		width = 240
	}
	if height == 0 {
		height = 80
	}

	if length == 0 {
		length = 4
	}

	var param = configJsonBody{
		CaptchaType: "",
		DriverDigit: &base64Captcha.DriverDigit{
			Length:   length,
			Height:   height,
			Width:    width,
			DotCount: 2,
		},
	}

	data, err := g.GetCaptcha(param)
	if err != nil {
		return nil, fmt.Errorf("DigitCaptcha errr")
	}

	return data, nil
}

func (g *GoCaptcha) StringCaptcha(width, height, length int) (*CaptchaData, error) {
	if width == 0 {
		width = 240
	}
	if height == 0 {
		height = 80
	}
	var param = configJsonBody{
		CaptchaType: "string",
		DriverString: &base64Captcha.DriverString{
			Length:          length,
			Height:          height,
			Width:           width,
			ShowLineOptions: 2,
			NoiseCount:      0,
			Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		},
	}

	data, err := g.GetCaptcha(param)
	if err != nil {
		return nil, fmt.Errorf("StringCaptcha errr")
	}

	return data, nil
}

func (g *GoCaptcha) AudioCaptcha(language string, length int) (*CaptchaData, error) {
	if language == "" {
		language = "123456"
	}
	var param = configJsonBody{
		CaptchaType: "audio",
		DriverAudio: &base64Captcha.DriverAudio{
			Length:   length,
			Language: language,
		},
	}

	data, err := g.GetCaptcha(param)
	if err != nil {
		return nil, fmt.Errorf("AudioCaptcha errr")
	}

	return data, nil
}

// source是中文英文字列表
func (g *GoCaptcha) ChineseCaptcha(width, height, length int, source string) (*CaptchaData, error) {
	if source == "" {
		source = "123456qwertyu你好adfkl在在载在饿工一ioplkjhgfdsazxcvbnm"
	}
	var param = configJsonBody{
		CaptchaType: "chinese",
		DriverChinese: &base64Captcha.DriverChinese{
			Length:          length,
			Height:          height,
			Width:           width,
			ShowLineOptions: 0,
			NoiseCount:      0,
			Source:          source,
		},
	}

	data, err := g.GetCaptcha(param)
	if err != nil {
		return nil, fmt.Errorf("ChineseCaptcha errr")
	}

	return data, nil
}

// 数学计算
func (g *GoCaptcha) MathCaptcha(width, height int) (*CaptchaData, error) {
	if width == 0 {
		width = 240
	}
	if height == 0 {
		height = 80
	}
	var param = configJsonBody{
		CaptchaType: "math",
		DriverMath: &base64Captcha.DriverMath{
			Height:          height,
			Width:           width,
			ShowLineOptions: 0,
			NoiseCount:      0,
		},
	}

	data, err := g.GetCaptcha(param)
	if err != nil {
		return nil, fmt.Errorf("MathCaptcha errr")
	}

	return data, nil
}
