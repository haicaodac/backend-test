package library

import (
	"log"
	"os"
)

// Logger ...
type Logger struct {
	file    *os.File
	logger  *log.Logger
	PathLog string
	Type    string
}

// Open ...
func (l *Logger) Open() *log.Logger {
	mode := os.Getenv("MODE")
	if l.PathLog == "" {
		l.PathLog = "logs/" + mode + ".log"
	}
	if l.Type == "" {
		l.Type = "ERROR"
	}

	file, err := os.OpenFile(l.PathLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	l.file = file
	log.SetOutput(l.file)
	l.logger = log.New(l.file, l.Type+": ", log.Ldate|log.Ltime|log.Lshortfile)
	return l.logger
}

// Close ...
func (l *Logger) Close() {
	l.file.Close()
}
