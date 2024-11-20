package main

import (
	gocontext "github.com/gif-gif/go.io/go-context"
	"github.com/gif-gif/go.io/goio"
	"testing"
)

func TestAdmobApps(t *testing.T) {
	goio.Init(goio.DEVELOPMENT)
	<-gocontext.Cancel().Done()
}
