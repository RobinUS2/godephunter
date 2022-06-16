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
}
