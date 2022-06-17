package service

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

type Instance struct {
	log *logrus.Logger
}

func (i Instance) Log() *logrus.Logger {
	return i.log
}

type Target string

type Opts struct {
	Target Target
}

func (o Opts) Validate() error {
	if len(strings.TrimSpace(string(o.Target))) < 1 {
		return errors.New(`target empty (pass via --find="github.com/my-org/dep@v0.0.0-20210615"`)
	}
	return nil
}

type Result struct {
	Query Opts
	Chain []Dep
}

func (res Result) String() string {
	lines := make([]string, 0)
	for i, elm := range res.Chain {
		prefix := strings.Repeat(" ", i)
		lines = append(lines, fmt.Sprintf("%s%s -> %s", prefix, elm.To, elm.From))
	}
	return fmt.Sprintf("locating %s\n", res.Query.Target) + strings.Join(lines, "\n")
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

func (i Instance) ReadStdIn() string {
	scanner := bufio.NewScanner(os.Stdin)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		i.log.Errorf("failed to read std in %s", scanner.Err())
	}
	return strings.Join(lines, "\n")
}

func (i Instance) Scan(goModOutput string, opts Opts) (*Result, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid options %s", err)
	}
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
		Query: opts,
		Chain: chain,
	}, nil
}

func New() *Instance {
	return &Instance{
		log: logrus.New(),
	}
}
