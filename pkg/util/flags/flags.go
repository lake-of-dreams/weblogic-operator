package flags

import (
	goflag "flag"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

// InitFlags parses, then logs the command line flags.
func InitFlags() {
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	pflag.Parse()
	pflag.VisitAll(func(flag *pflag.Flag) {
		glog.V(4).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
	})
}
