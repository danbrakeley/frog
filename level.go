package frog

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
