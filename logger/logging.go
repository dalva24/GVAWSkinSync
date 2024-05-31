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

func InitializeLoggerOnceToFile(name string) {
	logInitOnce.Do(func() {

		err := os.Remove(name)
		if err != nil {
			//nop
		}

		runLogFile, _ := os.OpenFile(
			name,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		multi := MultiLevelWriter(os.Stdout, runLogFile)

		//prep logger
		logOut := ConsoleWriter{Out: multi, TimeFormat: time.TimeOnly}
		Log = New(logOut).With().Caller().Timestamp().Logger()
		Log.Debug().Msg("INIT: Initializing logger.Log.")
	})
}

func DumpToFile(filename string, data string) {
	os.WriteFile(filename, []byte(data), 0600)
}
