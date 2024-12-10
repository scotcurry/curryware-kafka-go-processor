package logging

//import (
//	"log"
//	"os"
//)
//
//type Logger struct {
//	logger *log.Logger
//}
//
//func NewLogger(prefix string) *Logger {
//	return &Logger{
//		logger: log.New(os.Stdout, prefix, log.Ldate|log.Ltime|log.Lshortfile),
//	}
//}
//
//func (l *Logger) Info(msg string) {
//	l.logger.Println("INFO:", msg)
//}
//
//func (l *Logger) Warn(msg string) {
//	l.logger.Println("WARN:", msg)
//}
//
//func (l *Logger) Error(msg string) {
//	l.logger.Println("ERROR:", msg)
//}
