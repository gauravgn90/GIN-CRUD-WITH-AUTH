package logservice

import (
	"log"
	"os"
)

func Info(msg string, args ...interface{}) {
	// Log an info message
	log.SetOutput(os.Stdout)
	log.Printf("\033[1;34m"+msg+"\033[0m", args...)
}
func Error(msg string, args ...interface{}) {
	// Log an info message
	log.SetOutput(os.Stderr)
	log.Printf("\033[1;31m"+msg+"\033[0m", args...)
}
func Warning(msg string, args ...interface{}) {
	// Log an info message
	log.SetOutput(os.Stdout)
	log.Printf("\033[1;33m"+msg+"\033[0m", args...)
}
func DumpAndDie(msg interface{}) {
	// Log an info message
	log.SetOutput(os.Stdout)
	log.Printf("\033[1;33m%+v\033[0m", msg)
	os.Exit(1)
}
