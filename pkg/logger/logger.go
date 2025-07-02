package logger

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Logta/SurveyBot/types"
)

type logger struct {
	*log.Logger
}

// New creates a new logger instance
func New() types.Logger {
	return &logger{
		Logger: log.New(os.Stdout, "[SurveyBot] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *logger) Info(ctx context.Context, msg string, fields ...types.Field) {
	l.logWithFields("INFO", msg, fields...)
}

func (l *logger) Error(ctx context.Context, msg string, err error, fields ...types.Field) {
	allFields := append(fields, types.Field{Key: "error", Value: err.Error()})
	l.logWithFields("ERROR", msg, allFields...)
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...types.Field) {
	l.logWithFields("DEBUG", msg, fields...)
}

func (l *logger) logWithFields(level, msg string, fields ...types.Field) {
	output := fmt.Sprintf("[%s] %s", level, msg)
	if len(fields) > 0 {
		output += " |"
		for _, field := range fields {
			output += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		}
	}
	l.Println(output)
}
