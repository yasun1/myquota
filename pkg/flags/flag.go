package flags

import (
	"github.com/spf13/pflag"

	"github/yasun1/myquota/pkg/logs/debug"
)

// AddDebugFlag adds the '--debug' flag to the given set of command line flags.
func AddDebugFlag(fs *pflag.FlagSet) {
	debug.AddFlag(fs)
}
