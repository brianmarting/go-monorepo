package logging

import "github.com/rs/zerolog"

func EnableLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}
