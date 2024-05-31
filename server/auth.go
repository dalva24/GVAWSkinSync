package main

import (
	"context"
	"github.com/spf13/afero"
	"github.com/stephenafamo/kronika"
	"math/rand"
	"net.dalva.GvawSkinSync/data"
	"net.dalva.GvawSkinSync/logger"
	"time"
)

var authCode string
var passwords []string

func scheduleRefreshAuthCode() {
	ctx := context.Background()

	start, err := time.Parse(
		"2006-01-02 15:04:05",
		"2019-09-17 21:00:00",
	) // is a tuesday
	if err != nil {
		panic(err)
	}

	interval := time.Hour * 24 // 1 week

	for range kronika.Every(ctx, start, interval) {
		refreshAuthCode()
	}
}

func refreshAuthCode() {
	if passwords == nil || len(passwords) == 0 {
		passwords = data.LoadLines("passwords.txt", afero.NewBasePathFs(afero.NewOsFs(), "./res/"))
	}
	authCode = getRandomAuthCode()
	logger.Log.Info().Str("code", authCode).Msg("refreshAuthCode")
}

func getRandomAuthCode() string {
	length := len(passwords)
	return passwords[rand.Int()%length]
}
