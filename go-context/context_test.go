package gocontext

import (
	"fmt"
	"testing"
)

func TestOsExit(t *testing.T) {
	go func() {
		for {
			select {
			case <-WithCancel().Done():
				fmt.Println("---1---")
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-WithCancel().Done():
				fmt.Println("---2---")
				return
			}
		}
	}()

	fmt.Println("---3---")
	<-WithCancel().Done()
	fmt.Println("---4---")
}
