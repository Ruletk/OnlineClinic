package logging

type LoggerLevel string

const (
	Debug LoggerLevel = "debug"
	Info  LoggerLevel = "info"
	Warn  LoggerLevel = "warn"
	Error LoggerLevel = "error"
	Fatal LoggerLevel = "fatal"
)

type LoggerFormat string

const (
	JSON LoggerFormat = "json"
	Text LoggerFormat = "text"
)
