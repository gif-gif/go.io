package goetcd

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	gocontext "github.com/gif-gif/go.io/go-context"
	goutils "github.com/gif-gif/go.io/go-utils"
)

// 配置结构定义
type TestSt struct {
	Name string `json:"name"`
}

func TestInit(t *testing.T) {
	Init(Config{
		Endpoints: []string{"127.0.0.1:2379"},
		//Username:  "root",
		//Password:  "123456",
	})

	if _, err := Set("/xz/dsp/http-api/192.168.1.101:15001", "192.168.1.101:15001"); err != nil {
		log.Fatalln(err)
	}
	if _, err := Set("/xz/dsp/http-api/192.168.1.101:15002", "192.168.1.101:15002"); err != nil {
		log.Fatalln(err)
	}
	if _, err := SetTTL("/xz/dsp/http-api/192.168.1.101:15003", "192.168.1.101:15003", 3); err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < 2; i++ {
		fmt.Println(GetString("/xz/dsp/http-api/"))
		fmt.Println(GetArray("/xz/dsp/http-api/"))
		fmt.Println(GetMap("/xz/dsp/http-api/"))
		time.Sleep(5 * time.Second)
	}

	Del("/xz/dsp/http-api/192.168.1.101:15002")
}

func TestRegisterService(t *testing.T) {
	Init(Config{
		Endpoints: []string{"127.0.0.1:23790"},
		//Username:  "root",
		//Password:  "123456",
	})

	err := RegisterService("/xz/dsp/http-api/node-1", "192.168.1.101:15002")
	fmt.Println(err)

	<-gocontext.WithCancel().Done()
}

func TestWatch(t *testing.T) {
	Init(Config{
		Endpoints: []string{"127.0.0.1:2379"},
		//Username:  "root",
		//Password:  "123456",
	})

	go func() {
		for i := 0; i < 100; i++ {
			SetTTL(fmt.Sprintf("/xz/dsp/http-api/node-%d", 100), fmt.Sprintf("192.168.1.%d", i), 5)
			time.Sleep(time.Second)
		}
	}()

	ch := Watch("/xz/dsp/http-api")
	for i := range ch {
		fmt.Println(i)
	}
}

// get array value by prefix key
func TestWatch1(t *testing.T) {
	Init(Config{
		Endpoints: []string{"127.0.0.1:2379"},
		Username:  "root",
		//Password:  "123456",
	})

	startTime := time.Now()
	pool := goutils.NewErrorGroup(context.Background(), 1000)
	pool.Submit(func() error {
		for i := 0; i < 100000; i++ {
			SetTTL(fmt.Sprintf("/xz/node-server/node-%d", i), fmt.Sprintf("192.168.1.%d", i), 180)
		}
		return nil
	})

	err := pool.Wait()
	if err != nil {
		log.Fatalln(err)
	}
	elapsed := time.Since(startTime)
	fmt.Printf("SetTTL 10000 nodes cost: %v\n", elapsed)

	startTime = time.Now()
	GetArray("/xz/node-server/")
	elapsed = time.Since(startTime)
	fmt.Printf("GetArray 10000 nodes cost: %v\n", elapsed)
	//for _, v := range ch {
	//	fmt.Println(v)
	//}

}
