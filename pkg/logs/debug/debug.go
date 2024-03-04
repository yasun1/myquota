package debug

import (
	"github.com/spf13/pflag"
)

var debug bool

// AddFlag adds the log flag to the given set of command line flags.
func AddFlag(flags *pflag.FlagSet) {
	flags.BoolVarP(
		&debug,
		"debug",
		"d",
		false,
		"If the debug is true, will print out the rich logs",
	)
}

// DebugMode returns the value of debug
func DebugMode() bool {
	return debug
}
