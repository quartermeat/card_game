package singleton

import (
	"log"
	"os"
	"sync"
)

type errorLog struct {
	logger *log.Logger
}

var instance *errorLog
var once sync.Once

// GetInstance returns the singleton instance of errorLog
func GetInstance() *errorLog {
	once.Do(func() {
		file, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open error log file: %s", err)
		}
		instance = &errorLog{
			logger: log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		}
	})
	return instance
}

// LogError logs an error message
func (e *errorLog) LogError(err error) {
	e.logger.Println(err)
}
