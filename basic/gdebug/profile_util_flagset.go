package gdebug

import (
	"flag"
	"github.com/google/pprof/driver"
)

type flagset struct {
	*flag.FlagSet
	args []string
}

func (f *flagset) StringList(name string, def string, usage string) *[]*string {
	return &[]*string{f.FlagSet.String(name, def, usage)}
}

func (f *flagset) ExtraUsage() string {
	return ""
}

func (f *flagset) Parse(usage func()) []string {
	f.FlagSet.Usage = func() {}
	f.FlagSet.Parse(f.args)
	return f.FlagSet.Args()
}

func (f *flagset) AddExtraUsage(string) {
}

func newFlagSet(a ...string) driver.FlagSet {
	return &flagset{flag.NewFlagSet("ppf", flag.ContinueOnError), append(a, "")}
}
