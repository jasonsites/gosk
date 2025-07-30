package logger

type LoggerFormat struct {
	JSON   string
	Styled string
}

var Format = LoggerFormat{
	JSON:   "json",
	Styled: "styled",
}
