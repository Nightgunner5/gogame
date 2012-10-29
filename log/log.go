package log

import ("fmt";"log")

func Warning(message string, args ...interface{}) {
	log_("WARNING", fmt.Sprintf(message, args...))
}

func Panic(message string, args ...interface{}) {
	log_("PANIC", fmt.Sprintf(message, args...))
	panic(fmt.Sprintf(message, args...))
}

func log_(tag, message string) {
	// TODO: better logging method
	log.Printf("[%s] %s", tag, message)
}
