package service_test

import (
	"github.com/RobinUS2/godephunter/service"
	"testing"
)

func TestNew(t *testing.T) {
	svc := service.New()
	if svc == nil {
		t.FailNow()
	}
	opts := svc.NewScanOpts()
	opts.Target = "httpclient@v0.0.0-2021" // anything from 2021 june 15th or older
	res, err := svc.ScanFile("../privatetestdata/sample1.modfile", opts)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%s", res)
}
