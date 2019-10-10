package frog

import "fmt"

type Level byte

const (
	Transient Level = iota // strictly unimportant, ie progress bars, real-time byte counts, estimated time remaining, etc
	Verbose                // debugging info
	Info                   // normal message
	Warning                // something unusual happened
	Error                  // something bad happened
	Fatal                  // stop everything right now

	levelMax
	levelMin Level = 0
)

func init() {
	// ensure all levels have a valid string
	for l := levelMin; l < levelMax; l++ {
		if len(l.String()) == 0 {
			panic(fmt.Errorf("Empty String returned for level %d", int(l)))
		}
	}
}

func (l Level) String() string {
	switch l {
	case Transient:
		return "transient"
	case Verbose:
		return "verbose"
	case Info:
		return "info"
	case Warning:
		return "warning"
	case Error:
		return "error"
	case Fatal:
		return "fatal"
	}
	return ""
}
