package main

import (
	"fmt"
	"os"
	"sort"
)

type ArgParser struct {
	opts Options
	groups Groups
	usageFunc func()
}

type Option struct {
	opt string
	longOpt string
	help string
	required bool
	defaultValue interface{}
}

type Group struct {
	name string
	opts Options
}

type Args struct {
	indexMap map[string]int
}

type Options []Option
type Groups []Group

func (o Options) Len() int {
	return len(o)
}

func (o Options) Less(i, j int) bool {
	return o[i].opt < o[j].opt
}

func (o Options) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o Groups) Len() int {
	return len(o)
}

func (o Groups) Less(i, j int) bool {
	return o[i].name < o[j].name
}

func (o Groups) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func NewArgParser() (a *ArgParser) {
	return &ArgParser{
		opts:      make(Options, 0),
		groups:    make(Groups, 0),
		usageFunc: nil,
	}
}

func (a *ArgParser) Usage(call func()) {
	a.usageFunc = call
}

func (a *ArgParser) PrintUsage() {
	if a.usageFunc != nil {
		a.usageFunc()
	} else {
		a.PrintDefaultUsage()
	}
}

func (a *ArgParser) PrintDefaultUsage() {
	sort.Sort(a.opts)
	sort.Sort(a.groups)
	for _, opt := range a.opts {
		a.printOption(opt)
	}
	for _, group := range a.groups {
		fmt.Printf("  %s\n", group.name)
		for _, opt := range group.opts {
			fmt.Print("  ")
			a.printOption(opt)
		}
	}
}

func (a *ArgParser) printOption(opt Option) {
	// TODO: align help message
	fmt.Printf("  -%s", opt.opt)

	if opt.longOpt != "" {
		fmt.Printf(",%s", opt.longOpt)
	}

	fmt.Printf(" %s", opt.help)

	if _, ok := opt.defaultValue.(bool); !ok && !opt.required {
		fmt.Printf(" (default %v)", opt.defaultValue)
	}
	fmt.Printf("\n")
}

func (a *ArgParser) Options(opts ...Option) *ArgParser {
	a.opts = append(a.opts, opts...)
	return a
}

func (a *ArgParser) Group(name string, opts ...Option) *ArgParser {
	a.groups = append(a.groups, Group{name: name, opts: opts})
	return a
}

func (a *ArgParser) Parse() *Args {
	args := &Args{indexMap: make(map[string]int)}
	for idx, arg := range os.Args {
		if arg[0] == '-' {
			arg = arg[1:]
		}
		args.indexMap[arg] = idx
	}
	return args
}

func (a *Args) HasArg(opt string) (ok bool) {
	_, ok = a.indexMap[opt]
	return
}

func (a *Args) AsString(opt string) (string, bool) {
	if idx, ok := a.indexMap[opt]; ok {
		return os.Args[idx], true
	}
	return "", false
}