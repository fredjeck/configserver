package logging

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

// Creates a new logger based on the $CONFIGSERVER_ENV env variable to decide on the configuration
// If the variable is found AND contains "dev" creates a development logger (human readable output)
// If not generates a production optimized logger logging as JSON
func NewLogger() *zap.Logger {
	env := strings.ToLower(os.Getenv("CONFIGSERVER_ENV"))
	var logger *zap.Logger
	var err error

	if strings.Contains(env, "dev") {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}

	return logger
}
