package main

import (
	"fmt"
	"github.com/RobinUS2/godephunter/service"
	"time"
)

func main() {
	svc := service.New()
	if svc == nil {
		panic("nil")
	}
	var res chan string
	go func() {
		res <- svc.ReadStdIn()
	}()
	select {
	case <-res:
		fmt.Println(res)
	case <-time.After(5 * time.Second):
		// timeout
		svc.Log().Error(`failed to read from stdin. use as:` + "\n" + `go mod graph | godephunter --locate="github.com/Route42/golang-commons/httpclient@v0.0.0-20210615"` + "\n\n")
	}
	// @todo read from std in and pass to service, print output
}
