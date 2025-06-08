package logger

import (
	"log"
)

func Init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Info(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

func Warn(format string, v ...interface{}) {
	log.Printf("[WARN] "+format, v...)
}

func Error(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

func Fatal(format string, v ...interface{}) {
	log.Fatalf("[FATAL] "+format, v...)
}