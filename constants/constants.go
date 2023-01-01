package constants

const (
	NAME    string = "R6Influx"
	VERSION string = "v0.5.0"
)

// 0 = Debug, 1 = Info etc.
// Is set to Info in goreleaser ldflags. Needs to be var of type string for ldflags to work.
var LOG_LEVEL string = "0"
