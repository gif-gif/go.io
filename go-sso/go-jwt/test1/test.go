package main

import "flag"

var gopath string

func init() {
	println("init a")
}

func init() {
	println("init b")
}

func init() {
	println("init c")
	// gopath may be overridden by --gopath flag on command line.
	flag.StringVar(&gopath, "gopath", "/root/go", "override default GOPATH")
}

func main() {
	println("main")
	flag.Parse()
	println(gopath)
}
