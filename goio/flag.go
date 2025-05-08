package goio

import (
	"flag"
	"fmt"
	"os"
)

/**
  - deploy.sh

	// 定义
	version=v1.0.1
	buildVersion="${version}.$(date +%Y%m%d).$(date +%H%M)"

	// 编译
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X https://github.com/gif-gif/go.io/goio.Version=$buildVersion" -o ss

  - 执行

	./ss -v
*/

var (
	Version     string
	VersionFlag = flag.Bool("v", false, "version")

	HelpFlag = flag.Bool("h", false, "help")
)

func FlagInit() {
	if !flag.Parsed() {
		flag.Parse()
	}

	if *VersionFlag {
		fmt.Println(Version)
		os.Exit(0)
	}

	if *HelpFlag {
		flag.PrintDefaults()
		os.Exit(0)
	}
}
