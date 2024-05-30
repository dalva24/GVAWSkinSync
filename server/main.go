package main

import (
	"net.dalva.GvawSkinSync/conf"
	"net.dalva.GvawSkinSync/logger"
)

func init() {
	logger.InitializeLoggerOnce()
	conf.InitializeConfigOnce()
	//data.InitializeDbOnce()
}

func init() {
	refreshAuthCode()
}

func main() {

	go scheduleRefreshAuthCode()

	compileSkins()

	serveAlastor()

}
