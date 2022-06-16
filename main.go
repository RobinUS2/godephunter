package main

import "github.com/RobinUS2/godephunter/service"

func main() {
	svc := service.New()
	if svc == nil {
		panic("nil")
	}
	// @todo read from std in and pass to service, print output
}
