package main

import (
	"flag"
	"github.com/RobinUS2/godephunter/service"
	"time"
)

func main() {
	svc := service.New()
	if svc == nil {
		panic("nil")
	}

	var target string
	flag.StringVar(&target, "find", "", "the name of the dep to locate (e.g. github.com/Route42/golang-commons/httpclient@v0.0.0-20210615)")
	flag.Parse()

	var bodyChan = make(chan string, 1)
	go func() {
		bodyChan <- svc.ReadStdIn()
	}()
	select {
	case body := <-bodyChan:
		opts := svc.NewScanOpts()
		opts.Target = service.Target(target)
		res, err := svc.Scan(body, opts)
		if err != nil {
			svc.Log().Error(err)
			return
		}
		svc.Log().Info(res.String())
	case <-time.After(5 * time.Second):
		// timeout
		svc.Log().Error(`failed to read from stdin. use as:` + "\n" + `go mod graph | godephunter --find="github.com/Route42/golang-commons/httpclient@v0.0.0-20210615"` + "\n\n")
	}
}
