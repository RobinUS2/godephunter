package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

type Instance struct {
	log *logrus.Logger
}

type Target string

type Opts struct {
	Target Target
}

type Result struct {
	Chain []Dep
}

func (res Result) String() string {
	lines := make([]string, 0)
	for _, elm := range res.Chain {
		lines = append(lines, fmt.Sprintf("%s -> %s", elm.From, elm.To)) // @Todo check formatting
	}
	return strings.Join(lines, "\n")
}

func (i Instance) NewScanOpts() Opts {
	return Opts{}
}

func (i Instance) ScanFile(fileName string, opts Opts) (*Result, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file in %s: %s", fileName, err)
	}
	return i.Scan(string(b), opts)
}

type From string

type To string

type Dep struct {
	From From
	To   To
}

func (i Instance) Scan(goModOutput string, opts Opts) (*Result, error) {
	goModOutput = strings.TrimSpace(goModOutput)
	i.log.Debugf("%s", goModOutput)
	lines := strings.Split(goModOutput, "\n")
	deps := make([]Dep, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		pair := strings.Split(line, " ")
		if len(pair) != 2 {
			continue
		}
		dep := Dep{
			From: From(pair[0]),
			To:   To(pair[1]),
		}
		//i.log.Infof("%s %s", dep.From, dep.To)
		deps = append(deps, dep)
	}
	i.log.Debugf("%+v", deps)
	i.log.Debugf("opts %+v target %s", opts, string(opts.Target))
	chain := make([]Dep, 0)
	var currentTarget = opts.Target
outer:
	for {
		for _, dep := range deps {
			//log.Println(dep.From, dep.To)
			if strings.Contains(string(dep.To), string(currentTarget)) { // @semantic version support for both httpclient@v0.0.0-20210615 or httpclient@v1.0.2 formats
				i.log.Debugf("found %s in %+v", currentTarget, dep)
				chain = append(chain, dep)
				currentTarget = Target(dep.From) // shift to next one
				continue outer
			}
		}
		break outer
	}
	if len(chain) == 0 {
		return nil, fmt.Errorf("%s not found", opts.Target)
	}
	return &Result{
		Chain: chain,
	}, nil
}

func New() *Instance {
	return &Instance{
		log: logrus.New(),
	}
}
