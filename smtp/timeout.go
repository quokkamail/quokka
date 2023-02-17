package smtp

import "time"

const (
	DefaultDataBlockTimeout         = 3 * time.Minute
	DefaultDataInitiationTimeout    = 2 * time.Minute
	DefaultDataTerminationTimeout   = 10 * time.Minute
	DefaultInitial220MessageTimeout = 5 * time.Minute
	DefaultMAILCommandTimeout       = 5 * time.Minute
	DefaultRCPTCommandTimeout       = 5 * time.Minute
	DefaultServerTimeout            = 5 * time.Minute
)
