package data

import "os"

var (
	SigClose chan os.Signal
)

func init() {
	SigClose = make(chan os.Signal, 1)
}
