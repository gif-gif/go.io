// example of HTTP server that uses the captcha package.
package main

import (
	"context"
	"encoding/json"
	gocaptcha "github.com/gif-gif/go.io/go-captcha"
	goredisc "github.com/gif-gif/go.io/go-db/go-redis/go-redisc"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// configJsonBody json request body.
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

var store = base64Captcha.DefaultMemStore

// base64Captcha create http handler
func generateCaptchaHandler(w http.ResponseWriter, r *http.Request) {
	//parse request parameters
	//decoder := json.NewDecoder(r.Body)
	//var param configJsonBody
	//err := decoder.Decode(&param)
	//if err != nil {
	//	log.Println(err)
	//}

	var param configJsonBody = configJsonBody{
		Id:          "",
		CaptchaType: "audio",
		VerifyValue: "",
		DriverAudio: &base64Captcha.DriverAudio{
			Length:   4,
			Language: "123456",
		},
		DriverString: &base64Captcha.DriverString{
			Length:          4,
			Height:          60,
			Width:           240,
			ShowLineOptions: 2,
			NoiseCount:      0,
			Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		},
		DriverChinese: &base64Captcha.DriverChinese{
			Length:          4,
			Height:          60,
			Width:           240,
			ShowLineOptions: 2,
			NoiseCount:      0,
			Source:          "1234567890qwertyu你好adfkl在在载在饿工一ioplkjhgfdsazxcvbnm",
		},
		DriverMath: &base64Captcha.DriverMath{},
		DriverDigit: &base64Captcha.DriverDigit{
			Length:   4,
			Height:   60,
			Width:    240,
			DotCount: 2,
		},
	}

	defer r.Body.Close()
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
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, answer, err := c.Generate()
	body := map[string]interface{}{"code": 1, "data": b64s, "captchaId": id, "answer": answer, "msg": "success"}
	if err != nil {
		body = map[string]interface{}{"code": 0, "msg": err.Error()}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(body)
}

// base64Captcha verify http handler
func captchaVerifyHandle(w http.ResponseWriter, r *http.Request) {

	//parse request json body
	decoder := json.NewDecoder(r.Body)
	var param configJsonBody
	err := decoder.Decode(&param)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()
	//verify the captcha
	body := map[string]interface{}{"code": 0, "msg": "failed"}
	if store.Verify(param.Id, param.VerifyValue, true) {
		body = map[string]interface{}{"code": 1, "msg": "ok"}
	}

	//set json response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	json.NewEncoder(w).Encode(body)
}

// start a net/http server
func main() {
	config := goredisc.Config{
		Name:     "gocaptcha",
		Addrs:    []string{"127.0.0.1:6379"},
		Password: "",
		DB:       0,
		Prefix:   "gocaptcha",
		AutoPing: true,
	}

	err := gocaptcha.Init(gocaptcha.Config{
		RedisConfig: &config,
	})

	if err != nil {
		golog.Error(err.Error())
		return
	}

	data, err := gocaptcha.Default().DigitCaptcha(0, 0, 0)
	golog.WithTag("data").Info(data)

	simpleServer()
}

func simpleServer() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/getCaptcha", func(c *gin.Context) {
		//generateCaptchaHandler(c.Writer, c.Request)

		//c.JSON(http.StatusOK, goutils.CaptchaGet(240, 60))
		//c.JSON(http.StatusOK, goutils.CaptchaStringGet(240, 60))
		//c.JSON(http.StatusOK, goutils.CaptchaMathGet(240, 60))
		c.JSON(http.StatusOK, goutils.CaptchaAudioGet("123456"))
	})

	r.POST("/api/verifyCaptcha", func(c *gin.Context) {
		goutils.CaptchaVerify("yCjHYaNyAJdVT8yNER6r", "111")
	})

	srv := &http.Server{
		Addr:    ":1000",
		Handler: r,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器 (设置 5 秒的超时时间)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
