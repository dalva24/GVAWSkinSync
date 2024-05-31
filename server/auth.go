package main

import (
	"context"
	"github.com/spf13/afero"
	"github.com/stephenafamo/kronika"
	"math/rand"
	"net.dalva.GvawSkinSync/data"
	"net.dalva.GvawSkinSync/logger"
	"strings"
	"time"
)

var authCode string
var passwords []string
var blockedIPs []ipFails

type ipFails struct {
	ip    string
	fails int
}

func addIpFail(ipPort string) {
	fields := strings.Split(ipPort, ":")
	ip := strings.Join(fields[0:len(fields)-1], ":")
	found := false
	for i := range blockedIPs {
		if blockedIPs[i].ip == ip {
			blockedIPs[i].fails++
			logger.Log.Warn().Str("ip", ip).Int("fails", blockedIPs[i].fails).Msg("IP fails")
			found = true
			break
		}
	}
	if !found {
		blockedIPs = append(blockedIPs, ipFails{ip: ip, fails: 1})
		logger.Log.Warn().Str("ip", ip).Int("fails", 1).Msg("IP fails")
	}
}

func isIpBlocked(ipPort string) bool {
	fields := strings.Split(ipPort, ":")
	ip := strings.Join(fields[0:len(fields)-1], ":")
	found := false
	fails := 0
	for i := range blockedIPs {
		if blockedIPs[i].ip == ip {
			found = true
			fails = blockedIPs[i].fails
			break
		}
	}
	if found && fails >= 3 {
		logger.Log.Warn().Str("ip", ip).Int("fails", fails).Msg("IP blocked")
		return true
	} else {
		return false
	}
}

func resetIpBlock(ipPort string) {
	fields := strings.Split(ipPort, ":")
	ip := strings.Join(fields[0:len(fields)-1], ":")
	for i := range blockedIPs {
		if blockedIPs[i].ip == ip {
			blockedIPs[i].fails = 0
			break
		}
	}
}

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
