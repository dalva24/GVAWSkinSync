package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"log"
	"net.dalva.GvawSkinSync/Alastor"
	"os"
	"os/user"
	"strings"
)

func DisableInputs() {
	syncButton.Disable()
	authCode.Disable()
	syncServer.Disable()
	savedGamesDir.Disable()
	persSkin.Disable()
	skinSize.Disable()
	persSize.Disable()
}

func EnableInputs() {
	syncButton.Enable()
	authCode.Enable()
	syncServer.Enable()
	savedGamesDir.Enable()
	persSkin.Enable()
	skinSize.Enable()
	persSize.Enable()
}

func RefreshStatus() {
	statusMajor.Refresh()
	statusMinor.Refresh()
}

func UpdateStatus(status string) {

}

func browseDcs() {
	userObj, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// test multiple possible combinations
	dcsDir := ""
	dcsExist, _ := exists(userObj.HomeDir + "\\Saved Games\\DCS\\Config\\options.lua")
	dcsOBExist, _ := exists(userObj.HomeDir + "\\Saved Games\\DCS.openbeta\\Config\\options.lua")
	if dcsOBExist && dcsExist {
		dcsOptions, _ := os.Stat(userObj.HomeDir + "\\Saved Games\\DCS\\Config\\options.lua")
		dcsOBOptions, _ := os.Stat(userObj.HomeDir + "\\Saved Games\\DCS.openbeta\\Config\\options.lua")
		if dcsOptions.ModTime().After(dcsOBOptions.ModTime()) {
			dcsDir = userObj.HomeDir + "\\Saved Games\\DCS\\"
		} else {
			dcsDir = userObj.HomeDir + "\\Saved Games\\DCS.openbeta\\"
		}
	} else if dcsExist {
		dcsDir = userObj.HomeDir + "\\Saved Games\\DCS\\"
	} else if dcsOBExist {
		dcsDir = userObj.HomeDir + "\\Saved Games\\DCS.openbeta\\"
	} else {
		dcsDir = "not found!"
	}

	savedGamesDir.Text = dcsDir
	savedGamesDir.Refresh()
	/*
		w := fyneApp.NewWindow("Add File")
		d := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
			time.Sleep(time.Second) // this is a workaround for the second panic
			w.Close()
		}, w)
		d.Show()
		w.Resize(fyne.NewSize(1000, 500))
		w.Show()
	*/
}

func sync() {

	DisableInputs()

	addressPort = "gvaw.web.id:24003"
	if len(strings.Split(syncServer.Text, ":")) == 2 {
		addressPort = syncServer.Text
	} else if len(strings.Split(syncServer.Text, ":")) == 1 {
		addressPort = syncServer.Text + ":24003"
	}

	log.Println("SYNC" + addressPort)

	err := alastor.TestConnection(addressPort, authCode.Text)
	if err != nil {
		ShowInfo(err.Title, err.Err)
		EnableInputs()
	}

	log.Println("Conn OK" + addressPort)

	pStatus = statusConnected
	RefreshStatus()

	skinSync()

	EnableInputs()

}

func ShowUI() {
	fyneApp = app.New()
	window = fyneApp.NewWindow("GVAW SkinSync by Dalva")
	window.Resize(fyne.NewSize(400, 600))
	window.SetFixedSize(true)
	fyneApp.Settings().SetTheme(theme.DarkTheme())
	window.SetContent(initElements())
	browseDcs()
	window.ShowAndRun()
}

func ShowInfo(title string, content string) {
	wrapped := wordWrap(content, 45)
	dialog.ShowInformation(title, wrapped, window)
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func wordWrap(text string, lineWidth int) string {
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return text
	}
	wrapped := words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}

	return wrapped

}