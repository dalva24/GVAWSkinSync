package logger

import (
	. "github.com/rs/zerolog"
	"os"
	"sync"
	"time"
)

var (
	logInitOnce sync.Once
	Log         Logger
)

func init() {
	InitializeLoggerOnce()
}

func SetLevel(lvl string) error {
	level, err := ParseLevel(lvl)
	if err != nil {
		return err
	}
	Log.Level(level)
	SetGlobalLevel(level)
	return nil
}

func InitializeLoggerOnce() {
	logInitOnce.Do(func() {
		//prep logger
		logOut := ConsoleWriter{Out: os.Stdout, TimeFormat: time.TimeOnly}
		Log = New(logOut).With().Caller().Timestamp().Logger()
		Log.Debug().Msg("INIT: Initializing logger.Log.")
	})
}

func DumpToFile(filename string, data string) {
	os.WriteFile(filename, []byte(data), 0600)
}
